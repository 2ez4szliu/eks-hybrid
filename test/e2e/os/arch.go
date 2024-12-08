package os

import (
	"bytes"
	"context"
	"text/template"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ssm"
	"github.com/aws/eks-hybrid/test/e2e"
)

const (
	amd64Arch = "x86_64"
	arm64Arch = "arm64"
)

func populateBaseScripts(userDataInput *e2e.UserDataInput) error {
	logCollector, err := executeTemplate(logCollectorScript, userDataInput)
	if err != nil {
		return err
	}

	userDataInput.Files = append(userDataInput.Files,
		e2e.File{Content: string(nodeAdmInitScript), Path: "/tmp/nodeadm-init.sh", Permissions: "0755"},
		e2e.File{Content: string(logCollector), Path: "/tmp/log-collector.sh", Permissions: "0755"},
	)
	return nil
}

func executeTemplate(templateData []byte, values any) ([]byte, error) {
	tmpl, err := template.New("cloud-init").Funcs(templateFuncMap()).Parse(string(templateData))
	if err != nil {
		return nil, err
	}

	// Execute the template and write the result to a buffer
	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, values); err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

func getAmiIDFromSSM(ctx context.Context, client *ssm.SSM, amiName string) (*string, error) {
	getParameterInput := &ssm.GetParameterInput{
		Name:           aws.String(amiName),
		WithDecryption: aws.Bool(true),
	}

	output, err := client.GetParameterWithContext(ctx, getParameterInput)
	if err != nil {
		return nil, err
	}

	return output.Parameter.Value, nil
}

func getInstanceTypeFromRegionAndArch(region, arch string) string {
	if arch == amd64Arch {
		if region == "ap-southeast-5" {
			return "m6i.large"
		}
		return "m5.2xlarge"
	}
	return "t4g.2xlarge"
}
