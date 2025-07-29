package main

import (
	"fmt"
	"time"

	DeviceRegistrationService "github.com/manjunath-ravindra/go-edge-device/services/deviceRegistration"
	StatusPollingService "github.com/manjunath-ravindra/go-edge-device/services/statusPolling"
	DeviceTypes "github.com/manjunath-ravindra/go-edge-device/types/device"
	EnvTypes "github.com/manjunath-ravindra/go-edge-device/types/env"
	EnvVendors "github.com/manjunath-ravindra/go-edge-device/vendors/env"
)

func main() {
	var envVariables EnvTypes.EnvVariableStructTypes
	var err error
	for {
		envVariables, err = EnvVendors.LoadEnv()
		if err != nil {
			fmt.Println("Error loading environment variables:", err)
			time.Sleep(2 * time.Second)
			continue
		}
		break
	}

	for {
		status, err := StatusPollingService.PollDeviceStatus(envVariables.BaseURL, envVariables.DeviceID, envVariables.SecretKey)

		if err != nil {
			fmt.Println("Error fetching Status from the device")
		}

		var Status DeviceTypes.DeviceStatus = DeviceTypes.DeviceStatus(status.Data.Status)
		fmt.Println("Status of the Device: ", Status)

		switch Status {
		case DeviceTypes.Register:
			DeviceRegistrationService.ReRegistration(
				envVariables.BaseURL,
				envVariables.DeviceID,
				envVariables.SecretKey,
				envVariables.DeviceFrom,
				envVariables.EncryptionKey,
			)
		case DeviceTypes.DownloadComplete:
			//do nothing and exit the switch case
			//publish logs, records and weld params
		case DeviceTypes.Deregistered:
			//do nothing and exit the switch case
		case DeviceTypes.Failed:
			DeviceRegistrationService.InitialRegistration(
				envVariables.BaseURL,
				envVariables.DeviceID,
				envVariables.SecretKey,
				envVariables.DeviceFrom,
				envVariables.EncryptionKey,
			)
		default:

		}
		time.Sleep(15 * time.Second)
	}
}
