package http

import (
	"net/http"
	"web-studio-backend/internal/app/domain"
)

// GetUser godoc
// @Summary      Get user by identifier.
// @Description  Returns information about single user.
// @Tags         Users
// @Produce      json
// @Param        user_id path int64 true "User identifier."
// @Success      200  {object}  domain.GetUserResponse
// @Failure      400  {object}  errcore.CoreError
// @Failure      500  {object}  errcore.CoreError
// @Router       /api/v1/users/{user_id} [get]
func (s *server) GetUser(w http.ResponseWriter, r *http.Request) {
	userId := s.parseParamInt16("user_id", r)

	response, err := s.core.GetUser(r.Context(), &domain.GetUserRequest{UserId: userId})
	if err != nil {
		s.sendError(err, w)
		return
	}

	s.sendJSON(http.StatusOK, response, w)
}

// CreateUser godoc
// @Summary      Create user.
// @Description  Creates a new user. Returns an object with information about created user.
// @Tags         Users
// @Accept       json
// @Produce      json
// @Param        request body domain.CreateUserRequest true "Request body."
// @Success      200  {object}	domain.CreateUserResponse
// @Failure      400  {object}  errcore.CoreError
// @Failure      500  {object}  errcore.CoreError
// @Router       /api/v1/users [post]
func (s *server) CreateUser(w http.ResponseWriter, r *http.Request) {
	var req domain.CreateUserRequest
	if err := s.readJSON(&req, r); err != nil {
		s.sendError(err, w)
		return
	}

	response, err := s.core.CreateUser(r.Context(), &req)
	if err != nil {
		s.sendError(err, w)
		return
	}

	s.sendJSON(http.StatusOK, response, w)
}

// UpdateUser godoc
// @Summary      Update user.
// @Description  Updates a user. The request body must contain all required fields.
// @Tags         Users
// @Accept       json
// @Produce      json
// @Param        user_id path int64 true "User identifier."
// @Param        request body domain.UpdateUserRequest true "Request body."
// @Success      200  {object}	domain.UpdateUserResponse
// @Failure      404  {object}  errcore.CoreError
// @Failure      500  {object}  errcore.CoreError
// @Router       /api/v1/users/{user_id} [put]
func (s *server) UpdateUser(w http.ResponseWriter, r *http.Request) {
	userId := s.parseParamInt16("user_id", r)

	var req domain.UpdateUserRequest
	if err := s.readJSON(&req, r); err != nil {
		s.sendError(err, w)
		return
	}
	req.UserId = userId

	response, err := s.core.UpdateUser(r.Context(), &req)
	if err != nil {
		s.sendError(err, w)
		return
	}

	s.sendJSON(http.StatusOK, response, w)
}

// DeleteUser godoc
// @Summary      Delete user.
// @Description  Deletes a user.
// @Tags         Users
// @Accept       json
// @Produce      json
// @Param        user_id path int64 true "User identifier."
// @Success      200  {object}	nil
// @Failure      404  {object}  errcore.CoreError
// @Failure      500  {object}  errcore.CoreError
// @Router       /api/v1/users/{user_id} [delete]
func (s *server) DeleteUser(w http.ResponseWriter, r *http.Request) {
	userId := s.parseParamInt16("user_id", r)

	response, err := s.core.DeleteUser(r.Context(), &domain.DeleteUserRequest{UserId: userId})
	if err != nil {
		s.sendError(err, w)
		return
	}

	s.sendJSON(http.StatusOK, response, w)
}
