package postgres

import (
	"database/sql"
	"fmt"
	"strings"

	"github.com/tmrrwnxtsn/restaurant-table-booking-app/internal/apiserver/model"
	"github.com/tmrrwnxtsn/restaurant-table-booking-app/internal/apiserver/store"
)

// restaurantTable представляет название таблицы в БД, содержащей информацию о ресторанах.
const restaurantTable = "restaurants"

var _ store.RestaurantRepository = (*RestaurantRepository)(nil)

// RestaurantRepository представляет реализацю store.RestaurantRepository.
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

func (r *RestaurantRepository) GetAllAvailable(desiredDate, desiredTime string, peopleNumber int) ([]model.Restaurant, error) {
	getAllAvailableRestaurantsQuery := fmt.Sprintf(
		`SELECT r.id, r.name, r.average_waiting_time, r.average_check, SUM(t.seats_number) as available_seats_number
				FROM get_available_tables(date '%s', time '%s') t
				JOIN %s r ON r.id = t.restaurant_id
				GROUP BY r.id, r.name, r.average_waiting_time, r.average_check
				HAVING SUM(t.seats_number) > $1
				ORDER BY r.average_waiting_time, r.average_check`,
		desiredDate, desiredTime, restaurantTable,
	)

	rows, err := r.store.db.Query(getAllAvailableRestaurantsQuery, peopleNumber)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var restaurants []model.Restaurant

	for rows.Next() {
		var restaurant model.Restaurant
		if err = rows.Scan(
			&restaurant.ID, &restaurant.Name, &restaurant.AverageWaitingTime, &restaurant.AverageCheck, &restaurant.AvailableSeatsNumber,
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

func (r *RestaurantRepository) Get(id uint64) (*model.Restaurant, error) {
	getRestaurantQuery := fmt.Sprintf(
		"SELECT * FROM %s WHERE id = $1",
		restaurantTable,
	)

	restaurant := &model.Restaurant{}
	if err := r.store.db.QueryRow(
		getRestaurantQuery, id,
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
	// мы не можем удалить ресторан, если видим по оформленным броням, что клиенты посетят этот ресторан (сегодня или в будущем)
	// поэтому сначала смотрим, есть ли в будущем (или сегодня) брони в этом ресторане
	countBookingsWithThisRestaurantQuery := fmt.Sprintf(
		"SELECT COUNT(*) "+
			"FROM %s "+
			"JOIN %s bt on tables.id = bt.table_id "+
			"JOIN %s b on b.id = bt.booking_id "+
			"WHERE booked_date >= current_date AND restaurant_id = $1",
		tableTable, bookingsTablesTable, bookingTable,
	)
	var bookingsWithThisRestaurant int
	err := r.store.db.QueryRow(
		countBookingsWithThisRestaurantQuery, id,
	).Scan(&bookingsWithThisRestaurant)
	if err != nil {
		return err
	}

	// если нашлись брони в этом ресторане в будущем (или на сегодняшний день)
	if bookingsWithThisRestaurant > 0 {
		return fmt.Errorf("delete restaurant: %w", store.ErrRestaurantIsBooked)
	}

	// если все брони, которые связаны с этим рестораном, были в прошлом,
	// и в будущем (или на сегодняшний день) не ожидается клиентов, то его можно удалить
	deleteRestaurantQuery := fmt.Sprintf(
		"DELETE FROM %s WHERE id = $1",
		restaurantTable,
	)
	_, err = r.store.db.Exec(deleteRestaurantQuery, id)
	return err
}
