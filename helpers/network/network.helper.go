package NetworkHelper

import (
	"fmt"
	"net"
	"net/http"
	"time"
)

// NetworkStatus represents the current network connectivity status
type NetworkStatus struct {
	IsConnected bool
	LastCheck   time.Time
	Error       error
}

// NetworkChecker provides network connectivity checking functionality
type NetworkChecker struct {
	status        NetworkStatus
	checkURLs     []string
	timeout       time.Duration
	checkInterval time.Duration
	stopChan      chan bool
	isMonitoring  bool
}

// NewNetworkChecker creates a new network checker instance
func NewNetworkChecker() *NetworkChecker {
	return &NetworkChecker{
		checkURLs: []string{
			"https://www.google.com",
			"https://www.cloudflare.com",
			"https://www.amazon.com",
		},
		timeout:       3 * time.Second, // Reduced timeout for faster detection
		checkInterval: 5 * time.Second, // Much faster check interval for immediate detection
		status: NetworkStatus{
			IsConnected: false,
			LastCheck:   time.Time{},
		},
		stopChan:     make(chan bool),
		isMonitoring: false,
	}
}

// IsNetworkAvailable checks if the network is currently available
func (nc *NetworkChecker) IsNetworkAvailable() bool {
	// Always perform a fresh check for immediate detection
	nc.status = nc.checkConnectivity()
	return nc.status.IsConnected
}

// checkConnectivity performs the actual network connectivity check
func (nc *NetworkChecker) checkConnectivity() NetworkStatus {
	status := NetworkStatus{
		LastCheck:   time.Now(),
		IsConnected: false,
	}

	// Use concurrent checks for faster detection
	results := make(chan bool, len(nc.checkURLs)+2) // +2 for DNS and ping checks

	// Check URLs concurrently
	for _, url := range nc.checkURLs {
		go func(u string) {
			results <- nc.checkURL(u)
		}(url)
	}

	// Check DNS resolution concurrently
	go func() {
		results <- nc.checkDNSResolution()
	}()

	// Check ping concurrently
	go func() {
		results <- nc.checkPing()
	}()

	// Wait for any successful check or all to fail
	successCount := 0
	for i := 0; i < len(nc.checkURLs)+2; i++ {
		if <-results {
			successCount++
		}
	}

	if successCount > 0 {
		status.IsConnected = true
		return status
	}

	status.Error = fmt.Errorf("all connectivity checks failed")
	return status
}

// checkURL checks connectivity to a specific URL
func (nc *NetworkChecker) checkURL(url string) bool {
	client := &http.Client{
		Timeout: nc.timeout,
	}

	resp, err := client.Get(url)
	if err != nil {
		return false
	}
	defer resp.Body.Close()

	return resp.StatusCode >= 200 && resp.StatusCode < 400
}

// checkDNSResolution checks if DNS resolution is working
func (nc *NetworkChecker) checkDNSResolution() bool {
	_, err := net.LookupHost("google.com")
	return err == nil
}

// checkPing performs a lightweight ping check
func (nc *NetworkChecker) checkPing() bool {
	// Try to connect to a well-known port on a reliable host
	conn, err := net.DialTimeout("tcp", "8.8.8.8:53", 2*time.Second)
	if err != nil {
		return false
	}
	defer conn.Close()
	return true
}

// GetNetworkStatus returns the current network status
func (nc *NetworkChecker) GetNetworkStatus() NetworkStatus {
	return nc.status
}

// ForceCheck forces an immediate network connectivity check
func (nc *NetworkChecker) ForceCheck() NetworkStatus {
	nc.status = nc.checkConnectivity()
	return nc.status
}

// IsNetworkAvailableImmediate performs an immediate network check without any caching
func (nc *NetworkChecker) IsNetworkAvailableImmediate() bool {
	return nc.checkConnectivity().IsConnected
}

// WaitForNetwork waits for network connectivity to be restored
func (nc *NetworkChecker) WaitForNetwork(maxWait time.Duration) bool {
	startTime := time.Now()

	for time.Since(startTime) < maxWait {
		if nc.IsNetworkAvailable() {
			return true
		}
		time.Sleep(1 * time.Second) // Check every second for faster response
	}

	return false
}

// StartMonitoring starts continuous network monitoring in a goroutine
func (nc *NetworkChecker) StartMonitoring() {
	if nc.isMonitoring {
		return
	}

	nc.isMonitoring = true
	go func() {
		ticker := time.NewTicker(nc.checkInterval)
		defer ticker.Stop()

		for {
			select {
			case <-ticker.C:
				nc.ForceCheck()
			case <-nc.stopChan:
				return
			}
		}
	}()
}

// StopMonitoring stops the continuous network monitoring
func (nc *NetworkChecker) StopMonitoring() {
	if !nc.isMonitoring {
		return
	}

	nc.stopChan <- true
	nc.isMonitoring = false
}

// IsMonitoring returns whether continuous monitoring is active
func (nc *NetworkChecker) IsMonitoring() bool {
	return nc.isMonitoring
}
