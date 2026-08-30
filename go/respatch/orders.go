package respatch

import (
	"encoding/json"
	"fmt"
)

type OrdersService struct {
	client *Client
}

func (s *OrdersService) Quote(req *GetQuoteRequest) (map[string]interface{}, error) {
	resp, err := s.client.sendRequest("POST", "/api/v1/delivery-orders/quote", req)
	if err != nil {
		return nil, err
	}
	var result map[string]interface{}
	err = json.Unmarshal(resp, &result)
	return result, err
}

// CalculateFare is deprecated — use Quote.
func (s *OrdersService) CalculateFare(req *GetQuoteRequest) (map[string]interface{}, error) {
	resp, err := s.client.sendRequest("POST", "/api/v1/delivery-orders/calculate-fare", req)
	if err != nil {
		return nil, err
	}
	var result map[string]interface{}
	err = json.Unmarshal(resp, &result)
	return result, err
}

func (s *OrdersService) Create(req *CreateDeliveryOrderRequest) (map[string]interface{}, error) {
	resp, err := s.client.sendRequest("POST", "/api/v1/delivery-orders", req)
	if err != nil {
		return nil, err
	}
	var result map[string]interface{}
	err = json.Unmarshal(resp, &result)
	return result, err
}

// BookFromQuote creates an order from a prior quote (single HTTP call).
func (s *OrdersService) BookFromQuote(req *CreateDeliveryOrderRequest) (map[string]interface{}, error) {
	return s.Create(req)
}

func (s *OrdersService) List() ([]map[string]interface{}, error) {
	resp, err := s.client.sendRequest("GET", "/api/v1/delivery-orders", nil)
	if err != nil {
		return nil, err
	}
	var result []map[string]interface{}
	err = json.Unmarshal(resp, &result)
	return result, err
}

func (s *OrdersService) Retrieve(id string) (map[string]interface{}, error) {
	path := fmt.Sprintf("/api/v1/delivery-orders/%s", id)
	resp, err := s.client.sendRequest("GET", path, nil)
	if err != nil {
		return nil, err
	}
	var result map[string]interface{}
	err = json.Unmarshal(resp, &result)
	return result, err
}

func (s *OrdersService) Track(id string) (map[string]interface{}, error) {
	path := fmt.Sprintf("/api/v1/delivery-orders/%s/track", id)
	resp, err := s.client.sendRequest("GET", path, nil)
	if err != nil {
		return nil, err
	}
	var result map[string]interface{}
	err = json.Unmarshal(resp, &result)
	return result, err
}

func (s *OrdersService) Cancel(id string) error {
	path := fmt.Sprintf("/api/v1/delivery-orders/%s/cancel", id)
	_, err := s.client.sendRequest("POST", path, nil)
	return err
}
