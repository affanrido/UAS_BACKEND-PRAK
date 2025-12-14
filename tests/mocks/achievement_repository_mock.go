package mocks

import (
	model "UAS_BACKEND/domain/model"
	"context"
	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// MockAchievementRepository is a mock implementation of AchievementRepository
type MockAchievementRepository struct {
	mock.Mock
}

func (m *MockAchievementRepository) CreateAchievement(ctx context.Context, achievement *model.Achievement) (primitive.ObjectID, error) {
	args := m.Called(ctx, achievement)
	return args.Get(0).(primitive.ObjectID), args.Error(1)
}

func (m *MockAchievementRepository) CreateAchievementReference(ref *model.AchievementReference) error {
	args := m.Called(ref)
	return args.Error(0)
}

func (m *MockAchievementRepository) GetAchievementByID(ctx context.Context, id primitive.ObjectID) (*model.Achievement, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.Achievement), args.Error(1)
}

func (m *MockAchievementRepository) GetAchievementReferenceByID(id uuid.UUID) (*model.AchievementReference, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.AchievementReference), args.Error(1)
}

func (m *MockAchievementRepository) GetStudentByUserID(userID uuid.UUID) (*model.Student, error) {
	args := m.Called(userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.Student), args.Error(1)
}

func (m *MockAchievementRepository) GetStudentByID(studentID uuid.UUID) (*model.Student, error) {
	args := m.Called(studentID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.Student), args.Error(1)
}

func (m *MockAchievementRepository) GetLecturerByUserID(userID uuid.UUID) (*model.Lecturer, error) {
	args := m.Called(userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.Lecturer), args.Error(1)
}

func (m *MockAchievementRepository) GetLecturerByID(lecturerID uuid.UUID) (*model.Lecturer, error) {
	args := m.Called(lecturerID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.Lecturer), args.Error(1)
}

func (m *MockAchievementRepository) GetUserByID(userID uuid.UUID) (*model.Users, error) {
	args := m.Called(userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.Users), args.Error(1)
}

func (m *MockAchievementRepository) UpdateAchievementReferenceStatus(id uuid.UUID, status string, verifiedBy *uuid.UUID, rejectionNote *string) error {
	args := m.Called(id, status, verifiedBy, rejectionNote)
	return args.Error(0)
}

func (m *MockAchievementRepository) SoftDeleteAchievement(ctx context.Context, id primitive.ObjectID) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockAchievementRepository) SoftDeleteAchievementReference(id uuid.UUID) error {
	args := m.Called(id)
	return args.Error(0)
}

func (m *MockAchievementRepository) GetStudentsByAdvisorID(advisorID uuid.UUID) ([]model.Student, error) {
	args := m.Called(advisorID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]model.Student), args.Error(1)
}

func (m *MockAchievementRepository) GetAchievementReferencesByStudentIDs(studentIDs []uuid.UUID, limit, offset int) ([]model.AchievementReference, error) {
	args := m.Called(studentIDs, limit, offset)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]model.AchievementReference), args.Error(1)
}

func (m *MockAchievementRepository) CountAchievementReferencesByStudentIDs(studentIDs []uuid.UUID) (int, error) {
	args := m.Called(studentIDs)
	return args.Int(0), args.Error(1)
}

func (m *MockAchievementRepository) GetAchievementsByIDs(ctx context.Context, ids []primitive.ObjectID) ([]model.Achievement, error) {
	args := m.Called(ctx, ids)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]model.Achievement), args.Error(1)
}