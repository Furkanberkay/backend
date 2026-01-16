package dto

import (
	"time"
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
