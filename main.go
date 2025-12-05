package main

import (
	"UAS_BACKEND/domain/config"
	"UAS_BACKEND/domain/repository"
	"UAS_BACKEND/domain/route"
	"UAS_BACKEND/domain/service"
	"log"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
)

func main() {
	// Load configuration
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}
	defer cfg.DB.Close()

	// Initialize repository
	authRepo := repository.NewAuthRepository(cfg.DB)

	// Initialize service
	authService := service.NewAuthService(
		cfg.SecretKey,
		24*time.Hour, // Token TTL 24 jam
		authRepo,
	)

	// Initialize handler
	authHandler := route.NewAuthHandler(authService)

	// Setup Fiber app
	app := fiber.New(fiber.Config{
		ErrorHandler: func(c *fiber.Ctx, err error) error {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": err.Error(),
			})
		},
	})

	// Middleware
	app.Use(logger.New())
	app.Use(cors.New())

	// Setup routes
	route.SetupAuthRoutes(app, authHandler)

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