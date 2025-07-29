package rest

import (
	"context"
	"github.com/folivorra/ziper/internal/usecase"
	"log/slog"
	"net/http"
	"time"

	"github.com/folivorra/ziper/app"
	"github.com/gorilla/mux"
)

type Server struct {
	httpServer *http.Server
	logger     *slog.Logger
}

func NewServer(app *app.App, ts *usecase.TaskService, logger *slog.Logger, port string) *Server {
	r := mux.NewRouter()
	c := NewController(ts, logger)
	c.RegisterRoutes(r)

	srv := &http.Server{
		Addr:    ":" + port,
		Handler: r,
	}

	s := &Server{
		httpServer: srv,
		logger:     logger,
	}

	app.RegisterCleanup(func(ctx context.Context) {
		timeout, cancel := context.WithTimeout(ctx, 5*time.Second)
		defer cancel()
		if err := srv.Shutdown(timeout); err != nil {
			s.logger.Warn("HTTP server shutdown failed with error",
				slog.String("url", srv.Addr),
				slog.String("error", err.Error()),
			)
		}
	})

	return s
}

func (s *Server) Start() error {
	s.logger.Info("starting server",
		slog.String("listen", s.httpServer.Addr),
	)
	err := s.httpServer.ListenAndServe()
	if err != nil && err != http.ErrServerClosed {
		s.logger.Warn("HTTP server stopped with error",
			slog.String("listen", s.httpServer.Addr),
			slog.String("error", err.Error()),
		)
		return err
	}
	return nil
}

func (s *Server) Shutdown(ctx context.Context) error {
	s.logger.Info("shutting down server",
		slog.String("listen", s.httpServer.Addr),
	)
	return s.httpServer.Shutdown(ctx)
}
