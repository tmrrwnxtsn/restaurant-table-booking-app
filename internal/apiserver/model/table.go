package model

import "net/http"

// Table представляет столик в ресторане.
type Table struct {
	ID           uint64 `json:"id"`
	RestaurantID uint64 `json:"restaurant_id"`
	// SeatsNumber представляет вместимость столика.
	SeatsNumber int `json:"seats_number"`
}

// UpdateTableData содержит информацию о столике в ресторане и используется для обновления записи о нём в БД.
type UpdateTableData struct {
	SeatsNumber *int `json:"seats_number"`
}

// Bind осуществляет пост-обработку запроса UpdateTableData.
func (d *UpdateTableData) Bind(_ *http.Request) error {
	if d.SeatsNumber == nil {
		return ErrUpdateTableData
	}
	return nil
}
