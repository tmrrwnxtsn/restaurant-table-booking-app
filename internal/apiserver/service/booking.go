package service

import (
	"fmt"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/tmrrwnxtsn/restaurant-table-booking-app/internal/apiserver/model"
	"github.com/tmrrwnxtsn/restaurant-table-booking-app/internal/apiserver/store"
)

// BookingService представляет бизнес-логику работы с бронями.
type BookingService interface {
	// Create создаёт бронь в ресторане на выбранные дату, время и количество человек.
	Create(details model.BookingDetails) (uint64, error)
	// GetAll возвращает список всех броней ресторана.
	GetAll(restaurantID uint64) ([]model.Booking, error)
}

// BookingServiceImpl представляет реализацию BookingService.
type BookingServiceImpl struct {
	bookingRepo store.BookingRepository
	tableRepo   store.TableRepository
}

func NewBookingService(bookingRepo store.BookingRepository, tableRepo store.TableRepository) *BookingServiceImpl {
	return &BookingServiceImpl{bookingRepo: bookingRepo, tableRepo: tableRepo}
}

func (s *BookingServiceImpl) Create(details model.BookingDetails) (uint64, error) {
	desiredDateTime := strings.Split(details.DesiredDatetime, " ")

	// получаем доступные для брони столики в выбранном ресторане
	tables, err := s.tableRepo.GetAllAvailable(details.RestaurantID, desiredDateTime[0], desiredDateTime[1])
	if err != nil {
		return 0, err
	}

	peopleNum, err := strconv.Atoi(details.PeopleNumber)
	if err != nil {
		return 0, fmt.Errorf("%w: %s", ErrInvalidData, err.Error())
	}

	// подсчёт общего количество доступных мест в ресторане
	availableSeatsNum := 0
	for _, table := range tables {
		availableSeatsNum += table.SeatsNumber
	}

	// если суммарное количество доступных мест меньше, чем хочет прийти людей
	if availableSeatsNum < peopleNum {
		return 0, ErrNotEnoughSeatsInRestaurant
	}

	// алгоритм бронирования столиков:
	// 	1) доступные столики сортируются по возрастанию количества мест
	// 	2) пока суммарное количество мест у забронированных столиков не превысит (или будет равно) количество человек,
	//		которые желают прийти в ресторан, на каждой итерации бронируется столик

	// 1
	sort.SliceStable(tables, func(i, j int) bool {
		return tables[i].SeatsNumber < tables[j].SeatsNumber
	})

	// 2
	// столики, которые будут забронированы после создания брони (изначально пусто)
	bookedTables := make([]uint64, 0)
	// отслеживаем количество занятых мест
	bookedSeatsCurr := 0
	for i := 0; i < len(tables) && bookedSeatsCurr < peopleNum; i++ {
		bookedTables = append(bookedTables, tables[i].ID)
		bookedSeatsCurr += tables[i].SeatsNumber
	}

	dateTime, err := time.Parse("2006.01.02 15:04", details.DesiredDatetime)
	if err != nil {
		return 0, fmt.Errorf("%w: %s", ErrInvalidData, err.Error())
	}

	return s.bookingRepo.Create(details.ClientName, details.ClientPhone, dateTime, dateTime, bookedTables...)
}

func (s *BookingServiceImpl) GetAll(restaurantID uint64) ([]model.Booking, error) {
	return s.bookingRepo.GetAll(restaurantID)
}
