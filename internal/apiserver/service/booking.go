package service

import "github.com/tmrrwnxtsn/aero-table-booking-api/internal/apiserver/store"

// BookingService представляет бизнес-логику работы с бронями.
type BookingService interface {
}

type BookingServiceImpl struct {
	bookingRepo store.BookingRepository
}

func NewBookingService(bookingRepo store.BookingRepository) *BookingServiceImpl {
	return &BookingServiceImpl{bookingRepo: bookingRepo}
}
