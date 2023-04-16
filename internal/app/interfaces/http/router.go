package http

import (
	"net/http"
	"os"

	"github.com/go-chi/chi/v5/middleware"
)

func (s *server) initRouter() http.Handler {
	r := s.router

	if os.Getenv("ENV") == "dev" {
		r.Use(middleware.Logger)
	}
	r.Use(middleware.Recoverer)

	r.Use(corsMiddleware())

	r.Get(`/api/v1/users/{user_id}`, s.GetUser)
	r.Post(`/api/v1/users`, s.CreateUser)
	r.Put(`/api/v1/users/{user_id}`, s.UpdateUser)
	r.Delete(`/api/v1/users/{user_id}`, s.DeleteUser)

	r.Get(`/api/v1/projects/{project_id}`, s.GetProject)
	r.Post(`/api/v1/projects`, s.CreateProject)
	r.Put(`/api/v1/projects/{project_id}`, s.UpdateProject)
	r.Delete(`/api/v1/projects/{project_id}`, s.DeleteProject)
	r.Get(`/api/v1/projects/{project_id}/staffers`, s.GetProjectStaffers)

	r.Get(`/api/v1/staffers/{staffer_id}`, s.GetStaffer)
	r.Post(`/api/v1/staffers`, s.CreateStaffer)
	r.Put(`/api/v1/staffers/{staffer_id}`, s.UpdateStaffer)
	r.Delete(`/api/v1/staffers/{staffer_id}`, s.DeleteStaffer)

	r.Get(`/api/v1/docs`, s.GetApiDocs)

	return r
}
