package respatch

import (
	"encoding/json"
	"fmt"
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
