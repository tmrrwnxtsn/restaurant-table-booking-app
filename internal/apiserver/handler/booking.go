package handler

import (
	"github.com/go-chi/render"
	"github.com/tmrrwnxtsn/aero-table-booking-api/internal/apiserver/model"
	"net/http"
	"strconv"
)

// CreateBookingRequest представляет тело запроса на создание брони в ресторане.
type CreateBookingRequest struct {
	PeopleNumber int `json:"people_number"`
	// DesiredDatetime представляет дату и время посещения ресторана в рамках брони (строка вида "2022.06.16 17:03")
	DesiredDatetime string `json:"desired_datetime"`
	// ClientName имя клиента, оформляющего бронь.
	ClientName string `json:"client_name"`
	// ClientPhone телефон клиента, оформляющего бронь.
	ClientPhone string `json:"client_phone"`
}

// Bind осуществляет пост-обработку запроса CreateBookingRequest.
func (r *CreateBookingRequest) Bind(_ *http.Request) error {
	if r.PeopleNumber == 0 || r.DesiredDatetime == "" || r.ClientName == "" || r.ClientPhone == "" {
		return ErrBookingMissingFields
	}
	return nil
}

// CreateBookingResponse представляет тело ответа на создание брони.
type CreateBookingResponse struct {
	ID uint64 `json:"id"`
}

// Render осуществляет предобработку ответа CreateBookingResponse.
func (r *CreateBookingResponse) Render(_ http.ResponseWriter, _ *http.Request) error {
	return nil
}

// createBooking принимает запросы на создание брони в ресторане.
func (h *Handler) createBooking(w http.ResponseWriter, r *http.Request) {
	restaurant := r.Context().Value(restaurantCtxKey).(*model.Restaurant)

	data := &CreateBookingRequest{}
	if err := render.Bind(r, data); err != nil {
		_ = render.Render(w, r, ErrInvalidRequest(err))
		return
	}

	details := model.BookingDetails{
		RestaurantID:    restaurant.ID,
		PeopleNumber:    strconv.Itoa(data.PeopleNumber),
		DesiredDatetime: data.DesiredDatetime,
		ClientName:      data.ClientName,
		ClientPhone:     data.ClientPhone,
	}

	bookingID, err := h.service.BookingService.Create(details)
	if err != nil {
		_ = render.Render(w, r, ErrServiceFailure(err))
		return
	}

	render.Status(r, http.StatusCreated)
	_ = render.Render(w, r, &CreateBookingResponse{
		ID: bookingID,
	})
}

// ListBookingsResponse представляет тело ответа на получение списка броней ресторана.
type ListBookingsResponse struct {
	Data []model.Booking `json:"data"`
}

// Render осуществляет предобработку ответа ListBookingsResponse.
func (r *ListBookingsResponse) Render(_ http.ResponseWriter, _ *http.Request) error {
	return nil
}

// listBookings принимает запросы на получение списка броней ресторана.
func (h *Handler) listBookings(w http.ResponseWriter, r *http.Request) {
	restaurant := r.Context().Value(restaurantCtxKey).(*model.Restaurant)

	bookings, err := h.service.BookingService.GetAll(restaurant.ID)
	if err != nil {
		_ = render.Render(w, r, ErrServiceFailure(err))
		return
	}

	_ = render.Render(w, r, &ListBookingsResponse{
		Data: bookings,
	})
}
