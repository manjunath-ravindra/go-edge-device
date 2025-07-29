package main

import (
	"encoding/json"
	"fmt"
	"time"

	DeviceRegistrationService "github.com/manjunath-ravindra/go-edge-device/services/deviceRegistration"
	StatusPollingService "github.com/manjunath-ravindra/go-edge-device/services/statusPolling"
	DeviceTypes "github.com/manjunath-ravindra/go-edge-device/types/device"
	EnvVendors "github.com/manjunath-ravindra/go-edge-device/vendors/env"
)

func printAsJSON(v interface{}) {
	b, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		fmt.Println("Error marshaling to JSON:", err)
		return
	}
	fmt.Println(string(b), "Status printas json")
}

func main() {

	envVariables, err := EnvVendors.LoadEnv()
	if err != nil {
		fmt.Println("Error loading environment variables:", err)
		return
	}

	for {

		status, err := StatusPollingService.PollDeviceStatus(envVariables.BaseURL, envVariables.DeviceID, envVariables.SecretKey)

		if err != nil {
			fmt.Println("Error fetching Status from the device")
		}

		var Status DeviceTypes.DeviceStatus = DeviceTypes.DeviceStatus(status.Data.Status)

		printAsJSON(Status)
		switch Status {
		case DeviceTypes.Register:
			fmt.Println(Status, "Status in main")
			DeviceRegistrationService.ReRegistration(
				envVariables.BaseURL,
				envVariables.DeviceID,
				envVariables.SecretKey,
				envVariables.DeviceFrom,
				envVariables.EncryptionKey,
			)
		case DeviceTypes.DownloadComplete:
			//do nothing and exit the switch case
		case DeviceTypes.Deregistered:
			//do nothing and exit the switch case
		case DeviceTypes.Failed:
			fmt.Println(Status, "Status in main")
			DeviceRegistrationService.InitialRegistration(
				envVariables.BaseURL,
				envVariables.DeviceID,
				envVariables.SecretKey,
				envVariables.DeviceFrom,
				envVariables.EncryptionKey,
			)
		default:

		}
		time.Sleep(10 * time.Second)
	}
}
