package StatusPollingService

import (
	"fmt"
	DeviceRegistrationHelper "github.com/manjunath-ravindra/go-edge-device/helpers/deviceRegistration"
	DeviceTypes "github.com/manjunath-ravindra/go-edge-device/types/device"
)

func PollDeviceStatus(BASE_URL string, DEVICE_ID string, SECRET_KEY string) (*DeviceTypes.DeviceStatusResponse, error) {

	status, err := DeviceRegistrationHelper.CheckDeviceStatus(BASE_URL, DEVICE_ID, SECRET_KEY)

	if err != nil {
		fmt.Println("Error fetching the Status from DeviceStatusCheck")
		return nil, err
	}

	return status, nil
}
