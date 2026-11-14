package respatch

import (
	"encoding/json"
	"fmt"
	"net/url"
)

type WebhooksService struct {
	client *Client
}

func (s *WebhooksService) List() ([]map[string]interface{}, error) {
	resp, err := s.client.sendRequest("GET", "/api/v1/developers/webhooks", nil)
	if err != nil {
		return nil, err
	}

	var result []map[string]interface{}
	err = json.Unmarshal(resp, &result)
	return result, err
}

func (s *WebhooksService) Create(req *CreateWebhookRequest) (map[string]interface{}, error) {
	resp, err := s.client.sendRequest("POST", "/api/v1/developers/webhooks", req)
	if err != nil {
		return nil, err
	}

	var result map[string]interface{}
	err = json.Unmarshal(resp, &result)
	return result, err
}

func (s *WebhooksService) Retrieve(id string) (map[string]interface{}, error) {
	path := fmt.Sprintf("/api/v1/developers/webhooks/%s", id)
	resp, err := s.client.sendRequest("GET", path, nil)
	if err != nil {
		return nil, err
	}

	var result map[string]interface{}
	err = json.Unmarshal(resp, &result)
	return result, err
}

func (s *WebhooksService) Delete(id string) error {
	path := fmt.Sprintf("/api/v1/developers/webhooks/%s", id)
	_, err := s.client.sendRequest("DELETE", path, nil)
	return err
}

func (s *WebhooksService) DLQ(params map[string]string) (interface{}, error) {
	path := "/api/v1/developers/webhook-events"
	if len(params) > 0 {
		q := url.Values{}
		for k, v := range params {
			q.Set(k, v)
		}
		path += "?" + q.Encode()
	}
	resp, err := s.client.sendRequest("GET", path, nil)
	if err != nil {
		return nil, err
	}

	var result interface{}
	err = json.Unmarshal(resp, &result)
	return result, err
}

func (s *WebhooksService) ReplayDLQ(eventID string) error {
	path := fmt.Sprintf("/api/v1/developers/webhook-events/%s/retry", eventID)
	_, err := s.client.sendRequest("POST", path, nil)
	return err
}
