package PublishService

import (
	"crypto/tls"
	"crypto/x509"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math/rand"
	"os"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	NetworkHelper "github.com/manjunath-ravindra/go-edge-device/helpers/network"
)

func PublishMqttMessagesSerivce(IOT_ENDPOINT string) {
	endpoint := IOT_ENDPOINT
	clientID := "GO-TEST_CLIENT"
	topic := "test/topic"
	certFile := "certs/GO-TEST_certificate.pem.crt"
	keyFile := "certs/GO-TEST_private.pem.key"
	caFile := "certs/GO-TEST_AmazonRootCA1.pem"

	// Initialize network checker and start monitoring
	networkChecker := NetworkHelper.NewNetworkChecker()
	networkChecker.StartMonitoring()
	defer networkChecker.StopMonitoring()

	// Load device certificate and private key
	cert, err := tls.LoadX509KeyPair(certFile, keyFile)
	if err != nil {
		fmt.Printf("Error loading X509 key pair: %v\n", err)
		os.Exit(1)
	}

	// Load CA cert
	caCert, err := ioutil.ReadFile(caFile)
	if err != nil {
		fmt.Printf("Error reading CA file: %v\n", err)
		os.Exit(1)
	}
	caCertPool := x509.NewCertPool()
	caCertPool.AppendCertsFromPEM(caCert)

	tlsConfig := &tls.Config{
		Certificates:       []tls.Certificate{cert},
		RootCAs:            caCertPool,
		InsecureSkipVerify: false,
	}

	opts := mqtt.NewClientOptions()
	opts.AddBroker(fmt.Sprintf("tcps://%s", endpoint))
	opts.SetClientID(clientID)
	opts.SetTLSConfig(tlsConfig)
	opts.SetResumeSubs(true)
	opts.SetOrderMatters(true) // Maintain order of messages
	opts.SetMessageChannelDepth(1000)
	opts.SetAutoReconnect(true)

	client := mqtt.NewClient(opts)
	if token := client.Connect(); token.Wait() && token.Error() != nil {
		fmt.Printf("Error connecting to AWS IoT: %v\n", token.Error())
		os.Exit(1)
	}
	fmt.Println("Connected to AWS IoT!")
	fmt.Println("Network monitoring started. Publishing will pause when network is down.")

	type Payload struct {
		Message   string `json:"message"`
		Timestamp int64  `json:"timestamp"`
		DeviceID  string `json:"deviceId"`
	}

	rand.Seed(time.Now().UnixNano())

	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()

	fmt.Println("Publishing incremental data every second. Press Ctrl+C to exit.")
	counter := 0
	lastNetworkStatus := true // Track network status changes

	for t := range ticker.C {
		// Check network connectivity before publishing (immediate check)
		currentNetworkStatus := networkChecker.IsNetworkAvailableImmediate()

		// Detect network status changes for immediate response
		if currentNetworkStatus != lastNetworkStatus {
			if !currentNetworkStatus {
				fmt.Println("⚠️  Network connection lost! Pausing message publishing immediately...")
			} else {
				fmt.Println("✅ Network connection restored! Resuming message publishing...")
			}
			lastNetworkStatus = currentNetworkStatus
		}

		if !currentNetworkStatus {
			// Wait for network to be restored (max 5 minutes)
			if networkChecker.WaitForNetwork(5 * time.Minute) {
				fmt.Println("Network restored. Resuming message publishing...")
				lastNetworkStatus = true
			} else {
				fmt.Println("Network still down after 5 minutes. Continuing to wait...")
				continue
			}
		}

		payload := Payload{
			Message:   fmt.Sprintf("Incremental value: %d", counter),
			Timestamp: t.Unix(),
			DeviceID:  clientID,
		}

		payloadBytes, err := json.Marshal(payload)
		if err != nil {
			fmt.Printf("Error marshaling payload: %v\n", err)
			continue
		}

		// Double-check network before publishing (immediate check)
		if !networkChecker.IsNetworkAvailableImmediate() {
			fmt.Println("Network lost during message preparation. Skipping publish...")
			continue
		}

		token := client.Publish(topic, 0, false, payloadBytes)
		token.Wait()
		fmt.Printf("Published JSON message to topic '%s': %s\n", topic, string(payloadBytes))
		counter++
	}

	// Give time for the message to be sent before disconnecting
	time.Sleep(2 * time.Second)
	client.Disconnect(250)
	fmt.Println("Disconnected.")
}
