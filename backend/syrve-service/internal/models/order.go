package models

import (
	"time"

	"github.com/google/uuid"
)

// OrderEvent mirrors the JSON structure sent by order-service
type OrderEvent struct {
	ID                uuid.UUID    `json:"id"`
	UserID            *uuid.UUID   `json:"user_id,omitempty"`
	StatusID          string       `json:"status_id"`
	TotalPrice        float64      `json:"total_price"`
	CreatedAt         time.Time    `json:"created_at"`
	
	Address           string       `json:"address"`
	Entrance          string       `json:"entrance,omitempty"`
	Floor             string       `json:"floor,omitempty"`
	Flat              string       `json:"flat,omitempty"`
	AddressWishes     string       `json:"address_wishes,omitempty"`
	Name              string       `json:"name"`
	Phone             string       `json:"phone,omitempty"`
	Time              time.Time    `json:"time"`
	Email             string       `json:"email,omitempty"`
	Wishes            string       `json:"wishes,omitempty"`
	Promo             string       `json:"promo,omitempty"`
	Coords            string       `json:"coords,omitempty"`
	Cutlery           int          `json:"cutlery,omitempty"`
	DeliveryCost      float64      `json:"delivery_cost"`
	DeliveryDoor      bool         `json:"delivery_door"`
	DeliveryDoorPrice float64      `json:"delivery_door_price"`
	DeliveryTypeID    string       `json:"delivery_type_id"`
	PaymentMethod     string       `json:"payment_method,omitempty"`
	Zone              *string      `json:"zone,omitempty"`
	InvoiceID         *string      `json:"invoice_id,omitempty"`
	
	Items []OrderItem `json:"items"`
}

type OrderItem struct {
	ID                         uuid.UUID  `json:"id"`
	ProductID                  uuid.UUID  `json:"product_id"`
	ExternalProductID          string     `json:"external_product_id"`
	Name                       string     `json:"name"`
	Price                      float64    `json:"price"`
	Quantity                   int        `json:"quantity"`
	TotalPrice                 float64    `json:"total_price"`
	Weight                     float64    `json:"weight"`
	TotalWeight                float64    `json:"total_weight"`

	ProductVariationGroupID    *uuid.UUID `json:"product_variation_group_id,omitempty"`
	ProductVariationGroupName  *string    `json:"product_variation_group_name,omitempty"`
	ProductVariationID         *uuid.UUID `json:"product_variation_id,omitempty"`
	ProductVariationExternalID *string    `json:"product_variation_external_id,omitempty"`
	ProductVariationName       *string    `json:"product_variation_name,omitempty"`
}
