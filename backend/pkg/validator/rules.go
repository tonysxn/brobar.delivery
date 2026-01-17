package validator

import (
	"errors"
	"regexp"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/google/uuid"
)

// IsUUID checks if the value is a valid non-nil UUID.
var IsUUID = validation.By(func(value interface{}) error {
	id, ok := value.(uuid.UUID)
	if !ok {
		// Try to handle pointer to UUID if necessary, but typically we validate the value itself
		// if passed as value. If it is *uuid.UUID, casting to uuid.UUID fails.
		// For simplicity, let's assume we are validating uuid.UUID values.
		// If we encounter pointers, we might need a type switch.
		if ptr, ok := value.(*uuid.UUID); ok && ptr != nil {
			id = *ptr
		} else {
			// If it's a string, we could parse it, but usually we bind to UUID types.
			return nil // "Strictly speaking, nil interface is valid for Required to catch"
		}
	}

	// If it's the zero value (Nil UUID)
	if id == uuid.Nil {
		return errors.New("must be a valid UUID")
	}
	return nil
})

// IsNonNegative checks if a number is >= 0
var IsNonNegative = validation.By(func(value interface{}) error {
	switch v := value.(type) {
	case int:
		if v < 0 {
			return errors.New("must be >= 0")
		}
	case float64:
		if v < 0 {
			return errors.New("must be >= 0")
		}
	case float32:
		if v < 0 {
			return errors.New("must be >= 0")
		}
	case int64:
		if v < 0 {
			return errors.New("must be >= 0")
		}
	default:
		// nil or other types are considered valid by this rule; use Required for nil checks.
		return nil
	}
	return nil
})

var phoneRegex = regexp.MustCompile(`^\+?[0-9]{7,15}$`)

// IsPhone checks if the string is a valid phone number
var IsPhone = validation.Match(phoneRegex).Error("invalid phone number")

// Helper wrapper to validate a struct that implements the Validatable interface
type Validatable interface {
	Validate() error
}

func Validate(v Validatable) error {
	return v.Validate()
}
