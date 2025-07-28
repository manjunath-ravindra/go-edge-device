package envVendors

import (
	"fmt"
	"github.com/joho/godotenv"
	EnvTypes "github.com/manjunath-ravindra/go-edge-device/types/env"
	"os"
)

func LoadEnv() (EnvTypes.EnvVariableStructTypes, error) {
	var err error = godotenv.Load()
	if err != nil {
		fmt.Println("No .env file found or error loading .env file")
		return EnvTypes.EnvVariableStructTypes{}, err
	}

	// Now you can use os.Getenv to get your variables
	BASE_URL := os.Getenv("BASE_URL")
	ENCRYPTION_KEY := os.Getenv("ENCRYPTION_KEY")
	DEVICE_ID := os.Getenv("DEVICE_ID")
	SECRET_KEY := os.Getenv("SECRET_KEY")
	DEVICE_FROM := os.Getenv("DEVICE_FROM")

	return EnvTypes.EnvVariableStructTypes{
		BaseURL:       BASE_URL,
		EncryptionKey: ENCRYPTION_KEY,
		DeviceID:      DEVICE_ID,
		SecretKey:     SECRET_KEY,
		DeviceFrom:    DEVICE_FROM,
	}, nil
}
