package syrve

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"time"
)

const defaultTimeoutSec = 15

type Client struct {
	BaseURL        string
	ApiLogin       string
	HttpClient     *http.Client
	Timeout        time.Duration
	OrganizationID string
}

type AccessTokenRequest struct {
	APILogin string `json:"apiLogin"`
}

type AccessTokenResponse struct {
	CorrelationID string `json:"correlationId"`
	Token         string `json:"token"`
}

func NewClient(apiLogin string, organizationID string) *Client {
	return &Client{
		BaseURL:        "https://api-eu.syrve.live/api/1",
		ApiLogin:       apiLogin,
		OrganizationID: organizationID,
		HttpClient: &http.Client{
			Timeout: time.Second * defaultTimeoutSec,
		},
		Timeout: time.Second * defaultTimeoutSec,
	}
}

func (c *Client) WithTimeout(timeoutSec int) *Client {
	if timeoutSec > 0 {
		c.Timeout = time.Duration(timeoutSec) * time.Second
		c.HttpClient.Timeout = c.Timeout
	}
	return c
}

func (c *Client) GetAccessToken(ctx context.Context) (*AccessTokenResponse, error) {
	reqBody := AccessTokenRequest{APILogin: c.ApiLogin}
	bodyBytes, err := json.Marshal(reqBody)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, c.BaseURL+"/access_token", bytes.NewReader(bodyBytes))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := c.HttpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, errors.New("unexpected status code: " + resp.Status)
	}

	var tokenResp AccessTokenResponse
	err = json.NewDecoder(resp.Body).Decode(&tokenResp)
	if err != nil {
		return nil, err
	}

	if tokenResp.Token == "" || tokenResp.CorrelationID == "" {
		return nil, errors.New("empty token or correlationId in response")
	}

	return &tokenResp, nil
}
