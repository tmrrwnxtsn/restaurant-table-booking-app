package server

import (
	"context"
	"net/http"
	"time"
)

// Server представляет сервер, обрабатывающий HTTP-запросы.
type Server struct {
	bindAddr   string
	handler    http.Handler
	httpServer *http.Server
}

func NewServer(bindAddr string, handler http.Handler) *Server {
	return &Server{
		bindAddr: bindAddr,
		handler:  handler,
	}
}

func (s *Server) Run() error {
	s.httpServer = &http.Server{
		Addr:           s.bindAddr,
		Handler:        s.handler,
		MaxHeaderBytes: 1 << 20, // MB
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   15 * time.Second,
	}

	return s.httpServer.ListenAndServe()
}

func (s *Server) Shutdown(ctx context.Context) error {
	return s.httpServer.Shutdown(ctx)
}
