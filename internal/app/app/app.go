package app

import (
	"context"
	"fmt"
	"log/slog"
	stdhttp "net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/golang-migrate/migrate/v4"

	"web-studio-backend/internal/app/handler/http"
	"web-studio-backend/internal/app/infrastructure/repository/postgresql"
	"web-studio-backend/internal/app/service"
	"web-studio-backend/internal/pkg/config"
	"web-studio-backend/internal/pkg/wcrypto"
	"web-studio-backend/pkg/postgres"

	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

func Run(configPath string) error {
	config.Read(configPath)

	cfg := config.Get()

	// Logger initialization
	logLevel := slog.LevelDebug
	if cfg.App.Env == "prod" {
		logLevel = slog.LevelInfo
	}

	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		AddSource: false,
		Level:     logLevel,
	}))
	slog.SetDefault(logger)

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	// Database stuff
	user, password, err := wcrypto.DecodeUserPass(cfg.Database.User, cfg.Database.Password, config.Block)
	if err != nil {
		return fmt.Errorf("decoding database username: %w", err)
	}

	dbConnString := postgres.ConnectionString(user, password, cfg.Database.Host, cfg.Database.Database)

	pg, err := postgres.New(ctx, dbConnString)
	if err != nil {
		return fmt.Errorf("creating postgres: %w", err)
	}

	slog.Info("Connected to database")

	err = applyMigrations(dbConnString)
	if err != nil {
		return fmt.Errorf("applying migrations: %w", err)
	}

	// Repositories initialization
	userRepo := postgresql.NewUserRepository(pg.Pool)
	projectRepo := postgresql.NewProjectRepository(pg.Pool)

	// Services initialization
	userService := service.NewUserService(userRepo)
	projectService := service.NewProjectService(projectRepo)

	// Handler initialization
	handler := http.NewHandler(
		userService,
		projectService,
	)

	httpServer := &stdhttp.Server{
		Addr:    fmt.Sprintf("%s:%d", cfg.Http.Host, cfg.Http.Port),
		Handler: handler,
	}

	slog.Info("Server is started", slog.String("addr", httpServer.Addr))

	httpServerCh := make(chan error)
	if cfg.Http.HttpsEnabled {
		httpServerCh <- httpServer.ListenAndServeTLS(cfg.Http.CertFilePath, cfg.Http.KeyFilePath)
	} else {
		httpServerCh <- httpServer.ListenAndServe()
	}

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, syscall.SIGTERM)

	select {
	case s := <-interrupt:
		slog.Info("Interrupt signal: " + s.String())
	case err = <-httpServerCh:
		slog.Error("Server stop signal: " + err.Error())
	}

	// Shutdown
	err = httpServer.Shutdown(ctx)
	if err != nil {
		slog.Error("failed to shutdown the server", err)
	}
	slog.Info("Server has been shut down successfully")

	pg.Close()

	return nil
}

// applyMigrations applies migrations to database.
func applyMigrations(connString string) error {
	m, err := migrate.New("file://migrations", connString)
	if err != nil {
		return fmt.Errorf("creating migration: %w", err)
	}

	if err = m.Up(); err != nil {
		return fmt.Errorf("applying up: %w", err)
	}

	return nil
}
