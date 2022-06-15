package postgres

import (
	"fmt"
	"github.com/tmrrwnxtsn/aero-table-booking-api/internal/apiserver/model"
	"github.com/tmrrwnxtsn/aero-table-booking-api/internal/apiserver/store"
)

var _ store.TableRepository = (*TableRepository)(nil)

type TableRepository struct {
	store *Store
}

func NewTableRepository(store *Store) *TableRepository {
	return &TableRepository{store: store}
}

func (r *TableRepository) GetAllAvailableByRestaurant(restaurantID uint64, desiredDate, desiredTime string) ([]model.Table, error) {
	getAllTablesByRestaurantQuery := fmt.Sprintf(
		"SELECT * FROM get_available_tables(date '%s', time '%s') WHERE restaurant_id = $1",
		desiredDate, desiredTime,
	)

	rows, err := r.store.db.Query(getAllTablesByRestaurantQuery, restaurantID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tables []model.Table

	for rows.Next() {
		var table model.Table
		if err = rows.Scan(
			&table.ID, &table.RestaurantID, &table.SeatsNumber,
		); err != nil {
			return tables, err
		}
		tables = append(tables, table)
	}
	if err = rows.Err(); err != nil {
		return tables, err
	}
	return tables, nil
}
