package service

import (
	model "UAS_BACKEND/domain/Model"
	"UAS_BACKEND/domain/repository"
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
)

type StatisticsService struct {
	Repo *repository.StatisticsRepository
}

func NewStatisticsService(repo *repository.StatisticsRepository) *StatisticsService {
	return &StatisticsService{Repo: repo}
}

// StatisticsRequest - DTO untuk request statistics
type StatisticsRequest struct {
	StartDate *time.Time `json:"start_date,omitempty"`
	EndDate   *time.Time `json:"end_date,omitempty"`
	UserID    *uuid.UUID `json:"user_id,omitempty"`    // For student own stats
	AdvisorID *uuid.UUID `json:"advisor_id,omitempty"` // For lecturer advisee stats
}

// AchievementTypeStats - Statistik per tipe prestasi
type AchievementTypeStats struct {
	Type  string `json:"type"`
	Count int    `json:"count"`
	Total int    `json:"total"`
}

// AchievementPeriodStats - Statistik per periode
type AchievementPeriodStats struct {
	Period string `json:"period"` // Format: "2024-01" atau "2024"
	Count  int    `json:"count"`
	Total  int    `json:"total"`
}

// TopStudentStats - Top mahasiswa berprestasi
type TopStudentStats struct {
	StudentID    string  `json:"student_id"`
	StudentName  string  `json:"student_name"`
	ProgramStudy string  `json:"program_study"`
	AcademicYear string  `json:"academic_year"`
	TotalPoints  float64 `json:"total_points"`
	TotalCount   int     `json:"total_count"`
}

// CompetitionLevelStats - Distribusi tingkat kompetisi
type CompetitionLevelStats struct {
	Level string `json:"level"` // 'international', 'national', 'regional', 'local'
	Count int    `json:"count"`
	Total int    `json:"total"`
}

// StatisticsResponse - Response untuk statistics
type StatisticsResponse struct {
	TypeStats        []AchievementTypeStats    `json:"type_stats"`
	PeriodStats      []AchievementPeriodStats  `json:"period_stats"`
	TopStudents      []TopStudentStats         `json:"top_students"`
	CompetitionStats []CompetitionLevelStats   `json:"competition_stats"`
	Summary          *StatisticsSummary        `json:"summary"`
}

// StatisticsSummary - Summary statistics
type StatisticsSummary struct {
	TotalAchievements int     `json:"total_achievements"`
	TotalPoints       float64 `json:"total_points"`
	AveragePoints     float64 `json:"average_points"`
	VerifiedCount     int     `json:"verified_count"`
	PendingCount      int     `json:"pending_count"`
}

// GetStudentStatistics - Get statistics for student (own achievements)
func (s *StatisticsService) GetStudentStatistics(ctx context.Context, userID uuid.UUID, req *StatisticsRequest) (*StatisticsResponse, error) {
	// Validasi: User harus mahasiswa
	student, err := s.Repo.GetStudentByUserID(userID)
	if err != nil {
		return nil, errors.New("user is not a student")
	}

	// Set filter untuk mahasiswa ini saja
	req.UserID = &student.UserID

	return s.generateStatistics(ctx, req, "student")
}

// GetLecturerStatistics - Get statistics for lecturer (advisee achievements)
func (s *StatisticsService) GetLecturerStatistics(ctx context.Context, userID uuid.UUID, req *StatisticsRequest) (*StatisticsResponse, error) {
	// Validasi: User harus dosen
	lecturer, err := s.Repo.GetLecturerByUserID(userID)
	if err != nil {
		return nil, errors.New("user is not a lecturer")
	}

	// Set filter untuk mahasiswa bimbingan
	req.AdvisorID = &lecturer.ID

	return s.generateStatistics(ctx, req, "lecturer")
}

// GetAdminStatistics - Get statistics for admin (all achievements)
func (s *StatisticsService) GetAdminStatistics(ctx context.Context, req *StatisticsRequest) (*StatisticsResponse, error) {
	return s.generateStatistics(ctx, req, "admin")
}

// generateStatistics - Generate statistics based on request
func (s *StatisticsService) generateStatistics(ctx context.Context, req *StatisticsRequest, role string) (*StatisticsResponse, error) {
	// 1. Total prestasi per tipe
	typeStats, err := s.Repo.GetAchievementTypeStats(req)
	if err != nil {
		return nil, errors.New("failed to get type statistics: " + err.Error())
	}

	// 2. Total prestasi per periode
	periodStats, err := s.Repo.GetAchievementPeriodStats(req)
	if err != nil {
		return nil, errors.New("failed to get period statistics: " + err.Error())
	}

	// 3. Top mahasiswa berprestasi (hanya untuk lecturer dan admin)
	var topStudents []TopStudentStats
	if role != "student" {
		topStudents, err = s.Repo.GetTopStudentStats(req)
		if err != nil {
			return nil, errors.New("failed to get top student statistics: " + err.Error())
		}
	}

	// 4. Distribusi tingkat kompetisi
	competitionStats, err := s.Repo.GetCompetitionLevelStats(ctx, req)
	if err != nil {
		return nil, errors.New("failed to get competition statistics: " + err.Error())
	}

	// 5. Summary statistics
	summary, err := s.Repo.GetStatisticsSummary(req)
	if err != nil {
		return nil, errors.New("failed to get summary statistics: " + err.Error())
	}

	return &StatisticsResponse{
		TypeStats:        typeStats,
		PeriodStats:      periodStats,
		TopStudents:      topStudents,
		CompetitionStats: competitionStats,
		Summary:          summary,
	}, nil
}

// GetAchievementTrends - Get achievement trends over time
func (s *StatisticsService) GetAchievementTrends(ctx context.Context, userID uuid.UUID, role string, months int) (*TrendResponse, error) {
	if months <= 0 {
		months = 12 // Default 12 months
	}
	if months > 24 {
		months = 24 // Max 24 months
	}

	var req *StatisticsRequest
	switch role {
	case "student":
		student, err := s.Repo.GetStudentByUserID(userID)
		if err != nil {
			return nil, errors.New("user is not a student")
		}
		req = &StatisticsRequest{UserID: &student.UserID}
	case "lecturer":
		lecturer, err := s.Repo.GetLecturerByUserID(userID)
		if err != nil {
			return nil, errors.New("user is not a lecturer")
		}
		req = &StatisticsRequest{AdvisorID: &lecturer.ID}
	case "admin":
		req = &StatisticsRequest{}
	default:
		return nil, errors.New("invalid role")
	}

	trends, err := s.Repo.GetAchievementTrends(req, months)
	if err != nil {
		return nil, errors.New("failed to get achievement trends: " + err.Error())
	}

	return &TrendResponse{
		Trends: trends,
		Period: months,
	}, nil
}

// TrendData - Data untuk trend
type TrendData struct {
	Month       string `json:"month"`       // Format: "2024-01"
	Count       int    `json:"count"`       // Jumlah prestasi
	Points      float64 `json:"points"`     // Total points
	Verified    int    `json:"verified"`    // Jumlah verified
	Submitted   int    `json:"submitted"`   // Jumlah submitted
}

// TrendResponse - Response untuk trends
type TrendResponse struct {
	Trends []TrendData `json:"trends"`
	Period int         `json:"period"` // Jumlah bulan
}