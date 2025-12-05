package middleware

import (
	model "UAS_BACKEND/domain/Model"
	"UAS_BACKEND/domain/service"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

// RBACMiddleware - Middleware untuk RBAC (Role-Based Access Control)
type RBACMiddleware struct {
	AuthService *service.AuthService
	RBACService *service.RBACService
}

func NewRBACMiddleware(authService *service.AuthService, rbacService *service.RBACService) *RBACMiddleware {
	return &RBACMiddleware{
		AuthService: authService,
		RBACService: rbacService,
	}
}

// Authenticate - Middleware untuk autentikasi (ekstrak dan validasi JWT)
func (m *RBACMiddleware) Authenticate() fiber.Handler {
	return func(c *fiber.Ctx) error {
		// 1. Ekstrak JWT dari header
		authHeader := c.Get("Authorization")
		if authHeader == "" {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Missing authorization header",
			})
		}

		// Parse "Bearer <token>"
		tokenStr := authHeader
		if strings.HasPrefix(authHeader, "Bearer ") {
			tokenStr = strings.TrimPrefix(authHeader, "Bearer ")
		}

		// 2. Validasi token
		claims, err := m.AuthService.ParseToken(tokenStr)
		if err != nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Invalid or expired token",
			})
		}

		// Simpan claims di context untuk digunakan di handler
		c.Locals("user_id", claims.UserID)
		c.Locals("role_id", claims.RoleID)
		c.Locals("permissions", claims.Permissions)
		c.Locals("claims", claims)

		return c.Next()
	}
}

// RequirePermission - Middleware untuk check permission spesifik
func (m *RBACMiddleware) RequirePermission(requiredPermission string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Get permissions dari context (sudah di-set oleh Authenticate middleware)
		permissions, ok := c.Locals("permissions").([]string)
		if !ok {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Unauthorized: No permissions found",
			})
		}

		// 4. Check apakah user memiliki permission yang diperlukan
		hasPermission := false
		for _, perm := range permissions {
			if perm == requiredPermission {
				hasPermission = true
				break
			}
		}

		// 5. Allow/deny request
		if !hasPermission {
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
				"error":    "Forbidden: Insufficient permissions",
				"required": requiredPermission,
			})
		}

		return c.Next()
	}
}

// RequireAnyPermission - Middleware untuk check salah satu dari beberapa permissions
func (m *RBACMiddleware) RequireAnyPermission(requiredPermissions ...string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		permissions, ok := c.Locals("permissions").([]string)
		if !ok {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Unauthorized: No permissions found",
			})
		}

		// Check apakah user memiliki salah satu permission
		hasPermission := false
		for _, userPerm := range permissions {
			for _, reqPerm := range requiredPermissions {
				if userPerm == reqPerm {
					hasPermission = true
					break
				}
			}
			if hasPermission {
				break
			}
		}

		if !hasPermission {
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
				"error":    "Forbidden: Insufficient permissions",
				"required": requiredPermissions,
			})
		}

		return c.Next()
	}
}

// RequireAllPermissions - Middleware untuk check semua permissions harus ada
func (m *RBACMiddleware) RequireAllPermissions(requiredPermissions ...string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		permissions, ok := c.Locals("permissions").([]string)
		if !ok {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Unauthorized: No permissions found",
			})
		}

		// Check apakah user memiliki semua permissions
		permMap := make(map[string]bool)
		for _, perm := range permissions {
			permMap[perm] = true
		}

		missingPermissions := []string{}
		for _, reqPerm := range requiredPermissions {
			if !permMap[reqPerm] {
				missingPermissions = append(missingPermissions, reqPerm)
			}
		}

		if len(missingPermissions) > 0 {
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
				"error":   "Forbidden: Missing required permissions",
				"missing": missingPermissions,
			})
		}

		return c.Next()
	}
}

// RequireRole - Middleware untuk check role spesifik
func (m *RBACMiddleware) RequireRole(roleName string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		roleID, ok := c.Locals("role_id").(uuid.UUID)
		if !ok {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Unauthorized: No role found",
			})
		}

		// Get role name dari database
		role, err := m.RBACService.GetRoleByID(roleID)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Failed to get role information",
			})
		}

		if role.Name != roleName {
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
				"error":    "Forbidden: Insufficient role",
				"required": roleName,
				"current":  role.Name,
			})
		}

		return c.Next()
	}
}

// GetUserID - Helper untuk mendapatkan user ID dari context
func GetUserID(c *fiber.Ctx) (uuid.UUID, error) {
	userID, ok := c.Locals("user_id").(uuid.UUID)
	if !ok {
		return uuid.Nil, fiber.NewError(fiber.StatusUnauthorized, "User ID not found in context")
	}
	return userID, nil
}

// GetClaims - Helper untuk mendapatkan claims dari context
func GetClaims(c *fiber.Ctx) (*model.CustomClaims, error) {
	claims, ok := c.Locals("claims").(*model.CustomClaims)
	if !ok {
		return nil, fiber.NewError(fiber.StatusUnauthorized, "Claims not found in context")
	}
	return claims, nil
}

// GetPermissions - Helper untuk mendapatkan permissions dari context
func GetPermissions(c *fiber.Ctx) ([]string, error) {
	permissions, ok := c.Locals("permissions").([]string)
	if !ok {
		return nil, fiber.NewError(fiber.StatusUnauthorized, "Permissions not found in context")
	}
	return permissions, nil
}
