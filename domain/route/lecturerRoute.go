package route

import (
	"UAS_BACKEND/domain/middleware"
	model "UAS_BACKEND/domain/Model"
	"UAS_BACKEND/domain/service"
	"context"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
)

type LecturerHandler struct {
	AchievementService *service.AchievementService
	RBACMiddleware     *middleware.RBACMiddleware
}

func NewLecturerHandler(achievementService *service.AchievementService, rbacMiddleware *middleware.RBACMiddleware) *LecturerHandler {
	return &LecturerHandler{
		AchievementService: achievementService,
		RBACMiddleware:     rbacMiddleware,
	}
}

// ViewAdvisedStudentsAchievements - Handler untuk view prestasi mahasiswa bimbingan (FR-006)
func (h *LecturerHandler) ViewAdvisedStudentsAchievements(c *fiber.Ctx) error {
	// Get user ID dari context
	userID, err := middleware.GetUserID(c)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Unauthorized",
		})
	}

	// Parse pagination parameters
	page, _ := strconv.Atoi(c.Query("page", "1"))
	pageSize, _ := strconv.Atoi(c.Query("page_size", "10"))

	pagination := model.PaginationRequest{
		Page:     page,
		PageSize: pageSize,
	}

	// Create context with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Get advised students achievements
	response, err := h.AchievementService.ViewAdvisedStudentsAchievements(ctx, userID, pagination)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	// Return success response
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Advised students achievements retrieved successfully",
		"data":    response.Achievements,
		"pagination": response.Pagination,
	})
}

// SetupLecturerRoutes - Setup routes untuk lecturer/dosen
func SetupLecturerRoutes(app *fiber.App, handler *LecturerHandler, rbac *middleware.RBACMiddleware) {
	api := app.Group("/api")

	// Lecturer routes - require authentication and lecturer permissions
	lecturer := api.Group("/lecturer", rbac.Authenticate())
	{
		// View prestasi mahasiswa bimbingan - dosen wali melihat prestasi mahasiswa bimbingannya
		lecturer.Get("/advised-students/achievements",
			rbac.RequireAnyPermission("student.read", "achievement.read"),
			handler.ViewAdvisedStudentsAchievements,
		)
	}
}
