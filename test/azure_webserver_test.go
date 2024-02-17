package test

import (
	"testing"

	"github.com/gruntwork-io/terratest/modules/azure"
	"github.com/gruntwork-io/terratest/modules/terraform"
	"github.com/stretchr/testify/assert"
)

// You normally want to run this under a separate "Testing" subscription
// For lab purposes you will use your assigned subscription under the Cloud Dev/Ops program tenant
var subscriptionID string = "<removed subscription>"

func TestAzureLinuxVMCreation(t *testing.T) {
	terraformOptions := &terraform.Options{
		// The path to where our Terraform code is located
		TerraformDir: "../",
		// Override the default terraform variables
		Vars: map[string]interface{}{
			"labelPrefix": "song0138",
			"region":      "canadacentral",
		},
	}

	defer terraform.Destroy(t, terraformOptions)

	// Run `terraform init` and `terraform apply`. Fail the test if there are any errors.
	terraform.InitAndApply(t, terraformOptions)

	// Run `terraform output` to get the value of output variable
	vmName := terraform.Output(t, terraformOptions, "vm_name")
	resourceGroupName := terraform.Output(t, terraformOptions, "resource_group_name")

	// Confirm VM exists
	assert.True(t, azure.VirtualMachineExists(t, vmName, resourceGroupName, subscriptionID))
}
func TestAzureNICConnection(t *testing.T) {
	terraformOptions := &terraform.Options{
		TerraformDir: "../",
		Vars: map[string]interface{}{
			"labelPrefix": "song0138",
			"region":      "canadacentral",
		},
	}

	defer terraform.Destroy(t, terraformOptions)
	// 	// Run `terraform init` and `terraform apply`. Fail the test if there are any errors.
	terraform.InitAndApply(t, terraformOptions)
	vmName := terraform.Output(t, terraformOptions, "vm_name")
	resourceGroupName := terraform.Output(t, terraformOptions, "resource_group_name")
	nicName := terraform.Output(t, terraformOptions, "nic_name")

	// Get the list of NICs associated with the VM
	nicList := azure.GetVirtualMachineNics(t, vmName, resourceGroupName, subscriptionID)

	// Confirm the expected NIC is in the list of NICs associated with the VM
	var found bool
	for _, name := range nicList {
		if name == nicName {
			found = true
			break
		}
	}
	// Assert that the NIC with the specified name exists
	assert.True(t, found, "NIC with name %s should exist", nicName)
}

func TestAzureVMUbuntuVersion(t *testing.T) {
	terraformOptions := &terraform.Options{
		TerraformDir: "../",
		Vars: map[string]interface{}{
			"labelPrefix": "song0138",
			"region":      "canadacentral",
		},
	}

	defer terraform.Destroy(t, terraformOptions)

	// Run `terraform init` and `terraform apply`. Fail the test if there are any errors.
	terraform.InitAndApply(t, terraformOptions)
	vmName := terraform.Output(t, terraformOptions, "vm_name")
	resourceGroupName := terraform.Output(t, terraformOptions, "resource_group_name")

	// Get the image of the VM
	vmImage := azure.GetVirtualMachineImage(t, vmName, resourceGroupName, subscriptionID)

	// Check if the OS image is the expected Ubuntu version
	assert.Equal(t, "Canonical", vmImage.Publisher, "VM Publisher is not Canonical")
	assert.Equal(t, "0001-com-ubuntu-server-jammy", vmImage.Offer, "VM Offer is not Ubuntu Server Jammy")
	assert.Equal(t, "22_04-lts-gen2", vmImage.SKU, "VM SKU is not 22_04-lts-gen2")
}
