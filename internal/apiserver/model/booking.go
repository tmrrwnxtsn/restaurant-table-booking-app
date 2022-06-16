package model

import (
	"fmt"
	"time"
)

// Booking представляет бронь.
type Booking struct {
	ID uint64 `json:"id" example:"3"`
	// ClientName представляет имя клиента, оформляющего бронь.
	ClientName string `json:"client_name" example:"Павел"`
	// ClientPhone представляет телефон клиента, оформляющего бронь.
	ClientPhone string `json:"client_phone" example:"89485722648"`
	// BookedDate представляет дату посещения ресторана в рамках брони.
	BookedDate ShortFormattedDate `json:"booked_date" example:"2022.06.16"`
	// BookedTimeFrom представляет время начала брони.
	BookedTimeFrom ShortFormattedTime `json:"booked_time_from" example:"14:30"`
	// BookedTimeTo представляет время конца брони.
	BookedTimeTo ShortFormattedTime `json:"booked_time_to" example:"16:30"`
}

// ShortFormattedTime представляет время в формате "15:04".
type ShortFormattedTime time.Time

func (t ShortFormattedTime) MarshalJSON() ([]byte, error) {
	stamp := fmt.Sprintf("\"%s\"", time.Time(t).Format("15:04"))
	return []byte(stamp), nil
}

// ShortFormattedDate представляет дату в формате "2006.01.02".
type ShortFormattedDate time.Time

func (t ShortFormattedDate) MarshalJSON() ([]byte, error) {
	stamp := fmt.Sprintf("\"%s\"", time.Time(t).Format("2006.01.02"))
	return []byte(stamp), nil
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
	// DesiredDatetime представляет дату и время посещения ресторана в рамках брони (строка вида "2022.06.16 17:03")
	DesiredDatetime string
	// ClientName имя клиента, оформляющего бронь.
	ClientName string
	// ClientPhone телефон клиента, оформляющего бронь.
	ClientPhone string
}
