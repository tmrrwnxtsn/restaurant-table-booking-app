package model

import "net/http"

// Restaurant представляет ресторан.
type Restaurant struct {
	ID   uint64 `json:"id" example:"3"`
	Name string `json:"name" example:"Каравелла"`
	// AverageWaitingTime представляет среднее время ожидания заказа в минутах.
	AverageWaitingTime int `json:"average_waiting_time" example:"60"`
	// AverageCheck представляет средний чек на блюдо в ресторане.
	AverageCheck float64 `json:"average_check" example:"2500.00"`
	// AvailableSeatsNumber представляет актуальное количество свободных мест.
	AvailableSeatsNumber int `json:"available_seats_number,omitempty" example:"24"`
}

// UpdateRestaurantData содержит информацию о ресторане и используется для обновления записи о нём в БД.
type UpdateRestaurantData struct {
	Name               *string  `json:"name" example:"Каравелла"`
	AverageWaitingTime *int     `json:"average_waiting_time,string" example:"60"`
	AverageCheck       *float64 `json:"average_check,string" example:"2500.00"`
}

// Bind осуществляет пост-обработку запроса UpdateRestaurantData.
func (d *UpdateRestaurantData) Bind(_ *http.Request) error {
	if d.Name == nil && d.AverageWaitingTime == nil && d.AverageCheck == nil {
		return ErrUpdateRestaurantData
	}
	return nil
}
