package aws

import (
	"context"

	"github.com/aws/eks-hybrid/internal/util"

	"github.com/pkg/errors"
	"sigs.k8s.io/yaml"
)

// set build time
var manifestUrl string

type Manifest struct {
	SupportedEksReleases     []SupportedEksRelease     `json:"supported_eks_releases"`
	IamRolesAnywhereReleases []IamRolesAnywhereRelease `json:"iam_roles_anywhere_releases"`
	SsmReleases              []SsmRelease              `json:"ssm_releases"`
}

type SupportedEksRelease struct {
	MajorMinorVersion  string            `json:"major_minor_version"`
	LatestPatchVersion string            `json:"latest_patch_version"`
	PatchReleases      []EksPatchRelease `json:"patch_releases"`
}

type EksPatchRelease struct {
	Version      string     `json:"version"`
	PatchVersion string     `json:"patch_version"`
	ReleaseDate  string     `json:"release_date"`
	Artifacts    []Artifact `json:"artifacts"`
}

type IamRolesAnywhereRelease struct {
	Version   string     `json:"version"`
	Artifacts []Artifact `json:"artifacts"`
}

type SsmRelease struct {
	Version   string     `json:"version"`
	Artifacts []Artifact `json:"artifacts"`
}

type Artifact struct {
	Name        string `json:"name"`
	Arch        string `json:"arch"`
	OS          string `json:"os"`
	URI         string `json:"uri"`
	ChecksumURI string `json:"checksum_uri"`
}

// Read from the manifest file on s3 and parse into Manifest struct
func getReleaseManifest(ctx context.Context) (*Manifest, error) {
	yamlFileData, err := util.GetHttpFile(ctx, manifestUrl)
	if err != nil {
		return nil, err
	}
	var manifest Manifest
	err = yaml.Unmarshal(yamlFileData, &manifest)
	if err != nil {
		return nil, errors.Wrap(err, "invalid yaml data in release manifest")
	}
	return &manifest, nil
}
