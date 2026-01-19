package services

import (
	"context"
	"errors"
	"log/slog"
	"time"

	"github.com/Furkanberkay/ticket-booking-project-v1/models"
	"github.com/Furkanberkay/ticket-booking-project-v1/utils"
)

type AuthService struct {
	userRepo  models.AuthRepository
	tokenRepo models.TokenRepository
	jwtUtils  *utils.JWTWrapper
	logger    *slog.Logger
}

func NewAuthService(uRepo models.AuthRepository, tRepo models.TokenRepository, jwt *utils.JWTWrapper, l *slog.Logger) *AuthService {
	return &AuthService{
		userRepo:  uRepo,
		tokenRepo: tRepo,
		jwtUtils:  jwt,
		logger:    l,
	}
}

func (s *AuthService) Login(ctx context.Context, loginData *models.AuthCredentials) (*models.AuthResponse, error) {
	user, err := s.userRepo.GetUserByEmail(ctx, loginData.Email)
	if err != nil {
		if errors.Is(err, models.ErrRecordNotFound) {
			return nil, models.ErrRecordNotFound
		}
		s.logger.Error("login_user_lookup_failed", "error", err)
		return nil, models.InternalError
	}

	if !models.MatchesHash(loginData.Password, user.Password) {
		return nil, models.ErrInvalidCredentials
	}

	accessToken, refreshTokenStr, err := s.jwtUtils.GenerateTokens(user.ID, user.Email, string(user.Role))
	if err != nil {
		s.logger.Error("token_generation_failed", "error", err)
		return nil, models.InternalError
	}

	refreshTokenModel := &models.RefreshToken{
		UserID:    user.ID,
		Token:     refreshTokenStr,
		ExpiresAt: time.Now().Add(time.Hour * 24 * 7),
	}

	if err := s.tokenRepo.Create(ctx, refreshTokenModel); err != nil {
		s.logger.Error("redis_token_save_failed", "error", err)
		return nil, models.InternalError
	}

	return &models.AuthResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshTokenStr,
		User:         user,
	}, nil
}

func (s *AuthService) Register(ctx context.Context, registerData *models.User, plainPassword string) (*models.User, error) {
	_, err := s.userRepo.GetUserByEmail(ctx, registerData.Email)
	if err == nil {
		return nil, models.ErrUserAlreadyExist
	}

	hashedPassword, err := models.HashPassword(plainPassword)
	if err != nil {
		if errors.Is(err, models.ErrPasswordEmpty) {
			return nil, models.ErrPasswordEmpty
		}
		s.logger.Error("password_hash_failed", "error", err)
		return nil, models.InternalError
	}

	registerData.Password = hashedPassword
	user, err := s.userRepo.RegisterUser(ctx, registerData)
	if err != nil {
		s.logger.Error("register_user_failed", "error", err)
		return nil, models.InternalError
	}

	return user, nil
}

func (s *AuthService) RefreshToken(ctx context.Context, oldRefreshToken string) (*models.AuthResponse, error) {
	storedToken, err := s.tokenRepo.GetByToken(ctx, oldRefreshToken)
	if err != nil {
		return nil, models.ErrInvalidToken
	}

	if err := s.tokenRepo.Revoke(ctx, oldRefreshToken); err != nil {
		s.logger.Error("token_revoke_failed", "error", err)
		return nil, models.InternalError
	}

	user, err := s.userRepo.GetUserByUserID(ctx, storedToken.UserID)
	if err != nil {
		return nil, models.ErrRecordNotFound
	}

	newAccess, newRefresh, err := s.jwtUtils.GenerateTokens(user.ID, user.Email, string(user.Role))
	if err != nil {
		s.logger.Error("token_generation_failed", "error", err)
		return nil, models.InternalError
	}

	newTokenModel := &models.RefreshToken{
		UserID:    user.ID,
		Token:     newRefresh,
		ExpiresAt: time.Now().Add(time.Hour * 24 * 7),
	}

	if err := s.tokenRepo.Create(ctx, newTokenModel); err != nil {
		s.logger.Error("redis_token_save_failed", "error", err)
		return nil, models.InternalError
	}

	return &models.AuthResponse{
		AccessToken:  newAccess,
		RefreshToken: newRefresh,
		User:         user,
	}, nil
}
