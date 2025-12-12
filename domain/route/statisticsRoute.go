package route

import (
	"UAS_BACKEND/domain/middleware"
	"UAS_BACKEND/domain/service"
	"context"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type StatisticsHandler struct {
	StatisticsService *service.StatisticsService
	RBACMiddleware    *middleware.RBACMiddleware
}

func NewStatisticsHandler(statisticsService *service.StatisticsService, rbacMiddleware *middleware.RBACMiddleware) *StatisticsHandler {
	return &StatisticsHandler{
		StatisticsService: statisticsService,
		RBACMiddleware:    rbacMiddleware,
	}
}

// GetStudentStatistics - Handler untuk student statistics (own achievements)
func (h *StatisticsHandler) GetStudentStatistics(c *fiber.Ctx) error {
	// Get user ID dari context
	userID, err := middleware.GetUserID(c)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Unauthorized",
		})
	}

	// Parse query parameters
	req := &service.StatisticsRequest{}
	
	if startDate := c.Query("start_date"); startDate != "" {
		if t, err := time.Parse("2006-01-02", startDate); err == nil {
			req.StartDate = &t
		}
	}
	
	if endDate := c.Query("end_date"); endDate != "" {
		if t, err := time.Parse("2006-01-02", endDate); err == nil {
			req.EndDate = &t
		}
	}

	// Create context with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	// Get statistics
	stats, err := h.StatisticsService.GetStudentStatistics(ctx, userID, req)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Student statistics retrieved successfully",
		"data":    stats,
	})
}

// GetLecturerStatistics - Handler untuk lecturer statistics (advisee achievements)
func (h *StatisticsHandler) GetLecturerStatistics(c *fiber.Ctx) error {
	// Get user ID dari context
	userID, err := middleware.GetUserID(c)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Unauthorized",
		})
	}

	// Parse query parameters
	req := &service.StatisticsRequest{}
	
	if startDate := c.Query("start_date"); startDate != "" {
		if t, err := time.Parse("2006-01-02", startDate); err == nil {
			req.StartDate = &t
		}
	}
	
	if endDate := c.Query("end_date"); endDate != "" {
		if t, err := time.Parse("2006-01-02", endDate); err == nil {
			req.EndDate = &t
		}
	}

	// Create context with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	// Get statistics
	stats, err := h.StatisticsService.GetLecturerStatistics(ctx, userID, req)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Lecturer statistics retrieved successfully",
		"data":    stats,
	})
}

// GetAdminStatistics - Handler untuk admin statistics (all achievements)
func (h *StatisticsHandler) GetAdminStatistics(c *fiber.Ctx) error {
	// Parse query parameters
	req := &service.StatisticsRequest{}
	
	if startDate := c.Query("start_date"); startDate != "" {
		if t, err := time.Parse("2006-01-02", startDate); err == nil {
			req.StartDate = &t
		}
	}
	
	if endDate := c.Query("end_date"); endDate != "" {
		if t, err := time.Parse("2006-01-02", endDate); err == nil {
			req.EndDate = &t
		}
	}

	// Create context with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	// Get statistics
	stats, err := h.StatisticsService.GetAdminStatistics(ctx, req)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Admin statistics retrieved successfully",
		"data":    stats,
	})
}

// GetAchievementTrends - Handler untuk achievement trends
func (h *StatisticsHandler) GetAchievementTrends(c *fiber.Ctx) error {
	// Get user ID dan role dari context
	userID, err := middleware.GetUserID(c)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Unauthorized",
		})
	}

	// Get role dari claims
	claims, err := middleware.GetClaims(c)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Unauthorized",
		})
	}

	// Determine role
	role := "student" // default
	permissions := claims.Permissions
	for _, perm := range permissions {
		if perm == "admin" {
			role = "admin"
			break
		} else if perm == "achievement.verify" {
			role = "lecturer"
		}
	}

	// Parse months parameter
	months, _ := strconv.Atoi(c.Query("months", "12"))

	// Create context with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	// Get trends
	trends, err := h.StatisticsService.GetAchievementTrends(ctx, userID, role, months)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Achievement trends retrieved successfully",
		"data":    trends,
	})
}

// SetupStatisticsRoutes - Setup routes untuk statistics
func SetupStatisticsRoutes(app *fiber.App, handler *StatisticsHandler, rbac *middleware.RBACMiddleware) {
	api := app.Group("/api")

	// Student statistics - require authentication and student permissions
	student := api.Group("/student", rbac.Authenticate())
	{
		student.Get("/statistics", handler.GetStudentStatistics)
		student.Get("/trends", handler.GetAchievementTrends)
	}

	// Lecturer statistics - require authentication and lecturer permissions
	lecturer := api.Group("/lecturer", rbac.Authenticate())
	{
		lecturer.Get("/statistics", 
			rbac.RequirePermission("achievement.verify"),
			handler.GetLecturerStatistics,
		)
		lecturer.Get("/trends", 
			rbac.RequirePermission("achievement.verify"),
			handler.GetAchievementTrends,
		)
	}

	// Admin statistics - require authentication and admin role
	admin := api.Group("/admin", rbac.Authenticate(), rbac.RequireRole("admin"))
	{
		admin.Get("/statistics", handler.GetAdminStatistics)
		admin.Get("/trends", handler.GetAchievementTrends)
	}
}