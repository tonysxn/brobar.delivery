package errors

import "errors"

var (
	OrderNotFound    = errors.New("order not found")
	OrderInvalidData = errors.New("invalid order data")
)
