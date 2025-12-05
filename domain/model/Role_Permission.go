package model

import (
	"github.com/google/uuid"
)

// RolePermission - Tabel role_permissions (PostgreSQL)
type RolePermission struct {
	RoleID       uuid.UUID `json:"role_id" db:"role_id"`
	PermissionID uuid.UUID `json:"permission_id" db:"permission_id"`
}
