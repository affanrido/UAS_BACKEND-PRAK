package route

import (
	"UAS_BACKEND/domain/middleware"
	model "UAS_BACKEND/domain/Model"
	"UAS_BACKEND/domain/service"
	"context"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type AdminHandler struct {
	UserService             *service.UserService
	AdminAchievementService *service.AdminAchievementService
	RBACMiddleware          *middleware.RBACMiddleware
}

func NewAdminHandler(userService *service.UserService, adminAchievementService *service.AdminAchievementService, rbacMiddleware *middleware.RBACMiddleware) *AdminHandler {
	return &AdminHandler{
		UserService:             userService,
		AdminAchievementService: adminAchievementService,
		RBACMiddleware:          rbacMiddleware,
	}
}

// CreateUser - Handler untuk create user (FR-009)
func (h *AdminHandler) CreateUser(c *fiber.Ctx) error {
	var req service.CreateUserRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	response, err := h.UserService.CreateUser(&req)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message": "User created successfully",
		"data":    response,
	})
}

// GetAllUsers - Handler untuk get all users
func (h *AdminHandler) GetAllUsers(c *fiber.Ctx) error {
	page, _ := strconv.Atoi(c.Query("page", "1"))
	pageSize, _ := strconv.Atoi(c.Query("page_size", "10"))

	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 10
	}
	if pageSize > 100 {
		pageSize = 100
	}

	offset := (page - 1) * pageSize

	users, total, err := h.UserService.GetAllUsers(pageSize, offset)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	totalPages := total / pageSize
	if total%pageSize > 0 {
		totalPages++
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Users retrieved successfully",
		"data":    users,
		"pagination": fiber.Map{
			"page":        page,
			"page_size":   pageSize,
			"total_items": total,
			"total_pages": totalPages,
		},
	})
}

// GetUserByID - Handler untuk get user by ID
func (h *AdminHandler) GetUserByID(c *fiber.Ctx) error {
	userID, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid user ID",
		})
	}

	user, err := h.UserService.GetUserByID(userID)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "User retrieved successfully",
		"data":    user,
	})
}

// UpdateUser - Handler untuk update user (FR-009)
func (h *AdminHandler) UpdateUser(c *fiber.Ctx) error {
	userID, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid user ID",
		})
	}

	var req service.UpdateUserRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	response, err := h.UserService.UpdateUser(userID, &req)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "User updated successfully",
		"data":    response,
	})
}

// DeleteUser - Handler untuk delete user (FR-009)
func (h *AdminHandler) DeleteUser(c *fiber.Ctx) error {
	userID, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid user ID",
		})
	}

	err = h.UserService.DeleteUser(userID)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "User deleted successfully",
	})
}

// AssignRole - Handler untuk assign role (FR-009)
func (h *AdminHandler) AssignRole(c *fiber.Ctx) error {
	userID, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid user ID",
		})
	}

	var req struct {
		RoleID string `json:"role_id"`
	}
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	roleID, err := uuid.Parse(req.RoleID)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid role ID",
		})
	}

	response, err := h.UserService.AssignRole(userID, roleID)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Role assigned successfully",
		"data":    response,
	})
}

// SetStudentProfile - Handler untuk set student profile (FR-009)
func (h *AdminHandler) SetStudentProfile(c *fiber.Ctx) error {
	userID, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid user ID",
		})
	}

	var req service.StudentProfileRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	student, err := h.UserService.SetStudentProfile(userID, &req)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Student profile set successfully",
		"data":    student,
	})
}

// SetLecturerProfile - Handler untuk set lecturer profile (FR-009)
func (h *AdminHandler) SetLecturerProfile(c *fiber.Ctx) error {
	userID, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid user ID",
		})
	}

	var req service.LecturerProfileRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	lecturer, err := h.UserService.SetLecturerProfile(userID, &req)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Lecturer profile set successfully",
		"data":    lecturer,
	})
}

// SetAdvisor - Handler untuk set advisor (FR-009)
func (h *AdminHandler) SetAdvisor(c *fiber.Ctx) error {
	studentUserID, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid student user ID",
		})
	}

	var req struct {
		AdvisorID string `json:"advisor_id"`
	}
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	advisorID, err := uuid.Parse(req.AdvisorID)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid advisor ID",
		})
	}

	student, err := h.UserService.SetAdvisor(studentUserID, advisorID)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Advisor set successfully",
		"data":    student,
	})
}

// ViewAllAchievements - Handler untuk view all achievements (FR-010)
func (h *AdminHandler) ViewAllAchievements(c *fiber.Ctx) error {
	// Parse query parameters
	page, _ := strconv.Atoi(c.Query("page", "1"))
	size, _ := strconv.Atoi(c.Query("size", "10"))

	// Parse filters
	filter := &model.AdminAchievementFilter{
		Status:          c.Query("status"),
		AchievementType: c.Query("achievement_type"),
		StudentID:       c.Query("student_id"),
		AdvisorID:       c.Query("advisor_id"),
		ProgramStudy:    c.Query("program_study"),
	}

	// Parse sort
	sort := &model.AdminAchievementSort{
		Field: c.Query("sort_field", "created_at"),
		Order: c.Query("sort_order", "desc"),
	}

	// Create request
	req := &service.ViewAllAchievementsRequest{
		Page:   page,
		Size:   size,
		Filter: filter,
		Sort:   sort,
	}

	// Create context with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Get achievements
	response, err := h.AdminAchievementService.ViewAllAchievements(ctx, req)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Achievements retrieved successfully",
		"data":    response.Achievements,
		"pagination": response.Pagination,
		"summary": response.Summary,
	})
}

// GetAchievementDetail - Handler untuk get achievement detail by reference ID
func (h *AdminHandler) GetAchievementDetail(c *fiber.Ctx) error {
	referenceID, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid reference ID",
		})
	}

	// Create context with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Get achievement detail
	achievement, err := h.AdminAchievementService.GetAchievementByReferenceID(ctx, referenceID)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Achievement detail retrieved successfully",
		"data":    achievement,
	})
}

// GetRoles - Handler untuk get all roles
func (h *AdminHandler) GetRoles(c *fiber.Ctx) error {
	roles, err := h.UserService.GetAllRoles()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Roles retrieved successfully",
		"data":    roles,
	})
}

// SetupAdminRoutes - Setup routes untuk admin
func SetupAdminRoutes(app *fiber.App, handler *AdminHandler, rbac *middleware.RBACMiddleware) {
	api := app.Group("/api")

	// Admin routes - require authentication and admin role
	admin := api.Group("/admin", rbac.Authenticate(), rbac.RequireRole("admin"))
	{
		// User management
		admin.Post("/users", handler.CreateUser)                    // Create user
		admin.Get("/users", handler.GetAllUsers)                    // Get all users
		admin.Get("/users/:id", handler.GetUserByID)                // Get user by ID
		admin.Put("/users/:id", handler.UpdateUser)                 // Update user
		admin.Delete("/users/:id", handler.DeleteUser)              // Delete user
		admin.Post("/users/:id/assign-role", handler.AssignRole)    // Assign role
		admin.Post("/users/:id/student-profile", handler.SetStudentProfile)   // Set student profile
		admin.Post("/users/:id/lecturer-profile", handler.SetLecturerProfile) // Set lecturer profile
		admin.Post("/users/:id/set-advisor", handler.SetAdvisor)    // Set advisor

		// Achievement management
		admin.Get("/achievements", handler.ViewAllAchievements)      // View all achievements
		admin.Get("/achievements/:id", handler.GetAchievementDetail) // Get achievement detail

		// Utility endpoints
		admin.Get("/roles", handler.GetRoles)                       // Get all roles
	}
}
