package service_test

import (
	model "UAS_BACKEND/domain/model"
	"UAS_BACKEND/domain/service"
	"UAS_BACKEND/tests/mocks"
	"context"
	"errors"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func TestAchievementService_SubmitAchievement_Success(t *testing.T) {
	// Arrange
	mockRepo := new(mocks.MockAchievementRepository)
	achievementService := service.NewAchievementService(mockRepo)

	userID := uuid.New()
	studentID := uuid.New()
	mongoID := primitive.NewObjectID()

	student := &model.Student{
		ID:           studentID,
		UserID:       userID,
		StudentID:    "STD001",
		ProgramStudy: "Computer Science",
		AcademicYear: "2021",
		AdvisorID:    uuid.New(),
		CreatedAt:    time.Now(),
	}

	req := &service.SubmitAchievementRequest{
		AchievementType: "competition",
		Title:           "Test Achievement",
		Description:     "Test Description",
		Details: map[string]interface{}{
			"competitionName":  "Test Competition",
			"competitionLevel": "national",
			"rank":             1.0,
		},
		Attachments: []service.AttachmentRequest{
			{
				FileName: "certificate.pdf",
				FileURL:  "/uploads/certificate.pdf",
				FileType: "application/pdf",
			},
		},
		Tags:   []string{"competition", "programming"},
		Points: 100.0,
	}

	mockRepo.On("GetStudentByUserID", userID).Return(student, nil)
	mockRepo.On("CreateAchievement", context.Background(), mock.AnythingOfType("*model.Achievement")).Return(mongoID, nil)
	mockRepo.On("CreateAchievementReference", mock.AnythingOfType("*model.AchievementReference")).Return(nil)

	// Act
	result, err := achievementService.SubmitAchievement(context.Background(), userID, req)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, mongoID.Hex(), result.MongoAchievementID)
	assert.Equal(t, "draft", result.Status)
	assert.Equal(t, "Test Achievement", result.Achievement.Title)
	assert.Equal(t, "competition", result.Achievement.AchievementType)

	mockRepo.AssertExpectations(t)
}

func TestAchievementService_SubmitAchievement_UserNotStudent(t *testing.T) {
	// Arrange
	mockRepo := new(mocks.MockAchievementRepository)
	achievementService := service.NewAchievementService(mockRepo)

	userID := uuid.New()

	req := &service.SubmitAchievementRequest{
		AchievementType: "competition",
		Title:           "Test Achievement",
		Description:     "Test Description",
		Details:         map[string]interface{}{},
		Attachments:     []service.AttachmentRequest{},
		Tags:            []string{},
		Points:          100.0,
	}

	mockRepo.On("GetStudentByUserID", userID).Return(nil, errors.New("user is not a student"))

	// Act
	result, err := achievementService.SubmitAchievement(context.Background(), userID, req)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "user is not a student")

	mockRepo.AssertExpectations(t)
}

func TestAchievementService_SubmitAchievement_ValidationError(t *testing.T) {
	// Arrange
	mockRepo := new(mocks.MockAchievementRepository)
	achievementService := service.NewAchievementService(mockRepo)

	userID := uuid.New()
	studentID := uuid.New()

	student := &model.Student{
		ID:           studentID,
		UserID:       userID,
		StudentID:    "STD001",
		ProgramStudy: "Computer Science",
		AcademicYear: "2021",
		AdvisorID:    uuid.New(),
		CreatedAt:    time.Now(),
	}

	testCases := []struct {
		name        string
		req         *service.SubmitAchievementRequest
		expectedErr string
	}{
		{
			name: "Empty title",
			req: &service.SubmitAchievementRequest{
				AchievementType: "competition",
				Title:           "",
				Description:     "Test Description",
				Details:         map[string]interface{}{},
				Attachments:     []service.AttachmentRequest{},
				Tags:            []string{},
				Points:          100.0,
			},
			expectedErr: "title is required",
		},
		{
			name: "Empty achievement type",
			req: &service.SubmitAchievementRequest{
				AchievementType: "",
				Title:           "Test Achievement",
				Description:     "Test Description",
				Details:         map[string]interface{}{},
				Attachments:     []service.AttachmentRequest{},
				Tags:            []string{},
				Points:          100.0,
			},
			expectedErr: "achievement type is required",
		},
		{
			name: "Invalid achievement type",
			req: &service.SubmitAchievementRequest{
				AchievementType: "invalid_type",
				Title:           "Test Achievement",
				Description:     "Test Description",
				Details:         map[string]interface{}{},
				Attachments:     []service.AttachmentRequest{},
				Tags:            []string{},
				Points:          100.0,
			},
			expectedErr: "invalid achievement type",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			mockRepo.On("GetStudentByUserID", userID).Return(student, nil)

			// Act
			result, err := achievementService.SubmitAchievement(context.Background(), userID, tc.req)

			// Assert
			assert.Error(t, err)
			assert.Nil(t, result)
			assert.Contains(t, err.Error(), tc.expectedErr)
		})
	}

	mockRepo.AssertExpectations(t)
}

func TestAchievementService_SubmitForVerification_Success(t *testing.T) {
	// Arrange
	mockRepo := new(mocks.MockAchievementRepository)
	mockNotificationService := new(mocks.MockNotificationService)
	achievementService := service.NewAchievementService(mockRepo)

	userID := uuid.New()
	studentID := uuid.New()
	referenceID := uuid.New()
	advisorID := uuid.New()
	mongoID := primitive.NewObjectID()

	student := &model.Student{
		ID:           studentID,
		UserID:       userID,
		StudentID:    "STD001",
		ProgramStudy: "Computer Science",
		AcademicYear: "2021",
		AdvisorID:    advisorID,
		CreatedAt:    time.Now(),
	}

	reference := &model.AchievementReference{
		ID:                 referenceID,
		StudentID:          studentID,
		MongoAchievementID: mongoID.Hex(),
		Status:             "draft",
		CreatedAt:          time.Now(),
		UpdatedAt:          time.Now(),
	}

	achievement := &model.Achievement{
		ID:              mongoID,
		StudentID:       studentID,
		AchievementType: "competition",
		Title:           "Test Achievement",
		Description:     "Test Description",
		CreatedAt:       time.Now(),
		UpdatedAt:       time.Now(),
	}

	studentUser := &model.Users{
		ID:       userID,
		FullName: "Test Student",
		Email:    "student@example.com",
	}

	advisorInfo := &model.Lecturer{
		ID:         advisorID,
		UserID:     uuid.New(),
		LecturerID: "LEC001",
		Department: "Computer Science",
	}

	mockRepo.On("GetStudentByUserID", userID).Return(student, nil)
	mockRepo.On("GetAchievementReferenceByID", referenceID).Return(reference, nil)
	mockRepo.On("GetAchievementByID", context.Background(), mongoID).Return(achievement, nil)
	mockRepo.On("UpdateAchievementReferenceStatus", referenceID, "submitted", (*uuid.UUID)(nil), (*string)(nil)).Return(nil)
	mockRepo.On("GetUserByID", userID).Return(studentUser, nil)
	mockRepo.On("GetLecturerByID", advisorID).Return(advisorInfo, nil)

	mockNotificationService.On("CreateAchievementSubmittedNotification", 
		advisorInfo.UserID, studentUser.FullName, achievement.Title, referenceID).Return(nil)

	// Act
	result, err := achievementService.SubmitForVerification(context.Background(), userID, referenceID, mockNotificationService)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, referenceID, result.ReferenceID)
	assert.Equal(t, "submitted", result.Status)
	assert.Contains(t, result.Message, "submitted for verification successfully")

	mockRepo.AssertExpectations(t)
	mockNotificationService.AssertExpectations(t)
}

func TestAchievementService_SubmitForVerification_WrongStatus(t *testing.T) {
	// Arrange
	mockRepo := new(mocks.MockAchievementRepository)
	achievementService := service.NewAchievementService(mockRepo)

	userID := uuid.New()
	studentID := uuid.New()
	referenceID := uuid.New()

	student := &model.Student{
		ID:           studentID,
		UserID:       userID,
		StudentID:    "STD001",
		ProgramStudy: "Computer Science",
		AcademicYear: "2021",
		AdvisorID:    uuid.New(),
		CreatedAt:    time.Now(),
	}

	reference := &model.AchievementReference{
		ID:                 referenceID,
		StudentID:          studentID,
		MongoAchievementID: primitive.NewObjectID().Hex(),
		Status:             "submitted", // Already submitted
		CreatedAt:          time.Now(),
		UpdatedAt:          time.Now(),
	}

	mockRepo.On("GetStudentByUserID", userID).Return(student, nil)
	mockRepo.On("GetAchievementReferenceByID", referenceID).Return(reference, nil)

	// Act
	result, err := achievementService.SubmitForVerification(context.Background(), userID, referenceID, nil)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "achievement must be in 'draft' status")

	mockRepo.AssertExpectations(t)
}

func TestAchievementService_DeleteAchievement_Success(t *testing.T) {
	// Arrange
	mockRepo := new(mocks.MockAchievementRepository)
	achievementService := service.NewAchievementService(mockRepo)

	userID := uuid.New()
	studentID := uuid.New()
	referenceID := uuid.New()
	mongoID := primitive.NewObjectID()

	student := &model.Student{
		ID:           studentID,
		UserID:       userID,
		StudentID:    "STD001",
		ProgramStudy: "Computer Science",
		AcademicYear: "2021",
		AdvisorID:    uuid.New(),
		CreatedAt:    time.Now(),
	}

	reference := &model.AchievementReference{
		ID:                 referenceID,
		StudentID:          studentID,
		MongoAchievementID: mongoID.Hex(),
		Status:             "draft",
		IsDeleted:          false,
		CreatedAt:          time.Now(),
		UpdatedAt:          time.Now(),
	}

	mockRepo.On("GetStudentByUserID", userID).Return(student, nil)
	mockRepo.On("GetAchievementReferenceByID", referenceID).Return(reference, nil)
	mockRepo.On("SoftDeleteAchievement", context.Background(), mongoID).Return(nil)
	mockRepo.On("SoftDeleteAchievementReference", referenceID).Return(nil)

	// Act
	result, err := achievementService.DeleteAchievement(context.Background(), userID, referenceID)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, referenceID, result.ReferenceID)
	assert.Contains(t, result.Message, "deleted successfully")

	mockRepo.AssertExpectations(t)
}

func TestAchievementService_DeleteAchievement_WrongStatus(t *testing.T) {
	// Arrange
	mockRepo := new(mocks.MockAchievementRepository)
	achievementService := service.NewAchievementService(mockRepo)

	userID := uuid.New()
	studentID := uuid.New()
	referenceID := uuid.New()

	student := &model.Student{
		ID:           studentID,
		UserID:       userID,
		StudentID:    "STD001",
		ProgramStudy: "Computer Science",
		AcademicYear: "2021",
		AdvisorID:    uuid.New(),
		CreatedAt:    time.Now(),
	}

	reference := &model.AchievementReference{
		ID:                 referenceID,
		StudentID:          studentID,
		MongoAchievementID: primitive.NewObjectID().Hex(),
		Status:             "submitted", // Cannot delete submitted
		IsDeleted:          false,
		CreatedAt:          time.Now(),
		UpdatedAt:          time.Now(),
	}

	mockRepo.On("GetStudentByUserID", userID).Return(student, nil)
	mockRepo.On("GetAchievementReferenceByID", referenceID).Return(reference, nil)

	// Act
	result, err := achievementService.DeleteAchievement(context.Background(), userID, referenceID)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "only draft achievements can be deleted")

	mockRepo.AssertExpectations(t)
}