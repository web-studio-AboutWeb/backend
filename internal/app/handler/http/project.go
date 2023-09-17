package http

import (
	"context"
	"net/http"

	"web-studio-backend/internal/app/domain"
	"web-studio-backend/internal/app/handler/http/httphelp"
)

type ProjectService interface {
	GetProject(ctx context.Context, id int16) (*domain.Project, error)
	CreateProject(ctx context.Context, project *domain.Project) (*domain.Project, error)
	UpdateProject(ctx context.Context, project *domain.Project) (*domain.Project, error)
	GetProjectParticipants(ctx context.Context, projectID int16) ([]domain.User, error)
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
// @Param        project_id path int64 true "Project identifier."
// @Success      200  {object}  domain.Project
// @Failure      404  {object}  apperror.Error
// @Failure      500  {object}  apperror.Error
// @Router       /api/v1/projects/{project_id} [get]
func (h *projectHandler) getProject(w http.ResponseWriter, r *http.Request) {
	pid := httphelp.ParseParamInt16("project_id", r)

	response, err := h.projectService.GetProject(r.Context(), pid)
	if err != nil {
		httphelp.SendError(err, w)
		return
	}

	httphelp.SendJSON(http.StatusOK, response, w)
}

// createProject godoc
// @Summary      Created project
// @Description  Creates a new project. Returns an object with information about created project.
// @Tags         Projects
// @Accept       json
// @Produce      json
// @Param        request body domain.Project true "Request body."
// @Success      200  {object}	domain.Project
// @Failure      400  {object}  apperror.Error
// @Failure      500  {object}  apperror.Error
// @Router       /api/v1/projects [post]
func (h *projectHandler) createProject(w http.ResponseWriter, r *http.Request) {
	var project domain.Project
	if err := httphelp.ReadJSON(&project, r); err != nil {
		httphelp.SendError(err, w)
		return
	}

	response, err := h.projectService.CreateProject(r.Context(), &project)
	if err != nil {
		httphelp.SendError(err, w)
		return
	}

	httphelp.SendJSON(http.StatusOK, response, w)
}

// updateProject godoc
// @Summary      Update project
// @Description  Updates a project. The request body must contain all required fields.
// @Tags         Projects
// @Accept       json
// @Produce      json
// @Param        project_id path int64 true "Project identifier."
// @Param        request body domain.Project true "Request body."
// @Success      200  {object}	domain.Project
// @Failure      404  {object}  apperror.Error
// @Failure      500  {object}  apperror.Error
// @Router       /api/v1/projects/{project_id} [put]
func (h *projectHandler) updateProject(w http.ResponseWriter, r *http.Request) {
	projectID := httphelp.ParseParamInt16("project_id", r)

	var req domain.Project
	if err := httphelp.ReadJSON(&req, r); err != nil {
		httphelp.SendError(err, w)
		return
	}
	req.ID = projectID

	response, err := h.projectService.UpdateProject(r.Context(), &req)
	if err != nil {
		httphelp.SendError(err, w)
		return
	}

	httphelp.SendJSON(http.StatusOK, response, w)
}

// getProjectParticipants godoc
// @Summary      Get project participants
// @Description  Returns a list of project participants.
// @Tags         Projects
// @Produce      json
// @Param        project_id path int64 true "Project identifier."
// @Success      200  {array}   domain.User
// @Failure      404  {object}  apperror.Error
// @Failure      500  {object}  apperror.Error
// @Router       /api/v1/projects/{project_id}/participants [get]
func (h *projectHandler) getProjectParticipants(w http.ResponseWriter, r *http.Request) {
	projectID := httphelp.ParseParamInt16("project_id", r)

	response, err := h.projectService.GetProjectParticipants(r.Context(), projectID)
	if err != nil {
		httphelp.SendError(err, w)
		return
	}

	httphelp.SendJSON(http.StatusOK, response, w)
}
