package main

import (
	"crypto/aes"
	"crypto/rand"
	"encoding/hex"
	"fmt"

	"web-studio-backend/internal/pkg/wcrypto"
)

func main() {
	var buf [32]byte
	_, err := rand.Read(buf[:])
	if err != nil {
		panic(err)
	}

	key := hex.EncodeToString(buf[:])
	fmt.Println("Application key:", key)

	block, err := aes.NewCipher(buf[:])
	if err != nil {
		panic(err)
	}

	username, password, err := wcrypto.EncodeUserPass("webstudio", "webstudio", block)
	if err != nil {
		panic(err)
	}

	fmt.Println("Database credentials:")
	fmt.Println("Username:", username)
	fmt.Println("Password:", password)
}
