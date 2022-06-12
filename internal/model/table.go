package model

// Table представляет столик в ресторане.
type Table struct {
	ID           uint64 `json:"id"`
	RestaurantID uint64 `json:"restaurant_id"`
	// SeatsNumber представляет вместимость столика.
	SeatsNumber int `json:"seats_number"`
}
