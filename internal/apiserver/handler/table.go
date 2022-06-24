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

const tableCtxKey = "table"

// initTablesRouter подготавливает отдельный маршрутизатор для манипуляции столиками в ресторанах.
func (h *Handler) initTablesRouter() http.Handler {
	r := chi.NewRouter()
	r.Route("/{table_id}", func(r chi.Router) {
		r.Use(h.tableCtx)            // загрузить информацию о столике из контекста запроса
		r.Get("/", h.getTable)       // GET /tables/123/
		r.Patch("/", h.updateTable)  // PATCH /tables/123/
		r.Delete("/", h.deleteTable) // DELETE /tables/123/
	})
	return r
}

// createTableRequest представляет тело запроса на создание столика в ресторане.
type createTableRequest struct {
	SeatsNumber int `json:"seats_number" example:"3"`
}

// Bind осуществляет пост-обработку запроса.
func (r *createTableRequest) Bind(_ *http.Request) error {
	if r.SeatsNumber == 0 {
		return ErrTableMissingFields
	}
	return nil
}

// createTableResponse представляет тело ответа на создание столика в ресторане.
type createTableResponse struct {
	ID uint64 `json:"id" example:"2"`
}

// Render осуществляет предобработку ответа.
func (r *createTableResponse) Render(_ http.ResponseWriter, _ *http.Request) error {
	return nil
}

// createTable godoc
// @Summary  Создать столик в ресторане
// @Tags     tables
// @Accept   json
// @Produce  json
// @Param    restaurant_id  path      string               true  "ID ресторана"
// @Param    input          body      createTableRequest   true  "Информация о столике"
// @Success  201            {object}  createTableResponse  "ok"
// @Failure  400            {object}  errResponse          "Некорректные данные столика"
// @Failure  500            {object}  errResponse          "Ошибка на стороне сервера"
// @Router   /restaurants/{restaurant_id}/tables/ [post]
func (h *Handler) createTable(w http.ResponseWriter, r *http.Request) {
	restaurant := r.Context().Value(restaurantCtxKey).(*model.Restaurant)

	data := &createTableRequest{}
	if err := render.Bind(r, data); err != nil {
		_ = render.Render(w, r, errInvalidRequest(err))
		return
	}

	tableID, err := h.service.TableService.Create(restaurant.ID, data.SeatsNumber)
	if err != nil {
		_ = render.Render(w, r, errServiceFailure(err))
		return
	}

	render.Status(r, http.StatusCreated)
	_ = render.Render(w, r, &createTableResponse{
		ID: tableID,
	})
}

// listTablesResponse представляет тело ответа на получение списка столиков в ресторане.
type listTablesResponse struct {
	Data []model.Table `json:"data"`
}

// Render осуществляет предобработку ответа.
func (r *listTablesResponse) Render(_ http.ResponseWriter, _ *http.Request) error {
	return nil
}

// listTables godoc
// @Summary  Получить список столиков в ресторане
// @Tags     tables
// @Accept   json
// @Produce  json
// @Param    restaurant_id  path      string              true  "ID ресторана"
// @Success  200            {object}  listTablesResponse  "ok"
// @Failure  400            {object}  errResponse         "Некорректный restaurant_id"
// @Failure  500            {object}  errResponse         "Ошибка на стороне сервера"
// @Router   /restaurants/{restaurant_id}/tables/ [get]
func (h *Handler) listTables(w http.ResponseWriter, r *http.Request) {
	restaurant := r.Context().Value(restaurantCtxKey).(*model.Restaurant)

	tables, err := h.service.TableService.GetAll(restaurant.ID)
	if err != nil {
		_ = render.Render(w, r, errServiceFailure(err))
		return
	}

	_ = render.Render(w, r, &listTablesResponse{
		Data: tables,
	})
}

// tableCtx используется для загрузки столика (model.Table) из контекста запроса по table_id,
// переданному в параметрах URL запроса.
func (h *Handler) tableCtx(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if tableIDStr := chi.URLParam(r, "table_id"); tableIDStr != "" {
			tableID, err := strconv.ParseUint(tableIDStr, 10, 0)
			if err != nil {
				_ = render.Render(w, r, errInvalidRequest(err))
				return
			}

			table, err := h.service.TableService.Get(tableID)
			if err != nil {
				if errors.Is(err, store.ErrTableNotFound) {
					_ = render.Render(w, r, errNotFound(err))
					return
				}
				_ = render.Render(w, r, errServiceFailure(err))
				return
			}

			ctx := context.WithValue(r.Context(), tableCtxKey, table)
			next.ServeHTTP(w, r.WithContext(ctx))
		} else {
			_ = render.Render(w, r, errInvalidRequest(ErrTableMissingFields))
			return
		}
	})
}

// getTableResponse представляет тело ответа на получение столика в ресторане.
type getTableResponse struct {
	*model.Table
}

// Render осуществляет предобработку ответа.
func (r *getTableResponse) Render(_ http.ResponseWriter, _ *http.Request) error {
	return nil
}

// getTable godoc
// @Summary  Получить столик по его ID
// @Tags     tables
// @Accept   json
// @Produce  json
// @Param    table_id  path      string            true  "ID столика"
// @Success  200       {object}  getTableResponse  "ok"
// @Failure  400       {object}  errResponse       "Некорректный ID столика"
// @Failure  500       {object}  errResponse       "Ошибка на стороне сервера"
// @Router   /tables/{table_id}/ [get]
func (h *Handler) getTable(w http.ResponseWriter, r *http.Request) {
	table := r.Context().Value(tableCtxKey).(*model.Table)

	if err := render.Render(w, r, &getTableResponse{table}); err != nil {
		_ = render.Render(w, r, errRender(err))
		return
	}
}

// updateTableResponse представляет тело ответа на получение ресторана.
type updateTableResponse struct {
	Status string `json:"status" example:"ok"`
}

// Render осуществляет предобработку ответа.
func (r *updateTableResponse) Render(_ http.ResponseWriter, _ *http.Request) error {
	return nil
}

// updateTable godoc
// @Summary      Обновить информацию о столике по его ID
// @Description  Обновление происходит именно путём PATCH-запросов, чтобы была возможность изменять данные частично.
// @Tags         tables
// @Accept   	 json
// @Produce  	 json
// @Param        table_id  path      string                 true  "ID столика"
// @Param        input     body      model.UpdateTableData  true  "Информация о столике"
// @Success      200       {object}  updateTableResponse    "ok"
// @Failure      400       {object}  errResponse            "Некорректный данные запроса"
// @Failure      500       {object}  errResponse            "Ошибка на стороне сервера"
// @Router       /tables/{table_id}/ [patch]
func (h *Handler) updateTable(w http.ResponseWriter, r *http.Request) {
	table := r.Context().Value(tableCtxKey).(*model.Table)

	data := model.UpdateTableData{}
	if err := render.Bind(r, &data); err != nil {
		_ = render.Render(w, r, errInvalidRequest(err))
		return
	}

	if err := h.service.TableService.Update(table.ID, data); err != nil {
		_ = render.Render(w, r, errServiceFailure(err))
		return
	}

	_ = render.Render(w, r, &updateTableResponse{Status: "ok"})
}

// deleteTableResponse представляет тело ответа на удаление столика ресторана.
type deleteTableResponse struct {
	Status string `json:"status"`
}

// Render осуществляет предобработку ответа.
func (r *deleteTableResponse) Render(_ http.ResponseWriter, _ *http.Request) error {
	return nil
}

// deleteRestaurant godoc
// @Summary      Удаляет столик по его ID
// @Description  Если за столиком ещё будут сидеть клиенты (есть брони в будущем/сегодняшнем днях в этом ресторане), то его нельзя удалить.
// @Tags     	 tables
// @Accept   	 json
// @Produce  	 json
// @Param        table_id  path      string               true  "ID столика"
// @Success      200       {object}  deleteTableResponse  "ok"
// @Failure      400       {object}  errResponse          "Некорректный данные запроса"
// @Failure      500       {object}  errResponse          "Ошибка на стороне сервера"
// @Router       /tables/{table_id}/ [delete]
func (h *Handler) deleteTable(w http.ResponseWriter, r *http.Request) {
	table := r.Context().Value(tableCtxKey).(*model.Table)

	if err := h.service.TableService.Delete(table.ID); err != nil {
		_ = render.Render(w, r, errServiceFailure(err))
		return
	}

	_ = render.Render(w, r, &deleteTableResponse{Status: "ok"})
}
