package requests

import (
	"time"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/google/uuid"
	"github.com/tonysanin/brobar/order-service/internal/models"
	"github.com/tonysanin/brobar/pkg/validator"
)

// CreateOrderRequest - minimal data from frontend
type CreateOrderRequest struct {
	// Contact
	Name  string `json:"name"`
	Phone string `json:"phone"`
	Email string `json:"email,omitempty"`

	// Delivery
	DeliveryTypeID string `json:"delivery_type_id"` // "delivery" | "pickup"
	Address        string `json:"address,omitempty"`
	Zone           string `json:"zone,omitempty"`
	Entrance       string `json:"entrance,omitempty"`
	DeliveryDoor   bool   `json:"delivery_door,omitempty"`
	Coords         string `json:"coords,omitempty"`

	// Time
	Time string `json:"time"` // "ASAP" or "2026-01-18 14:30"

	// Payment & Other
	PaymentMethod string `json:"payment_method"` // "online" | "cash"
	Cutlery       int    `json:"cutlery,omitempty"`
	PromoCode     string `json:"promo_code,omitempty"`
	Wishes        string `json:"wishes,omitempty"`

	// Items (minimal - only IDs and quantities)
	Items []OrderItemRequest `json:"items"`

	// Client-calculated total for validation
	ClientTotal float64 `json:"client_total"`
}

// OrderItemRequest - minimal item data from frontend
type OrderItemRequest struct {
	ProductID          uuid.UUID  `json:"product_id"`
	ProductVariationID *uuid.UUID `json:"product_variation_id,omitempty"`
	Quantity           int        `json:"quantity"`
}

func (r CreateOrderRequest) Validate() error {
	return validation.ValidateStruct(&r,
		validation.Field(&r.Name, validation.Required, validation.Length(2, 100)),
		validation.Field(&r.Phone, validation.Required, validator.IsPhone, validation.Length(6, 32)),
		validation.Field(&r.DeliveryTypeID, validation.Required, validation.In(
			string(models.DeliveryTypeDelivery),
			string(models.DeliveryTypePickup),
			string(models.DeliveryTypeDine),
		)),
		validation.Field(&r.PaymentMethod, validation.Required, validation.In("online", "cash", "bank")),
		validation.Field(&r.Time, validation.Required),
		validation.Field(&r.Items, validation.Required, validation.Length(1, 100)),
		validation.Field(&r.ClientTotal, validation.Required, validator.IsNonNegative),
	)
}

func (r OrderItemRequest) Validate() error {
	return validation.ValidateStruct(&r,
		validation.Field(&r.ProductID, validation.Required, validator.IsUUID),
		validation.Field(&r.Quantity, validation.Required, validation.Min(1)),
	)
}

// UpdateOrderRequest - admin update (full data)
type UpdateOrderRequest struct {
	ID             uuid.UUID                `json:"-"`
	Name           string                   `json:"name"`
	Address        string                   `json:"address"`
	Phone          string                   `json:"phone"`
	StatusID       string                   `json:"status_id"`
	DeliveryTypeID string                   `json:"delivery_type_id"`
	PaymentMethod  string                   `json:"payment_method"`
	Time           time.Time                `json:"time"`
	Items          []UpdateOrderItemRequest `json:"items"`
}

type UpdateOrderItemRequest struct {
	ProductID                  uuid.UUID  `json:"product_id"`
	Quantity                   int        `json:"quantity"`
	Price                      float64    `json:"price"`
	Weight                     float64    `json:"weight"`
	Name                       string     `json:"name"`
	ExternalProductID          string     `json:"external_product_id"`
	ProductVariationGroupID    *uuid.UUID `json:"product_variation_group_id"`
	ProductVariationGroupName  *string    `json:"product_variation_group_name,omitempty"`
	ProductVariationID         *uuid.UUID `json:"product_variation_id"`
	ProductVariationExternalID *string    `json:"product_variation_external_id,omitempty"`
	ProductVariationName       *string    `json:"product_variation_name,omitempty"`
}

func (r UpdateOrderRequest) Validate() error {
	return validation.ValidateStruct(&r,
		validation.Field(&r.Name, validation.Required, validation.Length(2, 100)),
		validation.Field(&r.Address, validation.Required, validation.Length(5, 256)),
		validation.Field(&r.Phone, validation.Required, validator.IsPhone, validation.Length(6, 32)),
		validation.Field(&r.StatusID, validation.In(
			string(models.StatusPending),
			string(models.StatusPaid),
			string(models.StatusShipping),
			string(models.StatusCompleted),
			string(models.StatusCancelled))),
		validation.Field(&r.DeliveryTypeID, validation.In(
			string(models.DeliveryTypeDelivery),
			string(models.DeliveryTypePickup),
			string(models.DeliveryTypeDine),
		)),
		validation.Field(&r.PaymentMethod, validation.In("online", "cash", "bank")),
		validation.Field(&r.Time, validation.Required),
		validation.Field(&r.Items, validation.Required, validation.Length(1, 100)),
	)
}

func (r UpdateOrderItemRequest) Validate() error {
	return validation.ValidateStruct(&r,
		validation.Field(&r.ProductID, validation.Required, validator.IsUUID),
		validation.Field(&r.Quantity, validation.Required, validation.Min(1)),
		validation.Field(&r.Price, validation.Required, validator.IsNonNegative),
		validation.Field(&r.Weight, validation.Required, validator.IsNonNegative),
		validation.Field(&r.Name, validation.Required, validation.Length(1, 256)),
		validation.Field(&r.ExternalProductID, validation.Required, validation.Length(0, 100)),
		validation.Field(&r.ProductVariationGroupName, validation.Length(1, 255)),
		validation.Field(&r.ProductVariationExternalID, validation.Length(0, 100)),
		validation.Field(&r.ProductVariationName, validation.Length(1, 255)),
	)
}
