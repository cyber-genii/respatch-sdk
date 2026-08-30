package respatch

import (
	"encoding/json"
	"fmt"
)

type DeliveriesService struct {
	client *Client
}

// Deprecated: use OrdersService
func (s *DeliveriesService) Create(req *CreateDeliveryOrderRequest) (map[string]interface{}, error) {
	resp, err := s.client.sendRequest("POST", "/api/v1/delivery-orders", req)
	if err != nil {
		return nil, err
	}
	var result map[string]interface{}
	err = json.Unmarshal(resp, &result)
	return result, err
}

func (s *DeliveriesService) List() ([]map[string]interface{}, error) {
	resp, err := s.client.sendRequest("GET", "/api/v1/delivery-orders", nil)
	if err != nil {
		return nil, err
	}
	var result []map[string]interface{}
	err = json.Unmarshal(resp, &result)
	return result, err
}

func (s *DeliveriesService) Retrieve(id string) (map[string]interface{}, error) {
	path := fmt.Sprintf("/api/v1/delivery-orders/%s", id)
	resp, err := s.client.sendRequest("GET", path, nil)
	if err != nil {
		return nil, err
	}
	var result map[string]interface{}
	err = json.Unmarshal(resp, &result)
	return result, err
}

func (s *DeliveriesService) Quote(req *CalculateOrderFareRequest) (map[string]interface{}, error) {
	resp, err := s.client.sendRequest("POST", "/api/v1/delivery-orders/calculate-fare", req)
	if err != nil {
		return nil, err
	}
	var result map[string]interface{}
	err = json.Unmarshal(resp, &result)
	return result, err
}
