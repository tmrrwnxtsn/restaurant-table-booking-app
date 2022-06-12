package handler

import "errors"

// ошибки, возникающие при работе с маршрутами манипуляции с ресторанами.
var (
	// ErrRestaurantMissingFields возникает, когда в запросе на создание ресторана пропущены обязательные поля.
	ErrRestaurantMissingFields = errors.New("missing required restaurant fields")
)
