package auth

type RegisterRequest struct {
	Email                string `json:"email" binding:"required,email,min=3"`
	Password             string `json:"password" binding:"required,eqfield=PasswordConfirmation"`
	PasswordConfirmation string `json:"passwordConfirmation" binding:"required"`
}

type InternalLoginRequest struct {
	Email    string `json:"email" binding:"required,email,min=3"`
	Password string `json:"password" binding:"required"`
}

type RefreshTokenRequest struct {
	RefreshToken string `json:"refreshToken" binding:"required"`
}

type TokenResponse struct {
	Type         string `json:"type"`
	Token        string `json:"token"`
	RefreshToken string `json:"refreshToken"`
}

type RegisterResponse struct {
	Message string `json:"message"`
}

type SendPasswordResetRequest struct {
	Email string `json:"email" binding:"required,email,min=3"`
}

type ResetPasswordRequest struct {
	Token                string `json:"token" binding:"required,min=3"`
	Password             string `json:"password" binding:"required,eqfield=PasswordConfirmation"`
	PasswordConfirmation string `json:"passwordConfirmation" binding:"required"`
}

type OAuthCallbackData struct {
	Provider string `form:"-"`
	Code     string `form:"code" binding:"required,min=1"`
	State    string `form:"state" binding:"required,min=1"`
}

func NewTokenResp(token, refreshToken string) TokenResponse {
	return TokenResponse{
		Type:         "Bearer",
		Token:        token,
		RefreshToken: refreshToken,
	}
}
