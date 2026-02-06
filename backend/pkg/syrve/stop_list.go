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

type StopListRequest struct {
	OrganizationIDs []string `json:"organizationIds"`
}

type StopListResponse struct {
	CorrelationID          string                 `json:"correlationId"`
	TerminalGroupStopLists []OrganizationStopList `json:"terminalGroupStopLists"`
}

type OrganizationStopList struct {
	OrganizationID string                  `json:"organizationId"`
	Items          []TerminalGroupStopList `json:"items"`
}

type TerminalGroupStopList struct {
	TerminalGroupID string         `json:"terminalGroupId"`
	Items           []StopListItem `json:"items"`
}

type StopListItem struct {
	Balance   float64 `json:"balance"`
	ProductID string  `json:"productId"`
}

func (c *Client) GetStopLists(ctx context.Context, authToken, organizationID string) (*StopListResponse, error) {
	if authToken == "" {
		return nil, errors.New("authorization token is required")
	}
	if organizationID == "" {
		return nil, errors.New("organizationId is required")
	}

	// url := fmt.Sprintf("%s/stop_lists?organizationId=%s", c.BaseURL, organizationID)
	url := fmt.Sprintf("%s/stop_lists", c.BaseURL)
	
	reqBody := StopListRequest{
		OrganizationIDs: []string{organizationID},
	}
	bodyBytes, err := json.Marshal(reqBody)
	if err != nil {
		return nil, err
	}
	
	reqHTTP, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewReader(bodyBytes))
	if err != nil {
		return nil, err
	}

	reqHTTP.Header.Set("Content-Type", "application/json")
	reqHTTP.Header.Set("Authorization", "Bearer "+authToken)

	resp, err := c.HttpClient.Do(reqHTTP)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(resp.Body)
		log.Printf("GetStopLists failed: %s | URL: %s | Body: %s", resp.Status, url, string(bodyBytes))
		return nil, fmt.Errorf("unexpected status code: %s", resp.Status)
	}

	// Parse response
	respBodyBytes, _ := io.ReadAll(resp.Body)
	log.Printf("GetStopLists success: %s | Body length: %d", resp.Status, len(respBodyBytes))
	log.Printf("GetStopLists body: %s", string(respBodyBytes))

	var stopListResp StopListResponse
	if err := json.Unmarshal(respBodyBytes, &stopListResp); err != nil {
		log.Printf("Failed to unmarshal stop list response: %v | Body sample: %s", err, string(respBodyBytes))
		return nil, err
	}

	return &stopListResp, nil
}
