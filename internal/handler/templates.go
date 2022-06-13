package handler

import (
	"errors"
	"github.com/go-chi/render"
	"github.com/tmrrwnxtsn/aero-table-booking-api/internal/model"
	"github.com/tmrrwnxtsn/aero-table-booking-api/internal/service"
	"html/template"
	"net/http"
)

const templatesPattern = "web/templates/*.gohtml"

var tmpls *template.Template

func init() {
	tmpls = template.Must(template.ParseGlob(templatesPattern))
}

// SiteDetails представляет данные, которые помещаются в шаблоны для отображения в визуальном интерфейсе.
type SiteDetails struct {
	Title       string
	Restaurants []model.Restaurant
}

// ErrorDetails представляет данные, которые помещаются в шаблон, который оповещает о произошедшей ошибке.
type ErrorDetails struct {
	Title      string
	StatusCode int
	Text       string
}

// home отображает содержание стартовой страницы.
func (h *Handler) home(w http.ResponseWriter, r *http.Request) {
	renderTemplate(w, r,
		"index",
		&SiteDetails{
			Title: "Бронирование столиков в ресторанах",
		},
	)
}

// restaurants отображает содержание страницы с выбором ресторанов в зависимости от желаемых даты, времени и
// количества человек, на которых в ресторане бронируются столики.
func (h *Handler) restaurants(w http.ResponseWriter, r *http.Request) {
	desiredDateTime := r.URL.Query().Get("desired_datetime")
	peopleNumber := r.URL.Query().Get("people_number")

	if desiredDateTime == "" || peopleNumber == "" {
		renderTemplate(w, r,
			"error",
			&ErrorDetails{
				Title:      "Произошла ошибка",
				Text:       ErrFindAvailableRestaurants.Error(),
				StatusCode: http.StatusBadRequest,
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
		renderTemplate(w, r,
			"error",
			&ErrorDetails{
				Title:      "Произошла ошибка",
				Text:       err.Error(),
				StatusCode: statusCode,
			},
		)
		return
	}

	renderTemplate(w, r,
		"restaurants",
		&SiteDetails{
			Title:       "Выбор ресторана",
			Restaurants: restaurants,
		},
	)
}

func renderTemplate(w http.ResponseWriter, r *http.Request, name string, data interface{}) {
	if err := tmpls.ExecuteTemplate(w, name, data); err != nil {
		_ = render.Render(w, r, ErrServiceFailure(err))
		return
	}
}
