package handler

import (
	"context"
	"errors"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"github.com/tmrrwnxtsn/aero-table-booking-api/internal/apiserver/model"
	"github.com/tmrrwnxtsn/aero-table-booking-api/internal/apiserver/store"
	"net/http"
	"strconv"
)

const tableCtxKey = "table"

// CreateTableRequest представляет тело запроса на создание столика в ресторане.
type CreateTableRequest struct {
	SeatsNumber int `json:"seats_number"`
}

// Bind осуществляет пост-обработку запроса CreateTableRequest.
func (r *CreateTableRequest) Bind(_ *http.Request) error {
	if r.SeatsNumber == 0 {
		return ErrTableMissingFields
	}
	return nil
}

// CreateTableResponse представляет тело ответа на создание столика в ресторане.
type CreateTableResponse struct {
	ID uint64 `json:"id"`
}

// Render осуществляет предобработку ответа CreateTableResponse.
func (r *CreateTableResponse) Render(_ http.ResponseWriter, _ *http.Request) error {
	return nil
}

// createTable принимает запросы на создание столика в ресторане.
func (h *Handler) createTable(w http.ResponseWriter, r *http.Request) {
	restaurant := r.Context().Value(restaurantCtxKey).(*model.Restaurant)

	data := &CreateTableRequest{}
	if err := render.Bind(r, data); err != nil {
		_ = render.Render(w, r, ErrInvalidRequest(err))
		return
	}

	tableID, err := h.service.TableService.Create(restaurant.ID, data.SeatsNumber)
	if err != nil {
		_ = render.Render(w, r, ErrServiceFailure(err))
		return
	}

	render.Status(r, http.StatusCreated)
	_ = render.Render(w, r, &CreateTableResponse{
		ID: tableID,
	})
}

// ListTablesResponse представляет тело ответа на получение списка столиков в ресторане.
type ListTablesResponse struct {
	Data []model.Table `json:"data"`
}

// Render осуществляет предобработку ответа ListTablesResponse.
func (r *ListTablesResponse) Render(_ http.ResponseWriter, _ *http.Request) error {
	return nil
}

// listTables принимает запросы на получение списка столиков в ресторане.
func (h *Handler) listTables(w http.ResponseWriter, r *http.Request) {
	restaurant := r.Context().Value(restaurantCtxKey).(*model.Restaurant)

	tables, err := h.service.TableService.GetAll(restaurant.ID)
	if err != nil {
		_ = render.Render(w, r, ErrServiceFailure(err))
		return
	}

	_ = render.Render(w, r, &ListTablesResponse{
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
				_ = render.Render(w, r, ErrInvalidRequest(err))
				return
			}

			table, err := h.service.TableService.Get(tableID)
			if err != nil {
				if errors.Is(err, store.ErrTableNotFound) {
					_ = render.Render(w, r, ErrNotFound(err))
					return
				}
				_ = render.Render(w, r, ErrServiceFailure(err))
				return
			}

			ctx := context.WithValue(r.Context(), tableCtxKey, table)
			next.ServeHTTP(w, r.WithContext(ctx))
		} else {
			_ = render.Render(w, r, ErrInvalidRequest(ErrTableMissingFields))
			return
		}
	})
}

// GetTableResponse представляет тело ответа на получение столика в ресторане.
type GetTableResponse struct {
	*model.Table
}

// Render осуществляет предобработку ответа GetRestaurantResponse.
func (r *GetTableResponse) Render(_ http.ResponseWriter, _ *http.Request) error {
	return nil
}

// getRestaurant принимает запросы на получение столика в ресторане.
func (h *Handler) getTable(w http.ResponseWriter, r *http.Request) {
	table := r.Context().Value(tableCtxKey).(*model.Table)

	if err := render.Render(w, r, &GetTableResponse{table}); err != nil {
		_ = render.Render(w, r, ErrRender(err))
		return
	}
}

// UpdateTableResponse представляет тело ответа на получение ресторана.
type UpdateTableResponse struct {
	Status string `json:"status"`
}

// Render осуществляет предобработку ответа UpdateTableResponse.
func (r *UpdateTableResponse) Render(_ http.ResponseWriter, _ *http.Request) error {
	return nil
}

// updateRestaurant принимает запросы на обновление информации о столике в ресторане.
func (h *Handler) updateTable(w http.ResponseWriter, r *http.Request) {
	table := r.Context().Value(tableCtxKey).(*model.Table)

	data := model.UpdateTableData{}
	if err := render.Bind(r, &data); err != nil {
		_ = render.Render(w, r, ErrInvalidRequest(err))
		return
	}

	if err := h.service.TableService.Update(table.ID, data); err != nil {
		_ = render.Render(w, r, ErrServiceFailure(err))
		return
	}

	_ = render.Render(w, r, &UpdateTableResponse{Status: "ok"})
}

// DeleteTableResponse представляет тело ответа на удаление столика ресторана.
type DeleteTableResponse struct {
	Status string `json:"status"`
}

// Render осуществляет предобработку ответа DeleteTableResponse.
func (r *DeleteTableResponse) Render(_ http.ResponseWriter, _ *http.Request) error {
	return nil
}

// deleteRestaurant принимает запросы на удаление столика ресторана.
func (h *Handler) deleteTable(w http.ResponseWriter, r *http.Request) {
	table := r.Context().Value(tableCtxKey).(*model.Table)

	if err := h.service.TableService.Delete(table.ID); err != nil {
		_ = render.Render(w, r, ErrServiceFailure(err))
		return
	}

	_ = render.Render(w, r, &DeleteRestaurantResponse{Status: "ok"})
}
