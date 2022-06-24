package handler

import (
	"net/http"
	"net/http/pprof"

	"github.com/go-chi/chi/v5"
)

// initProfilerRouter подготавливает отдельный маршрутизатор для профилирования работы API сервера.
func (h *Handler) initProfilerRouter() http.Handler {
	r := chi.NewRouter()

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, r.RequestURI+"/pprof/", http.StatusMovedPermanently)
	})
	r.HandleFunc("/pprof", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, r.RequestURI+"/", http.StatusMovedPermanently)
	})

	// Получение списка всех профилей
	r.HandleFunc("/pprof/*", pprof.Index)
	// Отображение строки запуска (например: /go-observability-course/examples/caching/redis/__debug_bin)
	r.HandleFunc("/pprof/cmdline", pprof.Cmdline)
	// профиль ЦПУ, в query-параметрах можно указать seconds со значением времени в секундах для снимка (по-умолчанию 30с)
	r.HandleFunc("/pprof/profile", pprof.Profile)
	r.HandleFunc("/pprof/symbol", pprof.Symbol)
	// профиль для получения трассировки (последовательности инструкций) выполнения приложения за время seconds из query-параметров (по-умолчанию 1с)
	r.HandleFunc("/pprof/trace", pprof.Trace)

	return r
}
