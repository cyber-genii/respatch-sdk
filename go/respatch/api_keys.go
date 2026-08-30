package respatch

import (
	"encoding/json"
	"fmt"
)

type APIKeysService struct {
	client *Client
}

func (s *APIKeysService) List() ([]map[string]interface{}, error) {
	resp, err := s.client.sendRequest("GET", "/api/v1/developers/api-keys", nil)
	if err != nil {
		return nil, err
	}

	var result []map[string]interface{}
	err = json.Unmarshal(resp, &result)
	return result, err
}

func (s *APIKeysService) Create(req map[string]interface{}) (map[string]interface{}, error) {
	resp, err := s.client.sendRequest("POST", "/api/v1/developers/api-keys", req)
	if err != nil {
		return nil, err
	}

	var result map[string]interface{}
	err = json.Unmarshal(resp, &result)
	return result, err
}

func (s *APIKeysService) Revoke(id string) error {
	path := fmt.Sprintf("/api/v1/developers/api-keys/%s/revoke", id)
	_, err := s.client.sendRequest("POST", path, nil)
	return err
}

func (s *APIKeysService) Delete(id string) error {
	path := fmt.Sprintf("/api/v1/developers/api-keys/%s", id)
	_, err := s.client.sendRequest("DELETE", path, nil)
	return err
}
