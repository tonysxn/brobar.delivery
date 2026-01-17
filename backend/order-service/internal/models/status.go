package models

import (
	"database/sql/driver"
	"errors"
	"fmt"
)

type Status string

const (
	StatusPending   Status = "pending"
	StatusPaid      Status = "paid"
	StatusShipping  Status = "shipping"
	StatusCompleted Status = "completed"
	StatusCancelled Status = "cancelled"
)

func (r *Status) Scan(value interface{}) error {
	switch v := value.(type) {
	case string:
		*r = Status(v)
		return nil
	case []byte:
		*r = Status(string(v))
		return nil
	default:
		return errors.New("status should be a string or []byte")
	}
}

func (r Status) Value() (driver.Value, error) {
	switch r {
	case StatusPending, StatusPaid, StatusShipping, StatusCompleted, StatusCancelled:
		return string(r), nil
	default:
		return nil, fmt.Errorf("invalid status: %s", r)
	}
}
