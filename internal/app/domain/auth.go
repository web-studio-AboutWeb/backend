package domain

type (
	AccessTokenMetadata struct {
		Id       int64    `json:"id"`
		Username string   `json:"username"`
		Email    string   `json:"email"`
		Role     UserRole `json:"role"`
	}
	RefreshTokenMetadata struct {
		Id int64 `json:"id"`
	}

	SignInRequest struct {
		Login    string `json:"login"`
		Password string `json:"password"`
	}
	SignInResponse struct {
		AccessToken  string `json:"accessToken"`
		RefreshToken string `json:"refreshToken"`
	}

	SignUpRequest struct {
		Username string `json:"username"`
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	SignUpResponse struct {
		AccessToken  string `json:"accessToken"`
		RefreshToken string `json:"refreshToken"`
	}
)
