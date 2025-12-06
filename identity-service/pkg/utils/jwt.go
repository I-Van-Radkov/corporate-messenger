package utils

import (
	"fmt"
	"time"

	"github.com/I-Van-Radkov/corporate-messenger/identity-service/internal/models"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

func SignToken(userId uuid.UUID, role models.AccountRole, jwtSecret string, jwtExpiresIn time.Duration) (string, error) {
	now := time.Now()

	claims := jwt.MapClaims{
		"user_id":   userId,
		"user_role": role,
		"iat":       now.Unix(),
		"nbf":       now.Unix(),
		"exp":       now.Add(jwtExpiresIn).Unix(),
	}

	tokenString := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	signedToken, err := tokenString.SignedString(jwtSecret)
	if err != nil {
		return "", fmt.Errorf("failed to sign token: %w", err)
	}

	return signedToken, nil
}
