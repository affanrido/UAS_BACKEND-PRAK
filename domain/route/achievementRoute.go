package route

import (
	"UAS_BACKEND/domain/middleware"
	"UAS_BACKEND/domain/service"
	"context"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type AchievementHandler struct {
	AchievementService  *service.AchievementService
	NotificationService *service.NotificationService
	RBACMiddleware      *middleware.RBACMiddleware
}

func NewAchievementHandler(achievementService *service.AchievementService, notificationService *service.NotificationService, rbacMiddleware *middleware.RBACMiddleware) *AchievementHandler {
	return &AchievementHandler{
		AchievementService:  achievementService,
		NotificationService: notificationService,
		RBACMiddleware:      rbacMiddleware,
	}
}

// SubmitAchievement - Handler untuk submit prestasi (FR-003)
func (h *AchievementHandler) SubmitAchievement(c *fiber.Ctx) error {
	// Get user ID dari context (sudah di-set oleh Authenticate middleware)
	userID, err := middleware.GetUserID(c)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Unauthorized",
		})
	}

	// Parse request body
	var req service.SubmitAchievementRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	// Create context with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Submit achievement
	response, err := h.AchievementService.SubmitAchievement(ctx, userID, &req)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	// Return success response
	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message": "Achievement submitted successfully",
		"data":    response,
	})
}

// GetMyAchievements - Handler untuk get semua prestasi mahasiswa
func (h *AchievementHandler) GetMyAchievements(c *fiber.Ctx) error {
	userID, err := middleware.GetUserID(c)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Unauthorized",
		})
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	achievements, err := h.AchievementService.GetStudentAchievements(ctx, userID)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Achievements retrieved successfully",
		"data":    achievements,
		"count":   len(achievements),
	})
}

// GetAchievementByID - Handler untuk get detail prestasi
func (h *AchievementHandler) GetAchievementByID(c *fiber.Ctx) error {
	mongoID := c.Params("id")
	if mongoID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Achievement ID is required",
		})
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	achievement, err := h.AchievementService.GetAchievementByID(ctx, mongoID)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Achievement retrieved successfully",
		"data":    achievement,
	})
}

// SubmitForVerification - Handler untuk submit prestasi untuk verifikasi (FR-004)
func (h *AchievementHandler) SubmitForVerification(c *fiber.Ctx) error {
	// Get user ID dari context
	userID, err := middleware.GetUserID(c)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Unauthorized",
		})
	}

	// Parse request body
	var req service.SubmitForVerificationRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	// Validate reference_id
	if req.ReferenceID == uuid.Nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "reference_id is required",
		})
	}

	// Create context with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Submit for verification
	response, err := h.AchievementService.SubmitForVerification(ctx, userID, req.ReferenceID, h.NotificationService)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	// Return success response
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": response.Message,
		"data":    response,
	})
}

// DeleteAchievement - Handler untuk hapus prestasi draft (FR-005)
func (h *AchievementHandler) DeleteAchievement(c *fiber.Ctx) error {
	// Get user ID dari context
	userID, err := middleware.GetUserID(c)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Unauthorized",
		})
	}

	// Get reference_id from URL parameter
	referenceIDStr := c.Params("id")
	referenceID, err := uuid.Parse(referenceIDStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid reference ID",
		})
	}

	// Create context with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Delete achievement
	response, err := h.AchievementService.DeleteAchievement(ctx, userID, referenceID)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	// Return success response
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": response.Message,
		"data":    response,
	})
}

// SetupAchievementRoutes - Setup routes untuk achievement
func SetupAchievementRoutes(app *fiber.App, handler *AchievementHandler, rbac *middleware.RBACMiddleware) {
	api := app.Group("/api")

	// Achievement routes - require authentication
	achievements := api.Group("/achievements", rbac.Authenticate())
	{
		// Submit prestasi - hanya mahasiswa (butuh permission "achievement.write")
		achievements.Post("/",
			rbac.RequirePermission("achievement.write"),
			handler.SubmitAchievement,
		)

		// Get my achievements - mahasiswa melihat prestasi sendiri
		achievements.Get("/my",
			rbac.RequirePermission("achievement.read"),
			handler.GetMyAchievements,
		)

		// Get achievement by ID - siapa saja yang punya permission "achievement.read"
		achievements.Get("/:id",
			rbac.RequirePermission("achievement.read"),
			handler.GetAchievementByID,
		)

		// Submit for verification - mahasiswa submit prestasi draft untuk verifikasi
		achievements.Post("/submit-verification",
			rbac.RequirePermission("achievement.write"),
			handler.SubmitForVerification,
		)

		// Delete achievement - mahasiswa hapus prestasi draft
		achievements.Delete("/:id",
			rbac.RequirePermission("achievement.write"),
			handler.DeleteAchievement,
		)
	}
}
