package main

import (
	"fmt"
	DeviceRegistrationService "github.com/manjunath-ravindra/go-edge-device/services"
	EnvVendors "github.com/manjunath-ravindra/go-edge-device/vendors/env"
)

func main() {

	envVariables, err := EnvVendors.LoadEnv()
	if err != nil {
		fmt.Println("Error loading environment variables:", err)
		return
	}

	DeviceRegistrationService.InitialRegistration(
		envVariables.BaseURL,
		envVariables.DeviceID,
		envVariables.SecretKey,
		envVariables.DeviceFrom,
		envVariables.EncryptionKey,
	)

	fmt.Println("Main operation completed successfully")
}
