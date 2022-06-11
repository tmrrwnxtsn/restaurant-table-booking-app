package postgres

import "github.com/tmrrwnxtsn/aero-table-booking-api/internal/store"

var _ store.ClientRepository = (*ClientRepository)(nil)

type ClientRepository struct {
	store *Store
}

func NewClientRepository(store *Store) *ClientRepository {
	return &ClientRepository{store: store}
}
