package service

import (
	model "UAS_BACKEND/domain/Model"
	"UAS_BACKEND/domain/repository"
	"fmt"

	"github.com/google/uuid"
)

type NotificationService struct {
	Repo *repository.NotificationRepository
}

func NewNotificationService(repo *repository.NotificationRepository) *NotificationService {
	return &NotificationService{Repo: repo}
}

// CreateAchievementSubmittedNotification - Buat notifikasi untuk dosen wali saat mahasiswa submit prestasi
func (s *NotificationService) CreateAchievementSubmittedNotification(advisorID uuid.UUID, studentName string, achievementTitle string, referenceID uuid.UUID) error {
	notification := &model.Notification{
		UserID:    advisorID,
		Type:      "achievement_submitted",
		Title:     "Prestasi Baru Menunggu Verifikasi",
		Message:   fmt.Sprintf("Mahasiswa %s telah mengajukan prestasi '%s' untuk diverifikasi.", studentName, achievementTitle),
		RelatedID: &referenceID,
	}

	return s.Repo.CreateNotification(notification)
}

// CreateAchievementVerifiedNotification - Buat notifikasi untuk mahasiswa saat prestasi diverifikasi
func (s *NotificationService) CreateAchievementVerifiedNotification(studentUserID uuid.UUID, achievementTitle string, referenceID uuid.UUID) error {
	notification := &model.Notification{
		UserID:    studentUserID,
		Type:      "achievement_verified",
		Title:     "Prestasi Diverifikasi",
		Message:   fmt.Sprintf("Prestasi Anda '%s' telah diverifikasi dan disetujui.", achievementTitle),
		RelatedID: &referenceID,
	}

	return s.Repo.CreateNotification(notification)
}

// CreateAchievementRejectedNotification - Buat notifikasi untuk mahasiswa saat prestasi ditolak
func (s *NotificationService) CreateAchievementRejectedNotification(studentUserID uuid.UUID, achievementTitle string, rejectionNote string, referenceID uuid.UUID) error {
	notification := &model.Notification{
		UserID:    studentUserID,
		Type:      "achievement_rejected",
		Title:     "Prestasi Ditolak",
		Message:   fmt.Sprintf("Prestasi Anda '%s' ditolak. Alasan: %s", achievementTitle, rejectionNote),
		RelatedID: &referenceID,
	}

	return s.Repo.CreateNotification(notification)
}

// GetUserNotifications - Ambil notifikasi user
func (s *NotificationService) GetUserNotifications(userID uuid.UUID, limit int) ([]model.Notification, error) {
	if limit <= 0 {
		limit = 50 // Default limit
	}
	return s.Repo.GetUserNotifications(userID, limit)
}

// GetUnreadCount - Hitung notifikasi yang belum dibaca
func (s *NotificationService) GetUnreadCount(userID uuid.UUID) (int, error) {
	return s.Repo.GetUnreadCount(userID)
}

// MarkAsRead - Tandai notifikasi sebagai sudah dibaca
func (s *NotificationService) MarkAsRead(notificationID uuid.UUID) error {
	return s.Repo.MarkAsRead(notificationID)
}

// MarkAllAsRead - Tandai semua notifikasi sebagai sudah dibaca
func (s *NotificationService) MarkAllAsRead(userID uuid.UUID) error {
	return s.Repo.MarkAllAsRead(userID)
}
