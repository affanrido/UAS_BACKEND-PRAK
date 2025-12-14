package route

import (
	"UAS_BACKEND/domain/middleware"
	"UAS_BACKEND/domain/model"
	"UAS_BACKEND/domain/service"
	"context"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type V1ReportHandler struct {
	StatisticsService *service.StatisticsService
	RBACMiddleware    *middleware.RBACMiddleware
}

func NewV1ReportHandler(
	statisticsService *service.StatisticsService,
	rbacMiddleware *middleware.RBACMiddleware,
) *V1ReportHandler {
	return &V1ReportHandler{
		StatisticsService: statisticsService,
		RBACMiddleware:    rbacMiddleware,
	}
}

// SetupV1ReportRoutes - Setup report and analytics routes v1
func SetupV1ReportRoutes(app *fiber.App, handler *V1ReportHandler) {
	reports := app.Group("/api/v1/reports")
	reports.Use(handler.RBACMiddleware.RequireAuth())

	// 5.8 Reports & Analytics endpoints
	reports.Get("/statistics", handler.RBACMiddleware.RequireAnyPermission("admin.manage", "lecturer.read", "student.read"), handler.GetStatistics)
	reports.Get("/student/:id", handler.RBACMiddleware.RequireAnyPermission("admin.manage", "lecturer.read"), handler.GetStudentReport)
}

// GetStatistics - GET /api/v1/reports/statistics
func (h *V1ReportHandler) GetStatistics(c *fiber.Ctx) error {
	user := c.Locals("user").(*model.Claims)

	// Parse query parameters
	startDateStr := c.Query("start_date")
	endDateStr := c.Query("end_date")
	reportType := c.Query("type", "overview") // overview, detailed, trends

	// Parse dates if provided
	var startDate, endDate *time.Time
	if startDateStr != "" {
		if parsed, err := time.Parse("2006-01-02", startDateStr); err == nil {
			startDate = &parsed
		}
	}
	if endDateStr != "" {
		if parsed, err := time.Parse("2006-01-02", endDateStr); err == nil {
			endDate = &parsed
		}
	}

	req := &service.StatisticsRequest{
		StartDate: startDate,
		EndDate:   endDate,
	}

	var result *service.StatisticsResponse
	var err error

	// Role-based statistics
	switch {
	case h.hasPermission(user, "admin.manage"):
		result, err = h.StatisticsService.GetAdminStatistics(context.Background(), req)
	case h.hasPermission(user, "lecturer.read"):
		result, err = h.StatisticsService.GetLecturerStatistics(context.Background(), user.UserID, req)
	case h.hasPermission(user, "student.read"):
		result, err = h.StatisticsService.GetStudentStatistics(context.Background(), user.UserID, req)
	default:
		return c.Status(403).JSON(fiber.Map{
			"success": false,
			"error":   "Insufficient permissions",
		})
	}

	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"success": false,
			"error":   "Failed to retrieve statistics: " + err.Error(),
		})
	}

	// Add metadata based on report type
	response := fiber.Map{
		"success": true,
		"message": "Statistics retrieved successfully",
		"data":    result,
		"metadata": fiber.Map{
			"report_type":  reportType,
			"generated_at": time.Now(),
			"user_role":    h.getUserRole(user),
			"date_range": fiber.Map{
				"start_date": startDate,
				"end_date":   endDate,
			},
		},
	}

	// Add additional data based on report type
	switch reportType {
	case "detailed":
		response["additional_metrics"] = h.getDetailedMetrics(result)
	case "trends":
		trends, _ := h.StatisticsService.GetAchievementTrends(context.Background(), user.UserID, h.getUserRole(user), 12)
		response["trends"] = trends
	}

	return c.Status(200).JSON(response)
}

// GetStudentReport - GET /api/v1/reports/student/:id
func (h *V1ReportHandler) GetStudentReport(c *fiber.Ctx) error {
	studentID, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"success": false,
			"error":   "Invalid student ID format",
		})
	}

	// Parse query parameters
	startDateStr := c.Query("start_date")
	endDateStr := c.Query("end_date")
	includeDetails := c.Query("include_details", "false") == "true"

	// Parse dates if provided
	var startDate, endDate *time.Time
	if startDateStr != "" {
		if parsed, err := time.Parse("2006-01-02", startDateStr); err == nil {
			startDate = &parsed
		}
	}
	if endDateStr != "" {
		if parsed, err := time.Parse("2006-01-02", endDateStr); err == nil {
			endDate = &parsed
		}
	}

	req := &service.StatisticsRequest{
		StartDate: startDate,
		EndDate:   endDate,
		UserID:    &studentID,
	}

	// Get student statistics
	result, err := h.StatisticsService.GetStudentStatistics(context.Background(), studentID, req)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"success": false,
			"error":   "Failed to retrieve student report: " + err.Error(),
		})
	}

	response := fiber.Map{
		"success": true,
		"message": "Student report retrieved successfully",
		"data": fiber.Map{
			"student_id":  studentID,
			"statistics":  result,
			"report_date": time.Now(),
		},
		"metadata": fiber.Map{
			"include_details": includeDetails,
			"date_range": fiber.Map{
				"start_date": startDate,
				"end_date":   endDate,
			},
		},
	}

	// Add detailed information if requested
	if includeDetails {
		response["detailed_breakdown"] = h.getStudentDetailedBreakdown(result)
		
		// Get achievement trends for this student
		trends, _ := h.StatisticsService.GetAchievementTrends(context.Background(), studentID, "student", 6)
		response["trends"] = trends
	}

	return c.Status(200).JSON(response)
}

// Helper functions

func (h *V1ReportHandler) hasPermission(user *model.Claims, permission string) bool {
	for _, perm := range user.Permissions {
		if perm == permission {
			return true
		}
	}
	return false
}

func (h *V1ReportHandler) getUserRole(user *model.Claims) string {
	if h.hasPermission(user, "admin.manage") {
		return "admin"
	}
	if h.hasPermission(user, "lecturer.read") {
		return "lecturer"
	}
	if h.hasPermission(user, "student.read") {
		return "student"
	}
	return "unknown"
}

func (h *V1ReportHandler) getDetailedMetrics(stats *service.StatisticsResponse) fiber.Map {
	return fiber.Map{
		"performance_indicators": fiber.Map{
			"completion_rate": h.calculateCompletionRate(stats),
			"average_points_per_achievement": h.calculateAveragePoints(stats),
			"monthly_growth": h.calculateMonthlyGrowth(stats),
		},
		"quality_metrics": fiber.Map{
			"verification_rate": h.calculateVerificationRate(stats),
			"rejection_rate": h.calculateRejectionRate(stats),
		},
		"comparative_analysis": fiber.Map{
			"vs_previous_period": h.compareToPreviousPeriod(stats),
			"ranking": h.calculateRanking(stats),
		},
	}
}

func (h *V1ReportHandler) getStudentDetailedBreakdown(stats *service.StatisticsResponse) fiber.Map {
	return fiber.Map{
		"achievement_breakdown": fiber.Map{
			"by_type": stats.TypeStats,
			"by_period": stats.PeriodStats,
			"by_competition_level": stats.CompetitionStats,
		},
		"performance_metrics": fiber.Map{
			"total_points": stats.Summary.TotalPoints,
			"average_points": stats.Summary.AveragePoints,
			"completion_rate": h.calculateCompletionRate(stats),
		},
		"status_distribution": fiber.Map{
			"verified": stats.Summary.VerifiedCount,
			"pending": stats.Summary.PendingCount,
			"total": stats.Summary.TotalAchievements,
		},
	}
}

// Calculation helper functions (mock implementations)
func (h *V1ReportHandler) calculateCompletionRate(stats *service.StatisticsResponse) float64 {
	if stats.Summary.TotalAchievements == 0 {
		return 0.0
	}
	return float64(stats.Summary.VerifiedCount) / float64(stats.Summary.TotalAchievements) * 100
}

func (h *V1ReportHandler) calculateAveragePoints(stats *service.StatisticsResponse) float64 {
	return stats.Summary.AveragePoints
}

func (h *V1ReportHandler) calculateMonthlyGrowth(stats *service.StatisticsResponse) float64 {
	// Mock calculation - in real implementation, compare with previous month
	return 15.5 // 15.5% growth
}

func (h *V1ReportHandler) calculateVerificationRate(stats *service.StatisticsResponse) float64 {
	if stats.Summary.TotalAchievements == 0 {
		return 0.0
	}
	return float64(stats.Summary.VerifiedCount) / float64(stats.Summary.TotalAchievements) * 100
}

func (h *V1ReportHandler) calculateRejectionRate(stats *service.StatisticsResponse) float64 {
	// Mock calculation - in real implementation, calculate from rejected achievements
	rejectedCount := stats.Summary.TotalAchievements - stats.Summary.VerifiedCount - stats.Summary.PendingCount
	if stats.Summary.TotalAchievements == 0 {
		return 0.0
	}
	return float64(rejectedCount) / float64(stats.Summary.TotalAchievements) * 100
}

func (h *V1ReportHandler) compareToPreviousPeriod(stats *service.StatisticsResponse) fiber.Map {
	// Mock comparison data
	return fiber.Map{
		"total_achievements": fiber.Map{
			"current": stats.Summary.TotalAchievements,
			"previous": 18,
			"change": "+38.9%",
		},
		"total_points": fiber.Map{
			"current": stats.Summary.TotalPoints,
			"previous": 1800.0,
			"change": "+38.9%",
		},
	}
}

func (h *V1ReportHandler) calculateRanking(stats *service.StatisticsResponse) fiber.Map {
	// Mock ranking data
	return fiber.Map{
		"current_rank": 5,
		"total_students": 150,
		"percentile": 96.7,
		"category": "Top Performer",
	}
}