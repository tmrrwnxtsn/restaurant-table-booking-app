package postgres

import (
	"database/sql"
	"github.com/tmrrwnxtsn/aero-table-booking-api/internal/store"
)

var _ store.Store = (*Store)(nil)

type Store struct {
	db             *sql.DB
	restaurantRepo store.RestaurantRepository
	tableRepo      store.TableRepository
	clientRepo     store.ClientRepository
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

func (s *Store) Clients() store.ClientRepository {
	if s.clientRepo != nil {
		return s.clientRepo
	}

	s.clientRepo = NewClientRepository(s)

	return s.clientRepo
}

func (s *Store) Bookings() store.BookingRepository {
	if s.bookingRepo != nil {
		return s.bookingRepo
	}

	s.bookingRepo = NewBookingRepository(s)

	return s.bookingRepo
}
