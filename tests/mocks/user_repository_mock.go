package mocks

import (
	model "UAS_BACKEND/domain/model"
	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
)

// MockUserRepository is a mock implementation of UserRepository
type MockUserRepository struct {
	mock.Mock
}

func (m *MockUserRepository) CreateUser(user *model.Users) error {
	args := m.Called(user)
	return args.Error(0)
}

func (m *MockUserRepository) GetUserByID(userID uuid.UUID) (*model.Users, error) {
	args := m.Called(userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.Users), args.Error(1)
}

func (m *MockUserRepository) GetUserByUsername(username string) (*model.Users, error) {
	args := m.Called(username)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.Users), args.Error(1)
}

func (m *MockUserRepository) GetUserByEmail(email string) (*model.Users, error) {
	args := m.Called(email)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.Users), args.Error(1)
}

func (m *MockUserRepository) UpdateUser(user *model.Users) error {
	args := m.Called(user)
	return args.Error(0)
}

func (m *MockUserRepository) DeleteUser(userID uuid.UUID) error {
	args := m.Called(userID)
	return args.Error(0)
}

func (m *MockUserRepository) GetAllUsers(limit, offset int) ([]model.Users, int, error) {
	args := m.Called(limit, offset)
	if args.Get(0) == nil {
		return nil, args.Int(1), args.Error(2)
	}
	return args.Get(0).([]model.Users), args.Int(1), args.Error(2)
}

func (m *MockUserRepository) GetRoleByID(roleID uuid.UUID) (*model.Roles, error) {
	args := m.Called(roleID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.Roles), args.Error(1)
}

func (m *MockUserRepository) GetAllRoles() ([]model.Roles, error) {
	args := m.Called()
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]model.Roles), args.Error(1)
}

func (m *MockUserRepository) CreateStudent(student *model.Student) error {
	args := m.Called(student)
	return args.Error(0)
}

func (m *MockUserRepository) GetStudentByUserID(userID uuid.UUID) (*model.Student, error) {
	args := m.Called(userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.Student), args.Error(1)
}

func (m *MockUserRepository) UpdateStudent(student *model.Student) error {
	args := m.Called(student)
	return args.Error(0)
}

func (m *MockUserRepository) CreateLecturer(lecturer *model.Lecturer) error {
	args := m.Called(lecturer)
	return args.Error(0)
}

func (m *MockUserRepository) GetLecturerByUserID(userID uuid.UUID) (*model.Lecturer, error) {
	args := m.Called(userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.Lecturer), args.Error(1)
}

func (m *MockUserRepository) GetLecturerByID(lecturerID uuid.UUID) (*model.Lecturer, error) {
	args := m.Called(lecturerID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.Lecturer), args.Error(1)
}

func (m *MockUserRepository) UpdateLecturer(lecturer *model.Lecturer) error {
	args := m.Called(lecturer)
	return args.Error(0)
}