package postgres

import (
	"database/sql"

	"github.com/tmrrwnxtsn/restaurant-table-booking-app/internal/apiserver/store"
)

var _ store.Store = (*Store)(nil)

type Store struct {
	db             *sql.DB
	restaurantRepo store.RestaurantRepository
	tableRepo      store.TableRepository
	bookingRepo    store.BookingRepository
}

func NewStore(db *sql.DB) *Store {
	return &Store{db: db}
}

func (s *Store) Restaurants() store.RestaurantRepository {
	if s.restaurantRepo != nil {
		return s.restaurantRepo
	}

	s.restaurantRepo = NewRestaurantRepository(s)

	return s.restaurantRepo
}

func (s *Store) Tables() store.TableRepository {
	if s.tableRepo != nil {
		return s.tableRepo
	}

	s.tableRepo = NewTableRepository(s)

	return s.tableRepo
}

func (s *Store) Bookings() store.BookingRepository {
	if s.bookingRepo != nil {
		return s.bookingRepo
	}

	s.bookingRepo = NewBookingRepository(s)

	return s.bookingRepo
}
