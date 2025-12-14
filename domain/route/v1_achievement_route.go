package route

import (
	"UAS_BACKEND/domain/middleware"
	"UAS_BACKEND/domain/model"
	"UAS_BACKEND/domain/service"
	"context"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type V1AchievementHandler struct {
	AchievementService  *service.AchievementService
	NotificationService *service.NotificationService
	FileService         *service.FileService
	RBACMiddleware      *middleware.RBACMiddleware
}

func NewV1AchievementHandler(
	achievementService *service.AchievementService,
	notificationService *service.NotificationService,
	fileService *service.FileService,
	rbacMiddleware *middleware.RBACMiddleware,
) *V1AchievementHandler {
	return &V1AchievementHandler{
		AchievementService:  achievementService,
		NotificationService: notificationService,
		FileService:         fileService,
		RBACMiddleware:      rbacMiddleware,
	}
}

// SetupV1AchievementRoutes - Setup achievement routes v1
func SetupV1AchievementRoutes(app *fiber.App, handler *V1AchievementHandler) {
	achievements := app.Group("/api/v1/achievements")
	achievements.Use(handler.RBACMiddleware.RequireAuth())

	// 5.4 Achievements endpoints
	achievements.Get("/", handler.GetAchievements)                                                                    // List (filtered by role)
	achievements.Get("/:id", handler.GetAchievementDetail)                                                           // Detail
	achievements.Post("/", handler.RBACMiddleware.RequirePermission("achievement.write"), handler.CreateAchievement) // Create (Mahasiswa)
	achievements.Put("/:id", handler.RBACMiddleware.RequirePermission("achievement.write"), handler.UpdateAchievement) // Update (Mahasiswa)
	achievements.Delete("/:id", handler.RBACMiddleware.RequirePermission("achievement.write"), handler.DeleteAchievement) // Delete (Mahasiswa)
	achievements.Post("/:id/submit", handler.RBACMiddleware.RequirePermission("achievement.write"), handler.SubmitForVerification) // Submit for verification
	achievements.Post("/:id/verify", handler.RBACMiddleware.RequirePermission("achievement.verify"), handler.VerifyAchievement) // Verify (Dosen Wali)
	achievements.Post("/:id/reject", handler.RBACMiddleware.RequirePermission("achievement.verify"), handler.RejectAchievement) // Reject (Dosen Wali)
	achievements.Get("/:id/history", handler.GetAchievementHistory)                                                  // Status history
	achievements.Post("/:id/attachments", handler.RBACMiddleware.RequirePermission("achievement.write"), handler.UploadAttachments) // Upload files
}

// GetAchievements - GET /api/v1/achievements (List filtered by role)
func (h *V1AchievementHandler) GetAchievements(c *fiber.Ctx) error {
	user := c.Locals("user").(*model.Claims)
	
	// Parse query parameters
	page, _ := strconv.Atoi(c.Query("page", "1"))
	limit, _ := strconv.Atoi(c.Query("limit", "10"))
	status := c.Query("status")
	achievementType := c.Query("type")
	studentID := c.Query("student_id")

	offset := (page - 1) * limit

	// Role-based filtering logic would be implemented here
	// For now, return mock response
	return c.Status(200).JSON(fiber.Map{
		"success": true,
		"message": "Achievements retrieved successfully",
		"data":    []interface{}{},
		"pagination": fiber.Map{
			"page":        page,
			"limit":       limit,
			"total":       0,
			"total_pages": 0,
		},
		"filters": fiber.Map{
			"status":           status,
			"achievement_type": achievementType,
			"student_id":       studentID,
			"user_role":        user.RoleID,
		},
	})
}

// GetAchievementDetail - GET /api/v1/achievements/:id
func (h *V1AchievementHandler) GetAchievementDetail(c *fiber.Ctx) error {
	achievementID, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"success": false,
			"error":   "Invalid achievement ID format",
		})
	}

	// Get achievement detail logic would be implemented here
	return c.Status(200).JSON(fiber.Map{
		"success": true,
		"message": "Achievement detail retrieved successfully",
		"data": fiber.Map{
			"id":     achievementID,
			"status": "draft",
			// More fields would be added here
		},
	})
}

// CreateAchievement - POST /api/v1/achievements (Mahasiswa)
func (h *V1AchievementHandler) CreateAchievement(c *fiber.Ctx) error {
	user := c.Locals("user").(*model.Claims)

	var req service.SubmitAchievementRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"success": false,
			"error":   "Invalid JSON format",
		})
	}

	result, err := h.AchievementService.SubmitAchievement(context.Background(), user.UserID, &req)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"success": false,
			"error":   err.Error(),
		})
	}

	return c.Status(201).JSON(fiber.Map{
		"success": true,
		"message": "Achievement created successfully",
		"data":    result,
	})
}

// UpdateAchievement - PUT /api/v1/achievements/:id (Mahasiswa)
func (h *V1AchievementHandler) UpdateAchievement(c *fiber.Ctx) error {
	achievementID, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"success": false,
			"error":   "Invalid achievement ID format",
		})
	}

	var req service.SubmitAchievementRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"success": false,
			"error":   "Invalid JSON format",
		})
	}

	// Update achievement logic would be implemented here
	return c.Status(200).JSON(fiber.Map{
		"success": true,
		"message": "Achievement updated successfully",
		"data": fiber.Map{
			"id": achievementID,
		},
	})
}

// DeleteAchievement - DELETE /api/v1/achievements/:id (Mahasiswa)
func (h *V1AchievementHandler) DeleteAchievement(c *fiber.Ctx) error {
	user := c.Locals("user").(*model.Claims)
	achievementID, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"success": false,
			"error":   "Invalid achievement ID format",
		})
	}

	result, err := h.AchievementService.DeleteAchievement(context.Background(), user.UserID, achievementID)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"success": false,
			"error":   err.Error(),
		})
	}

	return c.Status(200).JSON(fiber.Map{
		"success": true,
		"message": "Achievement deleted successfully",
		"data":    result,
	})
}

// SubmitForVerification - POST /api/v1/achievements/:id/submit
func (h *V1AchievementHandler) SubmitForVerification(c *fiber.Ctx) error {
	user := c.Locals("user").(*model.Claims)
	achievementID, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"success": false,
			"error":   "Invalid achievement ID format",
		})
	}

	result, err := h.AchievementService.SubmitForVerification(context.Background(), user.UserID, achievementID, h.NotificationService)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"success": false,
			"error":   err.Error(),
		})
	}

	return c.Status(200).JSON(fiber.Map{
		"success": true,
		"message": "Achievement submitted for verification successfully",
		"data":    result,
	})
}

// VerifyAchievement - POST /api/v1/achievements/:id/verify (Dosen Wali)
func (h *V1AchievementHandler) VerifyAchievement(c *fiber.Ctx) error {
	user := c.Locals("user").(*model.Claims)
	achievementID, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"success": false,
			"error":   "Invalid achievement ID format",
		})
	}

	var req struct {
		Notes string `json:"notes"`
	}

	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"success": false,
			"error":   "Invalid JSON format",
		})
	}

	// Verify achievement logic would be implemented here
	return c.Status(200).JSON(fiber.Map{
		"success": true,
		"message": "Achievement verified successfully",
		"data": fiber.Map{
			"id":          achievementID,
			"status":      "verified",
			"verified_by": user.UserID,
			"verified_at": time.Now(),
			"notes":       req.Notes,
		},
	})
}

// RejectAchievement - POST /api/v1/achievements/:id/reject (Dosen Wali)
func (h *V1AchievementHandler) RejectAchievement(c *fiber.Ctx) error {
	user := c.Locals("user").(*model.Claims)
	achievementID, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"success": false,
			"error":   "Invalid achievement ID format",
		})
	}

	var req struct {
		RejectionNote string `json:"rejection_note"`
	}

	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"success": false,
			"error":   "Invalid JSON format",
		})
	}

	if req.RejectionNote == "" {
		return c.Status(400).JSON(fiber.Map{
			"success": false,
			"error":   "Rejection note is required",
		})
	}

	// Reject achievement logic would be implemented here
	return c.Status(200).JSON(fiber.Map{
		"success": true,
		"message": "Achievement rejected successfully",
		"data": fiber.Map{
			"id":             achievementID,
			"status":         "rejected",
			"rejected_by":    user.UserID,
			"rejected_at":    time.Now(),
			"rejection_note": req.RejectionNote,
		},
	})
}

// GetAchievementHistory - GET /api/v1/achievements/:id/history
func (h *V1AchievementHandler) GetAchievementHistory(c *fiber.Ctx) error {
	achievementID, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"success": false,
			"error":   "Invalid achievement ID format",
		})
	}

	// Mock history data
	history := []fiber.Map{
		{
			"status":     "draft",
			"timestamp":  time.Now().Add(-72 * time.Hour),
			"action":     "created",
			"actor":      "student",
			"notes":      "Achievement created",
		},
		{
			"status":     "submitted",
			"timestamp":  time.Now().Add(-24 * time.Hour),
			"action":     "submitted",
			"actor":      "student",
			"notes":      "Submitted for verification",
		},
	}

	return c.Status(200).JSON(fiber.Map{
		"success": true,
		"message": "Achievement history retrieved successfully",
		"data": fiber.Map{
			"achievement_id": achievementID,
			"history":        history,
		},
	})
}

// UploadAttachments - POST /api/v1/achievements/:id/attachments
func (h *V1AchievementHandler) UploadAttachments(c *fiber.Ctx) error {
	achievementID, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"success": false,
			"error":   "Invalid achievement ID format",
		})
	}

	// Handle file upload
	form, err := c.MultipartForm()
	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"success": false,
			"error":   "Failed to parse multipart form",
		})
	}

	files := form.File["files"]
	if len(files) == 0 {
		return c.Status(400).JSON(fiber.Map{
			"success": false,
			"error":   "No files uploaded",
		})
	}

	uploadedFiles := []fiber.Map{}
	for _, file := range files {
		// Save file using FileService
		savedFile, err := h.FileService.SaveFile(file)
		if err != nil {
			return c.Status(400).JSON(fiber.Map{
				"success": false,
				"error":   "Failed to save file: " + err.Error(),
			})
		}

		uploadedFiles = append(uploadedFiles, fiber.Map{
			"filename":     savedFile.FileName,
			"url":          savedFile.URL,
			"size":         savedFile.Size,
			"content_type": savedFile.ContentType,
		})
	}

	return c.Status(200).JSON(fiber.Map{
		"success": true,
		"message": "Files uploaded successfully",
		"data": fiber.Map{
			"achievement_id": achievementID,
			"files":          uploadedFiles,
		},
	})
}