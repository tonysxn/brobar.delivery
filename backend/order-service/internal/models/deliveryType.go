package models

import (
	"database/sql/driver"
	"errors"
	"fmt"
)

type DeliveryType string

const (
	DeliveryTypeDelivery DeliveryType = "delivery"
	DeliveryTypePickup   DeliveryType = "pickup"
	DeliveryTypeDine     DeliveryType = "dine"
)

func (r *DeliveryType) Scan(value interface{}) error {
	switch v := value.(type) {
	case string:
		*r = DeliveryType(v)
		return nil
	case []byte:
		*r = DeliveryType(string(v))
		return nil
	default:
		return errors.New("deliveryType should be a string or []byte")
	}
}

func (r DeliveryType) Value() (driver.Value, error) {
	switch r {
	case DeliveryTypeDelivery, DeliveryTypePickup, DeliveryTypeDine:
		return string(r), nil
	default:
		return nil, fmt.Errorf("invalid status: %s", r)
	}
}
