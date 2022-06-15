package model

import "time"

// Booking представляет бронь.
type Booking struct {
	ID uint64 `json:"id"`
	// ClientName представляет имя клиента, оформляющего бронь.
	ClientName string `json:"client_name"`
	// ClientPhone представляет телефон клиента, оформляющего бронь.
	ClientPhone string `json:"client_phone"`
	// BookedDate представляет дату посещения ресторана в рамках брони.
	BookedDate time.Time `json:"booked_date"`
	// BookedTimeFrom представляет время начала брони.
	BookedTimeFrom time.Time `json:"booked_time_from"`
	// BookedTimeTo представляет время конца брони.
	BookedTimeTo time.Time `json:"booked_time_to"`
}

// BookingsTables представляет таблицу в БД, в которой хранятся столики и брони, к которым они относятся.
type BookingsTables struct {
	ID        uint64
	BookingID uint64
	TableID   uint64
}

// BookingDetails представляет данные, необходимые для оформления брони.
type BookingDetails struct {
	// RestaurantID представляет ID ресторана, в котором оформляется бронь.
	RestaurantID uint64
	// PeopleNumber представляет количество человек, собирающихся посетит ресторан по брони.
	PeopleNumber string
	// DesiredDatetime представляет дату и время посещения ресторана в рамках брони (строка вида "16.06.2022 17:03")
	DesiredDatetime string
	// ClientName имя клиента, оформляющего бронь.
	ClientName string
	// ClientPhone телефон клиента, оформляющего бронь.
	ClientPhone string
}
