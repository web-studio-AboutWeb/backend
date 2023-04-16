package http

import (
	"net/http"
	_ "web-studio-backend/internal/app/core/shared/errors"
	staffer_dto "web-studio-backend/internal/app/core/staffer/dto"
)

// CreateStaffer godoc
// @Summary      Create staffer.
// @Description  Creates a new staffer. Returns an object with information about created staffer.
// @Tags         Staffers
// @Accept       json
// @Produce      json
// @Param        request body dto.StafferCreate true "Request body."
// @Success      200  {object}	dto.StafferObject
// @Failure      400  {object}  errors.CoreError
// @Failure      500  {object}  errors.CoreError
// @Router       /api/v1/staffers [post]
func (s *server) CreateStaffer(w http.ResponseWriter, r *http.Request) {
	var staffer staffer_dto.StafferCreate
	if err := s.readJSON(&staffer, r); err != nil {
		s.sendError(err, w)
		return
	}

	response, err := s.core.StafferHandlers.CreateStafferHandler.Execute(
		r.Context(), &staffer,
	)
	if err != nil {
		s.sendError(err, w)
		return
	}

	s.sendJSON(http.StatusOK, response, w)
}

// GetStaffer godoc
// @Summary      Get staffer by identifier.
// @Description  Returns information about single staffer.
// @Tags         Staffers
// @Produce      json
// @Param        staffer_id path int64 true "Staffer identifier."
// @Success      200  {object}  dto.StafferObject
// @Failure      400  {object}  errors.CoreError
// @Failure      500  {object}  errors.CoreError
// @Router       /api/v1/staffers/{staffer_id} [get]
func (s *server) GetStaffer(w http.ResponseWriter, r *http.Request) {
	stafferId := s.parseParamInt16("staffer_id", r)

	response, err := s.core.StafferHandlers.GetStafferHandler.Execute(
		r.Context(), &staffer_dto.StafferGet{StafferId: stafferId},
	)
	if err != nil {
		s.sendError(err, w)
		return
	}

	s.sendJSON(http.StatusOK, response, w)
}

// UpdateStaffer godoc
// @Summary      Update staffer.
// @Description  Updates a staffer. The request body must contain all required fields.
// @Tags         Staffers
// @Accept       json
// @Produce      json
// @Param        staffer_id path int64 true "Staffer identifier."
// @Param        request body dto.StafferUpdate true "Request body."
// @Success      200  {object}	dto.StafferObject
// @Failure      404  {object}  errors.CoreError
// @Failure      500  {object}  errors.CoreError
// @Router       /api/v1/staffers/{staffer_id} [put]
func (s *server) UpdateStaffer(w http.ResponseWriter, r *http.Request) {
	stafferId := s.parseParamInt16("staffer_id", r)

	var req staffer_dto.StafferUpdate
	if err := s.readJSON(&req, r); err != nil {
		s.sendError(err, w)
		return
	}
	req.StafferId = stafferId
	response, err := s.core.StafferHandlers.UpdateStafferHandler.Execute(
		r.Context(), &req,
	)
	if err != nil {
		s.sendError(err, w)
		return
	}

	s.sendJSON(http.StatusOK, response, w)
}

// DeleteStaffer godoc
// @Summary      Delete staffer.
// @Description  Deletes a staffer.
// @Tags         Staffers
// @Accept       json
// @Produce      json
// @Param        staffer_id path int64 true "Staffer identifier."
// @Success      200  {object}	nil
// @Failure      404  {object}  errors.CoreError
// @Failure      500  {object}  errors.CoreError
// @Router       /api/v1/staffers/{staffer_id} [delete]
func (s *server) DeleteStaffer(w http.ResponseWriter, r *http.Request) {
	stafferId := s.parseParamInt16("staffer_id", r)

	response, err := s.core.StafferHandlers.DeleteStafferHandler.Execute(
		r.Context(), &staffer_dto.StafferDelete{StafferId: stafferId},
	)
	if err != nil {
		s.sendError(err, w)
		return
	}

	s.sendJSON(http.StatusOK, response, w)
}
