package http

import (
	"context"
	"log/slog"
	"net/http"
	"time"

	"github.com/TTekmii/todo-list-app/internal/lib/logger/sl"
)

type Server struct {
	httpServer *http.Server
	logger     *slog.Logger
}

func NewServer(port string, handler http.Handler, logger *slog.Logger) *Server {
	return &Server{
		httpServer: &http.Server{
			Addr:           ":" + port,
			Handler:        handler,
			MaxHeaderBytes: 1 << 20, // 1 MB
			ReadTimeout:    15 * time.Second,
			WriteTimeout:   15 * time.Second,
			IdleTimeout:    60 * time.Second,
		},
		logger: logger,
	}
}

func (s *Server) Run() error {
	s.logger.Info("Starting HTTP server", slog.String("addr", s.httpServer.Addr))

	err := s.httpServer.ListenAndServe()
	if err != nil && err != http.ErrServerClosed {
		s.logger.Error("HTTP server listen error", sl.Err(err))
		return err
	}
	return nil
}

func (s *Server) Shutdown(ctx context.Context) error {
	s.logger.Info("Shutting down HTTP server...")
	return s.httpServer.Shutdown(ctx)
}
