package postgres

import "github.com/tmrrwnxtsn/aero-table-booking-api/internal/store"

var _ store.RestaurantRepository = (*RestaurantRepository)(nil)

type RestaurantRepository struct {
	store *Store
}

func NewRestaurantRepository(store *Store) *RestaurantRepository {
	return &RestaurantRepository{store: store}
}
