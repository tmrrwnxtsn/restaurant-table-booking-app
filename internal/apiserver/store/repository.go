package store

import (
	"github.com/tmrrwnxtsn/aero-table-booking-api/internal/apiserver/model"
	"time"
)

// RestaurantRepository представляет таблицу с информацией о ресторанах.
type RestaurantRepository interface {
	// Create создаёт новую запись о ресторане.
	Create(name string, averageWaitingTime int, averageCheck float64) (uint64, error)
	// GetAll возвращает список всех ресторанов.
	GetAll() ([]model.Restaurant, error)
	// GetAllAvailable возвращает список ресторанов, в которых можно забронировать столики на выбранные дату,
	// время и количество человек. Принимает desiredDate в формате "2006.01.02" и desiredTime - "15:04".
	GetAllAvailable(desiredDate, desiredTime string, peopleNumber int) ([]model.Restaurant, error)
	// GetByID возвращает ресторан по его ID.
	GetByID(id uint64) (*model.Restaurant, error)
	// Update обновляет информацию о ресторане по его ID.
	Update(id uint64, data model.UpdateRestaurantData) error
	// Delete удаляет запись о ресторане по его ID.
	Delete(id uint64) error
}

// TableRepository представляет таблицу с информацией о столиках в ресторанах.
type TableRepository interface {
	// GetAllAvailableByRestaurant возвращает список всех столиков, доступных для бронирования, в конкретном ресторане.
	// Принимает desiredDate в формате "2006.01.02" и desiredTime - "15:04".
	GetAllAvailableByRestaurant(restaurantID uint64, desiredDate, desiredTime string) ([]model.Table, error)
}

// BookingRepository представляет таблицу с информацией о совершённых клиентами бронях.
type BookingRepository interface {
	// Create создаёт новую запись о брони и связывает созданную бронь со столиками, которые бронируются в рамках неё.
	Create(clientName, clientPhone string, bookedDate, bookedTimeFrom time.Time, tableIDs ...uint64) (uint64, error)
}
