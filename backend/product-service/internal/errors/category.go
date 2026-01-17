package errors

import "errors"

var (
	CategoryNotFound    = errors.New("category not found")
	CategoryInvalidData = errors.New("invalid category data")
)
