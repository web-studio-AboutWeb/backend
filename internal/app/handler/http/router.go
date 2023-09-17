package http

import (
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func NewHandler(
	userService UserService,
	projectService ProjectService,
) http.Handler {
	uh := newUserHandler(userService)
	ph := newProjectHandler(projectService)

	r := chi.NewRouter()

	if os.Getenv("ENV") != "prod" {
		r.Use(middleware.Logger)
	}
	r.Use(middleware.Recoverer)

	r.Use(CorsMiddleware())

	r.Group(func(r chi.Router) {
		// TODO: router middleware

		r.Get(`/api/v1/users/{user_id}`, uh.getUser)
		r.Post(`/api/v1/users`, uh.createUser)
		r.Put(`/api/v1/users/{user_id}`, uh.updateUser)
		r.Delete(`/api/v1/users/{user_id}`, uh.removeUser)

		r.Get(`/api/v1/projects/{project_id}`, ph.getProject)
		r.Post(`/api/v1/projects`, ph.createProject)
		r.Put(`/api/v1/projects/{project_id}`, ph.updateProject)
		r.Get(`/api/v1/projects/{project_id}/participants`, ph.getProjectParticipants)
	})

	r.Get(`/api/v1/docs`, getApiDocs)

	return r
}
