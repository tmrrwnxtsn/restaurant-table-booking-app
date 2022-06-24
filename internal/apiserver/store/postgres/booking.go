package postgres

import (
	"context"
	"fmt"
	"time"

	"github.com/tmrrwnxtsn/restaurant-table-booking-app/internal/apiserver/model"
	"github.com/tmrrwnxtsn/restaurant-table-booking-app/internal/apiserver/store"
)

const (
	// bookingTable представляет название таблицы в БД, содержащей информацию о бронях.
	bookingTable = "bookings"
	// bookingsTablesTable представляет название таблицы в БД, содержащей информацию о связях между бронями и столиками ресторанов.
	bookingsTablesTable = "bookings_tables"
)

var _ store.BookingRepository = (*BookingRepository)(nil)

// BookingRepository представляет реализацю store.BookingRepository.
type BookingRepository struct {
	store *Store
}

func NewBookingRepository(store *Store) *BookingRepository {
	return &BookingRepository{store: store}
}

func (r *BookingRepository) Create(clientName, clientPhone string, bookedDate, bookedTimeFrom time.Time, tableIDs ...uint64) (uint64, error) {
	// хелпер-функция для выхода с ошибкой
	fail := func(err error) (uint64, error) {
		return 0, fmt.Errorf("create booking: %w", err)
	}

	// инициируем транзакцию
	ctx := context.Background()
	tx, err := r.store.db.BeginTx(ctx, nil)
	if err != nil {
		return fail(err)
	}
	defer tx.Rollback()

	// добавляем в таблицу с бронями новую бронь, возвращая её ID
	createBookingQuery := fmt.Sprintf(
		"INSERT INTO %s (client_name, client_phone, booked_date, booked_time_from, booked_time_to) VALUES ($1, $2, $3, $4, $5) RETURNING id",
		bookingTable,
	)
	var bookingID uint64
	if err = tx.QueryRowContext(ctx,
		createBookingQuery, clientName, clientPhone, bookedDate, bookedTimeFrom, bookedTimeFrom.Add(2*time.Hour),
	).Scan(&bookingID); err != nil {
		return fail(err)
	}

	// привязываем все столики, которые мы хотим забранировать, к только что созданной брони
	createBookingsTablesQuery := fmt.Sprintf(
		"INSERT INTO %s (booking_id, table_id) VALUES ($1, $2)",
		bookingsTablesTable,
	)
	for _, tableID := range tableIDs {
		_, err = tx.ExecContext(ctx, createBookingsTablesQuery, bookingID, tableID)
		if err != nil {
			return fail(err)
		}
	}

	// завершаем транзакцию
	if err = tx.Commit(); err != nil {
		return fail(err)
	}

	return bookingID, nil
}

func (r *BookingRepository) GetAll(restaurantID uint64) ([]model.Booking, error) {
	getAllBookingsQuery := fmt.Sprintf(
		"SELECT DISTINCT b.* "+
			"FROM %s b "+
			"JOIN bookings_tables bt on b.id = bt.booking_id "+
			"JOIN tables t on t.id = bt.table_id "+
			"WHERE restaurant_id = $1",
		bookingTable,
	)

	rows, err := r.store.db.Query(getAllBookingsQuery, restaurantID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var bookings []model.Booking

	for rows.Next() {
		var booking model.Booking
		if err = rows.Scan(
			&booking.ID, &booking.ClientName, &booking.ClientPhone, &booking.BookedDate, &booking.BookedTimeFrom, &booking.BookedTimeTo,
		); err != nil {
			return bookings, err
		}
		bookings = append(bookings, booking)
	}
	if err = rows.Err(); err != nil {
		return bookings, err
	}
	return bookings, nil
}
