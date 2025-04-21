package cluster

import (
	"context"
	"fmt"
	"os"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/cloudformation"
	"github.com/aws/aws-sdk-go-v2/service/cloudwatchlogs"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	"github.com/aws/aws-sdk-go-v2/service/eks"
	"github.com/aws/aws-sdk-go-v2/service/iam"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/ssm"
	"github.com/go-logr/logr"
	"gopkg.in/yaml.v2"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"

	"github.com/aws/eks-hybrid/test/e2e"
	"github.com/aws/eks-hybrid/test/e2e/addon"
	"github.com/aws/eks-hybrid/test/e2e/cni"
	"github.com/aws/eks-hybrid/test/e2e/constants"
)

const (
	clusterLogRetentionDays = 14
	clusterLogGroupName     = "/aws/eks/%s/cluster"
)

type TestResources struct {
	ClusterName       string        `yaml:"clusterName"`
	ClusterRegion     string        `yaml:"clusterRegion"`
	ClusterNetwork    NetworkConfig `yaml:"clusterNetwork"`
	HybridNetwork     NetworkConfig `yaml:"hybridNetwork"`
	KubernetesVersion string        `yaml:"kubernetesVersion"`
	Cni               string        `yaml:"cni"`
	EKS               EKSConfig     `yaml:"eks"`
}
type EKSConfig struct {
	Endpoint      string `yaml:"endpoint"`
	ClusterRoleSP string `yaml:"clusterRoleSP"`
	PodIdentitySP string `yaml:"podIdentitySP"`
}

type NetworkConfig struct {
	VpcCidr           string `yaml:"vpcCidr"`
	PublicSubnetCidr  string `yaml:"publicSubnetCidr"`
	PrivateSubnetCidr string `yaml:"privateSubnetCidr"`
	PodCidr           string `yaml:"podCidr"`
}

const (
	ciliumCni            = "cilium"
	calicoCni            = "calico"
	defaultClusterRoleSP = "eks.amazonaws.com"
	defaultPodIdentitySP = "pods.eks.amazonaws.com"
)

type Create struct {
	logger         logr.Logger
	eks            *eks.Client
	stack          *stack
	iam            *iam.Client
	s3             *s3.Client
	cloudWatchLogs *cloudwatchlogs.Client
}

// NewCreate creates a new workflow to create an EKS cluster. The EKS client will use
// the specified endpoint or the default endpoint if empty string is passed.
func NewCreate(aws aws.Config, logger logr.Logger, endpoint string) Create {
	return Create{
		logger: logger,
		eks:    e2e.NewEKSClient(aws, endpoint),
		stack: &stack{
			iamClient: iam.NewFromConfig(aws),
			cfn:       cloudformation.NewFromConfig(aws),
			ec2Client: ec2.NewFromConfig(aws),
			logger:    logger,
			ssmClient: ssm.NewFromConfig(aws),
		},
		iam:            iam.NewFromConfig(aws),
		s3:             s3.NewFromConfig(aws),
		cloudWatchLogs: cloudwatchlogs.NewFromConfig(aws),
	}
}

func (c *Create) Run(ctx context.Context, test TestResources) error {
	stackOut, err := c.stack.deploy(ctx, test)
	if err != nil {
		return fmt.Errorf("creating stack for cluster infra: %w", err)
	}

	hybridCluster := hybridCluster{
		Name:              test.ClusterName,
		Region:            test.ClusterRegion,
		KubernetesVersion: test.KubernetesVersion,
		SecurityGroup:     stackOut.clusterVpcConfig.securityGroup,
		SubnetIDs:         []string{stackOut.clusterVpcConfig.publicSubnet, stackOut.clusterVpcConfig.privateSubnet},
		Role:              stackOut.clusterRole,
		HybridNetwork:     test.HybridNetwork,
	}

	c.logger.Info("Creating EKS cluster..", "cluster", test.ClusterName)
	cluster, err := hybridCluster.create(ctx, c.eks, c.logger)
	if err != nil {
		return fmt.Errorf("creating %s EKS cluster: %w", test.KubernetesVersion, err)
	}

	err = c.tagClusterLogGroup(ctx, test.ClusterName)
	if err != nil {
		return fmt.Errorf("tagging cluster log group: %w", err)
	}

	err = c.setClusterLogRetention(ctx, test.ClusterName)
	if err != nil {
		return fmt.Errorf("setting cluster log retention: %w", err)
	}

	kubeconfig := KubeconfigPath(test.ClusterName)
	err = hybridCluster.UpdateKubeconfig(cluster, kubeconfig)
	if err != nil {
		return fmt.Errorf("saving kubeconfig for %s EKS cluster: %w", test.KubernetesVersion, err)
	}

	clientConfig, err := clientcmd.BuildConfigFromFlags("", kubeconfig)
	if err != nil {
		return fmt.Errorf("loading kubeconfig: %w", err)
	}

	k8sClient, err := kubernetes.NewForConfig(clientConfig)
	if err != nil {
		return fmt.Errorf("creating kubernetes client: %w", err)
	}

	podIdentityAddon := addon.NewPodIdentityAddon(hybridCluster.Name, stackOut.podIdentity.roleArn)

	err = podIdentityAddon.Create(ctx, c.logger, c.eks, k8sClient)
	if err != nil {
		return fmt.Errorf("creating add-on %s for EKS cluster: %w", podIdentityAddon.Name, err)
	}

	// upload test file to pod identity S3 bucket
	err = podIdentityAddon.UploadFileForVerification(ctx, c.logger, c.s3, stackOut.podIdentity.s3Bucket)
	if err != nil {
		return fmt.Errorf("uploading test file to s3 bucket: %s", stackOut.podIdentity.s3Bucket)
	}

	dynamicK8s, err := dynamic.NewForConfig(clientConfig)
	if err != nil {
		return fmt.Errorf("creating dynamic Kubernetes client: %w", err)
	}

	switch test.Cni {
	case ciliumCni:
		cilium := cni.NewCilium(dynamicK8s, test.HybridNetwork.PodCidr, test.ClusterRegion, test.KubernetesVersion)
		c.logger.Info("Installing cilium on cluster...", "cluster", test.ClusterName)
		if err = cilium.Deploy(ctx); err != nil {
			return fmt.Errorf("installing cilium for %s EKS cluster: %w", test.KubernetesVersion, err)
		}
		c.logger.Info("Cilium installed successfully.")
	case calicoCni:
		calico := cni.NewCalico(dynamicK8s, test.HybridNetwork.PodCidr, test.ClusterRegion)
		c.logger.Info("Installing calico on cluster...", "cluster", test.ClusterName)
		if err = calico.Deploy(ctx); err != nil {
			return fmt.Errorf("installing calico for %s EKS cluster: %w", test.KubernetesVersion, err)
		}
		c.logger.Info("Calico installed successfully.")
	}

	return nil
}

func KubeconfigPath(clusterName string) string {
	return fmt.Sprintf("/tmp/%s.kubeconfig", clusterName)
}

func LoadTestResources(path string) (TestResources, error) {
	file, err := os.ReadFile(path)
	if err != nil {
		return TestResources{}, fmt.Errorf("opening configuration file: %w", err)
	}

	var testResources TestResources
	err = yaml.Unmarshal(file, &testResources)
	if err != nil {
		return TestResources{}, fmt.Errorf("unmarshalling test resources: %w", err)
	}

	testResources = SetTestResourcesDefaults(testResources)

	return testResources, nil
}

func SetTestResourcesDefaults(testResources TestResources) TestResources {
	if testResources.EKS.ClusterRoleSP == "" {
		testResources.EKS.ClusterRoleSP = defaultClusterRoleSP
	}

	if testResources.EKS.PodIdentitySP == "" {
		testResources.EKS.PodIdentitySP = defaultPodIdentitySP
	}
	if testResources.ClusterNetwork == (NetworkConfig{}) {
		testResources.ClusterNetwork = NetworkConfig{
			VpcCidr:           "10.0.0.0/16",
			PublicSubnetCidr:  "10.0.10.0/24",
			PrivateSubnetCidr: "10.0.20.0/24",
		}
	}
	if testResources.HybridNetwork == (NetworkConfig{}) {
		testResources.HybridNetwork = NetworkConfig{
			VpcCidr:           "10.1.0.0/16",
			PublicSubnetCidr:  "10.1.1.0/24",
			PrivateSubnetCidr: "10.1.2.0/24",
			PodCidr:           "10.2.0.0/16",
		}
	}

	return testResources
}

func (c *Create) tagClusterLogGroup(ctx context.Context, clusterName string) error {
	describeLogGroups, err := c.cloudWatchLogs.DescribeLogGroups(ctx, &cloudwatchlogs.DescribeLogGroupsInput{
		LogGroupNamePrefix: aws.String(fmt.Sprintf(clusterLogGroupName, clusterName)),
	})
	if err != nil {
		return fmt.Errorf("describing log groups: %w", err)
	}
	if len(describeLogGroups.LogGroups) == 0 {
		return fmt.Errorf("log group not found")
	}

	_, err = c.cloudWatchLogs.TagResource(ctx, &cloudwatchlogs.TagResourceInput{
		ResourceArn: describeLogGroups.LogGroups[0].LogGroupArn,
		Tags: map[string]string{
			constants.TestClusterTagKey: clusterName,
		},
	})
	return err
}

func (c *Create) setClusterLogRetention(ctx context.Context, clusterName string) error {
	_, err := c.cloudWatchLogs.PutRetentionPolicy(ctx, &cloudwatchlogs.PutRetentionPolicyInput{
		LogGroupName:    aws.String(fmt.Sprintf(clusterLogGroupName, clusterName)),
		RetentionInDays: aws.Int32(clusterLogRetentionDays),
	})
	return err
}
