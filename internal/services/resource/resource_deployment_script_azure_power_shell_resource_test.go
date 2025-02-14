package resource_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/go-azure-sdk/resource-manager/resources/2020-10-01/deploymentscripts"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

type ResourceDeploymentScriptAzurePowerShellResource struct{}

func TestAccResourceDeploymentScriptAzurePowerShell_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_resource_deployment_script_azure_power_shell", "test")
	r := ResourceDeploymentScriptAzurePowerShellResource{}
	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccResourceDeploymentScriptAzurePowerShell_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_resource_deployment_script_azure_power_shell", "test")
	r := ResourceDeploymentScriptAzurePowerShellResource{}
	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.RequiresImportErrorStep(r.requiresImport),
	})
}

func TestAccResourceDeploymentScriptAzurePowerShell_complete(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_resource_deployment_script_azure_power_shell", "test")
	r := ResourceDeploymentScriptAzurePowerShellResource{}
	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.complete(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("environment_variable.0.secure_value", "storage_account.0.key"),
	})
}

func TestAccResourceDeploymentScriptAzurePowerShell_update(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_resource_deployment_script_azure_power_shell", "test")
	r := ResourceDeploymentScriptAzurePowerShellResource{}
	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.complete(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("environment_variable.0.secure_value", "storage_account.0.key"),
		{
			Config: r.update(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("environment_variable.0.secure_value", "storage_account.0.key"),
	})
}

func TestAccResourceDeploymentScriptAzurePowerShell_scriptUri(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_resource_deployment_script_azure_power_shell", "test")
	r := ResourceDeploymentScriptAzurePowerShellResource{}
	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.scriptUri(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("environment_variable.0.secure_value", "storage_account.0.key"),
	})
}

func (r ResourceDeploymentScriptAzurePowerShellResource) Exists(ctx context.Context, clients *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := deploymentscripts.ParseDeploymentScriptID(state.ID)
	if err != nil {
		return nil, err
	}

	client := clients.Resource.DeploymentScriptsClient
	resp, err := client.Get(ctx, *id)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			return utils.Bool(false), nil
		}
		return nil, fmt.Errorf("retrieving %s: %+v", id, err)
	}
	return utils.Bool(resp.Model != nil && *resp.Model != nil), nil
}

func (r ResourceDeploymentScriptAzurePowerShellResource) template(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctest-rg-%d"
  location = "%s"
}
`, data.RandomInteger, data.Locations.Primary)
}

func (r ResourceDeploymentScriptAzurePowerShellResource) basic(data acceptance.TestData) string {
	template := r.template(data)
	return fmt.Sprintf(`
				%s

resource "azurerm_resource_deployment_script_azure_power_shell" "test" {
  name                = "acctest-rdsaps-%d"
  resource_group_name = azurerm_resource_group.test.name
  location            = "%s"
  version             = "8.3"
  retention_interval  = "P1D"
  script_content      = <<EOF
		$output = 'Hello'
		Write-Output $output
		$DeploymentScriptOutputs = @{}
		$DeploymentScriptOutputs['text'] = $output
EOF
}
`, template, data.RandomInteger, data.Locations.Primary)
}

func (r ResourceDeploymentScriptAzurePowerShellResource) requiresImport(data acceptance.TestData) string {
	config := r.basic(data)
	return fmt.Sprintf(`
			%s

resource "azurerm_resource_deployment_script_azure_power_shell" "import" {
  name                = azurerm_resource_deployment_script_azure_power_shell.test.name
  resource_group_name = azurerm_resource_deployment_script_azure_power_shell.test.resource_group_name
  location            = azurerm_resource_deployment_script_azure_power_shell.test.location
  version             = azurerm_resource_deployment_script_azure_power_shell.test.version
  retention_interval  = azurerm_resource_deployment_script_azure_power_shell.test.retention_interval
  script_content      = azurerm_resource_deployment_script_azure_power_shell.test.script_content
}
`, config)
}

func (r ResourceDeploymentScriptAzurePowerShellResource) complete(data acceptance.TestData) string {
	template := r.template(data)
	return fmt.Sprintf(`
			%[1]s

resource "azurerm_user_assigned_identity" "test" {
  name                = "acctestUAI-%[2]d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
}

resource "azurerm_storage_account" "test" {
  name                = "unlikely23exst2acct%[4]s"
  resource_group_name = azurerm_resource_group.test.name

  location                 = azurerm_resource_group.test.location
  account_tier             = "Standard"
  account_replication_type = "LRS"
}


resource "azurerm_resource_deployment_script_azure_power_shell" "test" {
  name                = "acctest-rdsaps-%[2]d"
  resource_group_name = azurerm_resource_group.test.name
  location            = "%[3]s"
  version             = "8.3"
  retention_interval  = "P1D"
  command_line        = "-name \"John Dole\""
  cleanup_preference  = "OnSuccess"
  force_update_tag    = "1"
  timeout             = "PT30M"

  script_content = <<EOF
		param([string] $name)
      	$output = 'Hello {0}. The username is {1}, the password is {2}.' -f $name,$${Env:UserName},$${Env:Password}
      	Write-Output $output
      	$DeploymentScriptOutputs = @{}
      	$DeploymentScriptOutputs['text'] = $output
EOF

  supporting_script_uris = ["https://raw.githubusercontent.com/Azure/azure-docs-json-samples/master/deployment-script/create-cert.ps1"]

  container {
    container_group_name = "cgn-%[2]d"
  }

  environment_variable {
    name  = "UserName"
    value = "jdole"
  }

  environment_variable {
    name         = "Password"
    secure_value = "jDolePassword"
  }

  identity {
    type = "UserAssigned"
    identity_ids = [
      azurerm_user_assigned_identity.test.id
    ]
  }

  storage_account {
    name = azurerm_storage_account.test.name
    key  = azurerm_storage_account.test.primary_access_key
  }

  tags = {
    key = "value"
  }

}
`, template, data.RandomInteger, data.Locations.Primary, data.RandomString)
}

func (r ResourceDeploymentScriptAzurePowerShellResource) update(data acceptance.TestData) string {
	template := r.template(data)
	return fmt.Sprintf(`
			%[1]s

resource "azurerm_user_assigned_identity" "test" {
  name                = "acctestUAI-%[2]d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
}

resource "azurerm_storage_account" "test" {
  name                = "unlikely23exst2acct%[4]s"
  resource_group_name = azurerm_resource_group.test.name

  location                 = azurerm_resource_group.test.location
  account_tier             = "Standard"
  account_replication_type = "LRS"
}


resource "azurerm_resource_deployment_script_azure_power_shell" "test" {
  name                = "acctest-rdsaps-%[2]d"
  resource_group_name = azurerm_resource_group.test.name
  location            = "%[3]s"
  version             = "8.3"
  retention_interval  = "P1D"
  command_line        = "-name \"John Dole\""
  cleanup_preference  = "OnSuccess"
  force_update_tag    = "1"
  timeout             = "PT30M"

  script_content = <<EOF
		param([string] $name)
      	$output = 'Hello {0}. The username is {1}, the password is {2}.' -f $name,$${Env:UserName},$${Env:Password}
      	Write-Output $output
      	$DeploymentScriptOutputs = @{}
      	$DeploymentScriptOutputs['text'] = $output
EOF

  supporting_script_uris = ["https://raw.githubusercontent.com/Azure/azure-docs-json-samples/master/deployment-script/create-cert.ps1"]

  container {
    container_group_name = "cgn-%[2]d"
  }

  environment_variable {
    name  = "UserName"
    value = "jdole"
  }

  environment_variable {
    name         = "Password"
    secure_value = "jDolePassword"
  }

  identity {
    type = "UserAssigned"
    identity_ids = [
      azurerm_user_assigned_identity.test.id
    ]
  }

  storage_account {
    name = azurerm_storage_account.test.name
    key  = azurerm_storage_account.test.primary_access_key
  }

  tags = {
    key = "value2"
  }

}
`, template, data.RandomInteger, data.Locations.Primary, data.RandomString)
}

func (r ResourceDeploymentScriptAzurePowerShellResource) scriptUri(data acceptance.TestData) string {
	template := r.template(data)
	return fmt.Sprintf(`
			%[1]s

resource "azurerm_resource_deployment_script_azure_power_shell" "test" {
  name                = "acctest-rdsaps-%[2]d"
  resource_group_name = azurerm_resource_group.test.name
  location            = "%[3]s"
  version             = "8.3"
  retention_interval  = "P1D"
  command_line        = "-name \"John Dole\""
  primary_script_uri  = "https://raw.githubusercontent.com/Azure/azure-docs-json-samples/master/deployment-script/deploymentscript-helloworld.ps1"
}
`, template, data.RandomInteger, data.Locations.Primary, data.RandomString)
}
