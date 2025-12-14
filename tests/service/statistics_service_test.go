package service_test

import (
	model "UAS_BACKEND/domain/model"
	"UAS_BACKEND/domain/service"
	"UAS_BACKEND/tests/mocks"
	"context"
	"errors"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestStatisticsService_GetStudentStatistics_Success(t *testing.T) {
	// Arrange
	mockRepo := new(mocks.MockStatisticsRepository)
	statsService := service.NewStatisticsService(mockRepo)

	userID := uuid.New()
	studentID := uuid.New()

	student := &model.Student{
		ID:           studentID,
		UserID:       userID,
		StudentID:    "STD001",
		ProgramStudy: "Computer Science",
		AcademicYear: "2021",
	}

	req := &service.StatisticsRequest{
		StartDate: nil,
		EndDate:   nil,
	}

	typeStats := []service.AchievementTypeStats{
		{Type: "competition", Count: 5, Total: 5},
		{Type: "certification", Count: 3, Total: 3},
	}

	periodStats := []service.AchievementPeriodStats{
		{Period: "2024-01", Count: 2, Total: 2},
		{Period: "2024-02", Count: 3, Total: 3},
	}

	competitionStats := []service.CompetitionLevelStats{
		{Level: "national", Count: 3, Total: 3},
		{Level: "international", Count: 2, Total: 2},
	}

	summary := &service.StatisticsSummary{
		TotalAchievements: 8,
		TotalPoints:       800.0,
		AveragePoints:     100.0,
		VerifiedCount:     6,
		PendingCount:      2,
	}

	mockRepo.On("GetStudentByUserID", userID).Return(student, nil)
	mockRepo.On("GetAchievementTypeStats", req).Return(typeStats, nil)
	mockRepo.On("GetAchievementPeriodStats", req).Return(periodStats, nil)
	mockRepo.On("GetCompetitionLevelStats", context.Background(), req).Return(competitionStats, nil)
	mockRepo.On("GetStatisticsSummary", req).Return(summary, nil)

	// Act
	result, err := statsService.GetStudentStatistics(context.Background(), userID, req)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, 2, len(result.TypeStats))
	assert.Equal(t, "competition", result.TypeStats[0].Type)
	assert.Equal(t, 5, result.TypeStats[0].Count)
	assert.Equal(t, 2, len(result.PeriodStats))
	assert.Equal(t, "2024-01", result.PeriodStats[0].Period)
	assert.Equal(t, 2, len(result.CompetitionStats))
	assert.Equal(t, "national", result.CompetitionStats[0].Level)
	assert.Equal(t, summary, result.Summary)
	assert.Empty(t, result.TopStudents) // Student role should not have top students

	mockRepo.AssertExpectations(t)
}

func TestStatisticsService_GetStudentStatistics_UserNotStudent(t *testing.T) {
	// Arrange
	mockRepo := new(mocks.MockStatisticsRepository)
	statsService := service.NewStatisticsService(mockRepo)

	userID := uuid.New()
	req := &service.StatisticsRequest{}

	mockRepo.On("GetStudentByUserID", userID).Return(nil, errors.New("user is not a student"))

	// Act
	result, err := statsService.GetStudentStatistics(context.Background(), userID, req)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "user is not a student")

	mockRepo.AssertExpectations(t)
}

func TestStatisticsService_GetLecturerStatistics_Success(t *testing.T) {
	// Arrange
	mockRepo := new(mocks.MockStatisticsRepository)
	statsService := service.NewStatisticsService(mockRepo)

	userID := uuid.New()
	lecturerID := uuid.New()

	lecturer := &model.Lecturer{
		ID:         lecturerID,
		UserID:     userID,
		LecturerID: "LEC001",
		Department: "Computer Science",
	}

	req := &service.StatisticsRequest{
		StartDate: nil,
		EndDate:   nil,
	}

	typeStats := []service.AchievementTypeStats{
		{Type: "competition", Count: 15, Total: 15},
		{Type: "certification", Count: 10, Total: 10},
	}

	periodStats := []service.AchievementPeriodStats{
		{Period: "2024-01", Count: 8, Total: 8},
		{Period: "2024-02", Count: 12, Total: 12},
	}

	topStudents := []service.TopStudentStats{
		{
			StudentID:    "STD001",
			StudentName:  "John Doe",
			ProgramStudy: "Computer Science",
			AcademicYear: "2021",
			TotalPoints:  500.0,
			TotalCount:   5,
		},
		{
			StudentID:    "STD002",
			StudentName:  "Jane Smith",
			ProgramStudy: "Information Systems",
			AcademicYear: "2020",
			TotalPoints:  450.0,
			TotalCount:   4,
		},
	}

	competitionStats := []service.CompetitionLevelStats{
		{Level: "national", Count: 10, Total: 10},
		{Level: "international", Count: 5, Total: 5},
	}

	summary := &service.StatisticsSummary{
		TotalAchievements: 25,
		TotalPoints:       2500.0,
		AveragePoints:     100.0,
		VerifiedCount:     20,
		PendingCount:      5,
	}

	mockRepo.On("GetLecturerByUserID", userID).Return(lecturer, nil)
	mockRepo.On("GetAchievementTypeStats", req).Return(typeStats, nil)
	mockRepo.On("GetAchievementPeriodStats", req).Return(periodStats, nil)
	mockRepo.On("GetTopStudentStats", req).Return(topStudents, nil)
	mockRepo.On("GetCompetitionLevelStats", context.Background(), req).Return(competitionStats, nil)
	mockRepo.On("GetStatisticsSummary", req).Return(summary, nil)

	// Act
	result, err := statsService.GetLecturerStatistics(context.Background(), userID, req)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, 2, len(result.TypeStats))
	assert.Equal(t, "competition", result.TypeStats[0].Type)
	assert.Equal(t, 15, result.TypeStats[0].Count)
	assert.Equal(t, 2, len(result.TopStudents))
	assert.Equal(t, "John Doe", result.TopStudents[0].StudentName)
	assert.Equal(t, 500.0, result.TopStudents[0].TotalPoints)
	assert.Equal(t, summary, result.Summary)

	mockRepo.AssertExpectations(t)
}

func TestStatisticsService_GetLecturerStatistics_UserNotLecturer(t *testing.T) {
	// Arrange
	mockRepo := new(mocks.MockStatisticsRepository)
	statsService := service.NewStatisticsService(mockRepo)

	userID := uuid.New()
	req := &service.StatisticsRequest{}

	mockRepo.On("GetLecturerByUserID", userID).Return(nil, errors.New("user is not a lecturer"))

	// Act
	result, err := statsService.GetLecturerStatistics(context.Background(), userID, req)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "user is not a lecturer")

	mockRepo.AssertExpectations(t)
}

func TestStatisticsService_GetAdminStatistics_Success(t *testing.T) {
	// Arrange
	mockRepo := new(mocks.MockStatisticsRepository)
	statsService := service.NewStatisticsService(mockRepo)

	req := &service.StatisticsRequest{
		StartDate: nil,
		EndDate:   nil,
	}

	typeStats := []service.AchievementTypeStats{
		{Type: "competition", Count: 50, Total: 50},
		{Type: "certification", Count: 30, Total: 30},
		{Type: "research", Count: 20, Total: 20},
	}

	periodStats := []service.AchievementPeriodStats{
		{Period: "2024-01", Count: 25, Total: 25},
		{Period: "2024-02", Count: 35, Total: 35},
		{Period: "2024-03", Count: 40, Total: 40},
	}

	topStudents := []service.TopStudentStats{
		{
			StudentID:    "STD001",
			StudentName:  "John Doe",
			ProgramStudy: "Computer Science",
			AcademicYear: "2021",
			TotalPoints:  1000.0,
			TotalCount:   10,
		},
		{
			StudentID:    "STD002",
			StudentName:  "Jane Smith",
			ProgramStudy: "Information Systems",
			AcademicYear: "2020",
			TotalPoints:  900.0,
			TotalCount:   9,
		},
	}

	competitionStats := []service.CompetitionLevelStats{
		{Level: "international", Count: 20, Total: 20},
		{Level: "national", Count: 30, Total: 30},
		{Level: "regional", Count: 25, Total: 25},
		{Level: "local", Count: 25, Total: 25},
	}

	summary := &service.StatisticsSummary{
		TotalAchievements: 100,
		TotalPoints:       10000.0,
		AveragePoints:     100.0,
		VerifiedCount:     80,
		PendingCount:      20,
	}

	mockRepo.On("GetAchievementTypeStats", req).Return(typeStats, nil)
	mockRepo.On("GetAchievementPeriodStats", req).Return(periodStats, nil)
	mockRepo.On("GetTopStudentStats", req).Return(topStudents, nil)
	mockRepo.On("GetCompetitionLevelStats", context.Background(), req).Return(competitionStats, nil)
	mockRepo.On("GetStatisticsSummary", req).Return(summary, nil)

	// Act
	result, err := statsService.GetAdminStatistics(context.Background(), req)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, 3, len(result.TypeStats))
	assert.Equal(t, "competition", result.TypeStats[0].Type)
	assert.Equal(t, 50, result.TypeStats[0].Count)
	assert.Equal(t, 3, len(result.PeriodStats))
	assert.Equal(t, 2, len(result.TopStudents))
	assert.Equal(t, "John Doe", result.TopStudents[0].StudentName)
	assert.Equal(t, 4, len(result.CompetitionStats))
	assert.Equal(t, "international", result.CompetitionStats[0].Level)
	assert.Equal(t, summary, result.Summary)

	mockRepo.AssertExpectations(t)
}

func TestStatisticsService_GetAchievementTrends_Student_Success(t *testing.T) {
	// Arrange
	mockRepo := new(mocks.MockStatisticsRepository)
	statsService := service.NewStatisticsService(mockRepo)

	userID := uuid.New()
	studentID := uuid.New()

	student := &model.Student{
		ID:           studentID,
		UserID:       userID,
		StudentID:    "STD001",
		ProgramStudy: "Computer Science",
		AcademicYear: "2021",
	}

	trends := []service.TrendData{
		{
			Month:     "2024-01",
			Count:     3,
			Points:    300.0,
			Verified:  2,
			Submitted: 1,
		},
		{
			Month:     "2024-02",
			Count:     5,
			Points:    500.0,
			Verified:  4,
			Submitted: 1,
		},
	}

	mockRepo.On("GetStudentByUserID", userID).Return(student, nil)
	mockRepo.On("GetAchievementTrends", mock.AnythingOfType("*service.StatisticsRequest"), 12).Return(trends, nil)

	// Act
	result, err := statsService.GetAchievementTrends(context.Background(), userID, "student", 12)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, 2, len(result.Trends))
	assert.Equal(t, "2024-01", result.Trends[0].Month)
	assert.Equal(t, 3, result.Trends[0].Count)
	assert.Equal(t, 300.0, result.Trends[0].Points)
	assert.Equal(t, 12, result.Period)

	mockRepo.AssertExpectations(t)
}

func TestStatisticsService_GetAchievementTrends_Lecturer_Success(t *testing.T) {
	// Arrange
	mockRepo := new(mocks.MockStatisticsRepository)
	statsService := service.NewStatisticsService(mockRepo)

	userID := uuid.New()
	lecturerID := uuid.New()

	lecturer := &model.Lecturer{
		ID:         lecturerID,
		UserID:     userID,
		LecturerID: "LEC001",
		Department: "Computer Science",
	}

	trends := []service.TrendData{
		{
			Month:     "2024-01",
			Count:     10,
			Points:    1000.0,
			Verified:  8,
			Submitted: 2,
		},
		{
			Month:     "2024-02",
			Count:     15,
			Points:    1500.0,
			Verified:  12,
			Submitted: 3,
		},
	}

	mockRepo.On("GetLecturerByUserID", userID).Return(lecturer, nil)
	mockRepo.On("GetAchievementTrends", mock.AnythingOfType("*service.StatisticsRequest"), 6).Return(trends, nil)

	// Act
	result, err := statsService.GetAchievementTrends(context.Background(), userID, "lecturer", 6)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, 2, len(result.Trends))
	assert.Equal(t, "2024-01", result.Trends[0].Month)
	assert.Equal(t, 10, result.Trends[0].Count)
	assert.Equal(t, 1000.0, result.Trends[0].Points)
	assert.Equal(t, 6, result.Period)

	mockRepo.AssertExpectations(t)
}

func TestStatisticsService_GetAchievementTrends_InvalidRole(t *testing.T) {
	// Arrange
	mockRepo := new(mocks.MockStatisticsRepository)
	statsService := service.NewStatisticsService(mockRepo)

	userID := uuid.New()

	// Act
	result, err := statsService.GetAchievementTrends(context.Background(), userID, "invalid", 12)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "invalid role")

	mockRepo.AssertExpectations(t)
}

func TestStatisticsService_GetAchievementTrends_DefaultMonths(t *testing.T) {
	// Arrange
	mockRepo := new(mocks.MockStatisticsRepository)
	statsService := service.NewStatisticsService(mockRepo)

	userID := uuid.New()
	studentID := uuid.New()

	student := &model.Student{
		ID:           studentID,
		UserID:       userID,
		StudentID:    "STD001",
		ProgramStudy: "Computer Science",
		AcademicYear: "2021",
	}

	trends := []service.TrendData{}

	mockRepo.On("GetStudentByUserID", userID).Return(student, nil)
	mockRepo.On("GetAchievementTrends", mock.AnythingOfType("*service.StatisticsRequest"), 12).Return(trends, nil)

	// Act - Test with 0 months (should default to 12)
	result, err := statsService.GetAchievementTrends(context.Background(), userID, "student", 0)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, 12, result.Period)

	mockRepo.AssertExpectations(t)
}

func TestStatisticsService_GetAchievementTrends_MaxMonths(t *testing.T) {
	// Arrange
	mockRepo := new(mocks.MockStatisticsRepository)
	statsService := service.NewStatisticsService(mockRepo)

	userID := uuid.New()
	studentID := uuid.New()

	student := &model.Student{
		ID:           studentID,
		UserID:       userID,
		StudentID:    "STD001",
		ProgramStudy: "Computer Science",
		AcademicYear: "2021",
	}

	trends := []service.TrendData{}

	mockRepo.On("GetStudentByUserID", userID).Return(student, nil)
	mockRepo.On("GetAchievementTrends", mock.AnythingOfType("*service.StatisticsRequest"), 24).Return(trends, nil)

	// Act - Test with 30 months (should cap at 24)
	result, err := statsService.GetAchievementTrends(context.Background(), userID, "student", 30)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, 24, result.Period)

	mockRepo.AssertExpectations(t)
}