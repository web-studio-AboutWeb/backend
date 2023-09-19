package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/golang-migrate/migrate/v4"

	"web-studio-backend/internal/pkg/config"
	"web-studio-backend/internal/pkg/wcrypto"
	"web-studio-backend/pkg/postgres"

	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

func main() {
	var configPath string
	flag.StringVar(&configPath, "config-path", "config.yml", "Path to application config file.")
	flag.Parse()

	if len(os.Args) < 2 {
		fmt.Println("use up or down syntax to make migration: ./migrate up")
		return
	}

	config.Read(configPath)
	cfg := config.Get()

	user, password, err := wcrypto.DecodeUserPass(cfg.Database.User, cfg.Database.Password, config.Block)
	if err != nil {
		log.Fatalf("decondig db credentials: %v", err)
	}

	dbConnString := postgres.ConnectionString(user, password, cfg.Database.Host, cfg.Database.Database)

	m, err := migrate.New("file://migrations", dbConnString)
	if err != nil {
		log.Fatalf("creating migration: %v", err)
	}

	action := os.Args[1]

	switch action {
	case "up":
		err = m.Up()
	case "down":
		err = m.Down()
	case "drop":
		err = m.Drop()
	default:
		fmt.Println("unknown action:", action)
	}

	if err != nil {
		log.Fatalf("applying migration: %v", err)
	}

	fmt.Printf("Migrate %s done.\n", action)
}
