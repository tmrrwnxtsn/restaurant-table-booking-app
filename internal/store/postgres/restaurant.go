package postgres

import (
	"fmt"
	"github.com/tmrrwnxtsn/aero-table-booking-api/internal/model"
	"github.com/tmrrwnxtsn/aero-table-booking-api/internal/store"
)

const restaurantTable = "restaurants"

var _ store.RestaurantRepository = (*RestaurantRepository)(nil)

type RestaurantRepository struct {
	store *Store
}

func NewRestaurantRepository(store *Store) *RestaurantRepository {
	return &RestaurantRepository{store: store}
}

func (r *RestaurantRepository) Create(name string, averageWaitingTime int, averageCheck float64) (uint64, error) {
	createRestaurantQuery := fmt.Sprintf(
		"INSERT INTO %s (name, average_waiting_time, average_check) VALUES ($1, $2, $3) RETURNING id",
		restaurantTable,
	)

	var id uint64
	err := r.store.db.QueryRow(
		createRestaurantQuery,
		name, averageWaitingTime, averageCheck,
	).Scan(&id)
	if err != nil {
		return 0, err
	}
	return id, nil
}

func (r *RestaurantRepository) GetAll() ([]model.Restaurant, error) {
	listRestaurantsQuery := fmt.Sprintf(
		"SELECT * FROM %s ORDER BY average_waiting_time, average_check",
		restaurantTable,
	)

	rows, err := r.store.db.Query(listRestaurantsQuery)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var restaurants []model.Restaurant

	for rows.Next() {
		var restaurant model.Restaurant
		if err = rows.Scan(
			&restaurant.ID, &restaurant.Name, &restaurant.AverageWaitingTime, &restaurant.AverageCheck,
		); err != nil {
			return restaurants, err
		}
		restaurants = append(restaurants, restaurant)
	}
	if err = rows.Err(); err != nil {
		return restaurants, err
	}
	return restaurants, nil
}
