package route

import (
	"UAS_BACKEND/domain/middleware"

	"github.com/gofiber/fiber/v2"
)

type ProtectedHandler struct {
	RBACMiddleware *middleware.RBACMiddleware
}

func NewProtectedHandler(rbacMiddleware *middleware.RBACMiddleware) *ProtectedHandler {
	return &ProtectedHandler{RBACMiddleware: rbacMiddleware}
}

// Example handlers untuk demonstrasi RBAC

func (h *ProtectedHandler) GetProfile(c *fiber.Ctx) error {
	userID, err := middleware.GetUserID(c)
	if err != nil {
		return err
	}

	return c.JSON(fiber.Map{
		"message": "Profile accessed",
		"user_id": userID,
	})
}

func (h *ProtectedHandler) GetStudents(c *fiber.Ctx) error {
	permissions, _ := middleware.GetPermissions(c)

	return c.JSON(fiber.Map{
		"message":     "Students list",
		"permissions": permissions,
	})
}

func (h *ProtectedHandler) CreateStudent(c *fiber.Ctx) error {
	userID, _ := middleware.GetUserID(c)

	return c.JSON(fiber.Map{
		"message": "Student created",
		"by":      userID,
	})
}

func (h *ProtectedHandler) VerifyAchievement(c *fiber.Ctx) error {
	userID, _ := middleware.GetUserID(c)

	return c.JSON(fiber.Map{
		"message": "Achievement verified",
		"by":      userID,
	})
}

func (h *ProtectedHandler) AdminOnly(c *fiber.Ctx) error {
	return c.JSON(fiber.Map{
		"message": "Admin area accessed",
	})
}

// SetupProtectedRoutes - Setup routes dengan RBAC protection
func SetupProtectedRoutes(app *fiber.App, handler *ProtectedHandler, rbac *middleware.RBACMiddleware) {
	api := app.Group("/api")

	// Routes yang memerlukan autentikasi saja (tanpa permission check)
	authenticated := api.Group("", rbac.Authenticate())
	{
		authenticated.Get("/profile", handler.GetProfile)
	}

	// Routes dengan permission check
	students := api.Group("/students", rbac.Authenticate())
	{
		// Read students - butuh permission "student.read"
		students.Get("/", rbac.RequirePermission("student.read"), handler.GetStudents)

		// Create student - butuh permission "student.write"
		students.Post("/", rbac.RequirePermission("student.write"), handler.CreateStudent)
	}

	// Routes dengan multiple permissions (salah satu)
	achievements := api.Group("/achievements", rbac.Authenticate())
	{
		// Verify achievement - butuh permission "achievement.verify"
		achievements.Post("/:id/verify",
			rbac.RequirePermission("achievement.verify"),
			handler.VerifyAchievement,
		)
	}

	// Routes dengan role check
	admin := api.Group("/admin", rbac.Authenticate(), rbac.RequireRole("admin"))
	{
		admin.Get("/dashboard", handler.AdminOnly)
	}

	// Routes dengan multiple permissions (semua harus ada)
	management := api.Group("/management", rbac.Authenticate())
	{
		management.Get("/users",
			rbac.RequireAllPermissions("user.read", "user.write"),
			handler.AdminOnly,
		)
	}

	// Routes dengan any permission (salah satu dari beberapa)
	reports := api.Group("/reports", rbac.Authenticate())
	{
		reports.Get("/",
			rbac.RequireAnyPermission("student.read", "lecturer.read", "achievement.read"),
			handler.GetProfile,
		)
	}
}
