package product

import "errors"

var (
	ErrNotFound   = errors.New("product not found")
	ErrSKUExists  = errors.New("product with this SKU already exists")
	ErrInvalidInput = errors.New("invalid input")
)
