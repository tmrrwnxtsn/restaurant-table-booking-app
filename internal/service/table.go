package service

import "github.com/tmrrwnxtsn/aero-table-booking-api/internal/store"

// TableService представляет бизнес-логику работы со столиками.
type TableService interface {
}

type TableServiceImpl struct {
	tableRepo store.TableRepository
}

func NewTableService(tableRepo store.TableRepository) *TableServiceImpl {
	return &TableServiceImpl{tableRepo: tableRepo}
}
