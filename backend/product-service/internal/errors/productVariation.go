package errors

import "errors"

var (
	ProductVariationNotFound    = errors.New("product variation not found")
	ProductVariationInvalidData = errors.New("invalid product variation data")
)
