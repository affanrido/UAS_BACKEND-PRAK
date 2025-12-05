package route

import (
	model "UAS_BACKEND/domain/Model"
	"UAS_BACKEND/domain/service"

	"github.com/gofiber/fiber/v2"
)

type AuthHandler struct {
	AuthService *service.AuthService
}

func NewAuthHandler(authService *service.AuthService) *AuthHandler {
	return &AuthHandler{AuthService: authService}
}

// LoginHandler handles POST /api/auth/login
func (h *AuthHandler) LoginHandler(c *fiber.Ctx) error {
	var req model.LoginRequest

	// Parse request body
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	// Validasi input
	if req.Identifier == "" || req.Password == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Username/email and password are required",
		})
	}

	// Execute login
	resp, err := h.AuthService.Login(req.Identifier, req.Password)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	// Return success response
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Login successful",
		"data":    resp,
	})
}

// SetupAuthRoutes registers auth routes
func SetupAuthRoutes(app *fiber.App, handler *AuthHandler) {
	auth := app.Group("/api/auth")
	auth.Post("/login", handler.LoginHandler)
}
