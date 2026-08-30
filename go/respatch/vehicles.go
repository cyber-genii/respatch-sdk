package respatch

import "encoding/json"

type VehiclesService struct {
	client *Client
}

func (s *VehiclesService) GetTypes() ([]map[string]interface{}, error) {
	resp, err := s.client.sendRequest("GET", "/api/v1/vehicles/types", nil)
	if err != nil {
		return nil, err
	}

	var result []map[string]interface{}
	err = json.Unmarshal(resp, &result)
	return result, err
}
