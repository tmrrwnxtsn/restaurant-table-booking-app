package service

import (
	"github.com/tmrrwnxtsn/aero-table-booking-api/internal/model"
	"github.com/tmrrwnxtsn/aero-table-booking-api/internal/store"
)

// RestaurantService представляет бизнес-логику работы с ресторанами.
type RestaurantService interface {
	// Create создаёт ресторан.
	Create(name string, averageWaitingTime int, averageCheck float64) (uint64, error)
	// GetAll получает список всех ресторанов.
	GetAll() ([]model.Restaurant, error)
	// GetByID получает ресторан по его ID.
	GetByID(id uint64) (*model.Restaurant, error)
	// Update обновляет информацию о ресторане по его ID.
	Update(id uint64, data model.UpdateRestaurantData) error
	// Delete удаляет ресторан по его ID.
	Delete(id uint64) error
}

type RestaurantServiceImpl struct {
	restaurantRepo store.RestaurantRepository
}

func NewRestaurantService(restaurantRepo store.RestaurantRepository) RestaurantService {
	return &RestaurantServiceImpl{restaurantRepo: restaurantRepo}
}

func (r *RestaurantServiceImpl) Create(name string, averageWaitingTime int, averageCheck float64) (uint64, error) {
	return r.restaurantRepo.Create(name, averageWaitingTime, averageCheck)
}

func (r *RestaurantServiceImpl) GetAll() ([]model.Restaurant, error) {
	return r.restaurantRepo.GetAll()
}

func (r *RestaurantServiceImpl) GetByID(id uint64) (*model.Restaurant, error) {
	return r.restaurantRepo.GetByID(id)
}

func (r *RestaurantServiceImpl) Update(id uint64, data model.UpdateRestaurantData) error {
	return r.restaurantRepo.Update(id, data)
}

func (r *RestaurantServiceImpl) Delete(id uint64) error {
	return r.restaurantRepo.Delete(id)
}
