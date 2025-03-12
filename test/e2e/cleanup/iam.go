package cleanup

import (
	"context"
	"fmt"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/iam"
	"github.com/aws/aws-sdk-go-v2/service/iam/types"
	"github.com/go-logr/logr"

	"github.com/aws/eks-hybrid/test/e2e/constants"
	"github.com/aws/eks-hybrid/test/e2e/errors"
)

type IAMCleaner struct {
	iamClient *iam.Client
	logger    logr.Logger
}

func NewIAMCleaner(iamClient *iam.Client, logger logr.Logger) *IAMCleaner {
	return &IAMCleaner{
		iamClient: iamClient,
		logger:    logger,
	}
}

// TODO: introduce use of role prefix to clean this up

func (c *IAMCleaner) ListRoles(ctx context.Context, filterInput FilterInput) ([]string, error) {
	var roles []string

	// list-roles does not allow filtering by tags so we have to pull them all
	// We have the role =* checks to try and limit which roles we bother checking tags for
	// but we only delete those with the e2e cluster tag
	paginator := iam.NewListRolesPaginator(c.iamClient, &iam.ListRolesInput{})

	for paginator.HasMorePages() {
		page, err := paginator.NextPage(ctx)
		if err != nil {
			return nil, fmt.Errorf("listing IAM roles: %w", err)
		}
		for _, role := range page.Roles {
			if !strings.HasPrefix(*role.RoleName, constants.TestCredentialsStackNamePrefix) {
				continue
			}
			paginator := iam.NewListRoleTagsPaginator(c.iamClient, &iam.ListRoleTagsInput{
				RoleName: role.RoleName,
			})
			var tags []types.Tag
			for paginator.HasMorePages() {
				page, err := paginator.NextPage(ctx)
				if err != nil && errors.IsType(err, &types.NoSuchEntityException{}) {
					// skipping log since we are possiblying check rolesd we do not
					// intend to delete
					continue
				}
				if err != nil {
					return nil, fmt.Errorf("listing IAM role tags: %w", err)
				}
				tags = append(tags, page.Tags...)
			}
			role.Tags = tags
			if shouldDeleteRole(role, filterInput) {
				roles = append(roles, *role.RoleName)
			}
		}
	}
	return roles, nil
}

func shouldDeleteRole(role types.Role, input FilterInput) bool {
	var tags []Tag
	for _, tag := range role.Tags {
		tags = append(tags, Tag{
			Key:   *tag.Key,
			Value: *tag.Value,
		})
	}
	resource := ResourceWithTags{
		ID:           *role.RoleName,
		CreationTime: aws.ToTime(role.CreateDate),
		Tags:         tags,
	}
	return shouldDeleteResource(resource, input)
}

func (c *IAMCleaner) DeleteRole(ctx context.Context, roleName string) error {
	if err := c.detachRoleFromInstanceProfiles(ctx, roleName); err != nil {
		return fmt.Errorf("detaching role from instance profiles: %w", err)
	}

	if err := c.detachManagedPoliciesFromRole(ctx, roleName); err != nil {
		return fmt.Errorf("detaching managed policies from role: %w", err)
	}

	if err := c.deleteInlinePoliciesFromRole(ctx, roleName); err != nil {
		return fmt.Errorf("deleting inline policies from role: %w", err)
	}

	return c.deleteRoleEntity(ctx, roleName)
}

func (c *IAMCleaner) detachRoleFromInstanceProfiles(ctx context.Context, roleName string) error {
	instanceProfiles, err := c.iamClient.ListInstanceProfilesForRole(ctx, &iam.ListInstanceProfilesForRoleInput{
		RoleName: aws.String(roleName),
	})
	if err != nil && errors.IsType(err, &types.NoSuchEntityException{}) {
		c.logger.Info("IAM role already deleted", "roleName", roleName)
		return nil
	}
	if err != nil {
		return fmt.Errorf("listing IAM instance profiles for role %s: %w", roleName, err)
	}

	for _, instanceProfile := range instanceProfiles.InstanceProfiles {
		_, err := c.iamClient.RemoveRoleFromInstanceProfile(ctx, &iam.RemoveRoleFromInstanceProfileInput{
			InstanceProfileName: instanceProfile.InstanceProfileName,
			RoleName:            aws.String(roleName),
		})
		if err != nil {
			return fmt.Errorf("removing role from instance profile %s: %w", *instanceProfile.InstanceProfileName, err)
		}
	}

	return nil
}

func (c *IAMCleaner) detachManagedPoliciesFromRole(ctx context.Context, roleName string) error {
	paginator := iam.NewListAttachedRolePoliciesPaginator(c.iamClient, &iam.ListAttachedRolePoliciesInput{
		RoleName: aws.String(roleName),
	})

	for paginator.HasMorePages() {
		page, err := paginator.NextPage(ctx)
		if err != nil && errors.IsType(err, &types.NoSuchEntityException{}) {
			c.logger.Info("IAM role already deleted", "roleName", roleName)
			return nil
		}
		if err != nil {
			return fmt.Errorf("listing attached policies for role %s: %w", roleName, err)
		}

		for _, policy := range page.AttachedPolicies {
			_, err := c.iamClient.DetachRolePolicy(ctx, &iam.DetachRolePolicyInput{
				RoleName:  aws.String(roleName),
				PolicyArn: policy.PolicyArn,
			})
			if err != nil && !errors.IsType(err, &types.NoSuchEntityException{}) {
				return fmt.Errorf("detaching policy %s from role %s: %w", *policy.PolicyArn, roleName, err)
			}
		}
	}

	return nil
}

func (c *IAMCleaner) deleteInlinePoliciesFromRole(ctx context.Context, roleName string) error {
	paginator := iam.NewListRolePoliciesPaginator(c.iamClient, &iam.ListRolePoliciesInput{
		RoleName: aws.String(roleName),
	})

	for paginator.HasMorePages() {
		page, err := paginator.NextPage(ctx)
		if err != nil && errors.IsType(err, &types.NoSuchEntityException{}) {
			c.logger.Info("IAM role already deleted", "roleName", roleName)
			return nil
		}
		if err != nil {
			return fmt.Errorf("listing inline policies for role %s: %w", roleName, err)
		}

		for _, policyName := range page.PolicyNames {
			_, err := c.iamClient.DeleteRolePolicy(ctx, &iam.DeleteRolePolicyInput{
				RoleName:   aws.String(roleName),
				PolicyName: aws.String(policyName),
			})
			if err != nil && !errors.IsType(err, &types.NoSuchEntityException{}) {
				return fmt.Errorf("deleting inline policy %s from role %s: %w", policyName, roleName, err)
			}
		}
	}

	return nil
}

func (c *IAMCleaner) deleteRoleEntity(ctx context.Context, roleName string) error {
	_, err := c.iamClient.DeleteRole(ctx, &iam.DeleteRoleInput{
		RoleName: aws.String(roleName),
	})
	if err != nil && errors.IsType(err, &types.NoSuchEntityException{}) {
		c.logger.Info("IAM role already deleted", "roleName", roleName)
		return nil
	}
	if err != nil {
		return fmt.Errorf("deleting IAM role %s: %w", roleName, err)
	}
	c.logger.Info("Deleted IAM role", "roleName", roleName)
	return nil
}

func (c *IAMCleaner) ListInstanceProfiles(ctx context.Context, filterInput FilterInput) ([]string, error) {
	var instanceProfiles []string

	// list-instance-profiles does not allow filtering by tags so we have to pull them all
	// We have the role =* checks to try and limit which roles we bother checking tags for
	// but we only delete those with the e2e cluster tag
	paginator := iam.NewListInstanceProfilesPaginator(c.iamClient, &iam.ListInstanceProfilesInput{})

	for paginator.HasMorePages() {
		page, err := paginator.NextPage(ctx)
		if err != nil {
			return nil, fmt.Errorf("listing IAM instance profiles: %w", err)
		}
		for _, profile := range page.InstanceProfiles {
			if !strings.HasPrefix(*profile.InstanceProfileName, constants.TestCredentialsStackNamePrefix) {
				continue
			}
			instanceProfile, err := c.iamClient.GetInstanceProfile(ctx, &iam.GetInstanceProfileInput{
				InstanceProfileName: profile.InstanceProfileName,
			})
			if err != nil && errors.IsType(err, &types.NoSuchEntityException{}) {
				// skipping log since we are possiblying checking instance profiles we do not
				// intend to delete
				continue
			}
			if err != nil {
				return nil, fmt.Errorf("describing IAM instance profile %s: %w", *profile.InstanceProfileName, err)
			}

			if shouldDeleteInstanceProfile(instanceProfile.InstanceProfile, filterInput) {
				instanceProfiles = append(instanceProfiles, *profile.InstanceProfileName)
			}
		}
	}
	return instanceProfiles, nil
}

func shouldDeleteInstanceProfile(profile *types.InstanceProfile, input FilterInput) bool {
	var tags []Tag
	for _, tag := range profile.Tags {
		tags = append(tags, Tag{
			Key:   *tag.Key,
			Value: *tag.Value,
		})
	}
	resource := ResourceWithTags{
		ID:           *profile.InstanceProfileName,
		CreationTime: aws.ToTime(profile.CreateDate),
		Tags:         tags,
	}
	return shouldDeleteResource(resource, input)
}

func (c *IAMCleaner) ListRolesForInstanceProfile(ctx context.Context, profileName string) ([]string, error) {
	roles := []string{}

	describeInput := &iam.GetInstanceProfileInput{
		InstanceProfileName: aws.String(profileName),
	}
	profile, err := c.iamClient.GetInstanceProfile(ctx, describeInput)
	if err != nil && errors.IsType(err, &types.NoSuchEntityException{}) {
		c.logger.Info("IAM instance profile already deleted", "profileName", profileName)
		return roles, nil
	}
	if err != nil {
		return nil, fmt.Errorf("describing IAM instance profile %s: %w", profileName, err)
	}

	for _, role := range profile.InstanceProfile.Roles {
		roles = append(roles, *role.RoleName)
	}
	return roles, nil
}

func (c *IAMCleaner) RemoveRolesFromInstanceProfile(ctx context.Context, roles []string, profileName string) error {
	for _, role := range roles {
		_, err := c.iamClient.RemoveRoleFromInstanceProfile(ctx, &iam.RemoveRoleFromInstanceProfileInput{
			InstanceProfileName: aws.String(profileName),
			RoleName:            aws.String(role),
		})
		if err != nil {
			return fmt.Errorf("removing role from instance profile %s: %w", role, err)
		}
	}
	return nil
}

func (c *IAMCleaner) DeleteInstanceProfile(ctx context.Context, profileName string) error {
	_, err := c.iamClient.DeleteInstanceProfile(ctx, &iam.DeleteInstanceProfileInput{
		InstanceProfileName: aws.String(profileName),
	})
	if err != nil && errors.IsType(err, &types.NoSuchEntityException{}) {
		c.logger.Info("IAM instance profile already deleted", "profileName", profileName)
		return nil
	}
	if err != nil {
		return fmt.Errorf("deleting IAM instance profile %s: %w", profileName, err)
	}
	c.logger.Info("Deleted IAM instance profile", "profileName", profileName)
	return nil
}
