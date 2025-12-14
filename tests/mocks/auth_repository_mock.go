package mocks

import (
	model "UAS_BACKEND/domain/model"
	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
)

// MockAuthRepository is a mock implementation of AuthRepository
type MockAuthRepository struct {
	mock.Mock
}

func (m *MockAuthRepository) GetUserByIdentifier(identifier string) (*model.Users, error) {
	args := m.Called(identifier)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.Users), args.Error(1)
}

func (m *MockAuthRepository) GetUserByID(userID uuid.UUID) (*model.Users, error) {
	args := m.Called(userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.Users), args.Error(1)
}

func (m *MockAuthRepository) GetUserPermissions(userID uuid.UUID) ([]string, error) {
	args := m.Called(userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]string), args.Error(1)
}

func (m *MockAuthRepository) GetRoleByID(roleID uuid.UUID) (*model.Roles, error) {
	args := m.Called(roleID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.Roles), args.Error(1)
}