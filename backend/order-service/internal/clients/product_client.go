package clients

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/tonysanin/brobar/pkg/helpers"
)

type ProductClient struct {
	baseURL    string
	httpClient *http.Client
}

func NewProductClient() *ProductClient {
	return &ProductClient{
		baseURL: helpers.GetEnv("PRODUCT_SERVICE_URL", "http://product-service-dev:3000"),
		httpClient: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

// Product response
type ProductResponse struct {
	Success bool    `json:"success"`
	Data    Product `json:"data"`
}

type Product struct {
	ID         uuid.UUID `json:"id"`
	ExternalID string    `json:"external_id"`
	Name       string    `json:"name"`
	Price      float64   `json:"price"`
	Weight     float64   `json:"weight"`
	Stock      *float64  `json:"stock"`
}

// Variation response
type VariationResponse struct {
	Success bool      `json:"success"`
	Data    Variation `json:"data"`
}

type Variation struct {
	ID         uuid.UUID `json:"id"`
	GroupID    uuid.UUID `json:"group_id"`
	ExternalID string    `json:"external_id"`
	Name       string    `json:"name"`
}

// VariationGroup response
type VariationGroupResponse struct {
	Success bool           `json:"success"`
	Data    VariationGroup `json:"data"`
}

type VariationGroup struct {
	ID        uuid.UUID `json:"id"`
	ProductID uuid.UUID `json:"product_id"`
	Name      string    `json:"name"`
}

func (c *ProductClient) GetProduct(productID uuid.UUID) (*Product, error) {
	resp, err := c.httpClient.Get(fmt.Sprintf("%s/products/%s", c.baseURL, productID.String()))
	if err != nil {
		return nil, fmt.Errorf("failed to fetch product: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("product not found: %s", productID.String())
	}

	var productResp ProductResponse
	if err := json.NewDecoder(resp.Body).Decode(&productResp); err != nil {
		return nil, fmt.Errorf("failed to decode product response: %w", err)
	}

	if !productResp.Success {
		return nil, fmt.Errorf("product not found: %s", productID.String())
	}

	return &productResp.Data, nil
}

func (c *ProductClient) GetVariation(variationID uuid.UUID) (*Variation, error) {
	resp, err := c.httpClient.Get(fmt.Sprintf("%s/variations/%s", c.baseURL, variationID.String()))
	if err != nil {
		return nil, fmt.Errorf("failed to fetch variation: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("variation not found: %s", variationID.String())
	}

	var variationResp VariationResponse
	if err := json.NewDecoder(resp.Body).Decode(&variationResp); err != nil {
		return nil, fmt.Errorf("failed to decode variation response: %w", err)
	}

	if !variationResp.Success {
		return nil, fmt.Errorf("variation not found: %s", variationID.String())
	}

	return &variationResp.Data, nil
}

func (c *ProductClient) GetVariationGroup(groupID uuid.UUID) (*VariationGroup, error) {
	resp, err := c.httpClient.Get(fmt.Sprintf("%s/variation-groups/%s", c.baseURL, groupID.String()))
	if err != nil {
		return nil, fmt.Errorf("failed to fetch variation group: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("variation group not found: %s", groupID.String())
	}

	var groupResp VariationGroupResponse
	if err := json.NewDecoder(resp.Body).Decode(&groupResp); err != nil {
		return nil, fmt.Errorf("failed to decode variation group response: %w", err)
	}

	if !groupResp.Success {
		return nil, fmt.Errorf("variation group not found: %s", groupID.String())
	}

	return &groupResp.Data, nil
}
