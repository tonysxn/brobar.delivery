package models

import "github.com/google/uuid"

type OrderItem struct {
	ID                uuid.UUID `json:"id" db:"id"`
	OrderID           uuid.UUID `json:"order_id" db:"order_id"`
	ProductID         uuid.UUID `json:"product_id" db:"product_id"`
	ExternalProductID string    `json:"external_product_id" db:"external_product_id" validate:"max=100"`
	Name              string    `json:"name" db:"name"`
	Price             float64   `json:"price" db:"price"`
	Quantity          int       `json:"quantity" db:"quantity"`
	TotalPrice        float64   `json:"total_price" db:"total_price"`
	Weight            float64   `json:"weight" db:"weight"`
	TotalWeight       float64   `json:"total_weight" db:"total_weight"`

	ProductVariationGroupID    *uuid.UUID `json:"product_variation_group_id" db:"product_variation_group_id"`
	ProductVariationGroupName  *string    `json:"product_variation_group_name,omitempty" db:"product_variation_group_name" validate:"omitempty,min=1,max=255"`
	ProductVariationID         *uuid.UUID `json:"product_variation_id" db:"product_variation_id"`
	ProductVariationExternalID *string    `json:"product_variation_external_id,omitempty" db:"product_variation_external_id" validate:"max=100"`
	ProductVariationName       *string    `json:"product_variation_name,omitempty" db:"product_variation_name" validate:"omitempty,min=1,max=255"`
}
