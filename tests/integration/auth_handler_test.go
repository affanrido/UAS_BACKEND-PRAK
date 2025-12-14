package integration_test

import (
	"UAS_BACKEND/domain/model"
	"UAS_BACKEND/domain/route"
	"UAS_BACKEND/tests/mocks"
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestAuthHandler_Login_Success(t *testing.T) {
	// Arrange
	app := fiber.New()
	mockAuthService := new(mocks.MockAuthService)
	
	userID := uuid.New()
	roleID := uuid.New()
	
	loginReq := &model.LoginRequest{
		Identifier: "testuser",
		Password:   "password123",
	}
	
	expectedResponse := &model.LoginResponse{
		Token: "jwt-token-here",
		User: model.UserInfo{
			ID:       userID,
			Username: "testuser",
			Email:    "test@example.com",
			FullName: "Test User",
			RoleID:   roleID,
		},
		ExpiresAt: time.Now().Add(24 * time.Hour),
	}
	
	mockAuthService.On("Login", loginReq).Return(expectedResponse, nil)
	
	authHandler := route.NewAuthHandler(mockAuthService)
	authHandler.SetupRoutes(app)
	
	// Act
	reqBody, _ := json.Marshal(loginReq)
	req := httptest.NewRequest("POST", "/api/auth/login", bytes.NewBuffer(reqBody))
	req.Header.Set("Content-Type", "application/json")
	
	resp, err := app.Test(req)
	
	// Assert
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	
	var response map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&response)
	
	assert.Equal(t, "Login successful", response["message"])
	assert.NotNil(t, response["data"])
	
	mockAuthService.AssertExpectations(t)
}

func TestAuthHandler_Login_InvalidCredentials(t *testing.T) {
	// Arrange
	app := fiber.New()
	mockAuthService := new(mocks.MockAuthService)
	
	loginReq := &model.LoginRequest{
		Identifier: "wronguser",
		Password:   "wrongpassword",
	}
	
	mockAuthService.On("Login", loginReq).Return(nil, assert.AnError)
	
	authHandler := route.NewAuthHandler(mockAuthService)
	authHandler.SetupRoutes(app)
	
	// Act
	reqBody, _ := json.Marshal(loginReq)
	req := httptest.NewRequest("POST", "/api/auth/login", bytes.NewBuffer(reqBody))
	req.Header.Set("Content-Type", "application/json")
	
	resp, err := app.Test(req)
	
	// Assert
	assert.NoError(t, err)
	assert.Equal(t, http.StatusUnauthorized, resp.StatusCode)
	
	var response map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&response)
	
	assert.Contains(t, response["error"], "Invalid credentials")
	
	mockAuthService.AssertExpectations(t)
}

func TestAuthHandler_Login_ValidationError(t *testing.T) {
	// Arrange
	app := fiber.New()
	mockAuthService := new(mocks.MockAuthService)
	
	authHandler := route.NewAuthHandler(mockAuthService)
	authHandler.SetupRoutes(app)
	
	testCases := []struct {
		name    string
		request map[string]interface{}
	}{
		{
			name: "Empty identifier",
			request: map[string]interface{}{
				"identifier": "",
				"password":   "password123",
			},
		},
		{
			name: "Empty password",
			request: map[string]interface{}{
				"identifier": "testuser",
				"password":   "",
			},
		},
		{
			name: "Missing fields",
			request: map[string]interface{}{},
		},
	}
	
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Act
			reqBody, _ := json.Marshal(tc.request)
			req := httptest.NewRequest("POST", "/api/auth/login", bytes.NewBuffer(reqBody))
			req.Header.Set("Content-Type", "application/json")
			
			resp, err := app.Test(req)
			
			// Assert
			assert.NoError(t, err)
			assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
			
			var response map[string]interface{}
			json.NewDecoder(resp.Body).Decode(&response)
			
			assert.Contains(t, response["error"], "required")
		})
	}
}

func TestAuthHandler_Login_InvalidJSON(t *testing.T) {
	// Arrange
	app := fiber.New()
	mockAuthService := new(mocks.MockAuthService)
	
	authHandler := route.NewAuthHandler(mockAuthService)
	authHandler.SetupRoutes(app)
	
	// Act
	req := httptest.NewRequest("POST", "/api/auth/login", bytes.NewBuffer([]byte("invalid json")))
	req.Header.Set("Content-Type", "application/json")
	
	resp, err := app.Test(req)
	
	// Assert
	assert.NoError(t, err)
	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	
	var response map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&response)
	
	assert.Contains(t, response["error"], "Invalid JSON")
}