package test

import (
	"os"
	"testing"

	"github.com/gruntwork-io/terratest/modules/terraform"
	aksGitOpsIntegTests "github.com/microsoft/cobalt/infra/modules/providers/azure/aks-gitops/tests/integration"
	sbIntegTests "github.com/microsoft/cobalt/infra/modules/providers/azure/service-bus/tests/integration"
	esIntegTests "github.com/microsoft/cobalt/infra/modules/providers/elastic/elastic-cloud-enterprise/tests/integration"
	"github.com/microsoft/cobalt/test-harness/infratests"
)

var subscription = os.Getenv("ARM_SUBSCRIPTION_ID")
var kubeConfig = "../../output/bedrock_kube_config"
var tfOptions = &terraform.Options{
	TerraformDir: "../../",
	BackendConfig: map[string]interface{}{
		"storage_account_name": os.Getenv("TF_VAR_remote_state_account"),
		"container_name":       os.Getenv("TF_VAR_remote_state_container"),
	},
}

// Runs a suite of test assertions to validate that a provisioned data source environment
// is fully functional.
func TestDataEnvironment(t *testing.T) {
	testFixture := infratests.IntegrationTestFixture{
		GoTest:                t,
		TfOptions:             tfOptions,
		ExpectedTfOutputCount: 23,
		TfOutputAssertions: []infratests.TerraformOutputValidation{
			aksGitOpsIntegTests.BaselineClusterAssertions(
				kubeConfig,
				"contributor_service_principal_id",
				"mywebapp",
				"of your application is running on Kubernetes"),
			verifyServicePrincipalRoleAssignments,
			sbIntegTests.VerifySubscriptionsList(subscription,
				"resource_group",
				"sb_namespace_name",
				"sb_topics"),
			esIntegTests.ValidateElasticKvSecretValues("keyvault_secret_attributes", "elastic_cluster_properties"),
		},
	}
	infratests.RunIntegrationTests(&testFixture)
}
