package domain

type (
	SignInRequest struct {
		Login    string `json:"login"`
		Password string `json:"password"`
	}
	SignInResponse struct {
		SessionID string `json:"-"`
		CSRFToken string `json:"csrfToken"`
	}
)
