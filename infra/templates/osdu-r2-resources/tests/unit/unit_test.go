package test

import (
	"os"
	"testing"

	"github.com/gruntwork-io/terratest/modules/terraform"
	"github.com/microsoft/cobalt/test-harness/infratests"
)

var tfOptions = &terraform.Options{
	TerraformDir: "../../",
	Upgrade:      true,
	Vars: map[string]interface{}{
		"resource_group_location": region,
		"prefix":                  prefix,
		"app_services": []interface{}{
			map[string]interface{}{
				"app_name":         "tf-test-svc-1",
				"image":            *new(string),
				"linux_fx_version": "JAVA|8-jre8",
				"app_command_line": *new(string),
				"app_settings":     make(map[string]string, 0),
			},
		},
	},
	BackendConfig: map[string]interface{}{
		"storage_account_name": os.Getenv("TF_VAR_remote_state_account"),
		"container_name":       os.Getenv("TF_VAR_remote_state_container"),
	},
}

func TestTemplate(t *testing.T) {
	expectedAppDevResourceGroup := asMap(t, `{
		"location": "`+region+`"
	}`)

	expectedAppInsights := asMap(t, `{
		"application_type":    "Web"
	}`)

	resourceDescription := infratests.ResourceDescription{
		"azurerm_resource_group.app_rg":                                expectedAppDevResourceGroup,
		"module.app_insights.azurerm_application_insights.appinsights": expectedAppInsights,
	}

	appendAppServiceTests(t, resourceDescription)
	appendAutoScaleTests(t, resourceDescription)
	appendKeyVaultTests(t, resourceDescription)
	appendRedisTests(t, resourceDescription)
	appendStorageTests(t, resourceDescription)
	appendFunctionAppTests(t, resourceDescription)
	appendServicebusTests(t, resourceDescription)

	testFixture := infratests.UnitTestFixture{
		GoTest:                          t,
		TfOptions:                       tfOptions,
		Workspace:                       workspace,
		PlanAssertions:                  nil,
		ExpectedResourceCount:           104,
		ExpectedResourceAttributeValues: resourceDescription,
	}

	infratests.RunUnitTests(&testFixture)
}
