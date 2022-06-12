package handler

import (
	"context"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"github.com/tmrrwnxtsn/aero-table-booking-api/internal/model"
	"net/http"
	"strconv"
)

const (
	limitValueKey  = "limit"
	offsetValueKey = "offset"
)

type ListRestaurantsResponse struct {
	Data []model.Restaurant `json:"data"`
}

func (r *ListRestaurantsResponse) Render(_ http.ResponseWriter, _ *http.Request) error {
	return nil
}

// listRestaurants принимает запросы на отображение списка ресторанов.
func (h *Handler) listRestaurants(w http.ResponseWriter, r *http.Request) {
	restaurants, err := h.service.RestaurantService.GetAll()
	if err != nil {
		_ = render.Render(w, r, ErrServiceFailure(err))
		return
	}

	render.Status(r, http.StatusOK)
	_ = render.Render(w, r, &ListRestaurantsResponse{
		Data: restaurants,
	})
}

// paginate является middleware-компонентом для парсинга параметров (limit, offset) запроса к listRestaurants.
func (h *Handler) paginate(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var limit, offset uint64 = 5, 0

		if limitStr := chi.URLParam(r, limitValueKey); limitStr != "" {
			parseUint, err := strconv.ParseUint(limitStr, 10, 0)
			if err == nil {
				limit = parseUint
			}
		}

		if offsetStr := chi.URLParam(r, offsetValueKey); offsetStr != "" {
			parseUint, err := strconv.ParseUint(offsetStr, 10, 0)
			if err != nil {
				offset = parseUint
			}
		}

		ctx := context.WithValue(r.Context(), limitValueKey, limit)
		ctx = context.WithValue(ctx, offsetValueKey, offset)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// CreateRestaurantRequest представляет тело запроса на создание ресторана.
type CreateRestaurantRequest struct {
	Name               string  `json:"name"`
	AverageWaitingTime int     `json:"average_waiting_time,string"`
	AverageCheck       float64 `json:"average_check,string"`
}

// Bind осуществляет пост-обработку запроса.
func (r *CreateRestaurantRequest) Bind(_ *http.Request) error {
	if r.Name == "" || r.AverageWaitingTime == 0 || r.AverageCheck == 0.0 {
		return ErrRestaurantMissingFields
	}
	return nil
}

// CreateRestaurantResponse представляет тело ответа на создание ресторана.
type CreateRestaurantResponse struct {
	ID uint64 `json:"id"`
}

// Render осуществляет предобработку ответа.
func (r *CreateRestaurantResponse) Render(_ http.ResponseWriter, _ *http.Request) error {
	return nil
}

// createRestaurant принимает запросы на создание ресторана.
func (h *Handler) createRestaurant(w http.ResponseWriter, r *http.Request) {
	data := &CreateRestaurantRequest{}
	if err := render.Bind(r, data); err != nil {
		_ = render.Render(w, r, ErrInvalidRequest(err))
		return
	}

	restaurantID, err := h.service.RestaurantService.Create(data.Name, data.AverageWaitingTime, data.AverageCheck)
	if err != nil {
		_ = render.Render(w, r, ErrServiceFailure(err))
		return
	}

	render.Status(r, http.StatusCreated)
	_ = render.Render(w, r, &CreateRestaurantResponse{
		ID: restaurantID,
	})
}
