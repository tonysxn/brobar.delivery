package syrve

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"github.com/tonysanin/brobar/syrve-service/internal/services/syrve/types"
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

func (c *Client) GetProducts(ctx context.Context, authToken, organizationID string) ([]syrve.MenuItemDTO, error) {
	resp, err := c.GetNomenclature(ctx, authToken, syrve.NomenclatureRequest{OrganizationID: organizationID})
	if err != nil {
		return nil, err
	}

	productMap := make(map[string]syrve.MenuItem)
	for _, p := range resp.Products {
		if *p.Type == "Dish" {
			productMap[p.ID] = p
		}
	}

	modifiersMap := make(map[string]syrve.MenuItem)
	for _, p := range resp.Products {
		if *p.Type == "Modifier" {
			modifiersMap[p.ID] = p
		}
	}

	groupMap := make(map[string]syrve.Group)
	for _, group := range resp.Groups {
		groupMap[group.ID] = group
	}

	var result []syrve.MenuItemDTO
	for _, p := range productMap {
		full := syrve.MenuItemDTO{
			ID:             p.ID,
			Name:           p.Name,
			Modifiers:      []syrve.ModifierDTO{},
			GroupModifiers: []syrve.ModifierGroupDTO{},
		}

		for _, modifier := range p.Modifiers {
			modifierFull := syrve.ModifierDTO{ID: modifier.ID, Name: modifiersMap[modifier.ID].Name, DefaultAmount: modifier.DefaultAmount, Required: modifier.Required}

			full.Modifiers = append(full.Modifiers, modifierFull)
		}

		for _, groupModifier := range p.GroupModifiers {
			groupModifierFull := syrve.ModifierGroupDTO{
				ID:             groupModifier.ID,
				Name:           groupMap[groupModifier.ID].Name,
				Required:       groupModifier.Required,
				DefaultAmount:  groupModifier.DefaultAmount,
				ChildModifiers: []syrve.ModifierDTO{},
			}

			for _, cm := range groupModifier.ChildModifiers {
				m := syrve.ModifierDTO{
					ID:            cm.ID,
					Name:          modifiersMap[cm.ID].Name,
					DefaultAmount: cm.DefaultAmount,
					Required:      cm.Required,
				}
				groupModifierFull.ChildModifiers = append(groupModifierFull.ChildModifiers, m)
			}

			full.GroupModifiers = append(full.GroupModifiers, groupModifierFull)
		}

		result = append(result, full)
	}

	return result, nil
}
