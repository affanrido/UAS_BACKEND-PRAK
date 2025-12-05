package model

import (
	"time"

	"github.com/google/uuid"
)

// Lecturer - Tabel lecturers (PostgreSQL)
type Lecturer struct {
	ID          uuid.UUID `json:"id" db:"id"`
	UserID      uuid.UUID `json:"user_id" db:"user_id"`
	LecturerID  string    `json:"lecturer_id" db:"lecturer_id"`
	Department  string    `json:"department" db:"department"`
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
}
