package service

import "github.com/tmrrwnxtsn/aero-table-booking-api/internal/store"

// RestaurantService представляет бизнес-логику работы с ресторанами.
type RestaurantService interface {
}

type RestaurantServiceImpl struct {
	restaurantRepo store.RestaurantRepository
}

func NewRestaurantService(restaurantRepo store.RestaurantRepository) RestaurantService {
	return &RestaurantServiceImpl{restaurantRepo: restaurantRepo}
}
