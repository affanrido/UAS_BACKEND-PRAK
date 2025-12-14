package mocks

import (
	model "UAS_BACKEND/domain/model"
	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
)

// MockNotificationRepository is a mock implementation of NotificationRepository
type MockNotificationRepository struct {
	mock.Mock
}

func (m *MockNotificationRepository) CreateNotification(notification *model.Notification) error {
	args := m.Called(notification)
	return args.Error(0)
}

func (m *MockNotificationRepository) GetUserNotifications(userID uuid.UUID, limit int) ([]model.Notification, error) {
	args := m.Called(userID, limit)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]model.Notification), args.Error(1)
}

func (m *MockNotificationRepository) GetUnreadCount(userID uuid.UUID) (int, error) {
	args := m.Called(userID)
	return args.Int(0), args.Error(1)
}

func (m *MockNotificationRepository) MarkAsRead(notificationID uuid.UUID) error {
	args := m.Called(notificationID)
	return args.Error(0)
}

func (m *MockNotificationRepository) MarkAllAsRead(userID uuid.UUID) error {
	args := m.Called(userID)
	return args.Error(0)
}