package service

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/tonysanin/brobar/pkg/helpers"
)

type OrderDTO struct {
	ID            string `json:"id"`
	Name          string `json:"name"`
	Phone         string `json:"phone"`
	Address       string `json:"address"`
	Coords        string `json:"coords"`
	Entrance      string `json:"entrance"`
	Floor         string `json:"floor"`
	Flat          string `json:"flat"`
	Wishes        string `json:"wishes"`
	AddressWishes string `json:"address_wishes"`
	DeliveryDoor  bool   `json:"delivery_door"`
}

type OrderClient struct {
	baseURL string
	client  *http.Client
}

func NewOrderClient() *OrderClient {
	return &OrderClient{
		baseURL: helpers.GetEnv("ORDER_SERVICE_URL", "http://order-service:3001"),
		client:  &http.Client{Timeout: 10 * time.Second},
	}
}

func (c *OrderClient) GetOrder(id string) (*OrderDTO, error) {
	url := fmt.Sprintf("%s/orders/%s", c.baseURL, id)
	resp, err := c.client.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to get order: status %d", resp.StatusCode)
	}

	var result struct {
		Success bool     `json:"success"`
		Data    OrderDTO `json:"data"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}

	return &result.Data, nil
}
