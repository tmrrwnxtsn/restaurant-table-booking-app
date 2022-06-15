package model

import "errors"

// ErrUpdateRestaurantData возникает при попытке обновить данные о ресторане без передачи самих данных.
var ErrUpdateRestaurantData = errors.New("update restaurant request has no values")
