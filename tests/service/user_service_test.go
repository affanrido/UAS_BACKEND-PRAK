package service_test

import (
	model "UAS_BACKEND/domain/model"
	"UAS_BACKEND/domain/service"
	"UAS_BACKEND/tests/mocks"
	"errors"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"golang.org/x/crypto/bcrypt"
)

func TestUserService_CreateUser_Success(t *testing.T) {
	// Arrange
	mockRepo := new(mocks.MockUserRepository)
	userService := service.NewUserService(mockRepo)

	roleID := uuid.New()
	advisorID := uuid.New()

	req := &service.CreateUserRequest{
		Username: "testuser",
		Email:    "test@example.com",
		Password: "password123",
		FullName: "Test User",
		RoleID:   roleID,
		IsActive: true,
		StudentData: &service.StudentProfileRequest{
			StudentID:    "STD001",
			ProgramStudy: "Computer Science",
			AcademicYear: "2021",
			AdvisorID:    advisorID,
		},
	}

	role := &model.Roles{
		ID:   roleID,
		Name: "student",
	}

	mockRepo.On("GetUserByUsername", "testuser").Return(nil, errors.New("not found"))
	mockRepo.On("GetUserByEmail", "test@example.com").Return(nil, errors.New("not found"))
	mockRepo.On("CreateUser", mock.AnythingOfType("*model.Users")).Return(nil)
	mockRepo.On("GetRoleByID", roleID).Return(role, nil)
	mockRepo.On("CreateStudent", mock.AnythingOfType("*model.Student")).Return(nil)

	// Act
	result, err := userService.CreateUser(req)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, "testuser", result.User.Username)
	assert.Equal(t, "test@example.com", result.User.Email)
	assert.Equal(t, "Test User", result.User.FullName)
	assert.Equal(t, roleID, result.User.RoleID)
	assert.True(t, result.User.IsActive)
	assert.NotNil(t, result.Student)
	assert.Equal(t, "STD001", result.Student.StudentID)
	assert.Equal(t, "Computer Science", result.Student.ProgramStudy)
	assert.Equal(t, role, result.Role)

	mockRepo.AssertExpectations(t)
}

func TestUserService_CreateUser_ValidationErrors(t *testing.T) {
	// Arrange
	mockRepo := new(mocks.MockUserRepository)
	userService := service.NewUserService(mockRepo)

	testCases := []struct {
		name        string
		req         *service.CreateUserRequest
		expectedErr string
	}{
		{
			name: "Empty username",
			req: &service.CreateUserRequest{
				Username: "",
				Email:    "test@example.com",
				Password: "password123",
				FullName: "Test User",
				RoleID:   uuid.New(),
				IsActive: true,
			},
			expectedErr: "username is required",
		},
		{
			name: "Empty email",
			req: &service.CreateUserRequest{
				Username: "testuser",
				Email:    "",
				Password: "password123",
				FullName: "Test User",
				RoleID:   uuid.New(),
				IsActive: true,
			},
			expectedErr: "email is required",
		},
		{
			name: "Empty password",
			req: &service.CreateUserRequest{
				Username: "testuser",
				Email:    "test@example.com",
				Password: "",
				FullName: "Test User",
				RoleID:   uuid.New(),
				IsActive: true,
			},
			expectedErr: "password is required",
		},
		{
			name: "Short password",
			req: &service.CreateUserRequest{
				Username: "testuser",
				Email:    "test@example.com",
				Password: "123",
				FullName: "Test User",
				RoleID:   uuid.New(),
				IsActive: true,
			},
			expectedErr: "password must be at least 6 characters",
		},
		{
			name: "Empty full name",
			req: &service.CreateUserRequest{
				Username: "testuser",
				Email:    "test@example.com",
				Password: "password123",
				FullName: "",
				RoleID:   uuid.New(),
				IsActive: true,
			},
			expectedErr: "full name is required",
		},
		{
			name: "Empty role ID",
			req: &service.CreateUserRequest{
				Username: "testuser",
				Email:    "test@example.com",
				Password: "password123",
				FullName: "Test User",
				RoleID:   uuid.Nil,
				IsActive: true,
			},
			expectedErr: "role ID is required",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Act
			result, err := userService.CreateUser(tc.req)

			// Assert
			assert.Error(t, err)
			assert.Nil(t, result)
			assert.Contains(t, err.Error(), tc.expectedErr)
		})
	}
}

func TestUserService_CreateUser_UsernameExists(t *testing.T) {
	// Arrange
	mockRepo := new(mocks.MockUserRepository)
	userService := service.NewUserService(mockRepo)

	existingUser := &model.Users{
		ID:       uuid.New(),
		Username: "testuser",
		Email:    "existing@example.com",
	}

	req := &service.CreateUserRequest{
		Username: "testuser",
		Email:    "test@example.com",
		Password: "password123",
		FullName: "Test User",
		RoleID:   uuid.New(),
		IsActive: true,
	}

	mockRepo.On("GetUserByUsername", "testuser").Return(existingUser, nil)

	// Act
	result, err := userService.CreateUser(req)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "username already exists")

	mockRepo.AssertExpectations(t)
}

func TestUserService_CreateUser_EmailExists(t *testing.T) {
	// Arrange
	mockRepo := new(mocks.MockUserRepository)
	userService := service.NewUserService(mockRepo)

	existingUser := &model.Users{
		ID:       uuid.New(),
		Username: "existinguser",
		Email:    "test@example.com",
	}

	req := &service.CreateUserRequest{
		Username: "testuser",
		Email:    "test@example.com",
		Password: "password123",
		FullName: "Test User",
		RoleID:   uuid.New(),
		IsActive: true,
	}

	mockRepo.On("GetUserByUsername", "testuser").Return(nil, errors.New("not found"))
	mockRepo.On("GetUserByEmail", "test@example.com").Return(existingUser, nil)

	// Act
	result, err := userService.CreateUser(req)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "email already exists")

	mockRepo.AssertExpectations(t)
}

func TestUserService_UpdateUser_Success(t *testing.T) {
	// Arrange
	mockRepo := new(mocks.MockUserRepository)
	userService := service.NewUserService(mockRepo)

	userID := uuid.New()
	roleID := uuid.New()
	newRoleID := uuid.New()

	existingUser := &model.Users{
		ID:           userID,
		Username:     "olduser",
		Email:        "old@example.com",
		PasswordHash: "oldhash",
		FullName:     "Old Name",
		RoleID:       roleID,
		IsActive:     true,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}

	req := &service.UpdateUserRequest{
		Username: "newuser",
		Email:    "new@example.com",
		Password: "newpassword123",
		FullName: "New Name",
		RoleID:   newRoleID,
	}

	role := &model.Roles{
		ID:   newRoleID,
		Name: "admin",
	}

	mockRepo.On("GetUserByID", userID).Return(existingUser, nil)
	mockRepo.On("GetUserByUsername", "newuser").Return(nil, errors.New("not found"))
	mockRepo.On("GetUserByEmail", "new@example.com").Return(nil, errors.New("not found"))
	mockRepo.On("UpdateUser", mock.AnythingOfType("*model.Users")).Return(nil)
	mockRepo.On("GetRoleByID", newRoleID).Return(role, nil)
	mockRepo.On("GetStudentByUserID", userID).Return(nil, errors.New("not found"))
	mockRepo.On("GetLecturerByUserID", userID).Return(nil, errors.New("not found"))

	// Act
	result, err := userService.UpdateUser(userID, req)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, "newuser", result.User.Username)
	assert.Equal(t, "new@example.com", result.User.Email)
	assert.Equal(t, "New Name", result.User.FullName)
	assert.Equal(t, newRoleID, result.User.RoleID)

	mockRepo.AssertExpectations(t)
}

func TestUserService_UpdateUser_UserNotFound(t *testing.T) {
	// Arrange
	mockRepo := new(mocks.MockUserRepository)
	userService := service.NewUserService(mockRepo)

	userID := uuid.New()
	req := &service.UpdateUserRequest{
		Username: "newuser",
	}

	mockRepo.On("GetUserByID", userID).Return(nil, errors.New("user not found"))

	// Act
	result, err := userService.UpdateUser(userID, req)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "user not found")

	mockRepo.AssertExpectations(t)
}

func TestUserService_DeleteUser_Success(t *testing.T) {
	// Arrange
	mockRepo := new(mocks.MockUserRepository)
	userService := service.NewUserService(mockRepo)

	userID := uuid.New()
	user := &model.Users{
		ID:       userID,
		Username: "testuser",
		Email:    "test@example.com",
	}

	mockRepo.On("GetUserByID", userID).Return(user, nil)
	mockRepo.On("DeleteUser", userID).Return(nil)

	// Act
	err := userService.DeleteUser(userID)

	// Assert
	assert.NoError(t, err)

	mockRepo.AssertExpectations(t)
}

func TestUserService_DeleteUser_UserNotFound(t *testing.T) {
	// Arrange
	mockRepo := new(mocks.MockUserRepository)
	userService := service.NewUserService(mockRepo)

	userID := uuid.New()

	mockRepo.On("GetUserByID", userID).Return(nil, errors.New("user not found"))

	// Act
	err := userService.DeleteUser(userID)

	// Assert
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "user not found")

	mockRepo.AssertExpectations(t)
}

func TestUserService_AssignRole_Success(t *testing.T) {
	// Arrange
	mockRepo := new(mocks.MockUserRepository)
	userService := service.NewUserService(mockRepo)

	userID := uuid.New()
	oldRoleID := uuid.New()
	newRoleID := uuid.New()

	user := &model.Users{
		ID:       userID,
		Username: "testuser",
		Email:    "test@example.com",
		RoleID:   oldRoleID,
	}

	role := &model.Roles{
		ID:   newRoleID,
		Name: "admin",
	}

	mockRepo.On("GetUserByID", userID).Return(user, nil)
	mockRepo.On("GetRoleByID", newRoleID).Return(role, nil)
	mockRepo.On("UpdateUser", mock.AnythingOfType("*model.Users")).Return(nil)
	mockRepo.On("GetStudentByUserID", userID).Return(nil, errors.New("not found"))
	mockRepo.On("GetLecturerByUserID", userID).Return(nil, errors.New("not found"))

	// Act
	result, err := userService.AssignRole(userID, newRoleID)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, newRoleID, result.User.RoleID)
	assert.Equal(t, role, result.Role)

	mockRepo.AssertExpectations(t)
}

func TestUserService_SetStudentProfile_Success(t *testing.T) {
	// Arrange
	mockRepo := new(mocks.MockUserRepository)
	userService := service.NewUserService(mockRepo)

	userID := uuid.New()
	advisorID := uuid.New()

	user := &model.Users{
		ID:       userID,
		Username: "testuser",
	}

	req := &service.StudentProfileRequest{
		StudentID:    "STD001",
		ProgramStudy: "Computer Science",
		AcademicYear: "2021",
		AdvisorID:    advisorID,
	}

	mockRepo.On("GetUserByID", userID).Return(user, nil)
	mockRepo.On("GetStudentByUserID", userID).Return(nil, errors.New("not found"))
	mockRepo.On("CreateStudent", mock.AnythingOfType("*model.Student")).Return(nil)

	// Act
	result, err := userService.SetStudentProfile(userID, req)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, userID, result.UserID)
	assert.Equal(t, "STD001", result.StudentID)
	assert.Equal(t, "Computer Science", result.ProgramStudy)
	assert.Equal(t, "2021", result.AcademicYear)
	assert.Equal(t, advisorID, result.AdvisorID)

	mockRepo.AssertExpectations(t)
}

func TestUserService_GetAllUsers_Success(t *testing.T) {
	// Arrange
	mockRepo := new(mocks.MockUserRepository)
	userService := service.NewUserService(mockRepo)

	roleID := uuid.New()
	users := []model.Users{
		{
			ID:       uuid.New(),
			Username: "user1",
			Email:    "user1@example.com",
			RoleID:   roleID,
		},
		{
			ID:       uuid.New(),
			Username: "user2",
			Email:    "user2@example.com",
			RoleID:   roleID,
		},
	}

	role := &model.Roles{
		ID:   roleID,
		Name: "student",
	}

	mockRepo.On("GetAllUsers", 10, 0).Return(users, 2, nil)
	mockRepo.On("GetRoleByID", roleID).Return(role, nil).Times(2)
	mockRepo.On("GetStudentByUserID", users[0].ID).Return(nil, errors.New("not found"))
	mockRepo.On("GetLecturerByUserID", users[0].ID).Return(nil, errors.New("not found"))
	mockRepo.On("GetStudentByUserID", users[1].ID).Return(nil, errors.New("not found"))
	mockRepo.On("GetLecturerByUserID", users[1].ID).Return(nil, errors.New("not found"))

	// Act
	result, total, err := userService.GetAllUsers(10, 0)

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, 2, len(result))
	assert.Equal(t, 2, total)
	assert.Equal(t, "user1", result[0].User.Username)
	assert.Equal(t, "user2", result[1].User.Username)

	mockRepo.AssertExpectations(t)
}