package errors

import "errors"

var (
	ProductNotFound    = errors.New("product not found")
	ProductInvalidData = errors.New("invalid product data")
)
