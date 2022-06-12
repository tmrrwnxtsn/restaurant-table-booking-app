package postgres

import (
	"database/sql"
	"fmt"
	"github.com/tmrrwnxtsn/aero-table-booking-api/internal/model"
	"github.com/tmrrwnxtsn/aero-table-booking-api/internal/store"
	"strings"
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
	getAllRestaurantsQuery := fmt.Sprintf(
		"SELECT * FROM %s ORDER BY average_waiting_time, average_check",
		restaurantTable,
	)

	rows, err := r.store.db.Query(getAllRestaurantsQuery)
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

func (r *RestaurantRepository) GetByID(id uint64) (*model.Restaurant, error) {
	getRestaurantByIDQuery := fmt.Sprintf(
		"SELECT * FROM %s WHERE id = $1",
		restaurantTable,
	)

	restaurant := &model.Restaurant{}
	if err := r.store.db.QueryRow(
		getRestaurantByIDQuery, id,
	).Scan(&restaurant.ID, &restaurant.Name, &restaurant.AverageWaitingTime, &restaurant.AverageCheck); err != nil {
		if err == sql.ErrNoRows {
			return nil, store.ErrRestaurantNotFound
		}
		return nil, err
	}
	return restaurant, nil
}

func (r *RestaurantRepository) Update(id uint64, data model.UpdateRestaurantData) error {
	setValues := make([]string, 0, 3)
	args := make([]interface{}, 0, 3)
	argId := 1

	if data.Name != nil {
		setValues = append(setValues, fmt.Sprintf("name=$%d", argId))
		args = append(args, *data.Name)
		argId++
	}

	if data.AverageWaitingTime != nil {
		setValues = append(setValues, fmt.Sprintf("average_waiting_time=$%d", argId))
		args = append(args, *data.AverageWaitingTime)
		argId++
	}

	if data.AverageCheck != nil {
		setValues = append(setValues, fmt.Sprintf("average_check=$%d", argId))
		args = append(args, *data.AverageCheck)
		argId++
	}

	setQuery := strings.Join(setValues, ", ")

	updateRestaurantQuery := fmt.Sprintf("UPDATE %s SET %s WHERE id = $%d",
		restaurantTable, setQuery, argId,
	)

	args = append(args, id)

	_, err := r.store.db.Exec(updateRestaurantQuery, args...)
	return err
}

func (r *RestaurantRepository) Delete(id uint64) error {
	deleteRestaurantQuery := fmt.Sprintf(
		"DELETE FROM %s WHERE id = $1",
		restaurantTable,
	)

	_, err := r.store.db.Exec(deleteRestaurantQuery, id)
	return err
}
