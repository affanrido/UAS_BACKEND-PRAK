package route

import (
	"UAS_BACKEND/domain/middleware"
	"UAS_BACKEND/domain/model"
	"UAS_BACKEND/domain/service"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type V1StudentLecturerHandler struct {
	UserService         *service.UserService
	AchievementService  *service.AchievementService
	RBACMiddleware      *middleware.RBACMiddleware
}

func NewV1StudentLecturerHandler(
	userService *service.UserService,
	achievementService *service.AchievementService,
	rbacMiddleware *middleware.RBACMiddleware,
) *V1StudentLecturerHandler {
	return &V1StudentLecturerHandler{
		UserService:        userService,
		AchievementService: achievementService,
		RBACMiddleware:     rbacMiddleware,
	}
}

// SetupV1StudentLecturerRoutes - Setup student and lecturer routes v1
func SetupV1StudentLecturerRoutes(app *fiber.App, handler *V1StudentLecturerHandler) {
	// 5.5 Students endpoints
	students := app.Group("/api/v1/students")
	students.Use(handler.RBACMiddleware.RequireAuth())

	students.Get("/", handler.RBACMiddleware.RequireAnyPermission("student.read", "admin.manage"), handler.GetAllStudents)
	students.Get("/:id", handler.RBACMiddleware.RequireAnyPermission("student.read", "admin.manage"), handler.GetStudentByID)
	students.Get("/:id/achievements", handler.RBACMiddleware.RequireAnyPermission("student.read", "achievement.read"), handler.GetStudentAchievements)
	students.Put("/:id/advisor", handler.RBACMiddleware.RequirePermission("admin.manage"), handler.SetStudentAdvisor)

	// 5.5 Lecturers endpoints
	lecturers := app.Group("/api/v1/lecturers")
	lecturers.Use(handler.RBACMiddleware.RequireAuth())

	lecturers.Get("/", handler.RBACMiddleware.RequireAnyPermission("lecturer.read", "admin.manage"), handler.GetAllLecturers)
	lecturers.Get("/:id/advisees", handler.RBACMiddleware.RequireAnyPermission("lecturer.read", "admin.manage"), handler.GetLecturerAdvisees)
}

// GetAllStudents - GET /api/v1/students
func (h *V1StudentLecturerHandler) GetAllStudents(c *fiber.Ctx) error {
	// Parse query parameters
	page, _ := strconv.Atoi(c.Query("page", "1"))
	limit, _ := strconv.Atoi(c.Query("limit", "10"))
	programStudy := c.Query("program_study")
	academicYear := c.Query("academic_year")
	advisorID := c.Query("advisor_id")

	offset := (page - 1) * limit

	// Get all users with student role
	users, total, err := h.UserService.GetAllUsers(limit, offset)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"success": false,
			"error":   "Failed to retrieve students",
		})
	}

	// Filter only students
	students := []fiber.Map{}
	for _, user := range users {
		if user.Student != nil {
			studentData := fiber.Map{
				"id":            user.Student.ID,
				"user_id":       user.User.ID,
				"student_id":    user.Student.StudentID,
				"full_name":     user.User.FullName,
				"email":         user.User.Email,
				"program_study": user.Student.ProgramStudy,
				"academic_year": user.Student.AcademicYear,
				"advisor_id":    user.Student.AdvisorID,
				"is_active":     user.User.IsActive,
				"created_at":    user.Student.CreatedAt,
			}

			// Apply filters
			include := true
			if programStudy != "" && user.Student.ProgramStudy != programStudy {
				include = false
			}
			if academicYear != "" && user.Student.AcademicYear != academicYear {
				include = false
			}
			if advisorID != "" && user.Student.AdvisorID.String() != advisorID {
				include = false
			}

			if include {
				students = append(students, studentData)
			}
		}
	}

	totalPages := (len(students) + limit - 1) / limit

	return c.Status(200).JSON(fiber.Map{
		"success": true,
		"message": "Students retrieved successfully",
		"data":    students,
		"pagination": fiber.Map{
			"page":        page,
			"limit":       limit,
			"total":       len(students),
			"total_pages": totalPages,
		},
		"filters": fiber.Map{
			"program_study": programStudy,
			"academic_year": academicYear,
			"advisor_id":    advisorID,
		},
	})
}

// GetStudentByID - GET /api/v1/students/:id
func (h *V1StudentLecturerHandler) GetStudentByID(c *fiber.Ctx) error {
	studentID, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"success": false,
			"error":   "Invalid student ID format",
		})
	}

	// Get user by ID and check if it's a student
	user, err := h.UserService.GetUserByID(studentID)
	if err != nil {
		return c.Status(404).JSON(fiber.Map{
			"success": false,
			"error":   "Student not found",
		})
	}

	if user.Student == nil {
		return c.Status(404).JSON(fiber.Map{
			"success": false,
			"error":   "User is not a student",
		})
	}

	studentData := fiber.Map{
		"id":            user.Student.ID,
		"user_id":       user.User.ID,
		"student_id":    user.Student.StudentID,
		"full_name":     user.User.FullName,
		"email":         user.User.Email,
		"username":      user.User.Username,
		"program_study": user.Student.ProgramStudy,
		"academic_year": user.Student.AcademicYear,
		"advisor_id":    user.Student.AdvisorID,
		"is_active":     user.User.IsActive,
		"created_at":    user.Student.CreatedAt,
		"updated_at":    user.Student.UpdatedAt,
		"role":          user.Role,
	}

	return c.Status(200).JSON(fiber.Map{
		"success": true,
		"message": "Student retrieved successfully",
		"data":    studentData,
	})
}

// GetStudentAchievements - GET /api/v1/students/:id/achievements
func (h *V1StudentLecturerHandler) GetStudentAchievements(c *fiber.Ctx) error {
	studentID, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"success": false,
			"error":   "Invalid student ID format",
		})
	}

	// Parse query parameters
	page, _ := strconv.Atoi(c.Query("page", "1"))
	limit, _ := strconv.Atoi(c.Query("limit", "10"))
	status := c.Query("status")
	achievementType := c.Query("type")

	// Check if student exists
	user, err := h.UserService.GetUserByID(studentID)
	if err != nil || user.Student == nil {
		return c.Status(404).JSON(fiber.Map{
			"success": false,
			"error":   "Student not found",
		})
	}

	// Mock achievements data - in real implementation, this would fetch from AchievementService
	achievements := []fiber.Map{
		{
			"id":               uuid.New(),
			"title":            "Programming Competition Winner",
			"achievement_type": "competition",
			"status":           "verified",
			"points":           100.0,
			"created_at":       "2024-01-15T10:00:00Z",
		},
		{
			"id":               uuid.New(),
			"title":            "Research Publication",
			"achievement_type": "publication",
			"status":           "submitted",
			"points":           150.0,
			"created_at":       "2024-01-20T14:30:00Z",
		},
	}

	// Apply filters
	filteredAchievements := []fiber.Map{}
	for _, achievement := range achievements {
		include := true
		if status != "" && achievement["status"] != status {
			include = false
		}
		if achievementType != "" && achievement["achievement_type"] != achievementType {
			include = false
		}
		if include {
			filteredAchievements = append(filteredAchievements, achievement)
		}
	}

	totalPages := (len(filteredAchievements) + limit - 1) / limit

	return c.Status(200).JSON(fiber.Map{
		"success": true,
		"message": "Student achievements retrieved successfully",
		"data": fiber.Map{
			"student": fiber.Map{
				"id":            user.Student.ID,
				"student_id":    user.Student.StudentID,
				"full_name":     user.User.FullName,
				"program_study": user.Student.ProgramStudy,
				"academic_year": user.Student.AcademicYear,
			},
			"achievements": filteredAchievements,
		},
		"pagination": fiber.Map{
			"page":        page,
			"limit":       limit,
			"total":       len(filteredAchievements),
			"total_pages": totalPages,
		},
		"filters": fiber.Map{
			"status": status,
			"type":   achievementType,
		},
	})
}

// SetStudentAdvisor - PUT /api/v1/students/:id/advisor
func (h *V1StudentLecturerHandler) SetStudentAdvisor(c *fiber.Ctx) error {
	studentID, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"success": false,
			"error":   "Invalid student ID format",
		})
	}

	var req struct {
		AdvisorID uuid.UUID `json:"advisor_id"`
	}

	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"success": false,
			"error":   "Invalid JSON format",
		})
	}

	if req.AdvisorID == uuid.Nil {
		return c.Status(400).JSON(fiber.Map{
			"success": false,
			"error":   "Advisor ID is required",
		})
	}

	// Set advisor using UserService
	student, err := h.UserService.SetAdvisor(studentID, req.AdvisorID)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"success": false,
			"error":   err.Error(),
		})
	}

	return c.Status(200).JSON(fiber.Map{
		"success": true,
		"message": "Student advisor set successfully",
		"data": fiber.Map{
			"student_id": student.ID,
			"advisor_id": student.AdvisorID,
		},
	})
}

// GetAllLecturers - GET /api/v1/lecturers
func (h *V1StudentLecturerHandler) GetAllLecturers(c *fiber.Ctx) error {
	// Parse query parameters
	page, _ := strconv.Atoi(c.Query("page", "1"))
	limit, _ := strconv.Atoi(c.Query("limit", "10"))
	department := c.Query("department")

	offset := (page - 1) * limit

	// Get all users with lecturer role
	users, total, err := h.UserService.GetAllUsers(limit, offset)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"success": false,
			"error":   "Failed to retrieve lecturers",
		})
	}

	// Filter only lecturers
	lecturers := []fiber.Map{}
	for _, user := range users {
		if user.Lecturer != nil {
			lecturerData := fiber.Map{
				"id":           user.Lecturer.ID,
				"user_id":      user.User.ID,
				"lecturer_id":  user.Lecturer.LecturerID,
				"full_name":    user.User.FullName,
				"email":        user.User.Email,
				"department":   user.Lecturer.Department,
				"is_active":    user.User.IsActive,
				"created_at":   user.Lecturer.CreatedAt,
			}

			// Apply department filter
			if department == "" || user.Lecturer.Department == department {
				lecturers = append(lecturers, lecturerData)
			}
		}
	}

	totalPages := (len(lecturers) + limit - 1) / limit

	return c.Status(200).JSON(fiber.Map{
		"success": true,
		"message": "Lecturers retrieved successfully",
		"data":    lecturers,
		"pagination": fiber.Map{
			"page":        page,
			"limit":       limit,
			"total":       len(lecturers),
			"total_pages": totalPages,
		},
		"filters": fiber.Map{
			"department": department,
		},
	})
}

// GetLecturerAdvisees - GET /api/v1/lecturers/:id/advisees
func (h *V1StudentLecturerHandler) GetLecturerAdvisees(c *fiber.Ctx) error {
	lecturerID, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"success": false,
			"error":   "Invalid lecturer ID format",
		})
	}

	// Parse query parameters
	page, _ := strconv.Atoi(c.Query("page", "1"))
	limit, _ := strconv.Atoi(c.Query("limit", "10"))

	// Check if lecturer exists
	user, err := h.UserService.GetUserByID(lecturerID)
	if err != nil || user.Lecturer == nil {
		return c.Status(404).JSON(fiber.Map{
			"success": false,
			"error":   "Lecturer not found",
		})
	}

	// Get all students and filter by advisor_id
	allUsers, _, err := h.UserService.GetAllUsers(1000, 0) // Get all users
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"success": false,
			"error":   "Failed to retrieve advisees",
		})
	}

	advisees := []fiber.Map{}
	for _, userResp := range allUsers {
		if userResp.Student != nil && userResp.Student.AdvisorID == lecturerID {
			adviseeData := fiber.Map{
				"id":            userResp.Student.ID,
				"user_id":       userResp.User.ID,
				"student_id":    userResp.Student.StudentID,
				"full_name":     userResp.User.FullName,
				"email":         userResp.User.Email,
				"program_study": userResp.Student.ProgramStudy,
				"academic_year": userResp.Student.AcademicYear,
				"is_active":     userResp.User.IsActive,
				"created_at":    userResp.Student.CreatedAt,
			}
			advisees = append(advisees, adviseeData)
		}
	}

	// Apply pagination
	start := (page - 1) * limit
	end := start + limit
	if end > len(advisees) {
		end = len(advisees)
	}

	paginatedAdvisees := []fiber.Map{}
	if start < len(advisees) {
		paginatedAdvisees = advisees[start:end]
	}

	totalPages := (len(advisees) + limit - 1) / limit

	return c.Status(200).JSON(fiber.Map{
		"success": true,
		"message": "Lecturer advisees retrieved successfully",
		"data": fiber.Map{
			"lecturer": fiber.Map{
				"id":           user.Lecturer.ID,
				"lecturer_id":  user.Lecturer.LecturerID,
				"full_name":    user.User.FullName,
				"department":   user.Lecturer.Department,
			},
			"advisees": paginatedAdvisees,
		},
		"pagination": fiber.Map{
			"page":        page,
			"limit":       limit,
			"total":       len(advisees),
			"total_pages": totalPages,
		},
	})
}