package DeviceRegistrationService

import (
	"fmt"
	DeviceRegistrationHelper "github.com/manjunath-ravindra/go-edge-device/helpers/deviceRegistration"
)

func InitialRegistration(BASE_URL string, DEVICE_ID string, SECRET_KEY string, DEVICE_FROM string, ENCRYPTION_KEY string) {

	DeviceRegistrationHelper.RegisterDevice(BASE_URL, DEVICE_ID, SECRET_KEY, DEVICE_FROM, false)

	status, err := DeviceRegistrationHelper.CheckDeviceStatus(BASE_URL, DEVICE_ID, SECRET_KEY)

	if err != nil {
		fmt.Println("Error fetching the Status from DeviceStatusCheck")
		return
	}

	if status.Data.Status != "Certificate Available" {
		fmt.Println("Invalid Device status")
		return
	} else {
		// DOWNLOAD AND STORE THE CERTIFICATE
		err := DeviceRegistrationHelper.RemoveAllCertFiles()
		if err != nil {
			fmt.Println("Error removing cert files:", err)
		}
		response, err := DeviceRegistrationHelper.DownloadDeviceCertificate(BASE_URL, DEVICE_ID, SECRET_KEY, ENCRYPTION_KEY)
		if err != nil {
			fmt.Println("Error fetching the certificate from DownloadDeviceCertificate")
			return
		}

		if response != nil && response.StatusCode == 200 {
			DeviceRegistrationHelper.ReturnDownloadAcknowledgement(BASE_URL, DEVICE_ID, SECRET_KEY)
		} else {
			fmt.Println("No certificate found for the device")
			return
		}
	}

}

func ReRegistration(BASE_URL string, DEVICE_ID string, SECRET_KEY string, DEVICE_FROM string, ENCRYPTION_KEY string) {

	// FIRST REMOVE THE EXISTING CERTIFICATES FROM THE DEVICE
	err := DeviceRegistrationHelper.RemoveAllCertFiles()
	if err != nil {
		fmt.Println("Error removing cert files:", err)
	}

	// THEN RE-REGISTER THE DEVICE
	DeviceRegistrationHelper.RegisterDevice(BASE_URL, DEVICE_ID, SECRET_KEY, DEVICE_FROM, true)

	status, err := DeviceRegistrationHelper.CheckDeviceStatus(BASE_URL, DEVICE_ID, SECRET_KEY)

	if err != nil {
		fmt.Println("Error fetching the Status from DeviceStatusCheck")
		return
	}

	if status.Data.Status != "Certificate Available" {
		fmt.Println("Invalid Device status")
		return
	} else {
		// DOWNLOAD AND STORE THE CERTIFICATE
		response, err := DeviceRegistrationHelper.DownloadDeviceCertificate(BASE_URL, DEVICE_ID, SECRET_KEY, ENCRYPTION_KEY)
		if err != nil {
			fmt.Println("Error fetching the certificate from DownloadDeviceCertificate")
			return
		}

		if response != nil && response.StatusCode == 200 {
			DeviceRegistrationHelper.ReturnDownloadAcknowledgement(BASE_URL, DEVICE_ID, SECRET_KEY)
		} else {
			fmt.Println("No certificate found for the device")
			return
		}
	}

}

func DownloadCertificateAfterAdminApproval(BASE_URL string, DEVICE_ID string, SECRET_KEY string, DEVICE_FROM string, ENCRYPTION_KEY string) {

	status, err := DeviceRegistrationHelper.CheckDeviceStatus(BASE_URL, DEVICE_ID, SECRET_KEY)

	if err != nil {
		fmt.Println("Error fetching the Status from DeviceStatusCheck")
		return
	}

	if status.Data.Status != "Certificate Available" {
		fmt.Println("Invalid Device status")
		return
	} else {
		err := DeviceRegistrationHelper.RemoveAllCertFiles()
		if err != nil {
			fmt.Println("Error removing cert files:", err)
		}
		// DOWNLOAD AND STORE THE CERTIFICATE
		response, err := DeviceRegistrationHelper.DownloadDeviceCertificate(BASE_URL, DEVICE_ID, SECRET_KEY, ENCRYPTION_KEY)
		if err != nil {
			fmt.Println("Error fetching the certificate from DownloadDeviceCertificate")
			return
		}

		if response != nil && response.StatusCode == 200 {
			DeviceRegistrationHelper.ReturnDownloadAcknowledgement(BASE_URL, DEVICE_ID, SECRET_KEY)
		} else {
			fmt.Println("No certificate found for the device")
			return
		}
	}
}
