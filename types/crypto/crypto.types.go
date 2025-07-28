package CryptoTypes

// CertData represents the decrypted certificate data for a device
type CertData struct {
	DeviceId       string `json:"deviceId"`
	CertArn        string `json:"certArn"`
	CertificatePem string `json:"certificatePem"`
	CertificateKey string `json:"certificateKey"`
	PrivateCA      string `json:"privateCA"`
	AwsEndpointUrl string `json:"awsEndpointUrl"`
}
