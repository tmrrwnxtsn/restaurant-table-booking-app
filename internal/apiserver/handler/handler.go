package handler

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/sirupsen/logrus"
	"github.com/swaggo/http-swagger"
	_ "github.com/tmrrwnxtsn/aero-table-booking-api/docs"
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
		r.Mount("/restaurants", h.initRestaurantsRouter())
		// маршруты для манипуляции столиками ресторанов
		r.Mount("/tables", h.initTablesRouter())
	})

	// swagger-документация
	r.Get("/swagger/*", httpSwagger.WrapHandler)

	// профилирование
	r.Mount("/debug", h.initProfilerRouter())

	return r
}
