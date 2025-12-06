package model

import (
	"time"

	"github.com/google/uuid"
)

// Notification - Tabel notifications (PostgreSQL)
type Notification struct {
	ID        uuid.UUID  `json:"id" db:"id"`
	UserID    uuid.UUID  `json:"user_id" db:"user_id"`           // Penerima notifikasi
	Type      string     `json:"type" db:"type"`                 // 'achievement_submitted', 'achievement_verified', 'achievement_rejected'
	Title     string     `json:"title" db:"title"`               // Judul notifikasi
	Message   string     `json:"message" db:"message"`           // Isi notifikasi
	RelatedID *uuid.UUID `json:"related_id" db:"related_id"`     // ID terkait (achievement_reference_id)
	IsRead    bool       `json:"is_read" db:"is_read"`           // Status sudah dibaca
	ReadAt    *time.Time `json:"read_at" db:"read_at"`           // Waktu dibaca
	CreatedAt time.Time  `json:"created_at" db:"created_at"`     // Waktu dibuat
}
