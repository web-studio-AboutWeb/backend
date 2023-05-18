package http

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"

	"web-studio-backend/internal/app/core"
	"web-studio-backend/internal/app/infrastructure/config"
	"web-studio-backend/internal/app/infrastructure/logger"
)

type Server interface {
	Run() error
	Shutdown(ctx context.Context) error
}

type server struct {
	core       *core.Core
	router     chi.Router
	httpServer *http.Server
	config     *config.Config
}

func NewHttpServer(core *core.Core, config *config.Config) Server {
	httpPort := config.Interfaces.Http.Port
	httpHost := config.Interfaces.Http.Host

	router := chi.NewRouter()

	s := &server{
		core:   core,
		router: router,
		httpServer: &http.Server{
			Addr:    fmt.Sprintf("%s:%d", httpHost, httpPort),
			Handler: router,
		},
		config: config,
	}

	s.initRouter()

	return s
}

func (s *server) Run() error {
	logger.Logger.Info().Msgf("Starting HTTP server at %s", s.httpServer.Addr)
	err := s.httpServer.ListenAndServe()
	if err != nil && !errors.Is(err, http.ErrServerClosed) {
		logger.Logger.Error().Err(err).Msg("Starting HTTP server")
		return err
	}
	return nil
}

func (s *server) Shutdown(ctx context.Context) error {
	return s.httpServer.Shutdown(ctx)
}
