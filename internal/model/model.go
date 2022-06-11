package model

import "time"

// Restaurant представляет ресторан.
type Restaurant struct {
	ID   uint64 `json:"id"`
	Name string `json:"name"`
	// AverageWaitingTime среднее время ожидания заказа в минутах.
	AverageWaitingTime int `json:"average_waiting_time"`
	// AverageCheck средний чек на блюдо в ресторане.
	AverageCheck float64 `json:"average_check"`
}

// Table представляет столик в ресторане.
type Table struct {
	ID           uint64 `json:"id"`
	RestaurantID uint64 `json:"restaurant_id"`
	// SeatsNumber вместимость столика.
	SeatsNumber int `json:"seats_number"`
}

// Client представляет клиента ресторана, бронирующего столики.
type Client struct {
	ID    uint64 `json:"id"`
	Name  string `json:"name"`
	Phone string `json:"phone"`
}

// Booking представляет бронь.
type Booking struct {
	ID         uint64    `json:"id"`
	ClientID   uint64    `json:"client_id"`
	BookedFrom time.Time `json:"booked_from"`
	BookedTo   time.Time `json:"booked_to"`
}

// BookingsTables представляет таблицу в БД, в которой хранятся столики и брони, к которым они относятся.
type BookingsTables struct {
	ID        uint64
	BookingID uint64
	TableID   uint64
}
