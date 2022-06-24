package service

import (
	"github.com/tmrrwnxtsn/restaurant-table-booking-app/internal/apiserver/model"
	"github.com/tmrrwnxtsn/restaurant-table-booking-app/internal/apiserver/store"
)

// TableService представляет бизнес-логику работы со столиками.
type TableService interface {
	// Create создаёт столик в ресторане.
	Create(restaurantID uint64, seatsNumber int) (uint64, error)
	// GetAllAvailable возвращает список доступных для брони столиков конкретного ресторана.
	// Принимает desiredDate в формате "2006.01.02" и desiredTime - "15:04".
	GetAllAvailable(restaurantID uint64, desiredDate, desiredTime string) ([]model.Table, error)
	// GetAll возвращает список всех столиков ресторана.
	GetAll(restaurantID uint64) ([]model.Table, error)
	// Get получает столик ресторана по его ID.
	Get(id uint64) (*model.Table, error)
	// Update обновляет информацию о столике ресторана по его ID.
	Update(id uint64, data model.UpdateTableData) error
	// Delete удаляет столик из ресторана по его ID, ЕСЛИ ОН НЕ ЗАБРОНИРОВАН НА БУДУЩЕЕ ВРЕМЯ.
	Delete(id uint64) error
}

// TableServiceImpl представляет реализацю TableService.
type TableServiceImpl struct {
	tableRepo store.TableRepository
}

func NewTableService(tableRepo store.TableRepository) *TableServiceImpl {
	return &TableServiceImpl{tableRepo: tableRepo}
}

func (s *TableServiceImpl) Create(restaurantID uint64, seatsNumber int) (uint64, error) {
	return s.tableRepo.Create(restaurantID, seatsNumber)
}

func (s *TableServiceImpl) GetAllAvailable(restaurantID uint64, desiredDate, desiredTime string) ([]model.Table, error) {
	return s.tableRepo.GetAllAvailable(restaurantID, desiredDate, desiredTime)
}

func (s *TableServiceImpl) GetAll(restaurantID uint64) ([]model.Table, error) {
	return s.tableRepo.GetAll(restaurantID)
}

func (s *TableServiceImpl) Get(id uint64) (*model.Table, error) {
	return s.tableRepo.Get(id)
}

func (s *TableServiceImpl) Update(id uint64, data model.UpdateTableData) error {
	return s.tableRepo.Update(id, data)
}

func (s *TableServiceImpl) Delete(id uint64) error {
	return s.tableRepo.Delete(id)
}
