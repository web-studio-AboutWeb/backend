package http

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	"path/filepath"

	"web-studio-backend/internal/app/domain"
	"web-studio-backend/internal/app/handler/http/dto"
	"web-studio-backend/internal/app/handler/http/httphelp"
)

//go:generate mockgen -source=user.go -destination=./mocks/user.go -package=mocks
type UserService interface {
	GetUser(ctx context.Context, id int32) (*domain.User, error)
	GetUsers(ctx context.Context) ([]domain.User, error)
	CreateUser(ctx context.Context, user *domain.User) (*domain.User, error)
	UpdateUser(ctx context.Context, user *domain.User) (*domain.User, error)
	RemoveUser(ctx context.Context, id int32) error

	SetUserImage(ctx context.Context, userID int32, img []byte) error
	GetUserImage(ctx context.Context, userID int32) (*domain.User, error)
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
// @Param        user_id path int true "User identifier."
// @Success      200  {object}  domain.User
// @Failure      400  {object}  Error
// @Failure      500  {object}  Error
// @Router       /api/v1/users/{user_id} [get]
func (h *userHandler) getUser(w http.ResponseWriter, r *http.Request) {
	userID := httphelp.ParseParamInt32("user_id", r)

	response, err := h.userService.GetUser(r.Context(), userID)
	if err != nil {
		httphelp.SendError(err, w)
		return
	}

	httphelp.SendJSON(http.StatusOK, response, w)
}

// getUsers godoc
// @Summary      Get users
// @Description  Returns a list of users.
// @Tags         Users
// @Produce      json
// @Success      200  {array}  domain.User
// @Failure      400  {object}  Error
// @Failure      500  {object}  Error
// @Router       /api/v1/users [get]
func (h *userHandler) getUsers(w http.ResponseWriter, r *http.Request) {
	response, err := h.userService.GetUsers(r.Context())
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
// @Failure      400  {object}  Error
// @Failure      500  {object}  Error
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
// @Description  Updates a user.
// @Tags         Users
// @Accept       json
// @Produce      json
// @Param        user_id path int16 true "User identifier."
// @Param        request body dto.UpdateUserIn true "Request body."
// @Success      200  {object}	domain.User
// @Failure      404  {object}  Error
// @Failure      500  {object}  Error
// @Router       /api/v1/users/{user_id} [put]
func (h *userHandler) updateUser(w http.ResponseWriter, r *http.Request) {
	userID := httphelp.ParseParamInt32("user_id", r)

	var req dto.UpdateUserIn
	if err := httphelp.ReadJSON(&req, r); err != nil {
		httphelp.SendError(err, w)
		return
	}

	response, err := h.userService.UpdateUser(r.Context(), req.ToDomain(userID))
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
// @Param        user_id path int true "User identifier."
// @Success      200  {object}	nil
// @Failure      404  {object}  Error
// @Failure      500  {object}  Error
// @Router       /api/v1/users/{user_id} [delete]
func (h *userHandler) removeUser(w http.ResponseWriter, r *http.Request) {
	userID := httphelp.ParseParamInt32("user_id", r)

	err := h.userService.RemoveUser(r.Context(), userID)
	if err != nil {
		httphelp.SendError(err, w)
		return
	}

	w.WriteHeader(http.StatusOK)
}

// setUserImage godoc
// @Summary      Set user image
// @Description  Updated user image. Accepts `multipart/form-data`.
// @Description
// @Description  Note: if a user already has an image, it will be deleted automatically on success.
// @Tags         Users
// @Accept       mpfd
// @Param        user_id path int true "User identifier."
// @Param        file formData file true "Image file. MUST have one of the following mime types: [`image/jpeg`, `image/png`, `image/webp`]"
// @Success      200
// @Failure      400  {object}  Error
// @Failure      404  {object}  Error
// @Failure      500  {object}  Error
// @Router       /api/v1/users/{user_id}/image [post]
func (h *userHandler) setUserImage(w http.ResponseWriter, r *http.Request) {
	tid := httphelp.ParseParamInt32("user_id", r)

	file, _, err := r.FormFile("file")
	if err != nil {
		if errors.Is(err, http.ErrMissingFile) {
			// TODO: custom http error
			httphelp.SendError(fmt.Errorf("file is not presented"), w)
			return
		}
		httphelp.SendError(fmt.Errorf("parsing form file: %w", err), w)
		return
	}
	defer file.Close()

	content, err := io.ReadAll(file)
	if err != nil {
		httphelp.SendError(fmt.Errorf("reading file: %w", err), w)
		return
	}

	err = h.userService.SetUserImage(r.Context(), tid, content)
	if err != nil {
		httphelp.SendError(fmt.Errorf("setting team image: %w", err), w)
		return
	}

	w.WriteHeader(http.StatusOK)
}

// getUserImage godoc
// @Summary      Get user image content
// @Description  Returns user image.
// @Tags         Users
// @Produce      octet-stream
// @Param        user_id path int true "User identifier."
// @Success      200
// @Failure      404  {object}  Error
// @Failure      500  {object}  Error
// @Router       /api/v1/users/{user_id}/image [get]
func (h *userHandler) getUserImage(w http.ResponseWriter, r *http.Request) {
	tid := httphelp.ParseParamInt32("user_id", r)

	response, err := h.userService.GetUserImage(r.Context(), tid)
	if err != nil {
		httphelp.SendError(fmt.Errorf("getting team image: %w", err), w)
		return
	}

	fileName := fmt.Sprintf("%s.%s", response.Username, filepath.Ext(response.ImageID))

	http.ServeContent(w, r, fileName, response.UpdatedAt, bytes.NewReader(response.ImageContent))
}
