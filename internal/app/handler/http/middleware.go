package http

import (
	"net/http"

	"github.com/go-chi/cors"
)

func corsMiddleware() func(next http.Handler) http.Handler {
	return cors.Handler(cors.Options{
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: false,
	})
}
