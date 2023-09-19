package main

import (
	"context"
	"log"

	"web-studio-backend/internal/app/domain"
	"web-studio-backend/internal/app/infrastructure/repository/postgresql"
	"web-studio-backend/internal/app/service"
	"web-studio-backend/internal/pkg/config"
	"web-studio-backend/internal/pkg/wcrypto"
	"web-studio-backend/pkg/postgres"
)

func main() {
	config.Read("config.yml")

	cfg := config.Get()

	user, password, err := wcrypto.DecodeUserPass(cfg.Database.User, cfg.Database.Password, config.Block)
	if err != nil {
		log.Fatalf("decoding database username: %v", err)
	}

	dbConnString := postgres.ConnectionString(user, password, cfg.Database.Host, cfg.Database.Database)

	pg, err := postgres.New(context.Background(), dbConnString)
	if err != nil {
		log.Fatalf("creating postgres: %v", err)
	}

	log.Println("Connected to database")

	userRepo := postgresql.NewUserRepository(pg.Pool)
	userService := service.NewUserService(userRepo)

	_, err = userService.CreateUser(context.Background(), &domain.User{
		Name:            "test",
		Surname:         "testov",
		Username:        "test",
		Email:           "test@mail.com",
		EncodedPassword: "password123",
		Role:            domain.UserRoleGlobalAdmin,
	})
	if err != nil {
		log.Fatalf("Creating user: %v", err)
	}

	log.Println("Done!")
}
