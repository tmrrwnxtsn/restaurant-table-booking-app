package handler

import (
	"context"
	"errors"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"

	"github.com/tmrrwnxtsn/restaurant-table-booking-app/internal/apiserver/model"
	"github.com/tmrrwnxtsn/restaurant-table-booking-app/internal/apiserver/store"
)

const restaurantCtxKey = "restaurant"

// initRestaurantsRouter подготавливает отдельный маршрутизатор для манипуляции ресторанами.
func (h *Handler) initRestaurantsRouter() http.Handler {
	r := chi.NewRouter()
	r.Post("/", h.createRestaurant) // POST /restaurants/
	r.Get("/", h.listRestaurants)   // GET /restaurants/
	r.Route("/{restaurant_id}", func(r chi.Router) {
		r.Use(h.restaurantCtx)            // загрузить информацию о ресторане из контекста запроса
		r.Get("/", h.getRestaurant)       // GET /restaurants/123/
		r.Patch("/", h.updateRestaurant)  // PATCH /restaurants/123/
		r.Delete("/", h.deleteRestaurant) // DELETE /restaurants/123/
		r.Route("/tables", func(r chi.Router) { // работа со столиками ресторанов
			r.Post("/", h.createTable) // POST /restaurants/123/tables
			r.Get("/", h.listTables)   // GET /restaurants/123/tables
		})
		r.Route("/bookings", func(r chi.Router) { // работа со бронями ресторанов
			r.Post("/", h.createBooking) // POST /restaurants/123/bookings
			r.Get("/", h.listBookings)   // GET /restaurants/123/bookings
		})
	})
	return r
}

// createRestaurantRequest представляет тело запроса на создание ресторана.
type createRestaurantRequest struct {
	Name               string  `json:"name" example:"Каравелла"`
	AverageWaitingTime int     `json:"average_waiting_time,string" example:"60"`
	AverageCheck       float64 `json:"average_check,string" example:"2500.00"`
}

// Bind осуществляет пост-обработку запроса.
func (r *createRestaurantRequest) Bind(_ *http.Request) error {
	if r.Name == "" || r.AverageWaitingTime == 0 || r.AverageCheck == 0.0 {
		return ErrRestaurantMissingFields
	}
	return nil
}

// createRestaurantResponse представляет тело ответа на создание ресторана.
type createRestaurantResponse struct {
	ID uint64 `json:"id" example:"1"`
}

// Render осуществляет предобработку ответа.
func (r *createRestaurantResponse) Render(_ http.ResponseWriter, _ *http.Request) error {
	return nil
}

// createRestaurant godoc
// @Summary  Создать ресторан
// @Tags     restaurants
// @Accept   json
// @Produce  json
// @Param    input  body      createRestaurantRequest   true  "Информация о ресторане"
// @Success  201    {object}  createRestaurantResponse  "ok"
// @Failure  400    {object}  errResponse               "Некорректные данные ресторана"
// @Failure  500    {object}  errResponse               "Ошибка на стороне сервера"
// @Router   /restaurants/ [post]
func (h *Handler) createRestaurant(w http.ResponseWriter, r *http.Request) {
	data := &createRestaurantRequest{}
	if err := render.Bind(r, data); err != nil {
		_ = render.Render(w, r, errInvalidRequest(err))
		return
	}

	restaurantID, err := h.service.RestaurantService.Create(data.Name, data.AverageWaitingTime, data.AverageCheck)
	if err != nil {
		_ = render.Render(w, r, errServiceFailure(err))
		return
	}

	render.Status(r, http.StatusCreated)
	_ = render.Render(w, r, &createRestaurantResponse{
		ID: restaurantID,
	})
}

// listRestaurantsResponse представляет тело ответа на получение списка ресторанов.
type listRestaurantsResponse struct {
	Data []model.Restaurant `json:"data"`
}

// Render осуществляет предобработку ответа.
func (r *listRestaurantsResponse) Render(_ http.ResponseWriter, _ *http.Request) error {
	return nil
}

// listRestaurants godoc
// @Summary  Получить список всех ресторанов
// @Tags     restaurants
// @Accept   json
// @Produce  json
// @Success  200  {object}  listRestaurantsResponse  "ok"
// @Failure  500  {object}  errResponse              "Ошибка на стороне сервера"
// @Router   /restaurants/ [get]
func (h *Handler) listRestaurants(w http.ResponseWriter, r *http.Request) {
	restaurants, err := h.service.RestaurantService.GetAll()
	if err != nil {
		_ = render.Render(w, r, errServiceFailure(err))
		return
	}

	_ = render.Render(w, r, &listRestaurantsResponse{
		Data: restaurants,
	})
}

// restaurantCtx используется для загрузки ресторана (model.Restaurant) из контекста запроса по restaurant_id,
// переданному в параметрах URL запроса.
func (h *Handler) restaurantCtx(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if restaurantIDStr := chi.URLParam(r, "restaurant_id"); restaurantIDStr != "" {
			restaurantID, err := strconv.ParseUint(restaurantIDStr, 10, 0)
			if err != nil {
				_ = render.Render(w, r, errInvalidRequest(err))
				return
			}

			restaurant, err := h.service.RestaurantService.Get(restaurantID)
			if err != nil {
				if errors.Is(err, store.ErrRestaurantNotFound) {
					_ = render.Render(w, r, errNotFound(err))
					return
				}
				_ = render.Render(w, r, errServiceFailure(err))
				return
			}

			ctx := context.WithValue(r.Context(), restaurantCtxKey, restaurant)
			next.ServeHTTP(w, r.WithContext(ctx))
		} else {
			_ = render.Render(w, r, errInvalidRequest(ErrRestaurantMissingFields))
			return
		}
	})
}

// getRestaurantResponse представляет тело ответа на получение ресторана.
type getRestaurantResponse struct {
	*model.Restaurant
}

// Render осуществляет предобработку ответа.
func (r *getRestaurantResponse) Render(_ http.ResponseWriter, _ *http.Request) error {
	return nil
}

// getRestaurant godoc
// @Summary  Получить ресторан по его ID
// @Tags     restaurants
// @Accept   json
// @Produce  json
// @Param    restaurant_id  path      string                 true  "ID ресторана"
// @Success  200            {object}  getRestaurantResponse  "ok"
// @Failure  400            {object}  errResponse            "Некорректный ID ресторана"
// @Failure  500            {object}  errResponse            "Ошибка на стороне сервера"
// @Router   /restaurants/{restaurant_id}/ [get]
func (h *Handler) getRestaurant(w http.ResponseWriter, r *http.Request) {
	restaurant := r.Context().Value(restaurantCtxKey).(*model.Restaurant)

	if err := render.Render(w, r, &getRestaurantResponse{restaurant}); err != nil {
		_ = render.Render(w, r, errRender(err))
		return
	}
}

// updateRestaurantResponse представляет тело ответа на получение ресторана.
type updateRestaurantResponse struct {
	Status string `json:"status" example:"ok"`
}

// Render осуществляет предобработку ответа.
func (r *updateRestaurantResponse) Render(_ http.ResponseWriter, _ *http.Request) error {
	return nil
}

// updateRestaurant godoc
// @Summary      Обновить информацию о ресторане по его ID
// @Description  Обновление происходит именно путём PATCH-запросов, чтобы была возможность изменять данные частично.
// @Tags     	 restaurants
// @Accept   	 json
// @Produce  	 json
// @Param        restaurant_id  path      string                      true  "ID ресторана"
// @Param        input          body      model.UpdateRestaurantData  true  "Информация о ресторане"
// @Success      200            {object}  updateRestaurantResponse    "ok"
// @Failure      400            {object}  errResponse                 "Некорректный данные запроса"
// @Failure      500            {object}  errResponse                 "Ошибка на стороне сервера"
// @Router       /restaurants/{restaurant_id}/ [patch]
func (h *Handler) updateRestaurant(w http.ResponseWriter, r *http.Request) {
	restaurant := r.Context().Value(restaurantCtxKey).(*model.Restaurant)

	data := model.UpdateRestaurantData{}
	if err := render.Bind(r, &data); err != nil {
		_ = render.Render(w, r, errInvalidRequest(err))
		return
	}

	if err := h.service.RestaurantService.Update(restaurant.ID, data); err != nil {
		_ = render.Render(w, r, errServiceFailure(err))
		return
	}

	_ = render.Render(w, r, &updateRestaurantResponse{Status: "ok"})
}

// deleteRestaurantResponse представляет тело ответа на удаление ресторана.
type deleteRestaurantResponse struct {
	Status string `json:"status" example:"ok"`
}

// Render осуществляет предобработку ответа.
func (r *deleteRestaurantResponse) Render(_ http.ResponseWriter, _ *http.Request) error {
	return nil
}

// deleteRestaurant godoc
// @Summary      Удаляет ресторан по его ID
// @Description  Если ресторан ожидает гостей (есть брони в будущем/сегодняшнем днях в этом ресторане), то его нельзя удалить.
// @Tags     	 restaurants
// @Accept   	 json
// @Produce  	 json
// @Param        restaurant_id  path      string                    true  "ID ресторана"
// @Success      200            {object}  deleteRestaurantResponse  "ok"
// @Failure      400            {object}  errResponse               "Некорректный данные запроса"
// @Failure      500            {object}  errResponse               "Ошибка на стороне сервера"
// @Router       /restaurants/{restaurant_id}/ [delete]
func (h *Handler) deleteRestaurant(w http.ResponseWriter, r *http.Request) {
	restaurant := r.Context().Value(restaurantCtxKey).(*model.Restaurant)

	if err := h.service.RestaurantService.Delete(restaurant.ID); err != nil {
		_ = render.Render(w, r, errServiceFailure(err))
		return
	}

	_ = render.Render(w, r, &deleteRestaurantResponse{Status: "ok"})
}
