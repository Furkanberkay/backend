package repositories

import (
	"context"
	"encoding/json"
	"time"

	"github.com/Furkanberkay/ticket-booking-project-v1/models"
	"github.com/redis/go-redis/v9"
)

type RedisTokenRepository struct {
	client *redis.Client
}

func NewRedisTokenRepository(client *redis.Client) models.TokenRepository {
	return &RedisTokenRepository{client: client}
}

func getTokenKey(token string) string {
	return "rt:" + token
}

func (r *RedisTokenRepository) Create(ctx context.Context, token *models.RefreshToken) error {
	data, err := json.Marshal(token)
	if err != nil {
		return err
	}

	ttl := time.Until(token.ExpiresAt)

	return r.client.Set(ctx, getTokenKey(token.Token), data, ttl).Err()
}

func (r *RedisTokenRepository) GetByToken(ctx context.Context, tokenStr string) (*models.RefreshToken, error) {
	val, err := r.client.Get(ctx, getTokenKey(tokenStr)).Result()
	if err == redis.Nil {
		return nil, models.ErrInvalidToken
	}
	if err != nil {
		return nil, err
	}

	var token models.RefreshToken
	if err := json.Unmarshal([]byte(val), &token); err != nil {
		return nil, err
	}

	return &token, nil
}

func (r *RedisTokenRepository) Revoke(ctx context.Context, tokenStr string) error {
	return r.client.Del(ctx, getTokenKey(tokenStr)).Err()
}
