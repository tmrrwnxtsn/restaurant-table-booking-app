package store

import (
	"time"

	"github.com/tmrrwnxtsn/restaurant-table-booking-app/internal/apiserver/model"
)

// RestaurantRepository представляет методы работы с информацией о ресторанах.
type RestaurantRepository interface {
	// Create создаёт новую запись о ресторане.
	Create(name string, averageWaitingTime int, averageCheck float64) (uint64, error)
	// GetAll возвращает список всех ресторанов.
	GetAll() ([]model.Restaurant, error)
	// GetAllAvailable возвращает список ресторанов, в которых можно забронировать столики на выбранные дату,
	// время и количество человек. Принимает desiredDate в формате "2006.01.02" и desiredTime - "15:04".
	GetAllAvailable(desiredDate, desiredTime string, peopleNumber int) ([]model.Restaurant, error)
	// Get возвращает ресторан по его ID.
	Get(id uint64) (*model.Restaurant, error)
	// Update обновляет информацию о ресторане по его ID.
	Update(id uint64, data model.UpdateRestaurantData) error
	// Delete удаляет запись о ресторане по его ID.
	Delete(id uint64) error
}

// TableRepository представляет методы работы с информацией о столиках в ресторанах.
type TableRepository interface {
	// Create создаёт новую запись о столике в ресторане.
	Create(restaurantID uint64, seatsNumber int) (uint64, error)
	// GetAllAvailable возвращает список всех столиков, доступных для бронирования, в конкретном ресторане.
	// Принимает desiredDate в формате "2006.01.02" и desiredTime - "15:04".
	GetAllAvailable(restaurantID uint64, desiredDate, desiredTime string) ([]model.Table, error)
	// GetAll возвращает список всех столиков ресторана.
	GetAll(restaurantID uint64) ([]model.Table, error)
	// Get возвращает столик ресторана по его ID.
	Get(id uint64) (*model.Table, error)
	// Update обновляет информацию о столике ресторана по его ID.
	Update(id uint64, data model.UpdateTableData) error
	// Delete удаляет столик из ресторана по его ID, ЕСЛИ ОН НЕ ЗАБРОНИРОВАН НА БУДУЩЕЕ ВРЕМЯ.
	Delete(id uint64) error
}

// BookingRepository представляет методы работы с информацией о совершённых клиентами бронях.
type BookingRepository interface {
	// Create создаёт новую запись о брони и связывает созданную бронь со столиками, которые бронируются в рамках неё.
	Create(clientName, clientPhone string, bookedDate, bookedTimeFrom time.Time, tableIDs ...uint64) (uint64, error)
	// GetAll возвращает список всех броней ресторана.
	GetAll(restaurantID uint64) ([]model.Booking, error)
}
