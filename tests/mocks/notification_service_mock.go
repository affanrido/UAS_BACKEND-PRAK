package mocks

import (
	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
)

// MockNotificationService is a mock implementation of NotificationService
type MockNotificationService struct {
	mock.Mock
}

func (m *MockNotificationService) CreateAchievementSubmittedNotification(advisorID uuid.UUID, studentName string, achievementTitle string, referenceID uuid.UUID) error {
	args := m.Called(advisorID, studentName, achievementTitle, referenceID)
	return args.Error(0)
}

func (m *MockNotificationService) CreateAchievementVerifiedNotification(studentUserID uuid.UUID, lecturerName, achievementTitle string, referenceID uuid.UUID) error {
	args := m.Called(studentUserID, lecturerName, achievementTitle, referenceID)
	return args.Error(0)
}

func (m *MockNotificationService) CreateAchievementRejectedNotification(studentUserID uuid.UUID, lecturerName, achievementTitle string, referenceID uuid.UUID, rejectionNote string) error {
	args := m.Called(studentUserID, lecturerName, achievementTitle, referenceID, rejectionNote)
	return args.Error(0)
}

func (m *MockNotificationService) CreateAchievementVerifiedNotification(studentUserID uuid.UUID, lecturerName, achievementTitle string, referenceID uuid.UUID) error {
	args := m.Called(studentUserID, lecturerName, achievementTitle, referenceID)
	return args.Error(0)
}

func (m *MockNotificationService) CreateAchievementRejectedNotification(studentUserID uuid.UUID, lecturerName, achievementTitle string, referenceID uuid.UUID, rejectionNote string) error {
	args := m.Called(studentUserID, lecturerName, achievementTitle, referenceID, rejectionNote)
	return args.Error(0)
}