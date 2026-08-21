package transmit

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

type Client struct {
	BaseURL    string
	APIKey     string
	HTTPClient *http.Client

	Orders     *OrdersService
	Deliveries *DeliveriesService
	APIKeys    *APIKeysService
	Webhooks   *WebhooksService
	Sandbox    *SandboxService
	Vehicles   *VehiclesService
}

type Options struct {
	APIKey      string
	Environment string // "production" or "sandbox"
	BaseURL     string // Optional override
}

// NewClient creates a new Transmit API Client
func NewClient(options Options) (*Client, error) {
	if options.APIKey == "" {
		return nil, fmt.Errorf("API Key is required to initialize the Respatch SDK")
	}

	baseURL := options.BaseURL
	if baseURL == "" {
		// Public frontiers only (respatch-sandbox / respatch-api), not the platform backend.
		if options.Environment == "sandbox" {
			baseURL = "https://sandbox-api.respatch.com"
		} else {
			baseURL = "https://api.respatch.com"
		}
	}

	c := &Client{
		BaseURL: baseURL,
		APIKey:  options.APIKey,
		HTTPClient: &http.Client{
			Timeout: time.Second * 30,
		},
	}

	c.Orders = &OrdersService{client: c}
	c.Deliveries = &DeliveriesService{client: c}
	c.APIKeys = &APIKeysService{client: c}
	c.Webhooks = &WebhooksService{client: c}
	c.Sandbox = &SandboxService{client: c}
	c.Vehicles = &VehiclesService{client: c}

	return c, nil
}

func (c *Client) sendRequest(method, path string, body interface{}) ([]byte, error) {
	url := c.BaseURL + path
	var reqBody io.Reader

	if body != nil {
		jsonBody, err := json.Marshal(body)
		if err != nil {
			return nil, err
		}
		reqBody = bytes.NewBuffer(jsonBody)
	}

	req, err := http.NewRequest(method, url, reqBody)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", "Bearer "+c.APIKey)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("User-Agent", "Respatch-Go-SDK/1.0.0")

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, fmt.Errorf("transmit api error: %d - %s", resp.StatusCode, string(respBody))
	}

	return unwrapData(respBody), nil
}
