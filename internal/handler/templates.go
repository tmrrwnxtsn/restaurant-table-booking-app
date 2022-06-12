package handler

import (
	"github.com/go-chi/render"
	"github.com/tmrrwnxtsn/aero-table-booking-api/internal/model"
	"html/template"
	"net/http"
)

const templatesPattern = "web/templates/*.gohtml"

var tmpls *template.Template

func init() {
	tmpls = template.Must(template.ParseGlob(templatesPattern))
}

type SiteDetails struct {
	Metadata    Metadata
	Restaurants []model.Restaurant
}

type Metadata struct {
	Title string
}

// home отображает содержание стартовой страницы.
func (h *Handler) home(w http.ResponseWriter, r *http.Request) {
	renderTemplate(w, r,
		"index",
		&SiteDetails{
			Metadata: Metadata{
				Title: "Бронирование столиков в ресторанах",
			},
		},
	)
}

// restaurants отображает содержание страницы с выбором ресторанов.
func (h *Handler) restaurants(w http.ResponseWriter, r *http.Request) {
	restaurants, err := h.service.RestaurantService.GetAll()
	if err != nil {
		_ = render.Render(w, r, ErrServiceFailure(err))
		return
	}

	renderTemplate(w, r,
		"restaurants",
		&SiteDetails{
			Metadata: Metadata{
				Title: "Выбор ресторана",
			},
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
