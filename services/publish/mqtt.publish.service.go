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
)

func PublishMqttMessagesSerivce(IOT_ENDPOINT string) {
	endpoint := IOT_ENDPOINT
	clientID := "GO-TEST_CLIENT"
	topic := "test/topic"
	certFile := "certs/GO-TEST_certificate.pem.crt"
	keyFile := "certs/GO-TEST_private.pem.key"
	caFile := "certs/GO-TEST_AmazonRootCA1.pem"

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
	for t := range ticker.C {
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
