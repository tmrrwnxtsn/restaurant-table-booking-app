package store

import "github.com/tmrrwnxtsn/aero-table-booking-api/internal/model"

// RestaurantRepository представляет таблицу с информацией о ресторанах.
type RestaurantRepository interface {
	// Create создаёт новую запись о ресторане.
	Create(name string, averageWaitingTime int, averageCheck float64) (uint64, error)
	// GetAll возвращает список всех ресторанов.
	GetAll() ([]model.Restaurant, error)
	// GetByID возвращает ресторан по его ID.
	GetByID(id uint64) (*model.Restaurant, error)
	// Update обновляет информацию о ресторане по его ID.
	Update(id uint64, data model.UpdateRestaurantData) error
}

// TableRepository представляет таблицу с информацией о столиках в ресторанах.
type TableRepository interface {
}

// ClientRepository представляет таблицу с информацией о клиентах ресторанах.
type ClientRepository interface {
}

// BookingRepository представляет таблицу с информацией о совершённых клиентами бронях.
type BookingRepository interface {
}
