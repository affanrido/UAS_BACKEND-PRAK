package middleware

import (
	"UAS_BACKEND/domain/service"
	"errors"
	"strings"

	"github.com/google/uuid"
)

func AuthMiddleware(authSvc *service.AuthService, tokenStr string) (uuid.UUID, error) {
	if tokenStr == "" {
		return uuid.Nil, errors.New("missing token")
	}

	// parsing "Bearer xxx"
	if strings.HasPrefix(tokenStr, "Bearer ") {
		tokenStr = strings.TrimPrefix(tokenStr, "Bearer ")
	}

	claims, err := authSvc.ParseToken(tokenStr)
	if err != nil {
		return uuid.Nil, err
	}

	// Kembalikan UUID user
	return claims.UserID, nil
}
