package middleware_test

import (
	"UAS_BACKEND/domain/middleware"
	"UAS_BACKEND/domain/model"
	"UAS_BACKEND/tests/mocks"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestAuthMiddleware_ValidToken_Success(t *testing.T) {
	// Arrange
	app := fiber.New()
	mockAuthService := new(mocks.MockAuthService)
	
	userID := uuid.New()
	roleID := uuid.New()
	
	claims := &model.Claims{
		UserID:      userID,
		RoleID:      roleID,
		Permissions: []string{"user.read", "user.write"},
		ExpiresAt:   time.Now().Add(time.Hour),
	}
	
	mockAuthService.On("ValidateToken", "valid-token").Return(claims, nil)
	
	authMiddleware := middleware.NewAuthMiddleware(mockAuthService)
	
	app.Use(authMiddleware.RequireAuth())
	app.Get("/protected", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{"message": "success"})
	})
	
	// Act
	req := httptest.NewRequest("GET", "/protected", nil)
	req.Header.Set("Authorization", "Bearer valid-token")
	
	resp, err := app.Test(req)
	
	// Assert
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	
	mockAuthService.AssertExpectations(t)
}

func TestAuthMiddleware_MissingToken_Unauthorized(t *testing.T) {
	// Arrange
	app := fiber.New()
	mockAuthService := new(mocks.MockAuthService)
	
	authMiddleware := middleware.NewAuthMiddleware(mockAuthService)
	
	app.Use(authMiddleware.RequireAuth())
	app.Get("/protected", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{"message": "success"})
	})
	
	// Act
	req := httptest.NewRequest("GET", "/protected", nil)
	
	resp, err := app.Test(req)
	
	// Assert
	assert.NoError(t, err)
	assert.Equal(t, http.StatusUnauthorized, resp.StatusCode)
}

func TestAuthMiddleware_InvalidToken_Unauthorized(t *testing.T) {
	// Arrange
	app := fiber.New()
	mockAuthService := new(mocks.MockAuthService)
	
	mockAuthService.On("ValidateToken", "invalid-token").Return(nil, assert.AnError)
	
	authMiddleware := middleware.NewAuthMiddleware(mockAuthService)
	
	app.Use(authMiddleware.RequireAuth())
	app.Get("/protected", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{"message": "success"})
	})
	
	// Act
	req := httptest.NewRequest("GET", "/protected", nil)
	req.Header.Set("Authorization", "Bearer invalid-token")
	
	resp, err := app.Test(req)
	
	// Assert
	assert.NoError(t, err)
	assert.Equal(t, http.StatusUnauthorized, resp.StatusCode)
	
	mockAuthService.AssertExpectations(t)
}

func TestAuthMiddleware_MalformedAuthHeader_Unauthorized(t *testing.T) {
	// Arrange
	app := fiber.New()
	mockAuthService := new(mocks.MockAuthService)
	
	authMiddleware := middleware.NewAuthMiddleware(mockAuthService)
	
	app.Use(authMiddleware.RequireAuth())
	app.Get("/protected", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{"message": "success"})
	})
	
	testCases := []struct {
		name   string
		header string
	}{
		{"No Bearer prefix", "invalid-token"},
		{"Only Bearer", "Bearer"},
		{"Empty Bearer", "Bearer "},
	}
	
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Act
			req := httptest.NewRequest("GET", "/protected", nil)
			req.Header.Set("Authorization", tc.header)
			
			resp, err := app.Test(req)
			
			// Assert
			assert.NoError(t, err)
			assert.Equal(t, http.StatusUnauthorized, resp.StatusCode)
		})
	}
}