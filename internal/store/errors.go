package store

import "errors"

var (
	// ErrRestaurantNotFound возникает, когда по введённому ID в БД не находится искомого ресторана.
	ErrRestaurantNotFound = errors.New("restaurant not found")
)
