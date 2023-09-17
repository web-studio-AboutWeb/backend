package http

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"

	"web-studio-backend/internal/pkg/config"
)

func NewHandler(
	userService UserService,
	projectService ProjectService,
) http.Handler {
	uh := newUserHandler(userService)
	ph := newProjectHandler(projectService)

	r := chi.NewRouter()

	if config.Get().App.Env != "prod" {
		r.Use(middleware.Logger)
	}
	r.Use(middleware.Recoverer)

	r.Use(CorsMiddleware())

	r.Group(func(r chi.Router) {
		// TODO: auth middleware

		r.Get(`/api/v1/users/{user_id}`, uh.getUser)
		r.Post(`/api/v1/users`, uh.createUser)
		r.Put(`/api/v1/users/{user_id}`, uh.updateUser)
		r.Delete(`/api/v1/users/{user_id}`, uh.removeUser)

		r.Get(`/api/v1/projects/{project_id}`, ph.getProject)
		r.Post(`/api/v1/projects`, ph.createProject)
		r.Put(`/api/v1/projects/{project_id}`, ph.updateProject)
		r.Get(`/api/v1/projects/{project_id}/participants`, ph.getProjectParticipants)
	})

	r.Get("/static/*", getStatic)
	r.Get(`/api/v1/docs`, getApiDocs)
	r.Get(`/api/v1/docs/swagger.json`, getApiDocsSwagger)

	return r
}
