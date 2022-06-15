package postgres

import "github.com/tmrrwnxtsn/aero-table-booking-api/internal/apiserver/store"

var _ store.BookingRepository = (*BookingRepository)(nil)

type BookingRepository struct {
	store *Store
}

func NewBookingRepository(store *Store) *BookingRepository {
	return &BookingRepository{store: store}
}
