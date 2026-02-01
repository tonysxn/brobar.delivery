package payment

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/tonysanin/brobar/pkg/monobank"
)

type Client struct {
	baseURL    string
	httpClient *http.Client
}

func NewClient(baseURL string) *Client {
	return &Client{
		baseURL:    baseURL,
		httpClient: &http.Client{},
	}
}

type InitPaymentInput struct {
	Amount      int                    `json:"amount"`
	OrderID     string                 `json:"order_id"`
	RedirectURL string                 `json:"redirect_url"`
	WebhookURL  string                 `json:"webhook_url"`
	Basket      []monobank.BasketOrder `json:"basket"`
}

type InitPaymentOutput struct {
	PaymentURL string `json:"payment_url"`
	InvoiceID  string `json:"invoice_id"`
}

func (c *Client) InitPayment(input InitPaymentInput) (*InitPaymentOutput, error) {
	body, err := json.Marshal(input)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal input: %w", err)
	}

	req, err := http.NewRequest("POST", c.baseURL+"/api/v1/payments/init", bytes.NewBuffer(body))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		var errResp map[string]string
		_ = json.NewDecoder(resp.Body).Decode(&errResp)
		return nil, fmt.Errorf("payment service error: %s", errResp["error"])
	}

	var output InitPaymentOutput
	if err := json.NewDecoder(resp.Body).Decode(&output); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return &output, nil
}
