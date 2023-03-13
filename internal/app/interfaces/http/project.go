package http

import (
	"net/http"
	"web-studio-backend/internal/app/domain"
)

// GetProject godoc
// @Summary      Get project by identifier.
// @Description  Returns information about single user.
// @Tags         Projects
// @Produce      json
// @Param        project_id path int64 true "Project identifier."
// @Success      200  {object}  domain.GetProjectResponse
// @Failure      404  {object}  errcore.CoreError
// @Failure      500  {object}  errcore.CoreError
// @Router       /api/v1/projects/{project_id} [get]
func (s *server) GetProject(w http.ResponseWriter, r *http.Request) {
	ProjectId := s.parseParamInt16("project_id", r)

	response, err := s.core.GetProject(r.Context(), &domain.GetProjectRequest{ProjectId: ProjectId})
	if err != nil {
		s.sendError(err, w)
		return
	}

	s.sendJSON(http.StatusOK, response, w)
}

// CreateProject godoc
// @Summary      Created project.
// @Description  Creates a new project. Returns an object with information about created project.
// @Tags         Projects
// @Accept       json
// @Produce      json
// @Param        request body domain.CreateProjectRequest true "Request body."
// @Success      200  {object}	domain.CreateProjectResponse
// @Failure      400  {object}  errcore.CoreError
// @Failure      500  {object}  errcore.CoreError
// @Router       /api/v1/projects [post]
func (s *server) CreateProject(w http.ResponseWriter, r *http.Request) {
	var req domain.CreateProjectRequest
	if err := s.readJSON(&req, r); err != nil {
		s.sendError(err, w)
		return
	}

	response, err := s.core.CreateProject(r.Context(), &req)
	if err != nil {
		s.sendError(err, w)
		return
	}

	s.sendJSON(http.StatusOK, response, w)
}

// UpdateProject godoc
// @Summary      Update project.
// @Description  Updates a project. The request body must contain all required fields.
// @Tags         Projects
// @Accept       json
// @Produce      json
// @Param        project_id path int64 true "Project identifier."
// @Param        request body domain.UpdateProjectRequest true "Request body."
// @Success      200  {object}	domain.UpdateProjectResponse
// @Failure      404  {object}  errcore.CoreError
// @Failure      500  {object}  errcore.CoreError
// @Router       /api/v1/projects/{project_id} [put]
func (s *server) UpdateProject(w http.ResponseWriter, r *http.Request) {
	ProjectId := s.parseParamInt16("project_id", r)

	var req domain.UpdateProjectRequest
	if err := s.readJSON(&req, r); err != nil {
		s.sendError(err, w)
		return
	}
	req.ProjectId = ProjectId

	response, err := s.core.UpdateProject(r.Context(), &req)
	if err != nil {
		s.sendError(err, w)
		return
	}

	s.sendJSON(http.StatusOK, response, w)
}

// DeleteProject godoc
// @Summary      Delete project.
// @Description  Deletes a project.
// @Tags         Projects
// @Accept       json
// @Produce      json
// @Param        project_id path int64 true "Project identifier."
// @Success      200  {object}	nil
// @Failure      404  {object}  errcore.CoreError
// @Failure      500  {object}  errcore.CoreError
// @Router       /api/v1/projects/{project_id} [delete]
func (s *server) DeleteProject(w http.ResponseWriter, r *http.Request) {
	ProjectId := s.parseParamInt16("project_id", r)

	response, err := s.core.DeleteProject(r.Context(), &domain.DeleteProjectRequest{ProjectId: ProjectId})
	if err != nil {
		s.sendError(err, w)
		return
	}

	s.sendJSON(http.StatusOK, response, w)
}

// GetProjectParticipants godoc
// @Summary      Get project participants.
// @Description  Returns a list of project participants.
// @Tags         Projects
// @Produce      json
// @Param        project_id path int64 true "Project identifier."
// @Success      200  {object}  domain.GetProjectResponse
// @Failure      404  {object}  errcore.CoreError
// @Failure      500  {object}  errcore.CoreError
// @Router       /api/v1/projects/{project_id}/participants [get]
func (s *server) GetProjectParticipants(w http.ResponseWriter, r *http.Request) {
	ProjectId := s.parseParamInt16("project_id", r)

	response, err := s.core.GetProjectParticipants(r.Context(), &domain.GetProjectParticipantsRequest{ProjectId: ProjectId})
	if err != nil {
		s.sendError(err, w)
		return
	}

	s.sendJSON(http.StatusOK, response, w)
}
