package domain

type (
	AuthContext struct {
		UserID   int32
		Username string
		Email    string
		Role     UserRole
	}

	SignInRequest struct {
		Login    string `json:"login"`
		Password string `json:"password"`
	}

	SignInResponse struct {
		SessionID string `json:"-"`
		CSRFToken string `json:"csrfToken"`
		UserID    int32  `json:"userID"`
	}
)
