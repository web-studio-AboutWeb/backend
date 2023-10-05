package http

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"web-studio-backend/internal/app/domain"
	"web-studio-backend/internal/app/handler/http/httphelp"
	"web-studio-backend/internal/pkg/auth"
	"web-studio-backend/internal/pkg/auth/session"
)

//go:generate mockgen -source=auth.go -destination=./mocks/auth.go -package=mocks
type AuthService interface {
	SignIn(ctx context.Context, req *domain.SignInRequest) (*domain.SignInResponse, error)
	SignOut(ctx context.Context, sessionID string)
}

type authHandler struct {
	authService AuthService
	userService UserService
}

func newAuthHandler(authService AuthService, userService UserService) *authHandler {
	return &authHandler{authService: authService, userService: userService}
}

// signIn godoc
// @Summary      Sign in
// @Description  Starts user session. On success returns CSRF token and sets HTTP cookie.
// @Description
// @Description  All the following requests must contain the X-CSRF-token header for successful authorization.
// @Description
// @Description  Example: `X-CSRF-token: <token>`.
// @Description
// @Description  If authorization will fail, `401 Unauthorized` status code will be returned without any additional data.
// @Tags         Auth
// @Param        request body domain.SignInRequest true "Request body."
// @Success      200  {object}	domain.SignInResponse
// @Failure      401  {object}  Error
// @Failure      500  {object}  Error
// @Router       /api/v1/auth/sign-in [post]
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
		Name:     "session_id",
		Value:    response.SessionID,
		SameSite: http.SameSiteNoneMode,
		Path:     "/",
		HttpOnly: true,
		Secure:   true,
		Expires:  time.Now().Add(session.TTL),
	})

	httphelp.SendJSON(http.StatusOK, response, w)
}

// signOut godoc
// @Summary      Sign out
// @Description  Ends user session and logs out the user.
// @Tags         Auth
// @Success      200
// @Failure      401  {object}  Error
// @Failure      500  {object}  Error
// @Router       /api/v1/auth/sign-out [post]
func (h *authHandler) signOut(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("session_id")
	if err != nil {
		// Cookie is not provided, just skip.
		w.WriteHeader(http.StatusOK)
		return
	}

	h.authService.SignOut(r.Context(), cookie.Value)

	http.SetCookie(w, &http.Cookie{
		Name:   "session_id",
		Value:  "",
		Path:   "/",
		MaxAge: -1, // delete cookie
	})

	w.WriteHeader(http.StatusOK)
}

func (h *authHandler) authMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie("session_id")
		if err != nil {
			http.Error(w, "session_id does not exists", http.StatusUnauthorized)
			return
		}

		sess, err := session.GetSession(cookie.Value)
		if err != nil {
			http.Error(w, fmt.Sprintf("user session not found(session id: %s)", cookie.Value), http.StatusUnauthorized)
			return
		}

		token := r.Header.Get("X-CSRF-Token")
		if token != sess.CSRFToken {
			http.Error(w, "invalid csrf token("+token+")", http.StatusUnauthorized)
			return
		}

		user, err := h.userService.GetUser(r.Context(), sess.UserID)
		if err != nil {
			http.Error(w, "user not found", http.StatusUnauthorized)
			return
		}

		ctx := auth.NewContext(r.Context(), &domain.AuthContext{
			UserID:   user.ID,
			Username: user.Username,
			Email:    user.Email,
			Role:     user.Role,
		})

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
