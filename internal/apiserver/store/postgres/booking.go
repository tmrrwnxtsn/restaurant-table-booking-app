package postgres

import (
	"context"
	"fmt"
	"github.com/tmrrwnxtsn/aero-table-booking-api/internal/apiserver/store"
	"time"
)

const (
	bookingTable        = "bookings"
	bookingsTablesTable = "bookings_tables"
)

var _ store.BookingRepository = (*BookingRepository)(nil)

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
