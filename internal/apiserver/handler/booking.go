package handler

import (
	"net/http"
	"strconv"

	"github.com/go-chi/render"

	"github.com/tmrrwnxtsn/restaurant-table-booking-app/internal/apiserver/model"
)

// createBookingRequest представляет тело запроса на создание брони в ресторане.
type createBookingRequest struct {
	PeopleNumber int `json:"people_number" example:"3"`
	// DesiredDatetime представляет дату и время посещения ресторана в рамках брони
	DesiredDatetime string `json:"desired_datetime" example:"2022.06.16 17:03"`
	// ClientName имя клиента, оформляющего бронь.
	ClientName string `json:"client_name" example:"Павел"`
	// ClientPhone телефон клиента, оформляющего бронь.
	ClientPhone string `json:"client_phone" example:"89876545654"`
}

// Bind осуществляет пост-обработку запроса.
func (r *createBookingRequest) Bind(_ *http.Request) error {
	if r.PeopleNumber == 0 || r.DesiredDatetime == "" || r.ClientName == "" || r.ClientPhone == "" {
		return ErrBookingMissingFields
	}
	return nil
}

// createBookingResponse представляет тело ответа на создание брони.
type createBookingResponse struct {
	ID uint64 `json:"id" example:"1"`
}

// Render осуществляет предобработку ответа.
func (r *createBookingResponse) Render(_ http.ResponseWriter, _ *http.Request) error {
	return nil
}

// createBooking godoc
// @Summary  Оформить бронь в выбранном ресторане
// @Tags     bookings
// @Accept   json
// @Produce  json
// @Param    restaurant_id  path      string                 true  "ID ресторана"
// @Param    input          body      createBookingRequest   true  "Информация о брони"
// @Success  201            {object}  createBookingResponse  "ok"
// @Failure  400            {object}  errResponse            "Некорректные данные брони"
// @Failure  500            {object}  errResponse            "Ошибка на стороне сервера"
// @Router   /restaurants/{restaurant_id}/bookings/ [post]
func (h *Handler) createBooking(w http.ResponseWriter, r *http.Request) {
	restaurant := r.Context().Value(restaurantCtxKey).(*model.Restaurant)

	data := &createBookingRequest{}
	if err := render.Bind(r, data); err != nil {
		_ = render.Render(w, r, errInvalidRequest(err))
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
		_ = render.Render(w, r, errServiceFailure(err))
		return
	}

	render.Status(r, http.StatusCreated)
	_ = render.Render(w, r, &createBookingResponse{
		ID: bookingID,
	})
}

// listBookingsResponse представляет тело ответа на получение списка броней ресторана.
type listBookingsResponse struct {
	Data []model.Booking `json:"data"`
}

// Render осуществляет предобработку ответа.
func (r *listBookingsResponse) Render(_ http.ResponseWriter, _ *http.Request) error {
	return nil
}

// listBookings godoc
// @Summary  Получить список всех броней, совершённых в ресторане
// @Tags     bookings
// @Accept   json
// @Produce  json
// @Param    restaurant_id  path      string                true  "ID ресторана"
// @Success  200            {object}  listBookingsResponse  "ok"
// @Failure  400            {object}  errResponse           "Некорректный restaurant_id"
// @Failure  500            {object}  errResponse           "Ошибка на стороне сервера"
// @Router   /restaurants/{restaurant_id}/bookings/ [get]
func (h *Handler) listBookings(w http.ResponseWriter, r *http.Request) {
	restaurant := r.Context().Value(restaurantCtxKey).(*model.Restaurant)

	bookings, err := h.service.BookingService.GetAll(restaurant.ID)
	if err != nil {
		_ = render.Render(w, r, errServiceFailure(err))
		return
	}

	_ = render.Render(w, r, &listBookingsResponse{
		Data: bookings,
	})
}
