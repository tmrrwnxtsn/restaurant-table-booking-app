package handler

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/tmrrwnxtsn/aero-table-booking-api/internal/service"
	"net/http"
	"time"
)

// Handler представляет маршрутизатор.
type Handler struct {
	service *service.Services
}

func NewHandler(services *service.Services) *Handler {
	return &Handler{service: services}
}

// InitRoutes инициализирует маршруты.
func (h *Handler) InitRoutes() *chi.Mux {
	router := chi.NewRouter()

	// middleware
	router.Use(middleware.RequestID)
	router.Use(middleware.RealIP)
	router.Use(middleware.Logger)
	router.Use(middleware.Recoverer)

	// установка таймаута на обработку запроса
	router.Use(middleware.Timeout(60 * time.Second))

	router.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello!"))
	})

	return router
}