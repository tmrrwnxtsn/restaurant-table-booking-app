package service

import "errors"

// ErrInvalidData возникает, когда из формы выбора времени брони и ввода количества человек приходят некорректные данные.
var ErrInvalidData = errors.New("invalid input data")
