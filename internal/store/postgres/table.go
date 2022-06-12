package postgres

import "github.com/tmrrwnxtsn/aero-table-booking-api/internal/store"

var _ store.TableRepository = (*TableRepository)(nil)

type TableRepository struct {
	store *Store
}

func NewTableRepository(store *Store) *TableRepository {
	return &TableRepository{store: store}
}
