package handler

import (
	"errors"
	"net/http"

	"github.com/go-chi/render"
)

var (
	// ErrRestaurantMissingFields возникает, когда в запросе на создание/получение ресторана пропущены обязательные поля.
	ErrRestaurantMissingFields = errors.New("missing required restaurant fields")
	// ErrTableMissingFields возникает, когда в запросе на создание/получение столика в ресторане пропущены обязательные поля.
	ErrTableMissingFields = errors.New("missing required restaurant table fields")
	// ErrBookingMissingFields возникает, когда в запросе на создание/получение брони пропущены обязательные поля.
	ErrBookingMissingFields = errors.New("missing required booking fields")
	// ErrFindAvailableRestaurants возникает, когда в запросе на поиск доступных ресторанов пропущено либо кол-во человек,
	// либо дата и время.
	ErrFindAvailableRestaurants = errors.New("missing required datetime or people number")
	// ErrMakingBookingContentType возникает, когда в запросе на оформления брони Content Type отличный от
	// application/x-www-form-urlencoded.
	ErrMakingBookingContentType = errors.New("booking data with wrong content type")
)

// errResponse представляет ответ с ошибкой.
type errResponse struct {
	Err     error `json:"-"`
	AppCode int64 `json:"-"`

	HTTPStatusCode int    `json:"code,omitempty" example:"400"`
	StatusText     string `json:"status" example:"invalid request"`
	ErrorText      string `json:"error,omitempty" example:"missing required fields"`
}

// Render осуществляет предобработку ответа errResponse.
func (e *errResponse) Render(_ http.ResponseWriter, r *http.Request) error {
	render.Status(r, e.HTTPStatusCode)
	return nil
}

// errInvalidRequest вкладывает ошибку в кастомную структуру errResponse с кодом состояния http.StatusBadRequest.
// Создаётся при некорректном запросе.
func errInvalidRequest(err error) render.Renderer {
	return &errResponse{
		Err:            err,
		HTTPStatusCode: http.StatusBadRequest,
		StatusText:     "invalid request",
		ErrorText:      err.Error(),
	}
}

// errNotFound вкладывает ошибку в кастомную структуру errResponse с кодом состояния http.StatusNotFound.
// Создаётся при отсутствии искомого ресурса по указанному URL.
func errNotFound(err error) render.Renderer {
	return &errResponse{
		Err:            err,
		HTTPStatusCode: http.StatusNotFound,
		StatusText:     "resource not found",
		ErrorText:      err.Error(),
	}
}

// errRender вкладывает ошибку в кастомную структуру errResponse с кодом состояния http.StatusUnprocessableEntity.
// Создаётся при возникновении ошибки обработки ответа.
func errRender(err error) render.Renderer {
	return &errResponse{
		Err:            err,
		HTTPStatusCode: http.StatusUnprocessableEntity,
		StatusText:     "error rendering response",
		ErrorText:      err.Error(),
	}
}

// errServiceFailure вкладывает ошибку в кастомную структуру errResponse с кодом состояния http.StatusInternalServerError.
// Создаётся при возникновении ошибки на стороне сервера.
func errServiceFailure(err error) render.Renderer {
	return &errResponse{
		Err:            err,
		HTTPStatusCode: http.StatusInternalServerError,
		StatusText:     "service failure",
		ErrorText:      err.Error(),
	}
}
