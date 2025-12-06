package repository

import (
	model "UAS_BACKEND/domain/Model"
	"database/sql"
	"errors"
	"time"

	"github.com/google/uuid"
)

type NotificationRepository struct {
	DB *sql.DB
}

func NewNotificationRepository(db *sql.DB) *NotificationRepository {
	return &NotificationRepository{DB: db}
}

// CreateNotification - Buat notifikasi baru
func (r *NotificationRepository) CreateNotification(notification *model.Notification) error {
	query := `
		INSERT INTO notifications 
		(id, user_id, type, title, message, related_id, is_read, created_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
	`

	notification.ID = uuid.New()
	notification.IsRead = false
	notification.CreatedAt = time.Now()

	_, err := r.DB.Exec(query,
		notification.ID,
		notification.UserID,
		notification.Type,
		notification.Title,
		notification.Message,
		notification.RelatedID,
		notification.IsRead,
		notification.CreatedAt,
	)

	return err
}

// GetUserNotifications - Ambil semua notifikasi user
func (r *NotificationRepository) GetUserNotifications(userID uuid.UUID, limit int) ([]model.Notification, error) {
	query := `
		SELECT id, user_id, type, title, message, related_id, is_read, read_at, created_at
		FROM notifications
		WHERE user_id = $1
		ORDER BY created_at DESC
		LIMIT $2
	`

	rows, err := r.DB.Query(query, userID, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var notifications []model.Notification
	for rows.Next() {
		var notif model.Notification
		err := rows.Scan(
			&notif.ID,
			&notif.UserID,
			&notif.Type,
			&notif.Title,
			&notif.Message,
			&notif.RelatedID,
			&notif.IsRead,
			&notif.ReadAt,
			&notif.CreatedAt,
		)
		if err != nil {
			return nil, err
		}
		notifications = append(notifications, notif)
	}

	return notifications, nil
}

// GetUnreadCount - Hitung notifikasi yang belum dibaca
func (r *NotificationRepository) GetUnreadCount(userID uuid.UUID) (int, error) {
	query := `
		SELECT COUNT(*) 
		FROM notifications 
		WHERE user_id = $1 AND is_read = false
	`

	var count int
	err := r.DB.QueryRow(query, userID).Scan(&count)
	return count, err
}

// MarkAsRead - Tandai notifikasi sebagai sudah dibaca
func (r *NotificationRepository) MarkAsRead(notificationID uuid.UUID) error {
	query := `
		UPDATE notifications
		SET is_read = true, read_at = $1
		WHERE id = $2
	`

	now := time.Now()
	_, err := r.DB.Exec(query, now, notificationID)
	return err
}

// MarkAllAsRead - Tandai semua notifikasi user sebagai sudah dibaca
func (r *NotificationRepository) MarkAllAsRead(userID uuid.UUID) error {
	query := `
		UPDATE notifications
		SET is_read = true, read_at = $1
		WHERE user_id = $2 AND is_read = false
	`

	now := time.Now()
	_, err := r.DB.Exec(query, now, userID)
	return err
}

// GetNotificationByID - Ambil notifikasi berdasarkan ID
func (r *NotificationRepository) GetNotificationByID(id uuid.UUID) (*model.Notification, error) {
	query := `
		SELECT id, user_id, type, title, message, related_id, is_read, read_at, created_at
		FROM notifications
		WHERE id = $1
	`

	var notif model.Notification
	err := r.DB.QueryRow(query, id).Scan(
		&notif.ID,
		&notif.UserID,
		&notif.Type,
		&notif.Title,
		&notif.Message,
		&notif.RelatedID,
		&notif.IsRead,
		&notif.ReadAt,
		&notif.CreatedAt,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.New("notification not found")
		}
		return nil, err
	}

	return &notif, nil
}

// DeleteNotification - Hapus notifikasi
func (r *NotificationRepository) DeleteNotification(id uuid.UUID) error {
	query := `DELETE FROM notifications WHERE id = $1`
	_, err := r.DB.Exec(query, id)
	return err
}
