package route

import (
	"UAS_BACKEND/domain/middleware"
	"UAS_BACKEND/domain/service"

	"github.com/gofiber/fiber/v2"
)

type FileHandler struct {
	FileService    *service.FileService
	RBACMiddleware *middleware.RBACMiddleware
}

func NewFileHandler(fileService *service.FileService, rbacMiddleware *middleware.RBACMiddleware) *FileHandler {
	return &FileHandler{
		FileService:    fileService,
		RBACMiddleware: rbacMiddleware,
	}
}

// UploadFile - Handler untuk upload single file
func (h *FileHandler) UploadFile(c *fiber.Ctx) error {
	// Get file from form
	file, err := c.FormFile("file")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "No file uploaded",
		})
	}

	// Upload file
	uploaded, err := h.FileService.UploadFile(file)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "File uploaded successfully",
		"data":    uploaded,
	})
}

// UploadMultipleFiles - Handler untuk upload multiple files
func (h *FileHandler) UploadMultipleFiles(c *fiber.Ctx) error {
	// Get form
	form, err := c.MultipartForm()
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Failed to parse multipart form",
		})
	}

	// Get files
	files := form.File["files"]
	if len(files) == 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "No files uploaded",
		})
	}

	// Upload files
	uploaded, err := h.FileService.UploadMultipleFiles(files)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Files uploaded successfully",
		"data":    uploaded,
		"count":   len(uploaded),
	})
}

// SetupFileRoutes - Setup routes untuk file upload
func SetupFileRoutes(app *fiber.App, handler *FileHandler, rbac *middleware.RBACMiddleware) {
	api := app.Group("/api")

	// File upload routes - require authentication
	files := api.Group("/files", rbac.Authenticate())
	{
		// Upload single file
		files.Post("/upload", handler.UploadFile)

		// Upload multiple files
		files.Post("/upload-multiple", handler.UploadMultipleFiles)
	}

	// Serve static files (uploaded files)
	app.Static("/uploads", "./uploads")
}
