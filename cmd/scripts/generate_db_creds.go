package main

import (
	"fmt"

	"web-studio-backend/internal/pkg/config"
	"web-studio-backend/internal/pkg/wcrypto"
)

func main() {
	config.Read("./config.yml")

	username, password, err := wcrypto.EncodeUserPass("webstudio", "webstudio", config.Block)
	if err != nil {
		panic(err)
	}

	fmt.Println("Username:", username)
	fmt.Println("Password:", password)
}
