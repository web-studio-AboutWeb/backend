package http

import (
	"context"
	"net/http"

	"web-studio-backend/internal/app/domain"
	"web-studio-backend/internal/app/handler/http/dto"
	"web-studio-backend/internal/app/handler/http/httphelp"
)

//go:generate mockgen -source=project.go -destination=./mocks/project.go -package=mocks
type ProjectService interface {
	GetProject(ctx context.Context, id int32) (*domain.Project, error)
	GetProjects(ctx context.Context) ([]domain.Project, error)
	CreateProject(ctx context.Context, project *domain.Project) (*domain.Project, error)
	UpdateProject(ctx context.Context, project *domain.Project) (*domain.Project, error)

	GetParticipants(ctx context.Context, projectID int32) ([]domain.ProjectParticipant, error)
	GetParticipant(ctx context.Context, participantID, projectID int32) (*domain.ProjectParticipant, error)
	AddParticipant(ctx context.Context, participant *domain.ProjectParticipant) (*domain.ProjectParticipant, error)
	UpdateParticipant(ctx context.Context, participant *domain.ProjectParticipant) (*domain.ProjectParticipant, error)
	RemoveParticipant(ctx context.Context, participantID, projectID int32) error
}

type projectHandler struct {
	projectService ProjectService
}

func newProjectHandler(ps ProjectService) *projectHandler {
	return &projectHandler{ps}
}

// getProject godoc
// @Summary      Get project by identifier
// @Description  Returns information about single user.
// @Tags         Projects
// @Produce      json
// @Param        project_id path int true "Project identifier."
// @Success      200  {object}  domain.Project
// @Failure      404  {object}  apperror.Error
// @Failure      500  {object}  apperror.Error
// @Router       /api/v1/projects/{project_id} [get]
func (h *projectHandler) getProject(w http.ResponseWriter, r *http.Request) {
	pid := httphelp.ParseParamInt32("project_id", r)

	response, err := h.projectService.GetProject(r.Context(), pid)
	if err != nil {
		httphelp.SendError(err, w)
		return
	}

	httphelp.SendJSON(http.StatusOK, response, w)
}

// getProjects godoc
// @Summary      Get projects
// @Description  Returns list of projects.
// @Tags         Projects
// @Produce      json
// @Success      200  {array}  domain.Project
// @Failure      404  {object}  apperror.Error
// @Failure      500  {object}  apperror.Error
// @Router       /api/v1/projects [get]
func (h *projectHandler) getProjects(w http.ResponseWriter, r *http.Request) {
	response, err := h.projectService.GetProjects(r.Context())
	if err != nil {
		httphelp.SendError(err, w)
		return
	}

	httphelp.SendJSON(http.StatusOK, response, w)
}

// createProject godoc
// @Summary      Create project
// @Description  Creates a new project. Returns an object with information about created project.
// @Tags         Projects
// @Accept       json
// @Produce      json
// @Param        request body dto.CreateProjectRequest true "Request body."
// @Success      200  {object}	domain.Project
// @Failure      400  {object}  apperror.Error
// @Failure      500  {object}  apperror.Error
// @Router       /api/v1/projects [post]
func (h *projectHandler) createProject(w http.ResponseWriter, r *http.Request) {
	var req dto.CreateProjectRequest
	if err := httphelp.ReadJSON(&req, r); err != nil {
		httphelp.SendError(err, w)
		return
	}

	response, err := h.projectService.CreateProject(r.Context(), req.ToDomain())
	if err != nil {
		httphelp.SendError(err, w)
		return
	}

	httphelp.SendJSON(http.StatusOK, response, w)
}

// updateProject godoc
// @Summary      Update project
// @Description  Updates a project.
// @Tags         Projects
// @Accept       json
// @Produce      json
// @Param        project_id path int true "Project identifier."
// @Param        request body dto.UpdateProjectRequest true "Request body."
// @Success      200  {object}	domain.Project
// @Failure      404  {object}  apperror.Error
// @Failure      500  {object}  apperror.Error
// @Router       /api/v1/projects/{project_id} [put]
func (h *projectHandler) updateProject(w http.ResponseWriter, r *http.Request) {
	projectID := httphelp.ParseParamInt32("project_id", r)

	var req dto.UpdateProjectRequest
	if err := httphelp.ReadJSON(&req, r); err != nil {
		httphelp.SendError(err, w)
		return
	}

	response, err := h.projectService.UpdateProject(r.Context(), req.ToDomain(projectID))
	if err != nil {
		httphelp.SendError(err, w)
		return
	}

	httphelp.SendJSON(http.StatusOK, response, w)
}

// getParticipants godoc
// @Summary      Get project participants
// @Description  Returns a list of project participants.
// @Tags         Projects
// @Produce      json
// @Param        project_id path int true "Project identifier."
// @Success      200  {array}   domain.ProjectParticipant
// @Failure      404  {object}  apperror.Error
// @Failure      500  {object}  apperror.Error
// @Router       /api/v1/projects/{project_id}/participants [get]
func (h *projectHandler) getParticipants(w http.ResponseWriter, r *http.Request) {
	projectID := httphelp.ParseParamInt32("project_id", r)

	response, err := h.projectService.GetParticipants(r.Context(), projectID)
	if err != nil {
		httphelp.SendError(err, w)
		return
	}

	httphelp.SendJSON(http.StatusOK, response, w)
}

// getParticipant godoc
// @Summary      Get project participant
// @Description  Returns information about project participant.
// @Tags         Projects
// @Produce      json
// @Param        project_id path int true "Project identifier."
// @Param        user_id path int true "Participant identifier."
// @Success      200  {object}  domain.ProjectParticipant
// @Failure      404  {object}  apperror.Error
// @Failure      500  {object}  apperror.Error
// @Router       /api/v1/projects/{project_id}/participants/{user_id} [get]
func (h *projectHandler) getParticipant(w http.ResponseWriter, r *http.Request) {
	projectID := httphelp.ParseParamInt32("project_id", r)
	participantID := httphelp.ParseParamInt32("user_id", r)

	response, err := h.projectService.GetParticipant(r.Context(), projectID, participantID)
	if err != nil {
		httphelp.SendError(err, w)
		return
	}

	httphelp.SendJSON(http.StatusOK, response, w)
}

// addParticipant godoc
// @Summary      Add participant to project
// @Description  Adds user to project participants list.
// @Description
// @Description  On success returns information about added participant.
// @Tags         Projects
// @Produce      json
// @Param        project_id path int true "Project identifier."
// @Param        request body dto.AddProjectParticipantRequest true "Request body."
// @Success      200  {object}  domain.ProjectParticipant
// @Failure      404  {object}  apperror.Error
// @Failure      500  {object}  apperror.Error
// @Router       /api/v1/projects/{project_id}/participants [post]
func (h *projectHandler) addParticipant(w http.ResponseWriter, r *http.Request) {
	projectID := httphelp.ParseParamInt32("project_id", r)

	var req dto.AddProjectParticipantRequest
	if err := httphelp.ReadJSON(&req, r); err != nil {
		httphelp.SendError(err, w)
		return
	}

	response, err := h.projectService.AddParticipant(r.Context(), req.ToDomain(projectID))
	if err != nil {
		httphelp.SendError(err, w)
		return
	}

	httphelp.SendJSON(http.StatusOK, response, w)
}

// updateParticipant godoc
// @Summary      Update project participant
// @Description  Updates participant role and position.
// @Description
// @Description  On success returns information about updated participant.
// @Tags         Projects
// @Produce      json
// @Param        project_id path int true "Project identifier."
// @Param        user_id path int true "Participant identifier."
// @Param        request body dto.UpdateProjectParticipantRequest true "Request body."
// @Success      200  {object}  domain.ProjectParticipant
// @Failure      404  {object}  apperror.Error
// @Failure      500  {object}  apperror.Error
// @Router       /api/v1/projects/{project_id}/participants/{user_id} [put]
func (h *projectHandler) updateParticipant(w http.ResponseWriter, r *http.Request) {
	projectID := httphelp.ParseParamInt32("project_id", r)
	userID := httphelp.ParseParamInt32("user_id", r)

	var req dto.UpdateProjectParticipantRequest
	if err := httphelp.ReadJSON(&req, r); err != nil {
		httphelp.SendError(err, w)
		return
	}

	response, err := h.projectService.UpdateParticipant(r.Context(), req.ToDomain(projectID, userID))
	if err != nil {
		httphelp.SendError(err, w)
		return
	}

	httphelp.SendJSON(http.StatusOK, response, w)
}

// removeParticipant godoc
// @Summary      Remove project participant
// @Description  Deletes the user from project participants list.
// @Tags         Projects
// @Param        project_id path int true "Project identifier."
// @Param        user_id path int true "Participant identifier."
// @Success      200
// @Failure      404  {object}  apperror.Error
// @Failure      500  {object}  apperror.Error
// @Router       /api/v1/projects/{project_id}/participants/{user_id} [delete]
func (h *projectHandler) removeParticipant(w http.ResponseWriter, r *http.Request) {
	projectID := httphelp.ParseParamInt32("project_id", r)
	userID := httphelp.ParseParamInt32("user_id", r)

	err := h.projectService.RemoveParticipant(r.Context(), userID, projectID)
	if err != nil {
		httphelp.SendError(err, w)
		return
	}

	w.WriteHeader(http.StatusOK)
}
