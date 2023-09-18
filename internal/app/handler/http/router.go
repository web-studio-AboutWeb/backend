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
	authService AuthService,
) http.Handler {
	uh := newUserHandler(userService)
	ph := newProjectHandler(projectService)
	ah := newAuthHandler(authService)

	r := chi.NewRouter()

	if config.Get().App.Env != "prod" {
		r.Use(middleware.Logger)
	}
	r.Use(middleware.Recoverer)
	r.Use(corsMiddleware())

	r.Get("/static/*", getStatic)

	// TODO: use auth middleware here to keep API private
	r.Get(`/api/v1/docs`, getApiDocs)
	r.Get(`/api/v1/docs/swagger.json`, getApiDocsSwagger)

	r.Post(`/api/v1/auth/sign-in`, ah.signIn)
	r.Post(`/api/v1/auth/sign-out`, ah.signOut)

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

	return r
}
