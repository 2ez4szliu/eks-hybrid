package cluster

import (
	"context"
	_ "embed"
	"fmt"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/cloudformation"
	"github.com/aws/aws-sdk-go-v2/service/cloudformation/types"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	"github.com/aws/aws-sdk-go-v2/service/ssm"
	ssmTypes "github.com/aws/aws-sdk-go-v2/service/ssm/types"
	"github.com/go-logr/logr"
	"k8s.io/apimachinery/pkg/util/wait"

	"github.com/aws/eks-hybrid/test/e2e/constants"
	"github.com/aws/eks-hybrid/test/e2e/errors"
	"github.com/aws/eks-hybrid/test/e2e/peered"
	e2eSSM "github.com/aws/eks-hybrid/test/e2e/ssm"
)

//go:embed cfn-templates/setup-cfn.yaml
var setupTemplateBody []byte

const (
	stackWaitTimeout  = 5 * time.Minute
	stackWaitInterval = 10 * time.Second
)

type vpcConfig struct {
	vpcID         string
	publicSubnet  string
	privateSubnet string
	securityGroup string
}

type resourcesStackOutput struct {
	clusterRole      string
	clusterVpcConfig vpcConfig
}

type stack struct {
	logger    logr.Logger
	cfn       *cloudformation.Client
	ssmClient *ssm.Client
	ec2Client *ec2.Client
}

func (s *stack) deploy(ctx context.Context, test TestResources) (*resourcesStackOutput, error) {
	stackName := stackName(test.ClusterName)
	resp, err := s.cfn.DescribeStacks(ctx, &cloudformation.DescribeStacksInput{
		StackName: aws.String(stackName),
	})
	if err != nil && !errors.IsType(err, &types.StackInstanceNotFoundException{}) {
		return nil, fmt.Errorf("looking for hybrid nodes cfn stack: %w", err)
	}

	params := []types.Parameter{
		{
			ParameterKey:   aws.String("ClusterName"),
			ParameterValue: aws.String(test.ClusterName),
		},
		{
			ParameterKey:   aws.String("ClusterRegion"),
			ParameterValue: aws.String(test.ClusterRegion),
		},
		{
			ParameterKey:   aws.String("ClusterVPCCidr"),
			ParameterValue: aws.String(test.ClusterNetwork.VpcCidr),
		},
		{
			ParameterKey:   aws.String("ClusterPublicSubnetCidr"),
			ParameterValue: aws.String(test.ClusterNetwork.PublicSubnetCidr),
		},
		{
			ParameterKey:   aws.String("ClusterPrivateSubnetCidr"),
			ParameterValue: aws.String(test.ClusterNetwork.PrivateSubnetCidr),
		},
		{
			ParameterKey:   aws.String("HybridNodeVPCCidr"),
			ParameterValue: aws.String(test.HybridNetwork.VpcCidr),
		},
		{
			ParameterKey:   aws.String("HybridNodePublicSubnetCidr"),
			ParameterValue: aws.String(test.HybridNetwork.PublicSubnetCidr),
		},
		{
			ParameterKey:   aws.String("HybridNodePrivateSubnetCidr"),
			ParameterValue: aws.String(test.HybridNetwork.PrivateSubnetCidr),
		},
		{
			ParameterKey:   aws.String("HybridNodePodCidr"),
			ParameterValue: aws.String(test.HybridNetwork.PodCidr),
		},
		{
			ParameterKey:   aws.String("TestClusterTagKey"),
			ParameterValue: aws.String(constants.TestClusterTagKey),
		},
	}

	if resp == nil || resp.Stacks == nil {
		s.logger.Info("Creating hybrid nodes setup stack", "stackName", stackName)
		_, err = s.cfn.CreateStack(ctx, &cloudformation.CreateStackInput{
			StackName:    aws.String(stackName),
			TemplateBody: aws.String(string(setupTemplateBody)),
			Parameters:   params,
			Capabilities: []types.Capability{
				types.CapabilityCapabilityIam,
				types.CapabilityCapabilityNamedIam,
			},
			Tags: []types.Tag{{
				Key:   aws.String(constants.TestClusterTagKey),
				Value: aws.String(test.ClusterName),
			}},
		})
		if err != nil {
			return nil, fmt.Errorf("creating hybrid nodes setup cfn stack: %w", err)
		}

		s.logger.Info("Waiting for hybrid nodes setup stack to be created", "stackName", stackName)
		err = waitForStackOperation(ctx, s.cfn, stackName)
		if err != nil {
			return nil, fmt.Errorf("waiting for hybrid nodes setup cfn stack: %w", err)
		}
	} else if resp.Stacks[0].StackStatus == types.StackStatusCreateInProgress {
		s.logger.Info("Waiting for hybrid nodes setup stack to be created", "stackName", stackName)
		err = waitForStackOperation(ctx, s.cfn, stackName)
		if err != nil {
			return nil, fmt.Errorf("waiting for hybrid nodes setup cfn stack: %w", err)
		}
	} else {
		s.logger.Info("Updating hybrid nodes setup stack", "stackName", stackName)
		_, err = s.cfn.UpdateStack(ctx, &cloudformation.UpdateStackInput{
			StackName: aws.String(stackName),
			Capabilities: []types.Capability{
				types.CapabilityCapabilityIam,
				types.CapabilityCapabilityNamedIam,
			},
			TemplateBody: aws.String(string(setupTemplateBody)),
			Parameters:   params,
		})
		if err != nil {
			// Handle "No updates are to be performed"
			if strings.Contains(err.Error(), "No updates are to be performed") {
				s.logger.Info("No changes detected in the stack, no update necessary.")
			} else {
				return nil, fmt.Errorf("updating CloudFormation stack: %w", err)
			}
		}

		s.logger.Info("Waiting for hybrid nodes setup stack to be updated", "stackName", stackName)
		err = waitForStackOperation(ctx, s.cfn, stackName)
		if err != nil {
			return nil, fmt.Errorf("waiting for hybrid nodes setup cfn stack: %w", err)
		}
	}

	resp, err = s.cfn.DescribeStacks(ctx, &cloudformation.DescribeStacksInput{
		StackName: aws.String(stackName),
	})
	if err != nil {
		return nil, fmt.Errorf("describing hybrid nodes setup cfn stack: %w", err)
	}

	result := &resourcesStackOutput{}
	// extract relevant stack outputs
	for _, output := range resp.Stacks[0].Outputs {
		switch aws.ToString(output.OutputKey) {
		case "ClusterRole":
			result.clusterRole = *output.OutputValue
		case "ClusterVPC":
			result.clusterVpcConfig.vpcID = *output.OutputValue
		case "ClusterVPCPublicSubnet":
			result.clusterVpcConfig.publicSubnet = *output.OutputValue
		case "ClusterVPCPrivateSubnet":
			result.clusterVpcConfig.privateSubnet = *output.OutputValue
		case "ClusterSecurityGroup":
			result.clusterVpcConfig.securityGroup = *output.OutputValue
		}
	}

	// We explictly fetch the keypair/jumpbox instead of relying on outputs
	// from the cfn stack. We have seen cases across different regions
	// where the outputs have been unreliable
	keyPair, err := peered.KeyPair(ctx, s.ec2Client, test.ClusterName)
	if err != nil {
		return nil, err
	}
	_, err = s.ssmClient.AddTagsToResource(ctx, &ssm.AddTagsToResourceInput{
		ResourceId:   aws.String("/ec2/keypair/" + *keyPair.KeyPairId),
		ResourceType: ssmTypes.ResourceTypeForTaggingParameter,
		Tags: []ssmTypes.Tag{
			{
				Key:   aws.String(constants.TestClusterTagKey),
				Value: aws.String(test.ClusterName),
			},
		},
	})
	if err != nil {
		return nil, fmt.Errorf("tagging private key ssm parameter: %w", err)
	}

	jumpbox, err := peered.JumpboxInstance(ctx, s.ec2Client, test.ClusterName)
	if err != nil {
		return nil, err
	}

	if err := e2eSSM.WaitForInstance(ctx, s.ssmClient, *jumpbox.InstanceId, s.logger); err != nil {
		return nil, err
	}

	command := "/root/download-private-key.sh"
	output, err := e2eSSM.RunCommand(ctx, s.ssmClient, *jumpbox.InstanceId, command, s.logger)
	if err != nil {
		return nil, fmt.Errorf("jumpbox getting private key from ssm: %w", err)
	}
	if output.Status != "Success" {
		return nil, fmt.Errorf("jumpbox getting private key from ssm")
	}

	s.logger.Info("E2E resources stack deployed successfully", "stackName", stackName)
	return result, nil
}

func stackName(clusterName string) string {
	return fmt.Sprintf("EKSHybridCI-Arch-%s", clusterName)
}

func waitForStackOperation(ctx context.Context, client *cloudformation.Client, stackName string) error {
	err := wait.PollUntilContextTimeout(ctx, stackWaitInterval, stackWaitTimeout, true, func(ctx context.Context) (bool, error) {
		stackOutput, err := client.DescribeStacks(ctx, &cloudformation.DescribeStacksInput{
			StackName: aws.String(stackName),
		})
		if err != nil {
			if errors.IsType(err, &types.StackInstanceNotFoundException{}) {
				return true, nil
			}
			return false, err
		}

		stackStatus := stackOutput.Stacks[0].StackStatus
		switch stackStatus {
		case types.StackStatusCreateComplete, types.StackStatusUpdateComplete, types.StackStatusDeleteComplete:
			return true, nil
		case types.StackStatusCreateInProgress, types.StackStatusUpdateInProgress, types.StackStatusDeleteInProgress, types.StackStatusUpdateCompleteCleanupInProgress:
			return false, nil
		default:
			return false, fmt.Errorf("stack %s failed with status: %s", stackName, stackStatus)
		}
	})

	return err
}

func (s *stack) delete(ctx context.Context, clusterName string) error {
	stackName := stackName(clusterName)
	s.logger.Info("Deleting E2E test cluster stack", "stackName", stackName)
	_, err := s.cfn.DeleteStack(ctx, &cloudformation.DeleteStackInput{
		StackName: aws.String(stackName),
	})
	if err != nil {
		return fmt.Errorf("deleting hybrid nodes setup cfn stack: %w", err)
	}

	s.logger.Info("Waiting for stack to be deleted", "stackName", stackName)
	if err := waitForStackOperation(ctx, s.cfn, stackName); err != nil {
		return fmt.Errorf("waiting for hybrid nodes setup cfn stack to be deleted: %w", err)
	}

	s.logger.Info("E2E test cluster stack deleted successfully", "stackName", stackName)
	return nil
}
