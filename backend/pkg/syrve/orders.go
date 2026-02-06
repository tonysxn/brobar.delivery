package syrve

import (
	"context"
	"encoding/json"
	"net/http"
)

// -----------------------------
// Orders (Get by Tables)
// -----------------------------

type TableOrdersResponse struct {
	Orders []TableOrder `json:"orders"`
}

type TableOrder struct {
	OrderID   string      `json:"orderId"`
	Timestamp int64       `json:"timestamp"` // Last update time?
	Order     TableOrderInfo `json:"order"`
}

type TableOrderInfo struct {
	ID       string   `json:"id"`
	TableIDs []string `json:"tableIds"`
	// Add more fields if needed for "is free" check
}

type GetOrdersByTablesRequest struct {
	OrganizationIDs []string `json:"organizationIds"`
	TableIDs        []string `json:"tableIds,omitempty"`
	Statuses        []string `json:"statuses,omitempty"`
	DateFrom        string   `json:"dateFrom,omitempty"`
	DateTo          string   `json:"dateTo,omitempty"`
}

func (c *Client) GetOrdersByTables(ctx context.Context, authToken, organizationID string, tableIDs []string, dateFrom, dateTo string) (*TableOrdersResponse, error) {
	// Endpoint: /api/1/order/by_table
	// This endpoint allows filtering by tables to see active orders.

	reqBody := GetOrdersByTablesRequest{
		OrganizationIDs: []string{organizationID},
		TableIDs:        tableIDs,
		DateFrom:        dateFrom,
		DateTo:          dateTo,
	}

	respBytes, err := c.doRequest(ctx, http.MethodPost, "/order/by_table", authToken, reqBody)
	if err != nil {
		return nil, err
	}

	var resp TableOrdersResponse
	if err := json.Unmarshal(respBytes, &resp); err != nil {
		return nil, err
	}

	return &resp, nil
}

// -----------------------------
// Create Order
// -----------------------------

type CreateOrderRequest struct {
    OrganizationID string      `json:"organizationId"`
    TerminalID     string      `json:"terminalGroupId"`
    Order          OrderPayload `json:"order"`
}

type OrderPayload struct {
    ID             string           `json:"id,omitempty"` // External ID
    OrderTypeID    string           `json:"orderTypeId"`
    TableIDs       []string         `json:"tableIds"`
    Customer       *Customer        `json:"customer,omitempty"`
    Phone          string           `json:"phone,omitempty"`
    Items          []OrderItem      `json:"items"`
    Comment        string           `json:"comment,omitempty"`
}

type Customer struct {
    Name   string `json:"name"`
    Phone  string `json:"phone,omitempty"` // Sometimes needed inside customer too
    Type   string `json:"type"` // regular
}

type OrderItem struct {
    ProductID string   `json:"productId"`
    Amount    float64  `json:"amount"`
    Price     *float64 `json:"price,omitempty"` // Optional override
    Type      string   `json:"type"` // Product
    Modifiers []OrderModifier `json:"modifiers,omitempty"`
}

type OrderModifier struct {
    ProductID      string  `json:"productId,omitempty"`
    ProductGroupID string  `json:"productGroupId,omitempty"`
    Amount         float64 `json:"amount"`
}

type CreateOrderResponse struct {
    OrderInfo struct {
        ID string `json:"id"`
    } `json:"orderInfo"`
}

func (c *Client) CreateOrder(ctx context.Context, authToken string, req CreateOrderRequest) (*CreateOrderResponse, error) {
    respBytes, err := c.doRequest(ctx, http.MethodPost, "/order/create", authToken, req)
    if err != nil {
        return nil, err
    }

    var resp CreateOrderResponse
    if err := json.Unmarshal(respBytes, &resp); err != nil {
        return nil, err
    }

    return &resp, nil
}
