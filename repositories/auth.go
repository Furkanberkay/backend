package repositories

import (
	"context"
	"errors"
	"log/slog"

	"github.com/Furkanberkay/ticket-booking-project-v1/models"
	"gorm.io/gorm"
)

type AuthRepository struct {
	db     *gorm.DB
	logger *slog.Logger
}

func NewAuthRepository(db *gorm.DB, logger *slog.Logger) models.AuthRepository {
	if logger == nil {
		logger = slog.Default()
	}
	return &AuthRepository{
		db:     db,
		logger: logger,
	}
}

func (r *AuthRepository) RegisterUser(ctx context.Context, user *models.User) (*models.User, error) {
	err := r.db.WithContext(ctx).Create(user).Error
	if err != nil {
		r.logger.ErrorContext(ctx, "auth_repo_register_user_failed",
			"email", user.Email,
			"error", err,
		)
		return nil, err
	}

	return user, nil
}

func (r *AuthRepository) GetUserByEmail(ctx context.Context, email string) (*models.User, error) {
	user := new(models.User)

	err := r.db.WithContext(ctx).Where("email = ?", email).First(user).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			r.logger.DebugContext(ctx, "auth_repo_user_not_found_by_email", "email", email)
			return nil, models.ErrRecordNotFound
		}

		r.logger.ErrorContext(ctx, "auth_repo_get_user_by_email_failed", "email", email, "error", err)
		return nil, err
	}

	return user, nil
}

func (r *AuthRepository) GetUserByUserID(ctx context.Context, userID uint) (*models.User, error) {
	user := new(models.User)

	err := r.db.WithContext(ctx).First(user, userID).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			r.logger.DebugContext(ctx, "auth_repo_user_not_found_by_id", "user_id", userID)
			return nil, models.ErrRecordNotFound
		}

		r.logger.ErrorContext(ctx, "auth_repo_get_user_by_id_failed", "user_id", userID, "error", err)
		return nil, err
	}

	return user, nil
}
