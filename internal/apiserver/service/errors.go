package service

import "errors"

var (
	// ErrInvalidData возникает, когда из формы выбора времени брони и ввода количества человек приходят некорректные данные.
	ErrInvalidData = errors.New("invalid input data")
	// ErrNotEnoughSeatsInRestaurant возникает в процессе создания брони, когда в ресторане не достаточно свободных мест.
	ErrNotEnoughSeatsInRestaurant = errors.New("there are not enough seats in the restaurant to make a booking")
)
