package respatch

import (
	"encoding/json"
	"fmt"
)

type SandboxService struct {
	client *Client
}

func (s *SandboxService) CreateData(req map[string]interface{}) (map[string]interface{}, error) {
	resp, err := s.client.sendRequest("POST", "/api/v1/developers/sandbox", req)
	if err != nil {
		return nil, err
	}

	var result map[string]interface{}
	err = json.Unmarshal(resp, &result)
	return result, err
}

func (s *SandboxService) ListData() ([]map[string]interface{}, error) {
	resp, err := s.client.sendRequest("GET", "/api/v1/developers/sandbox", nil)
	if err != nil {
		return nil, err
	}

	var result []map[string]interface{}
	err = json.Unmarshal(resp, &result)
	return result, err
}

func (s *SandboxService) ClearData() error {
	_, err := s.client.sendRequest("DELETE", "/api/v1/developers/sandbox", nil)
	return err
}

func (s *SandboxService) SimulateEvent(orderID string, event string) (map[string]interface{}, error) {
	path := fmt.Sprintf("/api/v1/sandbox/orders/%s/simulate/%s", orderID, event)
	resp, err := s.client.sendRequest("POST", path, map[string]interface{}{})
	if err != nil {
		return nil, err
	}

	var result map[string]interface{}
	err = json.Unmarshal(resp, &result)
	return result, err
}
