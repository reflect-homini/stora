package user

import "github.com/reflect-homini/stora/internal/domain/dto"

type NewUserRequest struct {
	Email     string
	Password  string
	Name      string
	Avatar    string
	VerifyNow bool
}

type UserResponse struct {
	dto.BaseDTO
	Email   string          `json:"email"`
	Profile ProfileResponse `json:"profile"`
}
