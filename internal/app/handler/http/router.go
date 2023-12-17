package http

import (
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/httprate"
	"github.com/rs/cors"

	"web-studio-backend/internal/pkg/config"
)

func NewHandler(
	userService UserService,
	projectService ProjectService,
	authService AuthService,
	documentService DocumentService,
	teamService TeamService,
	projectCategoryService ProjectCategoryService,
) http.Handler {
	uh := newUserHandler(userService)
	ph := newProjectHandler(projectService)
	ah := newAuthHandler(authService, userService)
	dh := newDocumentHandler(documentService)
	th := newTeamHandler(teamService)
	pch := newProjectCategoryHandler(projectCategoryService)

	r := chi.NewRouter()

	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(httprate.LimitByIP(69, time.Minute))
	r.Use(middleware.Recoverer)
	r.Use(cors.New(cors.Options{
		AllowedOrigins:   config.Get().Http.AllowedOrigins,
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token", "X-Session-Id"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300,
	}).Handler)
	r.Use(
		middleware.SetHeader("X-Content-Type-Options", "nosniff"),
		middleware.SetHeader("X-Frame-Options", "deny"),
	)
	r.Use(middleware.NoCache)

	r.Get("/static/*", getStatic)

	// TODO: move to auth middleware handlers to keep API private
	r.Get(`/api/v1/docs`, getApiDocs)
	r.Get(`/api/v1/docs/swagger.json`, getApiDocsSwagger)

	r.Post(`/api/v1/auth/sign-in`, ah.signIn)
	r.Post(`/api/v1/auth/sign-out`, ah.signOut)

	r.Group(func(r chi.Router) {
		//r.Use(ah.authMiddleware)
		// TODO: role middlewares

		// Users
		r.Get(`/api/v1/users/{user_id}`, uh.getUser)
		r.Get(`/api/v1/users`, uh.getUsers)
		r.Post(`/api/v1/users`, uh.createUser)
		r.Put(`/api/v1/users/{user_id}`, uh.updateUser)
		r.Delete(`/api/v1/users/{user_id}`, uh.removeUser)
		r.Post(`/api/v1/users/{user_id}/image`, uh.setUserImage)
		r.Get(`/api/v1/users/{user_id}/image`, uh.getUserImage)

		// Projects
		r.Get(`/api/v1/projects/{project_id}`, ph.getProject)
		r.Get(`/api/v1/projects`, ph.getProjects)
		r.Group(func(r chi.Router) {
			// TODO: remove
			r.Use(ah.authMiddleware)
			r.Post(`/api/v1/projects`, ph.createProject)
		})
		r.Put(`/api/v1/projects/{project_id}`, ph.updateProject)
		r.Post(`/api/v1/projects/{project_id}/image`, ph.setProjectImage)
		r.Get(`/api/v1/projects/{project_id}/image`, ph.getProjectImage)
		r.Get(`/api/v1/projects/{project_id}/participants`, ph.getParticipants)
		r.Post(`/api/v1/projects/{project_id}/participants`, ph.addParticipant)
		r.Get(`/api/v1/projects/{project_id}/participants/{user_id}`, ph.getParticipant)
		r.Put(`/api/v1/projects/{project_id}/participants/{user_id}`, ph.updateParticipant)
		r.Delete(`/api/v1/projects/{project_id}/participants/{user_id}`, ph.removeParticipant)

		// Project categories
		r.Get(`/api/v1/projects/categories`, pch.getProjectCategories)
		r.Post(`/api/v1/projects/categories`, pch.createProjectCategory)
		r.Put(`/api/v1/projects/categories/{category_id}`, pch.updateProjectCategory)
		r.Delete(`/api/v1/projects/categories/{category_id}`, pch.deleteProjectCategory)

		// Documents
		r.Get(`/api/v1/projects/{project_id}/documents`, dh.getProjectDocuments)
		r.Post(`/api/v1/projects/{project_id}/documents`, dh.addDocumentToProject)
		r.Delete(`/api/v1/projects/{project_id}/documents/{document_id}`, dh.removeDocumentFromProject)
		r.Get(`/api/v1/documents/{document_id}`, dh.downloadDocument)

		// Teams
		r.Get(`/api/v1/teams/{team_id}`, th.getTeam)
		r.Get(`/api/v1/teams`, th.getTeams)
		r.Post(`/api/v1/teams`, th.createTeam)
		r.Put(`/api/v1/teams/{team_id}`, th.updateTeam)
		r.Post(`/api/v1/teams/{team_id}/image`, th.setTeamImage)
		r.Get(`/api/v1/teams/{team_id}/image`, th.getTeamImage)
		r.Post(`/api/v1/teams/{team_id}/disable`, th.disableTeam)
		r.Post(`/api/v1/teams/{team_id}/enable`, th.enableTeam)
	})

	return r
}
