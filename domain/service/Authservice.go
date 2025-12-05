package service

import (
	model "UAS_BACKEND/domain/Model"
	"UAS_BACKEND/domain/repository"
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

type AuthService struct {
	SecretKey []byte
	TokenTTL  time.Duration
	Repo      *repository.AuthRepository
}

func NewAuthService(secretKey string, tokenTTL time.Duration, repo *repository.AuthRepository) *AuthService {
	return &AuthService{
		SecretKey: []byte(secretKey),
		TokenTTL:  tokenTTL,
		Repo:      repo,
	}
}

func HashPassword(password string) (string, error) {
	b, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(b), err
}

func CheckPassword(hash, password string) error {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
}

func (a *AuthService) Login(identifier, password string) (*model.LoginResponse, error) {
	// 1. Validasi kredensial
	user, err := a.Repo.GetUserByIdentifier(identifier)
	if err != nil || user == nil {
		return nil, errors.New("invalid credentials")
	}

	// 2. Cek password
	if err := CheckPassword(user.PasswordHash, password); err != nil {
		return nil, errors.New("invalid credentials")
	}

	// 3. Cek status aktif user
	if !user.IsActive {
		return nil, errors.New("user account is inactive")
	}

	// 4. Load permissions dari RBAC
	permissions, err := a.Repo.GetUserPermissions(user.RoleID)
	if err != nil {
		permissions = []string{} // fallback jika error
	}

	// 5. Generate JWT token dengan role dan permissions
	claims := model.CustomClaims{
		UserID:      user.ID,
		RoleID:      user.RoleID,
		Permissions: permissions,
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

	// 6. Return token dan user profile
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
