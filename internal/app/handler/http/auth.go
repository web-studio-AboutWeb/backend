package http

import (
	"context"
	"net/http"
	"time"

	"web-studio-backend/internal/app/domain"
	"web-studio-backend/internal/app/domain/apperror"
	"web-studio-backend/internal/app/handler/http/httphelp"
	"web-studio-backend/internal/pkg/auth/session"
)

type AuthService interface {
	SignIn(ctx context.Context, req *domain.SignInRequest) (*domain.SignInResponse, error)
	SignOut(ctx context.Context, sessionID string)
}

type authHandler struct {
	authService AuthService
}

func newAuthHandler(authService AuthService) *authHandler {
	return &authHandler{authService: authService}
}

func (h *authHandler) signIn(w http.ResponseWriter, r *http.Request) {
	var req domain.SignInRequest
	if err := httphelp.ReadJSON(&req, r); err != nil {
		httphelp.SendError(err, w)
		return
	}

	response, err := h.authService.SignIn(r.Context(), &req)
	if err != nil {
		httphelp.SendError(err, w)
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:    "session_id",
		Value:   response.SessionID,
		Path:    "/",
		Expires: time.Now().Add(session.Timeout),
	})

	httphelp.SendJSON(http.StatusOK, response, w)
}

func (h *authHandler) signOut(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("session_id")
	if err != nil {
		httphelp.SendError(apperror.NewUnauthorized("Cookie is not provided."), w)
		return
	}

	h.authService.SignOut(r.Context(), cookie.Value)

	http.SetCookie(w, &http.Cookie{
		Name:   "session_id",
		Value:  "",
		Path:   "/",
		MaxAge: -1, // delete cookie
	})
}
