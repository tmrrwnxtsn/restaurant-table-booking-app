package service

import "github.com/tmrrwnxtsn/restaurant-table-booking-app/internal/apiserver/store"

// Services представляет слой бизнес-логики.
type Services struct {
	// BookingService представляет бизнес-логику работы с бронями.
	BookingService BookingService
	// RestaurantService представляет бизнес-логику работы с ресторанами.
	RestaurantService RestaurantService
	// TableService представляет бизнес-логику работы со столиками.
	TableService TableService
}

func NewServices(store store.Store) *Services {
	return &Services{
		BookingService:    NewBookingService(store.Bookings(), store.Tables()),
		RestaurantService: NewRestaurantService(store.Restaurants()),
		TableService:      NewTableService(store.Tables()),
	}
}
