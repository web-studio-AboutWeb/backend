package main

import (
	"flag"
	"log"

	"web-studio-backend/internal/app/app"
)

func main() {
	var configPath string
	flag.StringVar(&configPath, "config-path", "config.default.yml", "Path to application config file.")
	flag.Parse()

	if err := app.Run(configPath); err != nil {
		log.Fatal(err)
	}
}
