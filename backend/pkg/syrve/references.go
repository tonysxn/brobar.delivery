package syrve

import (
	"context"
	"encoding/json"
	"net/http"
)

// -----------------------------
// Terminal Groups
// -----------------------------

type TerminalGroupsResponse struct {
	TerminalGroups []TerminalGroup `json:"terminalGroups"`
}

type TerminalGroup struct {
	ID    string     `json:"id"`
	Name  string     `json:"name"`
	Items []Terminal `json:"items"`
}

type Terminal struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

func (c *Client) GetTerminalGroups(ctx context.Context, authToken, organizationID string) (*TerminalGroupsResponse, error) {
	// Endpoint: /api/1/terminal_groups
	// Body: { "organizationIds": ["..."] }

	body := map[string][]string{
		"organizationIds": {organizationID},
	}

	respBytes, err := c.doRequest(ctx, http.MethodPost, "/terminal_groups", authToken, body)
	if err != nil {
		return nil, err
	}

	var resp TerminalGroupsResponse
	if err := json.Unmarshal(respBytes, &resp); err != nil {
		return nil, err
	}

	return &resp, nil
}

// -----------------------------
// Order Types
// -----------------------------

type OrderTypesResponse struct {
	OrderTypes []OrderTypeGroup `json:"orderTypes"`
}

type OrderTypeGroup struct {
	Items []OrderType `json:"items"`
}

type OrderType struct {
	ID             string `json:"id"`
	Name           string `json:"name"`
	OrderServiceType string `json:"orderServiceType"` // Delivery, Common, etc
}

func (c *Client) GetOrderTypes(ctx context.Context, authToken, organizationID string) (*OrderTypesResponse, error) {
	// Endpoint: /api/1/deliveries/order_types
	// Body: { "organizationIds": ["..."] }
    // Legacy generic path: /api/1/corporation/order_types or just /api/1/order_types? 
    // Docs say /api/1/corporation/order_types needs structure. 
    // Let's try /api/1/deliveries/order_types which is common for delivery integrations, OR just replicate legacy logic.
    // Legacy PHP `Syrve::orderTypes` likely calls `/api/1/corporation/order_types`
    
    // Let's assume /api/1/corporation/order_types for now as it returns groups.
	body := map[string][]string{
		"organizationIds": {organizationID},
	}

	respBytes, err := c.doRequest(ctx, http.MethodPost, "/deliveries/order_types", authToken, body)
	if err != nil {
		return nil, err
	}

	var resp OrderTypesResponse
	if err := json.Unmarshal(respBytes, &resp); err != nil {
		return nil, err
	}

	return &resp, nil
}

// -----------------------------
// Restaurant Sections (Tables)
// -----------------------------

type RestaurantSectionsResponse struct {
	RestaurantSections []RestaurantSection `json:"restaurantSections"`
}

type RestaurantSection struct {
	ID   string  `json:"id"`
	Name string  `json:"name"`
	Tables []Table `json:"tables"`
}

type Table struct {
	ID       string `json:"id"`
	Number   int    `json:"number"`
	Name     string `json:"name"`
	SeatingCapacity int `json:"seatingCapacity"`
}

func (c *Client) GetRestaurantSections(ctx context.Context, authToken string, terminalGroupIDs []string) (*RestaurantSectionsResponse, error) {
	// Endpoint: /api/1/reserve/available_restaurant_sections
	// Body: { "terminalGroupIds": ["..."] }
	
	body := map[string][]string{
		"terminalGroupIds": terminalGroupIDs,
	}

	respBytes, err := c.doRequest(ctx, http.MethodPost, "/reserve/available_restaurant_sections", authToken, body)
	if err != nil {
		return nil, err
	}

	var resp RestaurantSectionsResponse
	if err := json.Unmarshal(respBytes, &resp); err != nil {
		return nil, err
	}

	return &resp, nil
}
