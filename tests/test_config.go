package tests

import (
	"time"

	"github.com/google/uuid"
)

// TestConfig contains common test configuration and utilities
type TestConfig struct {
	DefaultTimeout time.Duration
	TestUserID     uuid.UUID
	TestRoleID     uuid.UUID
}

// NewTestConfig creates a new test configuration
func NewTestConfig() *TestConfig {
	return &TestConfig{
		DefaultTimeout: 5 * time.Second,
		TestUserID:     uuid.New(),
		TestRoleID:     uuid.New(),
	}
}

// Common test data generators
func GenerateTestUserID() uuid.UUID {
	return uuid.New()
}

func GenerateTestRoleID() uuid.UUID {
	return uuid.New()
}

func GenerateTestStudentID() uuid.UUID {
	return uuid.New()
}

func GenerateTestLecturerID() uuid.UUID {
	return uuid.New()
}

// Test constants
const (
	TestUsername     = "testuser"
	TestEmail        = "test@example.com"
	TestPassword     = "password123"
	TestFullName     = "Test User"
	TestStudentID    = "STD001"
	TestLecturerID   = "LEC001"
	TestProgramStudy = "Computer Science"
	TestAcademicYear = "2021"
	TestDepartment   = "Computer Science"
)

// Test achievement data
const (
	TestAchievementTitle       = "Programming Competition Winner"
	TestAchievementDescription = "Won first place in national programming competition"
	TestAchievementType        = "competition"
	TestCompetitionLevel       = "national"
	TestAchievementPoints      = 100.0
)

// Test notification data
const (
	TestNotificationTitle   = "Test Notification"
	TestNotificationMessage = "This is a test notification"
	TestNotificationType    = "achievement_submitted"
)