package service_test

import (
	model "UAS_BACKEND/domain/model"
	"UAS_BACKEND/domain/service"
	"UAS_BACKEND/tests/mocks"
	"errors"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestNotificationService_CreateAchievementSubmittedNotification_Success(t *testing.T) {
	// Arrange
	mockRepo := new(mocks.MockNotificationRepository)
	notificationService := service.NewNotificationService(mockRepo)

	advisorID := uuid.New()
	referenceID := uuid.New()
	studentName := "John Doe"
	achievementTitle := "Programming Competition Winner"

	mockRepo.On("CreateNotification", mock.MatchedBy(func(notification *model.Notification) bool {
		return notification.UserID == advisorID &&
			notification.Type == "achievement_submitted" &&
			notification.Title == "Prestasi Baru Menunggu Verifikasi" &&
			notification.RelatedID != nil &&
			*notification.RelatedID == referenceID
	})).Return(nil)

	// Act
	err := notificationService.CreateAchievementSubmittedNotification(advisorID, studentName, achievementTitle, referenceID)

	// Assert
	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

func TestNotificationService_CreateAchievementSubmittedNotification_RepositoryError(t *testing.T) {
	// Arrange
	mockRepo := new(mocks.MockNotificationRepository)
	notificationService := service.NewNotificationService(mockRepo)

	advisorID := uuid.New()
	referenceID := uuid.New()
	studentName := "John Doe"
	achievementTitle := "Programming Competition Winner"

	mockRepo.On("CreateNotification", mock.AnythingOfType("*model.Notification")).Return(errors.New("database error"))

	// Act
	err := notificationService.CreateAchievementSubmittedNotification(advisorID, studentName, achievementTitle, referenceID)

	// Assert
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "database error")
	mockRepo.AssertExpectations(t)
}

func TestNotificationService_CreateAchievementVerifiedNotification_Success(t *testing.T) {
	// Arrange
	mockRepo := new(mocks.MockNotificationRepository)
	notificationService := service.NewNotificationService(mockRepo)

	studentUserID := uuid.New()
	referenceID := uuid.New()
	lecturerName := "Dr. Smith"
	achievementTitle := "Programming Competition Winner"

	mockRepo.On("CreateNotification", mock.MatchedBy(func(notification *model.Notification) bool {
		return notification.UserID == studentUserID &&
			notification.Type == "achievement_verified" &&
			notification.Title == "Prestasi Diverifikasi" &&
			notification.RelatedID != nil &&
			*notification.RelatedID == referenceID
	})).Return(nil)

	// Act
	err := notificationService.CreateAchievementVerifiedNotification(studentUserID, lecturerName, achievementTitle, referenceID)

	// Assert
	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

func TestNotificationService_CreateAchievementRejectedNotification_Success(t *testing.T) {
	// Arrange
	mockRepo := new(mocks.MockNotificationRepository)
	notificationService := service.NewNotificationService(mockRepo)

	studentUserID := uuid.New()
	referenceID := uuid.New()
	lecturerName := "Dr. Smith"
	achievementTitle := "Programming Competition Winner"
	rejectionNote := "Documentation incomplete"

	mockRepo.On("CreateNotification", mock.MatchedBy(func(notification *model.Notification) bool {
		return notification.UserID == studentUserID &&
			notification.Type == "achievement_rejected" &&
			notification.Title == "Prestasi Ditolak" &&
			notification.RelatedID != nil &&
			*notification.RelatedID == referenceID
	})).Return(nil)

	// Act
	err := notificationService.CreateAchievementRejectedNotification(studentUserID, lecturerName, achievementTitle, referenceID, rejectionNote)

	// Assert
	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

func TestNotificationService_CreateAchievementRejectedNotification_EmptyRejectionNote(t *testing.T) {
	// Arrange
	mockRepo := new(mocks.MockNotificationRepository)
	notificationService := service.NewNotificationService(mockRepo)

	studentUserID := uuid.New()
	referenceID := uuid.New()
	lecturerName := "Dr. Smith"
	achievementTitle := "Programming Competition Winner"
	rejectionNote := ""

	mockRepo.On("CreateNotification", mock.MatchedBy(func(notification *model.Notification) bool {
		return notification.UserID == studentUserID &&
			notification.Type == "achievement_rejected" &&
			notification.Title == "Prestasi Ditolak" &&
			notification.RelatedID != nil &&
			*notification.RelatedID == referenceID &&
			// Should not contain "Alasan:" when rejection note is empty
			!contains(notification.Message, "Alasan:")
	})).Return(nil)

	// Act
	err := notificationService.CreateAchievementRejectedNotification(studentUserID, lecturerName, achievementTitle, referenceID, rejectionNote)

	// Assert
	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

func TestNotificationService_GetUserNotifications_Success(t *testing.T) {
	// Arrange
	mockRepo := new(mocks.MockNotificationRepository)
	notificationService := service.NewNotificationService(mockRepo)

	userID := uuid.New()
	limit := 10

	expectedNotifications := []model.Notification{
		{
			ID:        uuid.New(),
			UserID:    userID,
			Type:      "achievement_submitted",
			Title:     "Prestasi Baru Menunggu Verifikasi",
			Message:   "Mahasiswa John Doe telah mengajukan prestasi 'Test Achievement' untuk diverifikasi.",
			IsRead:    false,
			CreatedAt: time.Now(),
		},
		{
			ID:        uuid.New(),
			UserID:    userID,
			Type:      "achievement_verified",
			Title:     "Prestasi Diverifikasi",
			Message:   "Prestasi Anda 'Test Achievement' telah diverifikasi dan disetujui oleh Dr. Smith.",
			IsRead:    true,
			CreatedAt: time.Now().Add(-1 * time.Hour),
		},
	}

	mockRepo.On("GetUserNotifications", userID, limit).Return(expectedNotifications, nil)

	// Act
	result, err := notificationService.GetUserNotifications(userID, limit)

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, 2, len(result))
	assert.Equal(t, expectedNotifications[0].ID, result[0].ID)
	assert.Equal(t, expectedNotifications[0].Type, result[0].Type)
	assert.Equal(t, expectedNotifications[1].ID, result[1].ID)
	assert.Equal(t, expectedNotifications[1].Type, result[1].Type)

	mockRepo.AssertExpectations(t)
}

func TestNotificationService_GetUserNotifications_DefaultLimit(t *testing.T) {
	// Arrange
	mockRepo := new(mocks.MockNotificationRepository)
	notificationService := service.NewNotificationService(mockRepo)

	userID := uuid.New()
	limit := 0 // Should default to 50

	mockRepo.On("GetUserNotifications", userID, 50).Return([]model.Notification{}, nil)

	// Act
	result, err := notificationService.GetUserNotifications(userID, limit)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, result)

	mockRepo.AssertExpectations(t)
}

func TestNotificationService_GetUserNotifications_RepositoryError(t *testing.T) {
	// Arrange
	mockRepo := new(mocks.MockNotificationRepository)
	notificationService := service.NewNotificationService(mockRepo)

	userID := uuid.New()
	limit := 10

	mockRepo.On("GetUserNotifications", userID, limit).Return(nil, errors.New("database error"))

	// Act
	result, err := notificationService.GetUserNotifications(userID, limit)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "database error")

	mockRepo.AssertExpectations(t)
}

func TestNotificationService_GetUnreadCount_Success(t *testing.T) {
	// Arrange
	mockRepo := new(mocks.MockNotificationRepository)
	notificationService := service.NewNotificationService(mockRepo)

	userID := uuid.New()
	expectedCount := 5

	mockRepo.On("GetUnreadCount", userID).Return(expectedCount, nil)

	// Act
	result, err := notificationService.GetUnreadCount(userID)

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, expectedCount, result)

	mockRepo.AssertExpectations(t)
}

func TestNotificationService_GetUnreadCount_RepositoryError(t *testing.T) {
	// Arrange
	mockRepo := new(mocks.MockNotificationRepository)
	notificationService := service.NewNotificationService(mockRepo)

	userID := uuid.New()

	mockRepo.On("GetUnreadCount", userID).Return(0, errors.New("database error"))

	// Act
	result, err := notificationService.GetUnreadCount(userID)

	// Assert
	assert.Error(t, err)
	assert.Equal(t, 0, result)
	assert.Contains(t, err.Error(), "database error")

	mockRepo.AssertExpectations(t)
}

func TestNotificationService_MarkAsRead_Success(t *testing.T) {
	// Arrange
	mockRepo := new(mocks.MockNotificationRepository)
	notificationService := service.NewNotificationService(mockRepo)

	notificationID := uuid.New()

	mockRepo.On("MarkAsRead", notificationID).Return(nil)

	// Act
	err := notificationService.MarkAsRead(notificationID)

	// Assert
	assert.NoError(t, err)

	mockRepo.AssertExpectations(t)
}

func TestNotificationService_MarkAsRead_RepositoryError(t *testing.T) {
	// Arrange
	mockRepo := new(mocks.MockNotificationRepository)
	notificationService := service.NewNotificationService(mockRepo)

	notificationID := uuid.New()

	mockRepo.On("MarkAsRead", notificationID).Return(errors.New("database error"))

	// Act
	err := notificationService.MarkAsRead(notificationID)

	// Assert
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "database error")

	mockRepo.AssertExpectations(t)
}

func TestNotificationService_MarkAllAsRead_Success(t *testing.T) {
	// Arrange
	mockRepo := new(mocks.MockNotificationRepository)
	notificationService := service.NewNotificationService(mockRepo)

	userID := uuid.New()

	mockRepo.On("MarkAllAsRead", userID).Return(nil)

	// Act
	err := notificationService.MarkAllAsRead(userID)

	// Assert
	assert.NoError(t, err)

	mockRepo.AssertExpectations(t)
}

func TestNotificationService_MarkAllAsRead_RepositoryError(t *testing.T) {
	// Arrange
	mockRepo := new(mocks.MockNotificationRepository)
	notificationService := service.NewNotificationService(mockRepo)

	userID := uuid.New()

	mockRepo.On("MarkAllAsRead", userID).Return(errors.New("database error"))

	// Act
	err := notificationService.MarkAllAsRead(userID)

	// Assert
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "database error")

	mockRepo.AssertExpectations(t)
}

// Helper function to check if a string contains a substring
func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(substr) == 0 || (len(s) > len(substr) && (s[:len(substr)] == substr || s[len(s)-len(substr):] == substr || containsAt(s, substr))))
}

func containsAt(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}