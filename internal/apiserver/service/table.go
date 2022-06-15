package service

import (
	"github.com/tmrrwnxtsn/aero-table-booking-api/internal/apiserver/model"
	"github.com/tmrrwnxtsn/aero-table-booking-api/internal/apiserver/store"
)

// TableService представляет бизнес-логику работы со столиками.
type TableService interface {
	// GetAllAvailableByRestaurant возвращает список доступных для брони столиков конкретного ресторана.
	// Принимает desiredDate в формате "2006.01.02" и desiredTime - "15:04".
	GetAllAvailableByRestaurant(restaurantID uint64, desiredDate, desiredTime string) ([]model.Table, error)
}

// TableServiceImpl представляет реализацю TableService.
type TableServiceImpl struct {
	tableRepo store.TableRepository
}

func NewTableService(tableRepo store.TableRepository) *TableServiceImpl {
	return &TableServiceImpl{tableRepo: tableRepo}
}

func (s *TableServiceImpl) GetAllAvailableByRestaurant(restaurantID uint64, desiredDate, desiredTime string) ([]model.Table, error) {
	return s.tableRepo.GetAllAvailableByRestaurant(restaurantID, desiredDate, desiredTime)
}
