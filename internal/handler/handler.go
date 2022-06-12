package handler

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/sirupsen/logrus"
	"github.com/tmrrwnxtsn/aero-table-booking-api/internal/service"
	"time"
)

// Handler представляет маршрутизатор.
type Handler struct {
	service *service.Services
	logger  *logrus.Logger
}

func NewHandler(services *service.Services, logger *logrus.Logger) *Handler {
	return &Handler{
		service: services,
		logger:  logger,
	}
}

// InitRoutes инициализирует маршруты.
func (h *Handler) InitRoutes() *chi.Mux {
	r := chi.NewRouter()

	// middleware
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(NewStructuredLogger(h.logger))
	r.Use(middleware.Recoverer)

	// установка таймаута на обработку запроса
	r.Use(middleware.Timeout(60 * time.Second))

	r.Get("/", h.home)
	r.Route("/restaurants", func(r chi.Router) {
		r.Get("/", h.restaurants)
	})

	r.Route("/api/v1", func(r chi.Router) {
		// маршруты для манипуляции ресторанами
		r.Route("/restaurants", func(r chi.Router) {
			r.Get("/", h.listRestaurants)   // GET /api/v1/restaurants
			r.Post("/", h.createRestaurant) // POST /api/v1/restaurants

			r.Route("/{restaurant_id}", func(r chi.Router) {
				r.Use(h.restaurantCtx)      // загрузить *model.Restaurant из контекста запроса
				r.Get("/", h.getRestaurant) // GET /api/v1/restaurants/123
				//r.Put("/", UpdateArticle)    // PUT /api/v1/restaurants/123
				//r.Delete("/", DeleteArticle) // DELETE /api/v1/restaurants/123
			})
		})
	})

	return r
}
