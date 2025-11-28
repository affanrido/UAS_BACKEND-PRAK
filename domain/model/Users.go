package model 

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

// Users struct (Sesuai Prompt)
type Users struct {
	ID           uuid.UUID `json:"id"`
	Username     string    `json:"username"`
	Email        string    `json:"email"`
	PasswordHash string    `json:"-"` // Tidak di-return di JSON
	FullName     string    `json:"full_name"`
	RoleID       uuid.UUID `json:"role_id"`
	ISActive     bool      `json:"is_active"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

// DTO untuk Input Login
type LoginRequest struct {
	Identifier string // Bisa Username atau Email
	Password   string
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