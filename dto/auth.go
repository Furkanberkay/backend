package dto

import (
	"time"

	"github.com/Furkanberkay/ticket-booking-project-v1/models"
)

type RegisterUserRequest struct {
	Name     string    `json:"name" validate:"required,min=2,max=100"`
	Surname  string    `json:"surname" validate:"required,min=2,max=100"`
	Email    string    `json:"email" validate:"required,email"`
	Password string    `json:"password" validate:"required,min=6,max=32"`
	Birthday time.Time `json:"birthday" validate:"required"`
}

type LoginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

type UpdateUserRequest struct {
	Name     *string    `json:"name" validate:"omitempty,min=2"`
	Surname  *string    `json:"surname" validate:"omitempty,min=2"`
	Birthday *time.Time `json:"birthday" validate:"omitempty"`
}

type LoginResponse struct {
	AccessToken  string        `json:"access_token"`
	RefreshToken string        `json:"refresh_token"`
	User         *UserResponse `json:"user"`
}

func ToUserResponse(user *models.User) *UserResponse {
	if user == nil {
		return nil
	}

	return &UserResponse{
		ID:       user.ID,
		Name:     user.Name,
		Surname:  user.Surname,
		Email:    user.Email,
		Role:     user.Role,
		Birthday: user.Birthday,
	}
}

type UserResponse struct {
	ID       uint            `json:"id"`
	Name     string          `json:"name"`
	Surname  string          `json:"surname"`
	Email    string          `json:"email"`
	Role     models.UserRole `json:"role"`
	Birthday time.Time       `json:"birthday"`
}

type RefreshTokenRequest struct {
	RefreshToken string `json:"refresh_token" validate:"required"`
}
