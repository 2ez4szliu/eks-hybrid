package imagecredentialprovider

import (
	"context"
	"os"
	"path"

	"github.com/pkg/errors"
	"go.uber.org/zap"

	"github.com/aws/eks-hybrid/internal/artifact"
	"github.com/aws/eks-hybrid/internal/tracker"
)

// BinPath is the path to the image-credential-provider binary.
const (
	BinPath = "/etc/eks/image-credential-provider/ecr-credential-provider"

	artifactName      = "image-credential-provider"
	artifactFilePerms = 0o755
)

// Source represents a source that serves an image-credential-provider binary.
type Source interface {
	GetImageCredentialProvider(context.Context) (artifact.Source, error)
}

// Install installs the image-credential-provider at BinPath.
func Install(ctx context.Context, tracker *tracker.Tracker, src Source) error {
	imageCredentialProvider, err := src.GetImageCredentialProvider(ctx)
	if err != nil {
		return errors.Wrap(err, "failed to image-credential-provider source")
	}
	defer imageCredentialProvider.Close()

	if err := artifact.InstallFile(BinPath, imageCredentialProvider, artifactFilePerms); err != nil {
		return errors.Wrap(err, "failed to install image-credential-provider")
	}

	if !imageCredentialProvider.VerifyChecksum() {
		return errors.Errorf("image-credential-provider checksum mismatch: %v", artifact.NewChecksumError(imageCredentialProvider))
	}
	if err = tracker.Add(artifact.ImageCredentialProvider); err != nil {
		return err
	}

	return nil
}

func Uninstall() error {
	return os.RemoveAll(path.Dir(BinPath))
}

func Upgrade(ctx context.Context, src Source, log *zap.Logger) error {
	imageCredentialProvider, err := src.GetImageCredentialProvider(ctx)
	if err != nil {
		return errors.Wrap(err, "getting image-credential-provider source")
	}
	defer imageCredentialProvider.Close()

	return artifact.Upgrade(artifactName, BinPath, imageCredentialProvider, artifactFilePerms, log)
}
