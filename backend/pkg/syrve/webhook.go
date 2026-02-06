package syrve

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
)

type WebhooksFilter struct {
	FilterEventTypes []string `json:"filterEventTypes"`
}

type WebhookSettings struct {
	WebHooksURI    string          `json:"webHooksUri"`
	AuthToken      string          `json:"authToken,omitempty"`
	OrganizationID string          `json:"organizationId"`
	WebhooksFilter *WebhooksFilter `json:"webhooksFilter,omitempty"`
}

func (c *Client) UpdateWebhook(ctx context.Context, authToken, organizationID, webhookURL string) error {
	if authToken == "" {
		return errors.New("authorization token is required")
	}

	// Endpoint: /api/1/webhooks/update_settings
	// Body: { "organizationId": "...", "webHooksUri": "...", "authToken": "..." }
	
	body := WebhookSettings{
		OrganizationID: organizationID,
		WebHooksURI:    webhookURL,
		AuthToken:      "brobar_secret_token_123",
		WebhooksFilter: &WebhooksFilter{
			FilterEventTypes: []string{
				"DeliveryOrderUpdate",
				"DeliveryOrderError",
				"TableOrderUpdate",
				"TableOrderError",
				"ReserveUpdate",
				"ReserveError",
				"StopListUpdate",
			},
		},
	}

	bodyBytes, err := json.Marshal(body)
	if err != nil {
		return err
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, c.BaseURL+"/webhooks/update_settings", bytes.NewReader(bodyBytes))
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+authToken)

	resp, err := c.HttpClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	respBytes, _ := io.ReadAll(resp.Body)
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("unexpected status code: %s, body: %s", resp.Status, string(respBytes))
	}
	
	log.Printf("Syrve UpdateWebhook success body: %s", string(respBytes))

	return nil
}

func (c *Client) GetWebhookSettings(ctx context.Context, authToken, organizationID string) (*WebhookSettings, error) {
	if authToken == "" {
		return nil, errors.New("authorization token is required")
	}

	body := map[string]string{"organizationId": organizationID}
	bodyBytes, _ := json.Marshal(body)

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, c.BaseURL+"/webhooks/settings", bytes.NewReader(bodyBytes))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+authToken)

	resp, err := c.HttpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		respBytes, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("unexpected status code: %s, body: %s", resp.Status, string(respBytes))
	}

	var settings WebhookSettings
	if err := json.NewDecoder(resp.Body).Decode(&settings); err != nil {
		return nil, err
	}

	return &settings, nil
}

func (c *Client) RemoveWebhook(ctx context.Context, authToken, organizationID string) error {
    // To remove, we set it to empty string or null.
    // Based on legacy logic, it might just valid if we unset it?
	// or update to empty string?
    return c.UpdateWebhook(ctx, authToken, organizationID, "")
}
