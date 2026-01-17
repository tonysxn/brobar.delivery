package errors

import "errors"

var (
	ProductVariationGroupNotFound    = errors.New("product variation group not found")
	ProductVariationGroupInvalidData = errors.New("invalid product variation group data")
)
