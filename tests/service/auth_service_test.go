package service_test

import (
	model "UAS_BACKEND/domain/model"
	"UAS_BACKEND/domain/service"
	"UAS_BACKEND/tests/mocks"
	"errors"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"golang.org/x/crypto/bcrypt"
)

func TestAuthService_Login_Success(t *testing.T) {
	// Arrange
	mockRepo := new(mocks.MockAuthRepository)
	authService := service.NewAuthService("test-secret", 24*time.Hour, mockRepo)

	userID := uuid.New()
	roleID := uuid.New()
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("password123"), bcrypt.DefaultCost)

	expectedUser := &model.Users{
		ID:           userID,
		Username:     "testuser",
		Email:        "test@example.com",
		PasswordHash: string(hashedPassword),
		FullName:     "Test User",
		RoleID:       roleID,
		IsActive:     true,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}

	expectedPermissions := []string{"user.read", "user.write"}

	mockRepo.On("GetUserByIdentifier", "testuser").Return(expectedUser, nil)
	mockRepo.On("GetUserPermissions", userID).Return(expectedPermissions, nil)

	loginReq := &model.LoginRequest{
		Identifier: "testuser",
		Password:   "password123",
	}

	// Act
	result, err := authService.Login(loginReq)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.NotEmpty(t, result.Token)
	assert.Equal(t, expectedUser.ID, result.User.ID)
	assert.Equal(t, expectedUser.Username, result.User.Username)
	assert.Equal(t, expectedUser.Email, result.User.Email)

	mockRepo.AssertExpectations(t)
}

func TestAuthService_Login_InvalidCredentials(t *testing.T) {
	// Arrange
	mockRepo := new(mocks.MockAuthRepository)
	authService := service.NewAuthService("test-secret", 24*time.Hour, mockRepo)

	mockRepo.On("GetUserByIdentifier", "invaliduser").Return(nil, errors.New("user not found"))

	loginReq := &model.LoginRequest{
		Identifier: "invaliduser",
		Password:   "wrongpassword",
	}

	// Act
	result, err := authService.Login(loginReq)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "invalid credentials")

	mockRepo.AssertExpectations(t)
}

func TestAuthService_Login_WrongPassword(t *testing.T) {
	// Arrange
	mockRepo := new(mocks.MockAuthRepository)
	authService := service.NewAuthService("test-secret", 24*time.Hour, mockRepo)

	userID := uuid.New()
	roleID := uuid.New()
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("correctpassword"), bcrypt.DefaultCost)

	expectedUser := &model.Users{
		ID:           userID,
		Username:     "testuser",
		Email:        "test@example.com",
		PasswordHash: string(hashedPassword),
		FullName:     "Test User",
		RoleID:       roleID,
		IsActive:     true,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}

	mockRepo.On("GetUserByIdentifier", "testuser").Return(expectedUser, nil)

	loginReq := &model.LoginRequest{
		Identifier: "testuser",
		Password:   "wrongpassword",
	}

	// Act
	result, err := authService.Login(loginReq)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "invalid credentials")

	mockRepo.AssertExpectations(t)
}

func TestAuthService_Login_InactiveUser(t *testing.T) {
	// Arrange
	mockRepo := new(mocks.MockAuthRepository)
	authService := service.NewAuthService("test-secret", 24*time.Hour, mockRepo)

	userID := uuid.New()
	roleID := uuid.New()
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("password123"), bcrypt.DefaultCost)

	expectedUser := &model.Users{
		ID:           userID,
		Username:     "testuser",
		Email:        "test@example.com",
		PasswordHash: string(hashedPassword),
		FullName:     "Test User",
		RoleID:       roleID,
		IsActive:     false, // User is inactive
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}

	mockRepo.On("GetUserByIdentifier", "testuser").Return(expectedUser, nil)

	loginReq := &model.LoginRequest{
		Identifier: "testuser",
		Password:   "password123",
	}

	// Act
	result, err := authService.Login(loginReq)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "user account is inactive")

	mockRepo.AssertExpectations(t)
}

func TestAuthService_ValidateToken_Success(t *testing.T) {
	// Arrange
	mockRepo := new(mocks.MockAuthRepository)
	authService := service.NewAuthService("test-secret", 24*time.Hour, mockRepo)

	userID := uuid.New()
	roleID := uuid.New()
	permissions := []string{"user.read", "user.write"}

	// First create a valid token
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("password123"), bcrypt.DefaultCost)
	user := &model.Users{
		ID:           userID,
		Username:     "testuser",
		Email:        "test@example.com",
		PasswordHash: string(hashedPassword),
		FullName:     "Test User",
		RoleID:       roleID,
		IsActive:     true,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}

	mockRepo.On("GetUserByIdentifier", "testuser").Return(user, nil)
	mockRepo.On("GetUserPermissions", userID).Return(permissions, nil)

	loginReq := &model.LoginRequest{
		Identifier: "testuser",
		Password:   "password123",
	}

	loginResult, _ := authService.Login(loginReq)

	// Act
	claims, err := authService.ValidateToken(loginResult.Token)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, claims)
	assert.Equal(t, userID, claims.UserID)
	assert.Equal(t, roleID, claims.RoleID)
	assert.Equal(t, permissions, claims.Permissions)

	mockRepo.AssertExpectations(t)
}

func TestAuthService_ValidateToken_InvalidToken(t *testing.T) {
	// Arrange
	mockRepo := new(mocks.MockAuthRepository)
	authService := service.NewAuthService("test-secret", 24*time.Hour, mockRepo)

	invalidToken := "invalid.token.here"

	// Act
	claims, err := authService.ValidateToken(invalidToken)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, claims)
}

func TestAuthService_ValidateToken_ExpiredToken(t *testing.T) {
	// Arrange
	mockRepo := new(mocks.MockAuthRepository)
	authService := service.NewAuthService("test-secret", -1*time.Hour, mockRepo) // Negative TTL for expired token

	userID := uuid.New()
	roleID := uuid.New()
	permissions := []string{"user.read"}

	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("password123"), bcrypt.DefaultCost)
	user := &model.Users{
		ID:           userID,
		Username:     "testuser",
		Email:        "test@example.com",
		PasswordHash: string(hashedPassword),
		FullName:     "Test User",
		RoleID:       roleID,
		IsActive:     true,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}

	mockRepo.On("GetUserByIdentifier", "testuser").Return(user, nil)
	mockRepo.On("GetUserPermissions", userID).Return(permissions, nil)

	loginReq := &model.LoginRequest{
		Identifier: "testuser",
		Password:   "password123",
	}

	loginResult, _ := authService.Login(loginReq)

	// Act
	claims, err := authService.ValidateToken(loginResult.Token)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, claims)
	assert.Contains(t, err.Error(), "token is expired")

	mockRepo.AssertExpectations(t)
}

func TestAuthService_Login_EmptyCredentials(t *testing.T) {
	// Arrange
	mockRepo := new(mocks.MockAuthRepository)
	authService := service.NewAuthService("test-secret", 24*time.Hour, mockRepo)

	testCases := []struct {
		name        string
		loginReq    *model.LoginRequest
		expectedErr string
	}{
		{
			name: "Empty identifier",
			loginReq: &model.LoginRequest{
				Identifier: "",
				Password:   "password123",
			},
			expectedErr: "Username/email and password are required",
		},
		{
			name: "Empty password",
			loginReq: &model.LoginRequest{
				Identifier: "testuser",
				Password:   "",
			},
			expectedErr: "Username/email and password are required",
		},
		{
			name: "Both empty",
			loginReq: &model.LoginRequest{
				Identifier: "",
				Password:   "",
			},
			expectedErr: "Username/email and password are required",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Act
			result, err := authService.Login(tc.loginReq)

			// Assert
			assert.Error(t, err)
			assert.Nil(t, result)
			assert.Contains(t, err.Error(), tc.expectedErr)
		})
	}
}