package model

import (
	"time"

	"github.com/google/uuid"
)

// Roles - Tabel roles (PostgreSQL)
type Roles struct {
	ID          uuid.UUID `json:"id" db:"id"`
	Name        string    `json:"name" db:"name"`
	Description string    `json:"description" db:"description"`
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
}
