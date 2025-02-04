package gcp

import (
	"strings"

	"github.com/litmuschaos/litmus-go/pkg/log"
	"github.com/pkg/errors"
	"golang.org/x/net/context"
	"google.golang.org/api/compute/v1"
	"google.golang.org/api/option"
)

// GetVMInstanceStatus returns the status of a VM instance
func GetVMInstanceStatus(instanceName string, gcpProjectID string, instanceZone string) (string, error) {
	var err error

	// create an empty context
	ctx := context.Background()

	// get service account credentials json
	json, err := GetServiceAccountJSONFromSecret()
	if err != nil {
		return "", err
	}

	// create a new GCP Compute Service client using the GCP service account credentials
	computeService, err := compute.NewService(ctx, option.WithCredentialsJSON(json))
	if err != nil {
		return "", err
	}

	// get information about the requisite VM instance
	response, err := computeService.Instances.Get(gcpProjectID, instanceZone, instanceName).Context(ctx).Do()
	if err != nil {
		return "", err
	}

	// return the VM status
	return response.Status, nil
}

//InstanceStatusCheckByName is used to check the status of all the VM instances under chaos
func InstanceStatusCheckByName(instanceNames string, gcpProjectId string, instanceZones string) error {
	instanceNamesList := strings.Split(instanceNames, ",")
	instanceZonesList := strings.Split(instanceZones, ",")
	if len(instanceNamesList) == 0 {
		return errors.Errorf("No instance name found to stop")
	}
	if len(instanceNamesList) != len(instanceZonesList) {
		return errors.Errorf("The number of instance names and the number of regions is not equal")
	}
	log.Infof("[Info]: The instances under chaos(IUC) are: %v", instanceNamesList)
	return InstanceStatusCheck(instanceNamesList, gcpProjectId, instanceZonesList)
}

//InstanceStatusCheck is used to check whether all VM instances under chaos are running or not
func InstanceStatusCheck(instanceNamesList []string, gcpProjectId string, instanceZonesList []string) error {
	for i := range instanceNamesList {
		instanceState, err := GetVMInstanceStatus(instanceNamesList[i], gcpProjectId, instanceZonesList[i])
		if err != nil {
			return err
		}
		if instanceState != "RUNNING" {
			return errors.Errorf("failed to get the vm instance '%v' in running sate, current state: %v", instanceNamesList[i], instanceState)
		}
	}
	return nil
}
