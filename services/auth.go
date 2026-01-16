package services

import (
	"context"

	"github.com/Furkanberkay/ticket-booking-project-v1/models"
)

type AuthService struct {
	repository models.AuthRepository
}

func NewAuthRepository(repository models.AuthRepository) models.AuthService {
	return &AuthService{
		repository: repository,
	}
}

func (s *AuthService) Login(ctx context.Context, loginData *models.AuthCredentials) (string, *models.User, error) {

}
func (s *AuthService) Register(ctx context.Context, registerData *models.User, plainPassword string) (string, *models.User, error) {

}
