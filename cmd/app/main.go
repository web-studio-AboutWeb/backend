package main

import (
	"flag"
	"log"

	"web-studio-backend/internal/app/app"
)

// @Title    AboutWeb API
// @Version  1.0
// @Accept   json
// @Produce  json

// @SecurityDefinitions.apikey CSRF
// @In                         header
// @Name                       X-CSRF-Token
// @Description                **Each endpoint that requires authorization MUST be used with CSRF token in header. **
// @Description                **Example usage: `X-CSRF-Token your-token`**.

func main() {
	var configPath string
	flag.StringVar(&configPath, "config-path", "config.yml", "Path to application config file.")
	flag.Parse()

	if err := app.Run(configPath); err != nil {
		log.Fatal(err)
	}
}
