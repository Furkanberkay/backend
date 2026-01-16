package models

import (
	"context"

	"golang.org/x/crypto/bcrypt"
)

type AuthCredentials struct {
	Email    string `json:"email" validate:"required"`
	Password string `json:"password" validate:"required"`
}

type AuthRepository interface {
	RegisterUser(ctx context.Context, registerData *User) (*User, error)
	GetUserByEmail(ctx context.Context, email string) (*User, error)
	GetUserByUserID(ctx context.Context, userID uint) (*User, error)
}

type AuthService interface {
	Login(ctx context.Context, loginData *AuthCredentials) (string, *User, error)
	Register(ctx context.Context, user *User, plainPassword string) (string, *User, error)
}

func MatchesHash(password string, hash string) bool {
	if err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password)); err != nil {
		return false
	}
	return true
}
