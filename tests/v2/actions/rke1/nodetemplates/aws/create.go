package nodetemplates

import (
	"github.com/imdario/mergo"
	"github.com/rancher/rancher/tests/v2/actions/rke1/nodetemplates"
	"github.com/rancher/shepherd/clients/rancher"
	management "github.com/rancher/shepherd/clients/rancher/generated/management/v3"
	"github.com/rancher/shepherd/extensions/cloudcredentials/aws"
	"github.com/rancher/shepherd/pkg/config"
)

const awsEC2NodeTemplateNameBase = "awsNodeConfig"

// CreateAWSNodeTemplate is a helper function that takes the rancher Client as a parameter and creates
// an AWS node template and returns the NodeTemplate response
func CreateAWSNodeTemplate(rancherClient *rancher.Client) (*nodetemplates.NodeTemplate, error) {
	var amazonEC2NodeTemplateConfig nodetemplates.AmazonEC2NodeTemplateConfig
	config.LoadConfig(nodetemplates.AmazonEC2NodeTemplateConfigurationFileKey, &amazonEC2NodeTemplateConfig)

	cloudCredential, err := aws.CreateAWSCloudCredentials(rancherClient)
	if err != nil {
		return nil, err
	}

	nodeTemplate := nodetemplates.NodeTemplate{
		EngineInstallURL:            "https://releases.rancher.com/install-docker/24.0.sh",
		Name:                        awsEC2NodeTemplateNameBase,
		AmazonEC2NodeTemplateConfig: &amazonEC2NodeTemplateConfig,
	}

	nodeTemplateConfig := &nodetemplates.NodeTemplate{
		CloudCredentialID: cloudCredential.ID,
	}

	config.LoadConfig(nodetemplates.NodeTemplateConfigurationFileKey, nodeTemplateConfig)

	err = mergo.Merge(&nodeTemplate, nodeTemplateConfig, mergo.WithOverride)
	if err != nil {
		return nil, err
	}

	resp := &nodetemplates.NodeTemplate{}
	err = rancherClient.Management.APIBaseClient.Ops.DoCreate(management.NodeTemplateType, nodeTemplate, resp)
	if err != nil {
		return nil, err
	}
	return resp, nil
}