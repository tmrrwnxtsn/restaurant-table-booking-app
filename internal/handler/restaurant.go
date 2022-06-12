package handler

import (
	"context"
	"errors"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"github.com/tmrrwnxtsn/aero-table-booking-api/internal/model"
	"github.com/tmrrwnxtsn/aero-table-booking-api/internal/store"
	"net/http"
	"strconv"
)

const (
	limitValueKey  = "limit"
	offsetValueKey = "offset"
)

const (
	restaurantCtxKey = "restaurant"
)

// ListRestaurantsResponse представляет тело ответа на получение списка всех ресторанов.
type ListRestaurantsResponse struct {
	Data []model.Restaurant `json:"data"`
}

// Render осуществляет предобработку ответа ListRestaurantsResponse.
func (r *ListRestaurantsResponse) Render(_ http.ResponseWriter, _ *http.Request) error {
	return nil
}

// listRestaurants принимает запросы на получение списка ресторанов.
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

// Bind осуществляет пост-обработку запроса CreateRestaurantRequest.
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

// Render осуществляет предобработку ответа CreateRestaurantResponse.
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

// restaurantCtx используется для загрузки ресторана (model.Restaurant) по restaurantID,
// переданному в параметрах URL запроса.
func (h *Handler) restaurantCtx(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if restaurantIDStr := chi.URLParam(r, "restaurant_id"); restaurantIDStr != "" {
			restaurantID, err := strconv.ParseUint(restaurantIDStr, 10, 0)
			if err != nil {
				_ = render.Render(w, r, ErrInvalidRequest(err))
				return
			}

			restaurant, err := h.service.RestaurantService.GetByID(restaurantID)
			if err != nil {
				if errors.Is(err, store.ErrRestaurantNotFound) {
					_ = render.Render(w, r, ErrNotFound(err))
					return
				}
				_ = render.Render(w, r, ErrServiceFailure(err))
				return
			}

			ctx := context.WithValue(r.Context(), restaurantCtxKey, restaurant)
			next.ServeHTTP(w, r.WithContext(ctx))
		} else {
			_ = render.Render(w, r, ErrInvalidRequest(ErrRestaurantMissingFields))
			return
		}
	})
}

// GetRestaurantResponse представляет тело ответа на получение ресторана.
type GetRestaurantResponse struct {
	*model.Restaurant
}

// Render осуществляет предобработку ответа GetRestaurantResponse.
func (r *GetRestaurantResponse) Render(_ http.ResponseWriter, _ *http.Request) error {
	return nil
}

// getRestaurant принимает запросы на получение ресторана.
func (h *Handler) getRestaurant(w http.ResponseWriter, r *http.Request) {
	restaurant := r.Context().Value(restaurantCtxKey).(*model.Restaurant)

	if err := render.Render(w, r, &GetRestaurantResponse{restaurant}); err != nil {
		_ = render.Render(w, r, ErrRender(err))
		return
	}
}
