package handler

import (
	"errors"
	"html/template"
	"net/http"

	"github.com/go-chi/render"

	"github.com/tmrrwnxtsn/restaurant-table-booking-app/internal/apiserver/model"
	"github.com/tmrrwnxtsn/restaurant-table-booking-app/internal/apiserver/service"
)

const templatesPattern = "website/templates/*.gohtml"

var tmpls *template.Template

func init() {
	tmpls = template.Must(template.ParseGlob(templatesPattern))
}

// TemplatesContext представляет данные, которые передаются в gohtml-шаблоны.
type TemplatesContext struct {
	PageTitle   string
	Restaurants []model.Restaurant
	BookingID   uint64

	ErrorCode int
	ErrorText string
}

// home отображает содержание стартовой страницы, где необходимо указать количество человек, дату и время посещения.
func (h *Handler) home(w http.ResponseWriter, r *http.Request) {
	renderTemplate(w, r, "home",
		&TemplatesContext{
			PageTitle: "Бронирование столиков в ресторанах",
		},
	)
}

// restaurants отображает содержание страницы с доступными для брони ресторанами.
func (h *Handler) restaurants(w http.ResponseWriter, r *http.Request) {
	desiredDateTime := r.URL.Query().Get("desired_datetime")
	peopleNumber := r.URL.Query().Get("people_number")

	if desiredDateTime == "" || peopleNumber == "" {
		renderTemplate(w, r, "error",
			&TemplatesContext{
				PageTitle: "Произошла ошибка",
				ErrorText: ErrFindAvailableRestaurants.Error(),
				ErrorCode: http.StatusBadRequest,
			},
		)
		return
	}

	restaurants, err := h.service.RestaurantService.GetAllAvailable(desiredDateTime, peopleNumber)
	if err != nil {
		statusCode := http.StatusInternalServerError
		if errors.Is(err, service.ErrInvalidData) {
			statusCode = http.StatusBadRequest
		}
		renderTemplate(w, r, "error",
			&TemplatesContext{
				PageTitle: "Произошла ошибка",
				ErrorText: err.Error(),
				ErrorCode: statusCode,
			},
		)
		return
	}

	renderTemplate(w, r, "restaurants",
		&TemplatesContext{
			PageTitle:   "Выбор ресторана",
			Restaurants: restaurants,
		},
	)
}

// makeBooking обрабатывает запрос на создание брони в выбранном ресторане.
func (h *Handler) makeBooking(w http.ResponseWriter, r *http.Request) {
	headerContentType := r.Header.Get("Content-Type")
	if headerContentType != "application/x-www-form-urlencoded" {
		renderTemplate(w, r, "error",
			&TemplatesContext{
				PageTitle: "Произошла ошибка",
				ErrorText: ErrMakingBookingContentType.Error(),
				ErrorCode: http.StatusUnsupportedMediaType,
			},
		)
		return
	}

	restaurant := r.Context().Value(restaurantCtxKey).(*model.Restaurant)

	details := model.BookingDetails{
		RestaurantID:    restaurant.ID,
		PeopleNumber:    r.FormValue("people_number"),
		DesiredDatetime: r.FormValue("desired_datetime"),
		ClientName:      r.FormValue("client_name"),
		ClientPhone:     r.FormValue("client_phone"),
	}

	bookingID, err := h.service.BookingService.Create(details)
	if err != nil {
		statusCode := http.StatusInternalServerError
		if errors.Is(err, service.ErrInvalidData) {
			statusCode = http.StatusBadRequest
		}
		renderTemplate(w, r, "error",
			&TemplatesContext{
				PageTitle: "Произошла ошибка",
				ErrorText: err.Error(),
				ErrorCode: statusCode,
			},
		)
		return
	}

	renderTemplate(w, r, "booking-created",
		&TemplatesContext{
			PageTitle: "Бронь успешно оформлена",
			BookingID: bookingID,
		},
	)
}

// renderTemplate обрабатывает шаблон страницы с переданными в него данными.
func renderTemplate(w http.ResponseWriter, r *http.Request, name string, data interface{}) {
	if err := tmpls.ExecuteTemplate(w, name, data); err != nil {
		_ = render.Render(w, r, errServiceFailure(err))
		return
	}
}
