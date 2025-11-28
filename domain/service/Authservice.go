package service

import (
	model "UAS_BACKEND/domain/Model"
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

type AuthService struct {
	SecretKey []byte
	TokenTTL  time.Duration

	GetUserByIdentifier func(identifier string) (*model.Users, error)
}

func HashPassword(password string) (string, error) {
	b, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(b), err
}

func CheckPassword(hash, password string) error {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
}

func (a *AuthService) Login(identifier, password string) (*model.LoginResponse, error) {
	user, err := a.GetUserByIdentifier(identifier)
	if err != nil || user == nil {
		return nil, errors.New("invalid credentials")
	}

	if err := CheckPassword(user.PasswordHash, password); err != nil {
		return nil, errors.New("invalid credentials")
	}

	claims := model.CustomClaims{
		UserID:      user.ID,
		RoleID:      user.RoleID,
		Permissions: []string{}, // nanti dari RBAC
		RegisteredClaims: jwt.RegisteredClaims{
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(a.TokenTTL)),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signed, err := token.SignedString(a.SecretKey)
	if err != nil {
		return nil, err
	}

	return &model.LoginResponse{
		Token: signed,
		User:  *user,
	}, nil
}

func (a *AuthService) ParseToken(tokenStr string) (*model.CustomClaims, error) {
	t, err := jwt.ParseWithClaims(tokenStr, &model.CustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		return a.SecretKey, nil
	})

	if err != nil {
		return nil, err
	}

	claims, ok := t.Claims.(*model.CustomClaims)
	if !ok || !t.Valid {
		return nil, errors.New("invalid token")
	}

	return claims, nil
}
