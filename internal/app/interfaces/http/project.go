package http

import (
	"net/http"
	project_dto "web-studio-backend/internal/app/core/project/dto"
	_ "web-studio-backend/internal/app/core/shared/errors"
)

// GetProject godoc
// @Summary      Get project by identifier.
// @Description  Returns information about single user.
// @Tags         Projects
// @Produce      json
// @Param        project_id path int64 true "Project identifier."
// @Success      200  {object}  dto.ProjectObject
// @Failure      404  {object}  errors.CoreError
// @Failure      500  {object}  errors.CoreError
// @Router       /api/v1/projects/{project_id} [get]
func (s *server) GetProject(w http.ResponseWriter, r *http.Request) {
	ProjectId := s.parseParamInt16("project_id", r)

	response, err := s.core.ProjectHandlers.GetProjectHandler.Execute(
		r.Context(), &project_dto.ProjectGet{ProjectId: ProjectId})
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
// @Param        request body dto.ProjectCreate true "Request body."
// @Success      200  {object}	dto.ProjectObject
// @Failure      400  {object}  errors.CoreError
// @Failure      500  {object}  errors.CoreError
// @Router       /api/v1/projects [post]
func (s *server) CreateProject(w http.ResponseWriter, r *http.Request) {
	var project project_dto.ProjectCreate
	if err := s.readJSON(&project, r); err != nil {
		s.sendError(err, w)
		return
	}

	response, err := s.core.ProjectHandlers.CreateProjectHandler.Execute(
		r.Context(), &project,
	)
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
// @Param        request body dto.ProjectUpdate true "Request body."
// @Success      200  {object}	dto.ProjectObject
// @Failure      404  {object}  errors.CoreError
// @Failure      500  {object}  errors.CoreError
// @Router       /api/v1/projects/{project_id} [put]
func (s *server) UpdateProject(w http.ResponseWriter, r *http.Request) {
	ProjectId := s.parseParamInt16("project_id", r)

	var project project_dto.ProjectUpdate
	if err := s.readJSON(&project, r); err != nil {
		s.sendError(err, w)
		return
	}
	project.ProjectId = ProjectId

	response, err := s.core.ProjectHandlers.UpdateProjectHandler.Execute(
		r.Context(), &project,
	)
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
// @Failure      404  {object}  errors.CoreError
// @Failure      500  {object}  errors.CoreError
// @Router       /api/v1/projects/{project_id} [delete]
func (s *server) DeleteProject(w http.ResponseWriter, r *http.Request) {
	ProjectId := s.parseParamInt16("project_id", r)

	response, err := s.core.ProjectHandlers.DeleteProjectHandler.Execute(
		r.Context(), &project_dto.ProjectDelete{ProjectId: ProjectId},
	)
	if err != nil {
		s.sendError(err, w)
		return
	}

	s.sendJSON(http.StatusOK, response, w)
}

// GetProjectStaffers godoc
// @Summary      Get project participants.
// @Description  Returns a list of project participants.
// @Tags         Projects
// @Produce      json
// @Param        project_id path int64 true "Project identifier."
// @Success      200  {object}  dto.ProjectParticipants
// @Failure      404  {object}  errors.CoreError
// @Failure      500  {object}  errors.CoreError
// @Router       /api/v1/projects/{project_id}/participants [get]
func (s *server) GetProjectStaffers(w http.ResponseWriter, r *http.Request) {
	ProjectId := s.parseParamInt16("project_id", r)

	response, err := s.core.ProjectHandlers.GetProjectStaffersHandler.Execute(
		r.Context(), &project_dto.ProjectStaffersGet{ProjectId: ProjectId},
	)
	if err != nil {
		s.sendError(err, w)
		return
	}

	s.sendJSON(http.StatusOK, response, w)
}
