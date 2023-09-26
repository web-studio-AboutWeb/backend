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

	"web-studio-backend/internal/app/handler/http"
	"web-studio-backend/internal/app/infrastructure/repository/filesystem"
	"web-studio-backend/internal/app/infrastructure/repository/postgresql"
	"web-studio-backend/internal/app/service"
	"web-studio-backend/internal/pkg/config"
	"web-studio-backend/internal/pkg/wcrypto"
	"web-studio-backend/pkg/postgres"

	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

const filesDir = "files"

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

	runCtx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	// Database stuff
	user, password, err := wcrypto.DecodeUserPass(cfg.Database.User, cfg.Database.Password, config.Block)
	if err != nil {
		return fmt.Errorf("decoding database username: %w", err)
	}

	dbConnString := postgres.ConnectionString(user, password, cfg.Database.Host, cfg.Database.Database)

	pg, err := postgres.New(runCtx, dbConnString)
	if err != nil {
		return fmt.Errorf("creating postgres: %w", err)
	}
	defer pg.Close()

	slog.Info("Connected to database")

	// Repositories initialization
	userRepo := postgresql.NewUserRepository(pg.Pool)
	projectRepo := postgresql.NewProjectRepository(pg.Pool)
	documentRepo := postgresql.NewDocumentRepository(pg.Pool)
	teamRepo := postgresql.NewTeamRepository(pg.Pool)

	// FS storage initialization
	filesFS, err := filesystem.New(filesDir)
	if err != nil {
		return fmt.Errorf("creating images fs: %w", err)
	}

	// Services initialization
	userService := service.NewUserService(userRepo, filesFS)
	projectService := service.NewProjectService(projectRepo, userRepo, teamRepo)
	authService := service.NewAuthService(userRepo)
	documentService := service.NewDocumentService(documentRepo, projectRepo, filesFS)
	teamService := service.NewTeamService(teamRepo, filesFS)

	// Handler initialization
	handler := http.NewHandler(
		userService,
		projectService,
		authService,
		documentService,
		teamService,
	)

	httpServer := &stdhttp.Server{
		Addr:    fmt.Sprintf("%s:%d", cfg.Http.Host, cfg.Http.Port),
		Handler: handler,
	}

	httpServerCh := make(chan error)
	go func() {
		if cfg.Http.HttpsEnabled {
			httpServerCh <- httpServer.ListenAndServeTLS(cfg.Http.CertFilePath, cfg.Http.KeyFilePath)
		} else {
			httpServerCh <- httpServer.ListenAndServe()
		}
	}()

	slog.Info(
		"Server is started",
		slog.String("addr", httpServer.Addr),
		slog.Bool("https", cfg.Http.HttpsEnabled),
	)

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, syscall.SIGTERM)

	select {
	case s := <-interrupt:
		slog.Info("Interrupt signal: " + s.String())
	case err = <-httpServerCh:
		slog.Error("Server stop signal: " + err.Error())
	}

	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer shutdownCancel()

	// Shutdown
	err = httpServer.Shutdown(shutdownCtx)
	if err != nil {
		slog.Error("failed to shutdown the server: " + err.Error())
	}
	slog.Info("Server has been shut down successfully")

	return nil
}
