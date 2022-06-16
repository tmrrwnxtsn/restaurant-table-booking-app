package model

import "errors"

var (
	// ErrUpdateRestaurantData возникает при попытке обновить данные о ресторане без передачи самих данных.
	ErrUpdateRestaurantData = errors.New("update restaurant data has no values")
	// ErrUpdateTableData возникает при попытке обновить данные о столике в ресторане без передачи самих данных.
	ErrUpdateTableData = errors.New("update table data has no values")
)
