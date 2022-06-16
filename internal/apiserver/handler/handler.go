package handler

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/sirupsen/logrus"
	"github.com/tmrrwnxtsn/aero-table-booking-api/internal/apiserver/service"
	"github.com/tmrrwnxtsn/aero-table-booking-api/pkg/logging"
	"net/http"
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
	r.Use(logging.NewStructuredLogger(h.logger))
	r.Use(middleware.Recoverer)

	// установка таймаута на обработку запроса
	r.Use(middleware.Timeout(60 * time.Second))

	// работа системы в визуальном оформлении
	r.Get("/", h.home) // GET / (начальная страница)
	r.Route("/restaurants", func(r chi.Router) {
		r.Get("/", h.restaurants)                                              // GET /restaurants/?people_num=...&desired_datetime=... (страница со всеми доступными ресторанами)
		r.With(h.restaurantCtx).Post("/{restaurant_id}/booked", h.makeBooking) // POST /restaurants/123/booked (забронировать места в ресторане)
	})

	// инициализируем FileServer, который будет обрабатывать HTTP-запросы к статическим файлам из папки "./website".
	fileServer := http.FileServer(http.Dir("./website/"))
	r.Handle("/static/*", http.StripPrefix("/static", fileServer))

	r.Route("/api/v1", func(r chi.Router) {
		// маршруты для манипуляции ресторанами
		r.Route("/restaurants", func(r chi.Router) {
			r.Post("/", h.createRestaurant) // POST /api/v1/restaurants/
			r.Get("/", h.listRestaurants)   // GET /api/v1/restaurants/
			r.Route("/{restaurant_id}", func(r chi.Router) {
				r.Use(h.restaurantCtx)            // загрузить информацию о ресторане из контекста запроса
				r.Get("/", h.getRestaurant)       // GET /api/v1/restaurants/123/
				r.Patch("/", h.updateRestaurant)  // PATCH /api/v1/restaurants/123/
				r.Delete("/", h.deleteRestaurant) // DELETE /api/v1/restaurants/123/
				r.Route("/tables", func(r chi.Router) { // работа со столиками ресторанов
					r.Post("/", h.createTable) // POST /api/v1/restaurants/123/tables
					r.Get("/", h.listTables)   // GET /api/v1/restaurants/123/tables
				})
			})
		})
		// маршруты для манипуляции столиками ресторанов
		r.Route("/tables/{table_id}", func(r chi.Router) {
			r.Use(h.tableCtx)            // загрузить информацию о столике из контекста запроса
			r.Get("/", h.getTable)       // GET /api/v1/tables/123/
			r.Patch("/", h.updateTable)  // PATCH /api/v1/tables/123/
			r.Delete("/", h.deleteTable) // DELETE /api/v1/tables/123/
		})
	})

	return r
}
