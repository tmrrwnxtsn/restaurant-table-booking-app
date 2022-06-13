package model

import "net/http"

// Restaurant представляет ресторан.
type Restaurant struct {
	ID   uint64 `json:"id"`
	Name string `json:"name"`
	// AverageWaitingTime представляет среднее время ожидания заказа в минутах.
	AverageWaitingTime int `json:"average_waiting_time"`
	// AverageCheck представляет средний чек на блюдо в ресторане.
	AverageCheck float64 `json:"average_check"`
	// AvailableSeatsNumber представляет актуальное количество свободных мест.
	AvailableSeatsNumber int `json:"available_seats_number,omitempty"`
}

// UpdateRestaurantData содержит информацию о ресторане и используется для обновления записи о нём в БД.
type UpdateRestaurantData struct {
	Name               *string  `json:"name"`
	AverageWaitingTime *int     `json:"average_waiting_time,string"`
	AverageCheck       *float64 `json:"average_check,string"`
}

// Bind осуществляет пост-обработку запроса UpdateRestaurantData.
func (d *UpdateRestaurantData) Bind(_ *http.Request) error {
	if d.Name == nil && d.AverageWaitingTime == nil && d.AverageCheck == nil {
		return ErrUpdateRestaurantRequest
	}
	return nil
}
