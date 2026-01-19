package utils

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type JWTWrapper struct {
	SecretKey       string
	Issuer          string
	ExpirationHours int64
}

func NewJWTWrapper(secretKey string) *JWTWrapper {
	return &JWTWrapper{
		SecretKey:       secretKey,
		Issuer:          "ticket-booking-app",
		ExpirationHours: 1,
	}
}

func (j *JWTWrapper) GenerateTokens(userID uint, email string, role string) (accessToken string, refreshToken string, err error) {
	claims := jwt.MapClaims{
		"sub":   userID,
		"email": email,
		"role":  role,
		"iss":   j.Issuer,
		"iat":   time.Now().Unix(),
		"exp":   time.Now().Add(time.Hour * time.Duration(j.ExpirationHours)).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	accessToken, err = token.SignedString([]byte(j.SecretKey))
	if err != nil {
		return "", "", fmt.Errorf("could not sign token: %w", err)
	}

	refreshToken, err = j.generateRandomHex(32)
	if err != nil {
		return "", "", fmt.Errorf("could not generate refresh token: %w", err)
	}

	return accessToken, refreshToken, nil
}

func (j *JWTWrapper) ValidateToken(tokenString string) (*jwt.Token, error) {
	return jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(j.SecretKey), nil
	})
}

func (j *JWTWrapper) generateRandomHex(n int) (string, error) {
	bytes := make([]byte, n)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}
