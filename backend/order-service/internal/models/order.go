package models

import (
	"time"

	"github.com/google/uuid"
)

type Order struct {
	ID         uuid.UUID  `json:"id" db:"id"`
	UserID     *uuid.UUID `json:"user_id,omitempty" db:"user_id"`
	StatusID   Status     `json:"status_id" db:"status_id"`
	TotalPrice float64    `json:"total_price" db:"total_price"`
	CreatedAt  time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt  time.Time  `json:"updated_at" db:"updated_at"`

	Address           string       `json:"address" db:"address"`
	Entrance          string       `json:"entrance,omitempty" db:"entrance"`
	Floor             string       `json:"floor,omitempty" db:"floor"`
	Flat              string       `json:"flat,omitempty" db:"flat"`
	AddressWishes     string       `json:"address_wishes,omitempty" db:"address_wishes"`
	Name              string       `json:"name" db:"name"`
	Phone             string       `json:"phone,omitempty" db:"phone"`
	Time              time.Time    `json:"time" db:"time"`
	Email             string       `json:"email,omitempty" db:"email"`
	Wishes            string       `json:"wishes,omitempty" db:"wishes"`
	Promo             string       `json:"promo,omitempty" db:"promo"`
	Coords            string       `json:"coords,omitempty" db:"coords"`
	Cutlery           int          `json:"cutlery,omitempty" db:"cutlery"`
	DeliveryCost      float64      `json:"delivery_cost" db:"delivery_cost"`
	DeliveryDoor      bool         `json:"delivery_door" db:"delivery_door"`
	DeliveryDoorPrice float64      `json:"delivery_door_price" db:"delivery_door_price"`
	DeliveryTypeID    DeliveryType `json:"delivery_type_id" db:"delivery_type_id"`
	PaymentMethod     string       `json:"payment_method,omitempty" db:"payment_method"`
	Zone              *string      `json:"zone,omitempty" db:"zone"`
	InvoiceID         *string      `json:"invoice_id,omitempty" db:"invoice_id"`
	SyrveNotified     bool         `json:"syrve_notified" db:"syrve_notified"`
	PaymentURL        string       `json:"payment_url,omitempty" db:"-"`

	Items []OrderItem `json:"items" db:"-"`
}
