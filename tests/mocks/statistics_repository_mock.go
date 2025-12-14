package mocks

import (
	model "UAS_BACKEND/domain/model"
	"UAS_BACKEND/domain/service"
	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
)

// MockStatisticsRepository is a mock implementation of StatisticsRepository
type MockStatisticsRepository struct {
	mock.Mock
}

func (m *MockStatisticsRepository) GetStudentByUserID(userID uuid.UUID) (*model.Student, error) {
	args := m.Called(userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.Student), args.Error(1)
}

func (m *MockStatisticsRepository) GetLecturerByUserID(userID uuid.UUID) (*model.Lecturer, error) {
	args := m.Called(userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.Lecturer), args.Error(1)
}

func (m *MockStatisticsRepository) GetAchievementTypeStats(req *service.StatisticsRequest) ([]service.AchievementTypeStats, error) {
	args := m.Called(req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]service.AchievementTypeStats), args.Error(1)
}

func (m *MockStatisticsRepository) GetAchievementPeriodStats(req *service.StatisticsRequest) ([]service.AchievementPeriodStats, error) {
	args := m.Called(req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]service.AchievementPeriodStats), args.Error(1)
}

func (m *MockStatisticsRepository) GetTopStudentStats(req *service.StatisticsRequest) ([]service.TopStudentStats, error) {
	args := m.Called(req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]service.TopStudentStats), args.Error(1)
}

func (m *MockStatisticsRepository) GetCompetitionLevelStats(ctx interface{}, req *service.StatisticsRequest) ([]service.CompetitionLevelStats, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]service.CompetitionLevelStats), args.Error(1)
}

func (m *MockStatisticsRepository) GetStatisticsSummary(req *service.StatisticsRequest) (*service.StatisticsSummary, error) {
	args := m.Called(req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*service.StatisticsSummary), args.Error(1)
}

func (m *MockStatisticsRepository) GetAchievementTrends(req *service.StatisticsRequest, months int) ([]service.TrendData, error) {
	args := m.Called(req, months)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]service.TrendData), args.Error(1)
}