package DeviceRegistrationHelper

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"

	CryptoTypes "github.com/manjunath-ravindra/go-edge-device/types/crypto"
	DeviceTypes "github.com/manjunath-ravindra/go-edge-device/types/device"
	CryptoVendors "github.com/manjunath-ravindra/go-edge-device/vendors/crypto"
)

func RegisterDevice(BASE_URL string, DEVICE_ID string, SECRET_KEY string, DEVICE_FROM string) {
	url := BASE_URL + "/device/registration"
	body := map[string]interface{}{
		"deviceId":   DEVICE_ID,
		"secretKey":  SECRET_KEY,
		"deviceFrom": DEVICE_FROM,
		"reRegister": false,
	}

	jsonBody, err := json.Marshal(body)
	if err != nil {
		fmt.Println("Error marshaling JSON: ", err)
		return
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonBody))
	if err != nil {
		fmt.Println("Error creating request:", err)
		return
	}

	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}

	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error making POST request:", err)
		return
	}
	defer resp.Body.Close()

	fmt.Println("Response Status: ", resp.Status)
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error reading response body:", err)
		return
	}
	fmt.Println("Response Body:", string(respBody))
}

func CheckDeviceStatus(BASE_URL string, DEVICE_ID string, SECRET_KEY string) (*DeviceTypes.DeviceStatusResponse, error) {
	url := BASE_URL + "/device/status"
	// Add query parameters
	url = url + "?deviceId=" + DEVICE_ID + "&secretKey=" + SECRET_KEY

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		fmt.Println("Error creating request: ", err)
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}

	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error making GET request: ", err)
		return nil, err
	}

	defer resp.Body.Close()

	fmt.Println("Response Status: ", resp.Status)
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error reading response body: ", err)
		return nil, err
	}
	var result DeviceTypes.DeviceStatusResponse
	if err := json.Unmarshal(respBody, &result); err != nil {
		fmt.Println("Error unmarshaling JSON: ", err)
		return nil, err
	}
	return &result, nil
}

func DownloadDeviceCertificate(BASE_URL string, DEVICE_ID string, SECRET_KEY string, ENCRYPTION_KEY string) (*DeviceTypes.DeviceCertificateResponse, error) {
	url := BASE_URL + "/device/certificate"
	// Add query parameters
	url = url + "?deviceId=" + DEVICE_ID + "&secretKey=" + SECRET_KEY

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		fmt.Println("Error creating request: ", err)
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}

	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error making GET request: ", err)
		return nil, err
	}

	defer resp.Body.Close()

	fmt.Println("Response Status: ", resp.Status)
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error reading response body:", err)
		return nil, err
	}
	var result DeviceTypes.DeviceCertificateResponse
	if err := json.Unmarshal(respBody, &result); err != nil {
		fmt.Println("Error unmarshaling JSON: ", err)
		return nil, err
	}

	if ENCRYPTION_KEY == "" {
		fmt.Println("ENCRYPTION_KEY is empty")
		return &result, fmt.Errorf("ENCRYPTION_KEY is empty")
	}

	decrypted, err := CryptoVendors.DecryptResponse(result.TxId, result.Data.EncryptedData, result.Data.IV, ENCRYPTION_KEY)
	if err != nil {
		fmt.Println("Decryption error:", err)
		return &result, err
	}

	// Unmarshal decrypted JSON
	var certData CryptoTypes.CertData
	err = json.Unmarshal([]byte(decrypted), &certData)
	if err != nil {
		fmt.Println("Error unmarshaling decrypted certificate data:", err)
		return &result, err
	}

	// Save files with device ID prefix in 'certs' directory
	certsDir := "certs"
	err = os.MkdirAll(certsDir, 0700)
	if err != nil {
		fmt.Println("Error creating certs directory:", err)
		return &result, err
	}
	certFile := certsDir + "/" + certData.DeviceId + "_certificate.pem"
	keyFile := certsDir + "/" + certData.DeviceId + "_private.key"
	caFile := certsDir + "/" + certData.DeviceId + "_AmazonRootCA1.pem"

	err = os.WriteFile(certFile, []byte(certData.CertificatePem), 0600)
	if err != nil {
		fmt.Println("Error writing certificate PEM:", err)
	}
	err = os.WriteFile(keyFile, []byte(certData.CertificateKey), 0600)
	if err != nil {
		fmt.Println("Error writing certificate key:", err)
	}
	err = os.WriteFile(caFile, []byte(certData.PrivateCA), 0600)
	if err != nil {
		fmt.Println("Error writing private CA:", err)
	}

	return &result, nil
}

func ReturnDownloadAcknowledgement(BASE_URL string, DEVICE_ID string, SECRET_KEY string) (*DeviceTypes.DownloadAcknowledgeReponse, error) {
	url := BASE_URL + "/device/cert_status"

	body := map[string]interface{}{
		"deviceId":  DEVICE_ID,
		"secretKey": SECRET_KEY,
	}

	jsonBody, err := json.Marshal(body)
	if err != nil {
		fmt.Println("Error marshaling JSON: ", err)
		return nil, err
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonBody))
	if err != nil {
		fmt.Println("Error creating request:", err)
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}

	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error making POST request: ", err)
		return nil, err
	}

	defer resp.Body.Close()

	fmt.Println("Response Status: ", resp.Status)
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error reading response body:", err)
		return nil, err
	}
	var result DeviceTypes.DownloadAcknowledgeReponse
	if err := json.Unmarshal(respBody, &result); err != nil {
		fmt.Println("Error unmarshaling JSON: ", err)
		return nil, err
	}

	return &result, nil
}
