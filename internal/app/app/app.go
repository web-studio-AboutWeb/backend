package app

import (
	"context"
	"os"
	"os/signal"
	"syscall"
	"time"
	"web-studio-backend/internal/app"
	"web-studio-backend/internal/app/core"
	"web-studio-backend/internal/app/infrastructure/config"
	"web-studio-backend/internal/app/infrastructure/logger"
	"web-studio-backend/internal/app/interfaces/http"
)

type application struct {
	core *core.Core
	config *config.Config
}

func New(configPath string) (app.App, error) {
	config, err := config.Init(configPath)
	if err != nil {
		return nil, err
	}

	loggerConfig := config.Logger

	logger.Init(&logger.Config{
		LogToConsole:     loggerConfig.LogToConsole,
		EncodeLogsAsJson: loggerConfig.EncodeLogsAsJson,
		LogToFile:        loggerConfig.LogToFile,
		Directory:        loggerConfig.Directory,
		Filename:         loggerConfig.Filename,
		MaxSize:          loggerConfig.MaxSize,
		MaxBackups:       loggerConfig.MaxBackups,
		MaxAge:           loggerConfig.MaxAge,
	})

	c, err := core.New(context.Background(), config)
	if err != nil {
		return nil, err
	}

	return &application{
		core: c,
		config: config,
	}, nil
}

func (app *application) Start() error {
	httpServer := http.NewHttpServer(app.core, app.config)

	go func() {
		err := httpServer.Run()
		if err != nil {
			logger.Logger.Error().Err(err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signals := []os.Signal{syscall.SIGABRT, syscall.SIGQUIT, syscall.SIGHUP, os.Interrupt, syscall.SIGTERM}
	signal.Notify(quit, signals...)

	<-quit
	logger.Logger.Info().Msg("Shutting down the server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := httpServer.Shutdown(ctx); err != nil {
		logger.Logger.Error().Msgf("server shutdown failed: %v", err)
	}

	logger.Logger.Info().Msg("Server has been shut down")

	return nil
}

func (app *application) Stop(ctx context.Context) error {
	return nil
}
