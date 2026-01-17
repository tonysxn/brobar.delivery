package syrve

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"github.com/tonysanin/brobar/syrve-service/internal/services/syrve/types"
	"net/http"
)

func (c *Client) GetNomenclature(ctx context.Context, authToken string, req syrve.NomenclatureRequest) (*syrve.NomenclatureResponse, error) {
	if authToken == "" {
		return nil, errors.New("authorization token is required")
	}
	if req.OrganizationID == "" {
		return nil, errors.New("organizationId is required")
	}

	bodyBytes, err := json.Marshal(req)
	if err != nil {
		return nil, err
	}

	reqHTTP, err := http.NewRequestWithContext(ctx, http.MethodPost, c.BaseURL+"/nomenclature", bytes.NewReader(bodyBytes))
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
		return nil, errors.New("unexpected status code: " + resp.Status)
	}

	var nomenResp syrve.NomenclatureResponse
	err = json.NewDecoder(resp.Body).Decode(&nomenResp)
	if err != nil {
		return nil, err
	}

	return &nomenResp, nil
}
