package http

import (
	"context"
	"errors"
	"fmt"
	"github.com/go-chi/chi/v5"
	"net/http"
	"os"
	"web-studio-backend/internal/app/core"
	"web-studio-backend/internal/app/infrastructure/config"
	"web-studio-backend/internal/app/infrastructure/logger"
)

type Server interface {
	Run() error
	Shutdown(ctx context.Context) error
}

type server struct {
	core                    core.Core
	router                  chi.Router
	httpServer, httpsServer *http.Server
}

func NewHttpServer(core core.Core) Server {
	httpPort := config.Config.Interfaces.Http.Port
	httpHost := config.Config.Interfaces.Http.Host

	router := chi.NewRouter()

	s := &server{
		core:   core,
		router: router,
		httpServer: &http.Server{
			Addr:    fmt.Sprintf("%s:%d", httpHost, httpPort),
			Handler: router,
		},
	}

	if os.Getenv("USE_HTTPS") != "" {
		httpsPort := config.Config.Interfaces.Https.Port
		httpsHost := config.Config.Interfaces.Https.Host

		s.httpsServer = &http.Server{
			Addr:    fmt.Sprintf("%s:%d", httpsHost, httpsPort),
			Handler: s.router,
		}
	}

	s.initRouter()

	return s
}

func (s *server) Run() error {
	if s.httpsServer != nil {
		go func() {
			httpsConfig := config.Config.Interfaces.Https

			logger.Logger.Info().Msgf("Starting HTTPS server at %s", s.httpsServer.Addr)
			err := s.httpsServer.ListenAndServeTLS(httpsConfig.CertFilePath, httpsConfig.KeyFilePath)
			if err != nil && !errors.Is(err, http.ErrServerClosed) {
				logger.Logger.Err(err)
			}
		}()
	}

	logger.Logger.Info().Msgf("Starting HTTP server at %s", s.httpServer.Addr)
	err := s.httpServer.ListenAndServe()
	if err != nil && !errors.Is(err, http.ErrServerClosed) {
		return err
	}

	return nil
}

func (s *server) Shutdown(ctx context.Context) error {
	return s.httpsServer.Shutdown(ctx)
}
