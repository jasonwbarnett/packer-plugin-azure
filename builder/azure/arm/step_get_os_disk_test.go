package arm

import (
	"context"
	"fmt"
	"testing"

	"github.com/Azure/azure-sdk-for-go/services/compute/mgmt/2018-04-01/compute"

	"github.com/hashicorp/packer-plugin-azure/builder/azure/common/constants"

	"github.com/hashicorp/packer-plugin-sdk/multistep"
)

func TestStepGetOSDiskShouldFailIfGetFails(t *testing.T) {
	var testSubject = &StepGetOSDisk{
		query: func(context.Context, string, string) (compute.VirtualMachine, error) {
			return createVirtualMachineFromUri("test.vhd"), fmt.Errorf("!! Unit Test FAIL !!")
		},
		say:   func(message string) {},
		error: func(e error) {},
	}

	stateBag := createTestStateBagStepGetOSDisk()

	var result = testSubject.Run(context.Background(), stateBag)
	if result != multistep.ActionHalt {
		t.Fatalf("Expected the step to return 'ActionHalt', but got '%d'.", result)
	}

	if _, ok := stateBag.GetOk(constants.Error); ok == false {
		t.Fatalf("Expected the step to set stateBag['%s'], but it was not.", constants.Error)
	}
}

func TestStepGetOSDiskShouldPassIfGetPasses(t *testing.T) {
	var testSubject = &StepGetOSDisk{
		query: func(context.Context, string, string) (compute.VirtualMachine, error) {
			return createVirtualMachineFromUri("test.vhd"), nil
		},
		say:   func(message string) {},
		error: func(e error) {},
	}

	stateBag := createTestStateBagStepGetOSDisk()

	var result = testSubject.Run(context.Background(), stateBag)
	if result != multistep.ActionContinue {
		t.Fatalf("Expected the step to return 'ActionContinue', but got '%d'.", result)
	}

	if _, ok := stateBag.GetOk(constants.Error); ok == true {
		t.Fatalf("Expected the step to not set stateBag['%s'], but it was.", constants.Error)
	}
}

func TestStepGetOSDiskShouldTakeValidateArgumentsFromStateBag(t *testing.T) {
	var actualResourceGroupName string
	var actualComputeName string

	var testSubject = &StepGetOSDisk{
		query: func(ctx context.Context, resourceGroupName string, computeName string) (compute.VirtualMachine, error) {
			actualResourceGroupName = resourceGroupName
			actualComputeName = computeName

			return createVirtualMachineFromUri("test.vhd"), nil
		},
		say:   func(message string) {},
		error: func(e error) {},
	}

	stateBag := createTestStateBagStepGetOSDisk()
	var result = testSubject.Run(context.Background(), stateBag)

	if result != multistep.ActionContinue {
		t.Fatalf("Expected the step to return 'ActionContinue', but got '%d'.", result)
	}

	var expectedComputeName = stateBag.Get(constants.ArmComputeName).(string)
	var expectedResourceGroupName = stateBag.Get(constants.ArmResourceGroupName).(string)

	if actualComputeName != expectedComputeName {
		t.Fatal("Expected the step to source 'constants.ArmResourceGroupName' from the state bag, but it did not.")
	}

	if actualResourceGroupName != expectedResourceGroupName {
		t.Fatal("Expected the step to source 'constants.ArmResourceGroupName' from the state bag, but it did not.")
	}

	expectedOSDiskVhd, ok := stateBag.GetOk(constants.ArmOSDiskVhd)
	if !ok {
		t.Fatalf("Expected the state bag to have a value for '%s', but it did not.", constants.ArmOSDiskVhd)
	}

	if expectedOSDiskVhd != "test.vhd" {
		t.Fatalf("Expected the value of stateBag[%s] to be 'test.vhd', but got '%s'.", constants.ArmOSDiskVhd, expectedOSDiskVhd)
	}
}

func createTestStateBagStepGetOSDisk() multistep.StateBag {
	stateBag := new(multistep.BasicStateBag)

	stateBag.Put(constants.ArmComputeName, "Unit Test: ComputeName")
	stateBag.Put(constants.ArmResourceGroupName, "Unit Test: ResourceGroupName")

	return stateBag
}

func createVirtualMachineFromUri(vhdUri string) compute.VirtualMachine {
	vm := compute.VirtualMachine{
		VirtualMachineProperties: &compute.VirtualMachineProperties{
			StorageProfile: &compute.StorageProfile{
				OsDisk: &compute.OSDisk{
					Vhd: &compute.VirtualHardDisk{
						URI: &vhdUri,
					},
				},
			},
		},
	}

	return vm
}
