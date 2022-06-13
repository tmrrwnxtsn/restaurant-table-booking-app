package service

import (
	"fmt"
	"github.com/tmrrwnxtsn/aero-table-booking-api/internal/model"
	"github.com/tmrrwnxtsn/aero-table-booking-api/internal/store"
	"strconv"
	"time"
)

// RestaurantService представляет бизнес-логику работы с ресторанами.
type RestaurantService interface {
	// Create создаёт ресторан.
	Create(name string, averageWaitingTime int, averageCheck float64) (uint64, error)
	// GetAll получает список всех ресторанов.
	GetAll() ([]model.Restaurant, error)
	// GetAllAvailable возвращает список ресторанов, в которых можно забронировать столики.
	GetAllAvailable(desiredDateTime, peopleNumber string) ([]model.Restaurant, error)
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

func NewRestaurantService(restaurantRepo store.RestaurantRepository) *RestaurantServiceImpl {
	return &RestaurantServiceImpl{restaurantRepo: restaurantRepo}
}

func (r *RestaurantServiceImpl) Create(name string, averageWaitingTime int, averageCheck float64) (uint64, error) {
	return r.restaurantRepo.Create(name, averageWaitingTime, averageCheck)
}

func (r *RestaurantServiceImpl) GetAll() ([]model.Restaurant, error) {
	return r.restaurantRepo.GetAll()
}

func (r *RestaurantServiceImpl) GetAllAvailable(desiredDateTime, peopleNumber string) ([]model.Restaurant, error) {
	peopleNum, err := strconv.Atoi(peopleNumber)
	if err != nil {
		return nil, fmt.Errorf("%w: %s", ErrInvalidData, err.Error())
	}

	if peopleNum < 1 {
		return nil, fmt.Errorf("%w: the number of people cannot be less than 1", ErrInvalidData)
	}

	dateTime, err := time.Parse("2006-01-02T15:04", desiredDateTime)
	if err != nil {
		return nil, fmt.Errorf("%w: %s", ErrInvalidData, err.Error())
	}

	// если указанное время пользователем в прошлом, то ругаемся на некорректный ввод
	if time.Now().After(dateTime) {
		return nil, fmt.Errorf("%w: the date and time of booking cannot be in the past", ErrInvalidData)
	}

	// рестораны работают с 9:00 до 23:00 (последнюю бронь можно создать на 21:00)
	dateTimeHour := dateTime.Hour()
	dateTimeMinute := dateTime.Minute()
	if dateTimeHour < 9 || dateTimeHour >= 22 || dateTimeHour == 21 && dateTimeMinute > 0 {
		return nil, fmt.Errorf("%w: the restaurant is closed", ErrInvalidData)
	}

	desiredDate := dateTime.Format("2006.01.02")
	desiredTime := dateTime.Format("15:04")

	return r.restaurantRepo.GetAllAvailable(desiredDate, desiredTime, peopleNum)
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
