package model

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

// Users - Tabel users (PostgreSQL)
type Users struct {
	ID           uuid.UUID `json:"id" db:"id"`
	Username     string    `json:"username" db:"username"`
	Email        string    `json:"email" db:"email"`
	PasswordHash string    `json:"-" db:"password_hash"`
	FullName     string    `json:"full_name" db:"full_name"`
	RoleID       uuid.UUID `json:"role_id" db:"role_id"`
	IsActive     bool      `json:"is_active" db:"is_active"`
	CreatedAt    time.Time `json:"created_at" db:"created_at"`
	UpdatedAt    time.Time `json:"updated_at" db:"updated_at"`
}

// DTO untuk Input Login
type LoginRequest struct {
	Identifier string `json:"identifier"`
	Password   string `json:"password"`
}

// DTO untuk Output Login
type LoginResponse struct {
	Token string `json:"token"`
	User  Users  `json:"user"`
}

// Claims untuk JWT
type CustomClaims struct {
	UserID      uuid.UUID `json:"user_id"`
	RoleID      uuid.UUID `json:"role_id"`
	Permissions []string  `json:"permissions"`
	jwt.RegisteredClaims
}