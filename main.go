package main

import (
	"UAS_BACKEND/domain/config"
	"UAS_BACKEND/domain/middleware"
	"UAS_BACKEND/domain/repository"
	"UAS_BACKEND/domain/route"
	"UAS_BACKEND/domain/service"
	"context"
	"log"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
)

func main() {
	// Load PostgreSQL configuration
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}
	defer cfg.DB.Close()

	// Load MongoDB configuration
	mongoCfg, err := config.LoadMongoConfig()
	if err != nil {
		log.Fatalf("Failed to load MongoDB config: %v", err)
	}
	defer mongoCfg.Client.Disconnect(context.Background())

	// Initialize repositories
	authRepo := repository.NewAuthRepository(cfg.DB)
	rbacRepo := repository.NewRBACRepository(cfg.DB)
	achievementRepo := repository.NewAchievementRepository(cfg.DB, mongoCfg.Database)
	notificationRepo := repository.NewNotificationRepository(cfg.DB)

	// Initialize services
	authService := service.NewAuthService(
		cfg.SecretKey,
		24*time.Hour, // Token TTL 24 jam
		authRepo,
	)
	rbacService := service.NewRBACService(rbacRepo)
	achievementService := service.NewAchievementService(achievementRepo)
	notificationService := service.NewNotificationService(notificationRepo)
	fileService := service.NewFileService("./uploads", 10) // Max 10MB per file

	// Initialize middleware
	rbacMiddleware := middleware.NewRBACMiddleware(authService, rbacService)

	// Initialize handlers
	authHandler := route.NewAuthHandler(authService)
	protectedHandler := route.NewProtectedHandler(rbacMiddleware)
	achievementHandler := route.NewAchievementHandler(achievementService, notificationService, rbacMiddleware)
	notificationHandler := route.NewNotificationHandler(notificationService, rbacMiddleware)
	lecturerHandler := route.NewLecturerHandler(achievementService, notificationService, rbacMiddleware)
	fileHandler := route.NewFileHandler(fileService, rbacMiddleware)

	// Setup Fiber app
	app := fiber.New(fiber.Config{
		ErrorHandler: func(c *fiber.Ctx, err error) error {
			code := fiber.StatusInternalServerError
			if e, ok := err.(*fiber.Error); ok {
				code = e.Code
			}
			return c.Status(code).JSON(fiber.Map{
				"error": err.Error(),
			})
		},
		BodyLimit: 10 * 1024 * 1024, // 10MB body limit
	})

	// Global middleware
	app.Use(logger.New())
	app.Use(cors.New())

	// Setup routes
	route.SetupAuthRoutes(app, authHandler)
	route.SetupProtectedRoutes(app, protectedHandler, rbacMiddleware)
	route.SetupAchievementRoutes(app, achievementHandler, rbacMiddleware)
	route.SetupNotificationRoutes(app, notificationHandler, rbacMiddleware)
	route.SetupLecturerRoutes(app, lecturerHandler, rbacMiddleware)
	route.SetupFileRoutes(app, fileHandler, rbacMiddleware)

	// Health check
	app.Get("/health", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{"status": "ok"})
	})

	// Start server
	port := ":8080"
	log.Printf("Server running on http://localhost%s", port)
	if err := app.Listen(port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}