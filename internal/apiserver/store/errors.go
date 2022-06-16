package store

import "errors"

var (
	// ErrRestaurantNotFound возникает, когда по введённому ID в БД не находится искомого ресторана.
	ErrRestaurantNotFound = errors.New("restaurant not found")
	// ErrTableNotFound возникает, когда по введённому ID в БД не находится искомого ресторана.
	ErrTableNotFound = errors.New("table not found")
	// ErrRestaurantIsBooked возникает при попытке удалить ресторан, в который ещё придут клиенты.
	ErrRestaurantIsBooked = errors.New("clients are expected in the restaurant today or in the future")
	// ErrTableIsBooked возникает при попытке удалить столик, за которым должны будут сидеть клиенты.
	ErrTableIsBooked = errors.New("the table is booked for today or in the future")
)
