package handler

import (
	"errors"
	"github.com/go-chi/render"
	"net/http"
)

// Ошибки, возникающие при работе с маршрутами манипуляции с ресторанами.
var (
	// ErrRestaurantMissingFields возникает, когда в запросе на создание ресторана пропущены обязательные поля.
	ErrRestaurantMissingFields = errors.New("missing required restaurant fields")
	// ErrFindAvailableRestaurants возникает, когда в запросе на поиск доступных ресторанов пропущено либо кол-во человек,
	// либо дата и время.
	ErrFindAvailableRestaurants = errors.New("missing required datetime or people number")
)

// ErrResponse представляет ответ-ошибку.
type ErrResponse struct {
	Err            error `json:"-"`
	HTTPStatusCode int   `json:"-"`

	StatusText string `json:"status"`
	AppCode    int64  `json:"code,omitempty"`
	ErrorText  string `json:"error,omitempty"`
}

// Render осуществляет предобработку ответа ErrResponse.
func (e *ErrResponse) Render(_ http.ResponseWriter, r *http.Request) error {
	render.Status(r, e.HTTPStatusCode)
	return nil
}

// ErrInvalidRequest вкладывает ошибку в кастомную структуру ErrResponse с кодом состояния 400.
// Создаётся при некорректном запросе.
func ErrInvalidRequest(err error) render.Renderer {
	return &ErrResponse{
		Err:            err,
		HTTPStatusCode: http.StatusBadRequest,
		StatusText:     "invalid request",
		ErrorText:      err.Error(),
	}
}

// ErrNotFound вкладывает ошибку в кастомную структуру ErrResponse с кодом состояния 404.
// Создаётся при отсутствии искомого ресурса по указанному URL.
func ErrNotFound(err error) render.Renderer {
	return &ErrResponse{
		Err:            err,
		HTTPStatusCode: http.StatusNotFound,
		StatusText:     "resource not found",
		ErrorText:      err.Error(),
	}
}

// ErrRender вкладывает ошибку в кастомную структуру ErrResponse с кодом состояния 422.
// Создаётся при возникновении ошибки обработки ответа.
func ErrRender(err error) render.Renderer {
	return &ErrResponse{
		Err:            err,
		HTTPStatusCode: http.StatusUnprocessableEntity,
		StatusText:     "error rendering response",
		ErrorText:      err.Error(),
	}
}

// ErrServiceFailure вкладывает ошибку в кастомную структуру ErrResponse с кодом состояния 500.
// Создаётся при возникновении ошибки на стороне сервера.
func ErrServiceFailure(err error) render.Renderer {
	return &ErrResponse{
		Err:            err,
		HTTPStatusCode: http.StatusInternalServerError,
		StatusText:     "service failure",
		ErrorText:      err.Error(),
	}
}
