package http

import (
	"github.com/go-chi/chi/v5/middleware"
	"net/http"
	"os"
)

func (s *server) initRouter() http.Handler {
	r := s.router

	if os.Getenv("ENV") == "dev" {
		r.Use(middleware.Logger)
	}
	r.Use(middleware.Recoverer)

	r.Use(corsMiddleware())

	r.Post("/api/v1/auth/sign-in", s.SignIn)
	r.Post(`/api/v1/auth/sign-up`, s.SignUp)

	r.Get(`/api/v1/users/{user_id}`, s.GetUser)

	return r
}
