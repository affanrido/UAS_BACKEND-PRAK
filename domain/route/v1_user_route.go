package route

import (
	"UAS_BACKEND/domain/middleware"
	"UAS_BACKEND/domain/model"
	"UAS_BACKEND/domain/service"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type V1UserHandler struct {
	UserService    *service.UserService
	RBACMiddleware *middleware.RBACMiddleware
}

func NewV1UserHandler(userService *service.UserService, rbacMiddleware *middleware.RBACMiddleware) *V1UserHandler {
	return &V1UserHandler{
		UserService:    userService,
		RBACMiddleware: rbacMiddleware,
	}
}

// SetupV1UserRoutes - Setup user management routes v1
func SetupV1UserRoutes(app *fiber.App, handler *V1UserHandler) {
	users := app.Group("/api/v1/users")
	users.Use(handler.RBACMiddleware.RequireAuth())

	// 5.2 Users (Admin) endpoints
	users.Get("/", handler.RBACMiddleware.RequirePermission("admin.manage"), handler.GetAllUsers)
	users.Get("/:id", handler.RBACMiddleware.RequirePermission("admin.manage"), handler.GetUserByID)
	users.Post("/", handler.RBACMiddleware.RequirePermission("admin.manage"), handler.CreateUser)
	users.Put("/:id", handler.RBACMiddleware.RequirePermission("admin.manage"), handler.UpdateUser)
	users.Delete("/:id", handler.RBACMiddleware.RequirePermission("admin.manage"), handler.DeleteUser)
	users.Put("/:id/role", handler.RBACMiddleware.RequirePermission("admin.manage"), handler.AssignRole)
}

// GetAllUsers - GET /api/v1/users
func (h *V1UserHandler) GetAllUsers(c *fiber.Ctx) error {
	// Parse query parameters
	page, _ := strconv.Atoi(c.Query("page", "1"))
	limit, _ := strconv.Atoi(c.Query("limit", "10"))
	role := c.Query("role")
	isActiveStr := c.Query("is_active")

	// Calculate offset
	offset := (page - 1) * limit

	// Get users
	users, total, err := h.UserService.GetAllUsers(limit, offset)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"success": false,
			"error":   "Failed to retrieve users",
		})
	}

	// Filter by role if specified
	if role != "" {
		filteredUsers := []service.UserResponse{}
		for _, user := range users {
			if user.Role != nil && user.Role.Name == role {
				filteredUsers = append(filteredUsers, user)
			}
		}
		users = filteredUsers
		total = len(users)
	}

	// Filter by active status if specified
	if isActiveStr != "" {
		isActive := isActiveStr == "true"
		filteredUsers := []service.UserResponse{}
		for _, user := range users {
			if user.User.IsActive == isActive {
				filteredUsers = append(filteredUsers, user)
			}
		}
		users = filteredUsers
		total = len(users)
	}

	// Calculate pagination
	totalPages := (total + limit - 1) / limit

	return c.Status(200).JSON(fiber.Map{
		"success": true,
		"message": "Users retrieved successfully",
		"data":    users,
		"pagination": fiber.Map{
			"page":        page,
			"limit":       limit,
			"total":       total,
			"total_pages": totalPages,
		},
	})
}

// GetUserByID - GET /api/v1/users/:id
func (h *V1UserHandler) GetUserByID(c *fiber.Ctx) error {
	userID, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"success": false,
			"error":   "Invalid user ID format",
		})
	}

	user, err := h.UserService.GetUserByID(userID)
	if err != nil {
		return c.Status(404).JSON(fiber.Map{
			"success": false,
			"error":   "User not found",
		})
	}

	return c.Status(200).JSON(fiber.Map{
		"success": true,
		"message": "User retrieved successfully",
		"data":    user,
	})
}

// CreateUser - POST /api/v1/users
func (h *V1UserHandler) CreateUser(c *fiber.Ctx) error {
	var req service.CreateUserRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"success": false,
			"error":   "Invalid JSON format",
		})
	}

	user, err := h.UserService.CreateUser(&req)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"success": false,
			"error":   err.Error(),
		})
	}

	return c.Status(201).JSON(fiber.Map{
		"success": true,
		"message": "User created successfully",
		"data":    user,
	})
}

// UpdateUser - PUT /api/v1/users/:id
func (h *V1UserHandler) UpdateUser(c *fiber.Ctx) error {
	userID, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"success": false,
			"error":   "Invalid user ID format",
		})
	}

	var req service.UpdateUserRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"success": false,
			"error":   "Invalid JSON format",
		})
	}

	user, err := h.UserService.UpdateUser(userID, &req)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"success": false,
			"error":   err.Error(),
		})
	}

	return c.Status(200).JSON(fiber.Map{
		"success": true,
		"message": "User updated successfully",
		"data":    user,
	})
}

// DeleteUser - DELETE /api/v1/users/:id
func (h *V1UserHandler) DeleteUser(c *fiber.Ctx) error {
	userID, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"success": false,
			"error":   "Invalid user ID format",
		})
	}

	err = h.UserService.DeleteUser(userID)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"success": false,
			"error":   err.Error(),
		})
	}

	return c.Status(200).JSON(fiber.Map{
		"success": true,
		"message": "User deleted successfully",
		"data": fiber.Map{
			"user_id": userID,
		},
	})
}

// AssignRole - PUT /api/v1/users/:id/role
func (h *V1UserHandler) AssignRole(c *fiber.Ctx) error {
	userID, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"success": false,
			"error":   "Invalid user ID format",
		})
	}

	var req struct {
		RoleID uuid.UUID `json:"role_id"`
	}

	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"success": false,
			"error":   "Invalid JSON format",
		})
	}

	if req.RoleID == uuid.Nil {
		return c.Status(400).JSON(fiber.Map{
			"success": false,
			"error":   "Role ID is required",
		})
	}

	user, err := h.UserService.AssignRole(userID, req.RoleID)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"success": false,
			"error":   err.Error(),
		})
	}

	return c.Status(200).JSON(fiber.Map{
		"success": true,
		"message": "Role assigned successfully",
		"data":    user,
	})
}