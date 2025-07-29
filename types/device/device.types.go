package DeviceTypes

// DeviceStatusResponse represents the response from the device status check API
type DeviceStatusResponse struct {
	StatusCode int `json:"statusCode"`
	Data       struct {
		Status string `json:"status"`
	} `json:"data"`
	TxId string `json:"txId"`
}

// DeviceCertificateResponse represents the response from the device certificate download API
type DeviceCertificateResponse struct {
	StatusCode int `json:"statusCode"`
	Data       struct {
		IV            string `json:"iv"`
		EncryptedData string `json:"encryptedData"`
	} `json:"data"`
	TxId string `json:"txId"`
}

type DownloadAcknowledgeReponse struct {
	StatusCode int `json:"statusCode"`
	Data       struct {
		Status  string  `json:"status"`
		Message *string `json:"message,omitempty"`
	} `json:"data"`
	TxId string `json:"txId"`
}

type DeviceStatus string

const (
	Failed               DeviceStatus = "Failed"
	CertificateAvailable DeviceStatus = "Certificate Available"
	DownloadComplete     DeviceStatus = "Download Complete"
	AdminApprovalPending DeviceStatus = "Admin Approval Pending"
	AdminApproved        DeviceStatus = "Admin Approved"
	AdminRejected        DeviceStatus = "Admin Rejected"
	Deregistered         DeviceStatus = "Deregistered"
	Register             DeviceStatus = "Register"
	Pending              DeviceStatus = "Pending"
)
