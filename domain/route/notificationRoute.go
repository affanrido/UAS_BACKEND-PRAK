package route

import (
	"UAS_BACKEND/domain/middleware"
	"UAS_BACKEND/domain/service"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type NotificationHandler struct {
	NotificationService *service.NotificationService
	RBACMiddleware      *middleware.RBACMiddleware
}

func NewNotificationHandler(notificationService *service.NotificationService, rbacMiddleware *middleware.RBACMiddleware) *NotificationHandler {
	return &NotificationHandler{
		NotificationService: notificationService,
		RBACMiddleware:      rbacMiddleware,
	}
}

// GetMyNotifications - Handler untuk get notifikasi user
func (h *NotificationHandler) GetMyNotifications(c *fiber.Ctx) error {
	userID, err := middleware.GetUserID(c)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Unauthorized",
		})
	}

	// Get limit from query parameter (default 50)
	limitStr := c.Query("limit", "50")
	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit <= 0 {
		limit = 50
	}

	notifications, err := h.NotificationService.GetUserNotifications(userID, limit)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to get notifications",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Notifications retrieved successfully",
		"data":    notifications,
		"count":   len(notifications),
	})
}

// GetUnreadCount - Handler untuk get jumlah notifikasi yang belum dibaca
func (h *NotificationHandler) GetUnreadCount(c *fiber.Ctx) error {
	userID, err := middleware.GetUserID(c)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Unauthorized",
		})
	}

	count, err := h.NotificationService.GetUnreadCount(userID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to get unread count",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"unread_count": count,
	})
}

// MarkAsRead - Handler untuk tandai notifikasi sebagai sudah dibaca
func (h *NotificationHandler) MarkAsRead(c *fiber.Ctx) error {
	notificationIDStr := c.Params("id")
	notificationID, err := uuid.Parse(notificationIDStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid notification ID",
		})
	}

	err = h.NotificationService.MarkAsRead(notificationID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to mark as read",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Notification marked as read",
	})
}

// MarkAllAsRead - Handler untuk tandai semua notifikasi sebagai sudah dibaca
func (h *NotificationHandler) MarkAllAsRead(c *fiber.Ctx) error {
	userID, err := middleware.GetUserID(c)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Unauthorized",
		})
	}

	err = h.NotificationService.MarkAllAsRead(userID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to mark all as read",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "All notifications marked as read",
	})
}

// SetupNotificationRoutes - Setup routes untuk notification
func SetupNotificationRoutes(app *fiber.App, handler *NotificationHandler, rbac *middleware.RBACMiddleware) {
	api := app.Group("/api")

	// Notification routes - require authentication
	notifications := api.Group("/notifications", rbac.Authenticate())
	{
		// Get my notifications
		notifications.Get("/", handler.GetMyNotifications)

		// Get unread count
		notifications.Get("/unread-count", handler.GetUnreadCount)

		// Mark as read
		notifications.Put("/:id/read", handler.MarkAsRead)

		// Mark all as read
		notifications.Put("/read-all", handler.MarkAllAsRead)
	}
}
