package model

import "errors"

var (
	// ErrUpdateRestaurantRequest возникает при попытке обновить данные о ресторане без передачи самих данных.
	ErrUpdateRestaurantRequest = errors.New("update restaurant request has no values")
)
