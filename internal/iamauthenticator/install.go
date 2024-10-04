package iamauthenticator

import (
	"context"
	"fmt"
	"os"

	"github.com/aws/eks-hybrid/internal/artifact"
	"github.com/aws/eks-hybrid/internal/tracker"
)

// IAMAuthenticatorBinPath is the path the IAM Authenticator is installed to.
const IAMAuthenticatorBinPath = "/usr/local/bin/aws-iam-authenticator"

// IAMAuthenticatorSource retrieves the aws-iam-authenticator binary.
type IAMAuthenticatorSource interface {
	GetIAMAuthenticator(context.Context) (artifact.Source, error)
}

// Install installs the aws_signing_helper and aws-iam-authenticator on the system at
// SigningHelperBinPath and IAMAuthenticatorBinPath respectively.
func Install(ctx context.Context, tracker *tracker.Tracker, iamAuthSrc IAMAuthenticatorSource) error {
	authenticator, err := iamAuthSrc.GetIAMAuthenticator(ctx)
	if err != nil {
		return fmt.Errorf("failed to get aws-iam-authenticator source: %w", err)
	}
	defer authenticator.Close()

	if err := artifact.InstallFile(IAMAuthenticatorBinPath, authenticator, 0755); err != nil {
		return fmt.Errorf("failed to install aws-iam-authenticator: %w", err)
	}

	if !authenticator.VerifyChecksum() {
		return fmt.Errorf("aws-iam-authenticator checksum mismatch: %w", artifact.NewChecksumError(authenticator))
	}
	if err = tracker.Add(artifact.IamAuthenticator); err != nil {
		return err
	}

	return nil
}

func Uninstall() error {
	return os.RemoveAll(IAMAuthenticatorBinPath)
}
