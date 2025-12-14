package route

import (
	"UAS_BACKEND/domain/middleware"
	"UAS_BACKEND/domain/service"

	"github.com/gofiber/fiber/v2"
)

// SetupV1Routes - Setup all v1 API routes
func SetupV1Routes(
	app *fiber.App,
	authService *service.AuthService,
	userService *service.UserService,
	achievementService *service.AchievementService,
	notificationService *service.NotificationService,
	fileService *service.FileService,
	statisticsService *service.StatisticsService,
	rbacMiddleware *middleware.RBACMiddleware,
) {
	// Initialize handlers
	v1AuthHandler := NewV1AuthHandler(authService, rbacMiddleware)
	v1UserHandler := NewV1UserHandler(userService, rbacMiddleware)
	v1AchievementHandler := NewV1AchievementHandler(achievementService, notificationService, fileService, rbacMiddleware)
	v1StudentLecturerHandler := NewV1StudentLecturerHandler(userService, achievementService, rbacMiddleware)
	v1ReportHandler := NewV1ReportHandler(statisticsService, rbacMiddleware)

	// Setup routes
	SetupV1AuthRoutes(app, v1AuthHandler)
	SetupV1UserRoutes(app, v1UserHandler)
	SetupV1AchievementRoutes(app, v1AchievementHandler)
	SetupV1StudentLecturerRoutes(app, v1StudentLecturerHandler)
	SetupV1ReportRoutes(app, v1ReportHandler)

	// API v1 info endpoint
	app.Get("/api/v1", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"success": true,
			"message": "UAS Backend API v1",
			"version": "1.0.0",
			"endpoints": fiber.Map{
				"authentication": "/api/v1/auth",
				"users":          "/api/v1/users",
				"achievements":   "/api/v1/achievements",
				"students":       "/api/v1/students",
				"lecturers":      "/api/v1/lecturers",
				"reports":        "/api/v1/reports",
			},
			"documentation": "/swagger/",
		})
	})
}