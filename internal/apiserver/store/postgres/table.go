package postgres

import (
	"database/sql"
	"fmt"
	"strings"

	"github.com/tmrrwnxtsn/restaurant-table-booking-app/internal/apiserver/model"
	"github.com/tmrrwnxtsn/restaurant-table-booking-app/internal/apiserver/store"
)

// tableTable представляет название таблицы в БД, содержащей записи о столиках в ресторанах.
const tableTable = "tables"

var _ store.TableRepository = (*TableRepository)(nil)

// TableRepository представляет реализацю store.TableRepository.
type TableRepository struct {
	store *Store
}

func NewTableRepository(store *Store) *TableRepository {
	return &TableRepository{store: store}
}

func (r *TableRepository) Create(restaurantID uint64, seatsNumber int) (uint64, error) {
	createTableQuery := fmt.Sprintf(
		"INSERT INTO %s (restaurant_id, seats_number) VALUES ($1, $2) RETURNING id",
		tableTable,
	)

	var id uint64
	err := r.store.db.QueryRow(
		createTableQuery, restaurantID, seatsNumber,
	).Scan(&id)
	if err != nil {
		return 0, err
	}
	return id, nil
}

func (r *TableRepository) GetAllAvailable(restaurantID uint64, desiredDate, desiredTime string) ([]model.Table, error) {
	getAllAvailableTablesQuery := fmt.Sprintf(
		"SELECT * FROM get_available_tables(date '%s', time '%s') WHERE restaurant_id = $1",
		desiredDate, desiredTime,
	)

	rows, err := r.store.db.Query(getAllAvailableTablesQuery, restaurantID)
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

func (r *TableRepository) GetAll(restaurantID uint64) ([]model.Table, error) {
	getAllTablesQuery := fmt.Sprintf(
		"SELECT * FROM %s WHERE restaurant_id = $1",
		tableTable,
	)

	rows, err := r.store.db.Query(getAllTablesQuery, restaurantID)
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

func (r *TableRepository) Get(id uint64) (*model.Table, error) {
	getTableQuery := fmt.Sprintf(
		"SELECT * FROM %s WHERE id = $1",
		tableTable,
	)

	table := &model.Table{}
	if err := r.store.db.QueryRow(
		getTableQuery, id,
	).Scan(&table.ID, &table.RestaurantID, &table.SeatsNumber); err != nil {
		if err == sql.ErrNoRows {
			return nil, store.ErrTableNotFound
		}
		return nil, err
	}
	return table, nil
}

func (r *TableRepository) Update(id uint64, data model.UpdateTableData) error {
	setValues := make([]string, 0, 1)
	args := make([]interface{}, 0, 1)
	argId := 1

	if data.SeatsNumber != nil {
		setValues = append(setValues, fmt.Sprintf("seats_number=$%d", argId))
		args = append(args, *data.SeatsNumber)
		argId++
	}

	setQuery := strings.Join(setValues, ", ")

	updateTableQuery := fmt.Sprintf("UPDATE %s SET %s WHERE id = $%d",
		tableTable, setQuery, argId,
	)

	args = append(args, id)

	_, err := r.store.db.Exec(updateTableQuery, args...)
	return err
}

func (r *TableRepository) Delete(id uint64) error {
	// мы не можем удалить столик из ресторана, если видим, что клиенты в будущем придут и сядут за него,
	// поэтому сначала смотрим, есть ли в будущем или сегодня брони с этим столиком
	countBookingsWithThisTableQuery := fmt.Sprintf(
		"SELECT COUNT(*) "+
			"FROM %s "+
			"JOIN %s bt on tables.id = bt.table_id "+
			"JOIN %s b on b.id = bt.booking_id "+
			"WHERE booked_date >= current_date AND tables.id = $1",
		tableTable, bookingsTablesTable, bookingTable,
	)
	var bookingsWithThisTable int
	err := r.store.db.QueryRow(
		countBookingsWithThisTableQuery, id,
	).Scan(&bookingsWithThisTable)
	if err != nil {
		return err
	}

	// если нашлись брони на сегодняшний день или день в будущем
	if bookingsWithThisTable > 0 {
		return fmt.Errorf("delete table: %w", store.ErrTableIsBooked)
	}

	// если столик был забронирован в прошлом, и в будущем не ожидается использований этого столика, можно его удалить
	deleteTableQuery := fmt.Sprintf(
		"DELETE FROM %s WHERE id = $1",
		tableTable,
	)
	_, err = r.store.db.Exec(deleteTableQuery, id)
	return err
}
