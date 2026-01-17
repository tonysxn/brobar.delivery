package syrve

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
)

type OrganizationsRequest struct {
	OrganizationIDs      []string `json:"organizationIds,omitempty"`
	ReturnAdditionalInfo bool     `json:"returnAdditionalInfo,omitempty"`
	IncludeDisabled      bool     `json:"includeDisabled,omitempty"`
	ReturnExternalData   []string `json:"returnExternalData,omitempty"`
}

type Organization struct {
	ID                string  `json:"id"`
	Name              string  `json:"name"`
	RMSVersion        *string `json:"rmsVersion,omitempty"`
	Country           *string `json:"country,omitempty"`
	RestaurantAddress *string `json:"restaurantAddress,omitempty"`
}

type OrganizationsResponse struct {
	Organizations []Organization `json:"organizations"`
}

func (c *Client) GetOrganizations(ctx context.Context, authToken string, reqBody OrganizationsRequest) (*OrganizationsResponse, error) {
	if authToken == "" {
		return nil, errors.New("authorization token is required")
	}

	bodyBytes, err := json.Marshal(reqBody)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, c.BaseURL+"/organizations", bytes.NewReader(bodyBytes))
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
		return nil, errors.New("unexpected status code: " + resp.Status)
	}

	var orgsResp OrganizationsResponse
	err = json.NewDecoder(resp.Body).Decode(&orgsResp)
	if err != nil {
		return nil, err
	}

	return &orgsResp, nil
}
