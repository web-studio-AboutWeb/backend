package http

import (
	"context"
	"net/http"

	"web-studio-backend/internal/app/domain"
	"web-studio-backend/internal/app/handler/http/dto"
	"web-studio-backend/internal/app/handler/http/httphelp"
)

//go:generate mockgen -source=user.go -destination=./mocks/user.go -package=mocks
type UserService interface {
	GetUser(ctx context.Context, id int16) (*domain.User, error)
	CreateUser(ctx context.Context, user *domain.User) (*domain.User, error)
	UpdateUser(ctx context.Context, user *domain.User) (*domain.User, error)
	RemoveUser(ctx context.Context, id int16) error
}

type userHandler struct {
	userService UserService
}

func newUserHandler(us UserService) *userHandler {
	return &userHandler{us}
}

// getUser godoc
// @Summary      Get user by identifier
// @Description  Returns information about single user.
// @Tags         Users
// @Produce      json
// @Param        user_id path int64 true "User identifier."
// @Success      200  {object}  domain.User
// @Failure      400  {object}  apperror.Error
// @Failure      500  {object}  apperror.Error
// @Router       /api/v1/users/{user_id} [get]
func (h *userHandler) getUser(w http.ResponseWriter, r *http.Request) {
	userID := httphelp.ParseParamInt16("user_id", r)

	response, err := h.userService.GetUser(r.Context(), userID)
	if err != nil {
		httphelp.SendError(err, w)
		return
	}

	httphelp.SendJSON(http.StatusOK, response, w)
}

// createUser godoc
// @Summary      Create user
// @Description  Creates a new user. Returns an object with information about created user.
// @Tags         Users
// @Accept       json
// @Produce      json
// @Param        request body dto.CreateUserIn true "Request body."
// @Success      200  {object}	domain.User
// @Failure      400  {object}  apperror.Error
// @Failure      500  {object}  apperror.Error
// @Router       /api/v1/users [post]
func (h *userHandler) createUser(w http.ResponseWriter, r *http.Request) {
	var req dto.CreateUserIn
	if err := httphelp.ReadJSON(&req, r); err != nil {
		httphelp.SendError(err, w)
		return
	}

	response, err := h.userService.CreateUser(r.Context(), req.ToDomain())
	if err != nil {
		httphelp.SendError(err, w)
		return
	}

	httphelp.SendJSON(http.StatusOK, response, w)
}

// updateUser godoc
// @Summary      Update user
// @Description  Updates a user. The request body must contain all required fields.
// @Tags         Users
// @Accept       json
// @Produce      json
// @Param        user_id path int16 true "User identifier."
// @Param        request body dto.UpdateUserIn true "Request body."
// @Success      200  {object}	domain.User
// @Failure      404  {object}  apperror.Error
// @Failure      500  {object}  apperror.Error
// @Router       /api/v1/users/{user_id} [put]
func (h *userHandler) updateUser(w http.ResponseWriter, r *http.Request) {
	userID := httphelp.ParseParamInt16("user_id", r)

	var req dto.UpdateUserIn
	if err := httphelp.ReadJSON(&req, r); err != nil {
		httphelp.SendError(err, w)
		return
	}

	response, err := h.userService.UpdateUser(r.Context(), &domain.User{
		ID:       userID,
		Name:     req.Name,
		Surname:  req.Surname,
		Role:     req.Role,
		Position: req.Position,
	})
	if err != nil {
		httphelp.SendError(err, w)
		return
	}

	httphelp.SendJSON(http.StatusOK, response, w)
}

// removeUser godoc
// @Summary      Remove user
// @Description  Marks user as inactive.
// @Tags         Users
// @Accept       json
// @Produce      json
// @Param        user_id path int64 true "User identifier."
// @Success      200  {object}	nil
// @Failure      404  {object}  apperror.Error
// @Failure      500  {object}  apperror.Error
// @Router       /api/v1/users/{user_id} [delete]
func (h *userHandler) removeUser(w http.ResponseWriter, r *http.Request) {
	userID := httphelp.ParseParamInt16("user_id", r)

	err := h.userService.RemoveUser(r.Context(), userID)
	if err != nil {
		httphelp.SendError(err, w)
		return
	}

	w.WriteHeader(http.StatusOK)
}
