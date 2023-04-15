package http

import (
	"net/http"
	_ "web-studio-backend/internal/app/core/shared/errors"
	user_dto "web-studio-backend/internal/app/core/user/dto"
)

// CreateUser godoc
// @Summary      Create user.
// @Description  Creates a new user. Returns an object with information about created user.
// @Tags         Users
// @Accept       json
// @Produce      json
// @Param        request body dto.UserCreate true "Request body."
// @Success      200  {object}	dto.UserObject
// @Failure      400  {object}  errors.CoreError
// @Failure      500  {object}  errors.CoreError
// @Router       /api/v1/users [post]
func (s *server) CreateUser(w http.ResponseWriter, r *http.Request) {
	var user user_dto.UserCreate
	if err := s.readJSON(&user, r); err != nil {
		s.sendError(err, w)
		return
	}

	response, err := s.core.UserHandlers.CreateUserHandler.Execute(
		r.Context(), &user,
	)
	if err != nil {
		s.sendError(err, w)
		return
	}

	s.sendJSON(http.StatusOK, response, w)
}

// GetUser godoc
// @Summary      Get user by identifier.
// @Description  Returns information about single user.
// @Tags         Users
// @Produce      json
// @Param        user_id path int64 true "User identifier."
// @Success      200  {object}  dto.UserObject
// @Failure      400  {object}  errors.CoreError
// @Failure      500  {object}  errors.CoreError
// @Router       /api/v1/users/{user_id} [get]
func (s *server) GetUser(w http.ResponseWriter, r *http.Request) {
	userId := s.parseParamInt16("user_id", r)

	response, err := s.core.UserHandlers.GetUserHandler.Execute(
		r.Context(), &user_dto.UserGet{UserId: userId},
	)
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
// @Param        request body dto.UserUpdate true "Request body."
// @Success      200  {object}	dto.UserObject
// @Failure      404  {object}  errors.CoreError
// @Failure      500  {object}  errors.CoreError
// @Router       /api/v1/users/{user_id} [put]
func (s *server) UpdateUser(w http.ResponseWriter, r *http.Request) {
	userId := s.parseParamInt16("user_id", r)

	var req user_dto.UserUpdate
	if err := s.readJSON(&req, r); err != nil {
		s.sendError(err, w)
		return
	}
	req.UserId = userId
	response, err := s.core.UserHandlers.UpdateUserHandler.Execute(
		r.Context(), &req,
	)
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
// @Failure      404  {object}  errors.CoreError
// @Failure      500  {object}  errors.CoreError
// @Router       /api/v1/users/{user_id} [delete]
func (s *server) DeleteUser(w http.ResponseWriter, r *http.Request) {
	userId := s.parseParamInt16("user_id", r)

	response, err := s.core.UserHandlers.DeleteUserHandler.Execute(
		r.Context(), &user_dto.UserDelete{UserId: userId},
	)
	if err != nil {
		s.sendError(err, w)
		return
	}

	s.sendJSON(http.StatusOK, response, w)
}
