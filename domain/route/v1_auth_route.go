package route

import (
	"UAS_BACKEND/domain/middleware"
	"UAS_BACKEND/domain/model"
	"UAS_BACKEND/domain/service"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type V1AuthHandler struct {
	AuthService     *service.AuthService
	RBACMiddleware  *middleware.RBACMiddleware
}

func NewV1AuthHandler(authService *service.AuthService, rbacMiddleware *middleware.RBACMiddleware) *V1AuthHandler {
	return &V1AuthHandler{
		AuthService:    authService,
		RBACMiddleware: rbacMiddleware,
	}
}

// SetupV1AuthRoutes - Setup authentication routes v1
func SetupV1AuthRoutes(app *fiber.App, handler *V1AuthHandler) {
	auth := app.Group("/api/v1/auth")

	// 5.1 Authentication endpoints
	auth.Post("/login", handler.Login)
	auth.Post("/refresh", handler.RefreshToken)
	auth.Post("/logout", handler.RBACMiddleware.RequireAuth(), handler.Logout)
	auth.Get("/profile", handler.RBACMiddleware.RequireAuth(), handler.GetProfile)
}

// Login - POST /api/v1/auth/login
func (h *V1AuthHandler) Login(c *fiber.Ctx) error {
	var req model.LoginRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": "Invalid JSON format",
		})
	}

	// Validate required fields
	if req.Identifier == "" || req.Password == "" {
		return c.Status(400).JSON(fiber.Map{
			"error": "Username/email and password are required",
		})
	}

	// Authenticate user
	result, err := h.AuthService.Login(&req)
	if err != nil {
		return c.Status(401).JSON(fiber.Map{
			"error": "Invalid credentials",
		})
	}

	return c.Status(200).JSON(fiber.Map{
		"success": true,
		"message": "Login successful",
		"data": fiber.Map{
			"token":      result.Token,
			"user":       result.User,
			"expires_at": result.ExpiresAt,
		},
	})
}

// RefreshToken - POST /api/v1/auth/refresh
func (h *V1AuthHandler) RefreshToken(c *fiber.Ctx) error {
	var req struct {
		RefreshToken string `json:"refresh_token"`
	}

	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": "Invalid JSON format",
		})
	}

	if req.RefreshToken == "" {
		return c.Status(400).JSON(fiber.Map{
			"error": "Refresh token is required",
		})
	}

	// Validate refresh token
	claims, err := h.AuthService.ValidateToken(req.RefreshToken)
	if err != nil {
		return c.Status(401).JSON(fiber.Map{
			"error": "Invalid refresh token",
		})
	}

	// Generate new access token
	newToken, err := h.AuthService.GenerateToken(claims.UserID, claims.RoleID, claims.Permissions)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": "Failed to generate new token",
		})
	}

	return c.Status(200).JSON(fiber.Map{
		"success": true,
		"message": "Token refreshed successfully",
		"data": fiber.Map{
			"token":      newToken,
			"expires_at": time.Now().Add(24 * time.Hour),
		},
	})
}

// Logout - POST /api/v1/auth/logout
func (h *V1AuthHandler) Logout(c *fiber.Ctx) error {
	// Get user from context
	user := c.Locals("user").(*model.Claims)
	
	// In a real implementation, you would:
	// 1. Add token to blacklist
	// 2. Clear refresh token from database
	// 3. Log the logout event
	
	return c.Status(200).JSON(fiber.Map{
		"success": true,
		"message": "Logged out successfully",
		"data": fiber.Map{
			"user_id": user.UserID,
		},
	})
}

// GetProfile - GET /api/v1/auth/profile
func (h *V1AuthHandler) GetProfile(c *fiber.Ctx) error {
	// Get user from context
	user := c.Locals("user").(*model.Claims)

	return c.Status(200).JSON(fiber.Map{
		"success": true,
		"message": "Profile retrieved successfully",
		"data": fiber.Map{
			"user_id":     user.UserID,
			"role_id":     user.RoleID,
			"permissions": user.Permissions,
			"expires_at":  user.ExpiresAt,
		},
	})
}