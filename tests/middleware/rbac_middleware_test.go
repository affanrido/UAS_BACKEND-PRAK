package middleware_test

import (
	"UAS_BACKEND/domain/middleware"
	"UAS_BACKEND/domain/model"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestRBACMiddleware_HasPermission_Success(t *testing.T) {
	// Arrange
	app := fiber.New()
	
	userID := uuid.New()
	roleID := uuid.New()
	
	// Mock user with required permission
	app.Use(func(c *fiber.Ctx) error {
		claims := &model.Claims{
			UserID:      userID,
			RoleID:      roleID,
			Permissions: []string{"user.read", "user.write", "admin.manage"},
			ExpiresAt:   time.Now().Add(time.Hour),
		}
		c.Locals("user", claims)
		return c.Next()
	})
	
	rbacMiddleware := middleware.NewRBACMiddleware()
	
	app.Get("/admin", rbacMiddleware.RequirePermission("admin.manage"), func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{"message": "success"})
	})
	
	// Act
	req := httptest.NewRequest("GET", "/admin", nil)
	resp, err := app.Test(req)
	
	// Assert
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
}

func TestRBACMiddleware_MissingPermission_Forbidden(t *testing.T) {
	// Arrange
	app := fiber.New()
	
	userID := uuid.New()
	roleID := uuid.New()
	
	// Mock user without required permission
	app.Use(func(c *fiber.Ctx) error {
		claims := &model.Claims{
			UserID:      userID,
			RoleID:      roleID,
			Permissions: []string{"user.read", "user.write"}, // Missing admin.manage
			ExpiresAt:   time.Now().Add(time.Hour),
		}
		c.Locals("user", claims)
		return c.Next()
	})
	
	rbacMiddleware := middleware.NewRBACMiddleware()
	
	app.Get("/admin", rbacMiddleware.RequirePermission("admin.manage"), func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{"message": "success"})
	})
	
	// Act
	req := httptest.NewRequest("GET", "/admin", nil)
	resp, err := app.Test(req)
	
	// Assert
	assert.NoError(t, err)
	assert.Equal(t, http.StatusForbidden, resp.StatusCode)
}

func TestRBACMiddleware_NoUserInContext_Unauthorized(t *testing.T) {
	// Arrange
	app := fiber.New()
	
	rbacMiddleware := middleware.NewRBACMiddleware()
	
	app.Get("/admin", rbacMiddleware.RequirePermission("admin.manage"), func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{"message": "success"})
	})
	
	// Act
	req := httptest.NewRequest("GET", "/admin", nil)
	resp, err := app.Test(req)
	
	// Assert
	assert.NoError(t, err)
	assert.Equal(t, http.StatusUnauthorized, resp.StatusCode)
}

func TestRBACMiddleware_RequireRole_Success(t *testing.T) {
	// Arrange
	app := fiber.New()
	
	userID := uuid.New()
	adminRoleID := uuid.New()
	
	// Mock user with admin role
	app.Use(func(c *fiber.Ctx) error {
		claims := &model.Claims{
			UserID:      userID,
			RoleID:      adminRoleID,
			Permissions: []string{"admin.manage"},
			ExpiresAt:   time.Now().Add(time.Hour),
		}
		c.Locals("user", claims)
		c.Locals("userRole", "admin") // Mock role name
		return c.Next()
	})
	
	rbacMiddleware := middleware.NewRBACMiddleware()
	
	app.Get("/admin", rbacMiddleware.RequireRole("admin"), func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{"message": "success"})
	})
	
	// Act
	req := httptest.NewRequest("GET", "/admin", nil)
	resp, err := app.Test(req)
	
	// Assert
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
}

func TestRBACMiddleware_RequireRole_WrongRole_Forbidden(t *testing.T) {
	// Arrange
	app := fiber.New()
	
	userID := uuid.New()
	studentRoleID := uuid.New()
	
	// Mock user with student role
	app.Use(func(c *fiber.Ctx) error {
		claims := &model.Claims{
			UserID:      userID,
			RoleID:      studentRoleID,
			Permissions: []string{"user.read"},
			ExpiresAt:   time.Now().Add(time.Hour),
		}
		c.Locals("user", claims)
		c.Locals("userRole", "student") // Mock role name
		return c.Next()
	})
	
	rbacMiddleware := middleware.NewRBACMiddleware()
	
	app.Get("/admin", rbacMiddleware.RequireRole("admin"), func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{"message": "success"})
	})
	
	// Act
	req := httptest.NewRequest("GET", "/admin", nil)
	resp, err := app.Test(req)
	
	// Assert
	assert.NoError(t, err)
	assert.Equal(t, http.StatusForbidden, resp.StatusCode)
}

func TestRBACMiddleware_MultiplePermissions_Success(t *testing.T) {
	// Arrange
	app := fiber.New()
	
	userID := uuid.New()
	roleID := uuid.New()
	
	// Mock user with multiple permissions
	app.Use(func(c *fiber.Ctx) error {
		claims := &model.Claims{
			UserID:      userID,
			RoleID:      roleID,
			Permissions: []string{"user.read", "user.write", "achievement.manage", "admin.view"},
			ExpiresAt:   time.Now().Add(time.Hour),
		}
		c.Locals("user", claims)
		return c.Next()
	})
	
	rbacMiddleware := middleware.NewRBACMiddleware()
	
	// Test multiple permission requirements
	testCases := []struct {
		name       string
		permission string
		expected   int
	}{
		{"Has user.read", "user.read", http.StatusOK},
		{"Has achievement.manage", "achievement.manage", http.StatusOK},
		{"Missing admin.manage", "admin.manage", http.StatusForbidden},
	}
	
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			app.Get("/test/"+tc.name, rbacMiddleware.RequirePermission(tc.permission), func(c *fiber.Ctx) error {
				return c.JSON(fiber.Map{"message": "success"})
			})
			
			// Act
			req := httptest.NewRequest("GET", "/test/"+tc.name, nil)
			resp, err := app.Test(req)
			
			// Assert
			assert.NoError(t, err)
			assert.Equal(t, tc.expected, resp.StatusCode)
		})
	}
}