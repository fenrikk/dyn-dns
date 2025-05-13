package iplocator

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

// IPifyLocator implements the IPLocator interface using the ipify.org API
type IPifyResponse struct {
	IP string `json:"ip"`
}

type IPifyLocator struct {
	apiURL     string
	httpClient *http.Client
}

func NewIPifyLocator() *IPifyLocator {
	return &IPifyLocator{
		apiURL: "https://api.ipify.org?format=json",
		httpClient: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

func (i *IPifyLocator) GetCurrentIP() (string, error) {
	resp, err := i.httpClient.Get(i.apiURL)
	if err != nil {
		return "", fmt.Errorf("error while executing HTTP request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("unexpected response code: %d", resp.StatusCode)
	}

	var ipResp IPifyResponse
	if err := json.NewDecoder(resp.Body).Decode(&ipResp); err != nil {
		return "", fmt.Errorf("error while decoding JSON: %w", err)
	}

	if ipResp.IP == "" {
		return "", fmt.Errorf("empty IP address received")
	}

	return ipResp.IP, nil
}
