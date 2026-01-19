package models

import (
	"context"
	"errors"

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

type TokenRepository interface {
	Create(ctx context.Context, token *RefreshToken) error
	GetByToken(ctx context.Context, token string) (*RefreshToken, error)
	Revoke(ctx context.Context, tokenStr string) error
}

type AuthService interface {
	Login(ctx context.Context, loginData *AuthCredentials) (*AuthResponse, error)
	Register(ctx context.Context, registerData *User, plainPassword string) (*User, error)
	RefreshToken(ctx context.Context, oldRefreshToken string) (*AuthResponse, error)
}

type AuthResponse struct {
	AccessToken  string
	RefreshToken string
	User         *User
}

func MatchesHash(password string, hash string) bool {
	if err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password)); err != nil {
		return false
	}
	return true
}

func HashPassword(plain string) (string, error) {
	if plain == "" {
		return "", ErrPasswordEmpty
	}
	hashBytes, err := bcrypt.GenerateFromPassword([]byte(plain), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashBytes), nil
}

var ErrUserAlreadyExist = errors.New("user already exist")
var ErrPasswordEmpty = errors.New("password cannot be empty")
var ErrInvalidCredentials = errors.New("invalid credentials")
var ErrInvalidToken = errors.New("invalid token")
