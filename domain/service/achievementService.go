package service

import (
	model "UAS_BACKEND/domain/Model"
	"UAS_BACKEND/domain/repository"
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type AchievementService struct {
	Repo *repository.AchievementRepository
}

func NewAchievementService(repo *repository.AchievementRepository) *AchievementService {
	return &AchievementService{Repo: repo}
}

// SubmitAchievementRequest - DTO untuk submit prestasi
type SubmitAchievementRequest struct {
	AchievementType string                 `json:"achievementType"` // 'academic', 'competition', 'organization', 'publication', 'certification', 'other'
	Title           string                 `json:"title"`
	Description     string                 `json:"description"`
	Details         map[string]interface{} `json:"details"`
	CustomFields    map[string]interface{} `json:"customFields,omitempty"`
	Attachments     []AttachmentRequest    `json:"attachments"`
	Tags            []string               `json:"tags"`
	Points          float64                `json:"points"`
}

type AttachmentRequest struct {
	FileName string `json:"fileName"`
	FileURL  string `json:"fileUrl"`
	FileType string `json:"fileType"`
}

// SubmitAchievementResponse - DTO untuk response
type SubmitAchievementResponse struct {
	ReferenceID        uuid.UUID              `json:"reference_id"`
	MongoAchievementID string                 `json:"mongo_achievement_id"`
	Status             string                 `json:"status"`
	Achievement        *model.Achievement     `json:"achievement"`
	CreatedAt          time.Time              `json:"created_at"`
}

// SubmitAchievement - Flow FR-003: Submit Prestasi
func (s *AchievementService) SubmitAchievement(ctx context.Context, userID uuid.UUID, req *SubmitAchievementRequest) (*SubmitAchievementResponse, error) {
	// Validasi: User harus mahasiswa
	student, err := s.Repo.GetStudentByUserID(userID)
	if err != nil {
		return nil, errors.New("user is not a student")
	}

	// 1. Mahasiswa mengisi data prestasi (dari request)
	// Validasi input
	if err := s.validateAchievementRequest(req); err != nil {
		return nil, err
	}

	// 2. Mahasiswa upload dokumen pendukung (attachments dari request)
	attachments := make([]model.Attachment, len(req.Attachments))
	for i, att := range req.Attachments {
		attachments[i] = model.Attachment{
			FileName:   att.FileName,
			FileURL:    att.FileURL,
			FileType:   att.FileType,
			UploadedAt: time.Now(),
		}
	}

	// Parse details ke AchievementDetails
	details := s.parseDetails(req.AchievementType, req.Details)

	// 3. Sistem simpan ke MongoDB (achievement)
	achievement := &model.Achievement{
		StudentID:       student.ID,
		AchievementType: req.AchievementType,
		Title:           req.Title,
		Description:     req.Description,
		Details:         details,
		CustomFields:    req.CustomFields,
		Attachments:     attachments,
		Tags:            req.Tags,
		Points:          req.Points,
	}

	mongoID, err := s.Repo.CreateAchievement(ctx, achievement)
	if err != nil {
		return nil, errors.New("failed to save achievement to MongoDB: " + err.Error())
	}

	// 3. Sistem simpan ke PostgreSQL (reference)
	referenceID := uuid.New()
	reference := &model.AchievementReference{
		ID:                 referenceID,
		StudentID:          student.ID,
		MongoAchievementID: mongoID.Hex(),
		Status:             "draft", // 4. Status awal: 'draft'
	}

	err = s.Repo.CreateAchievementReference(reference)
	if err != nil {
		return nil, errors.New("failed to save achievement reference to PostgreSQL: " + err.Error())
	}

	// Set ID yang sudah di-generate
	achievement.ID = mongoID

	// 5. Return achievement data
	return &SubmitAchievementResponse{
		ReferenceID:        referenceID,
		MongoAchievementID: mongoID.Hex(),
		Status:             "draft",
		Achievement:        achievement,
		CreatedAt:          achievement.CreatedAt,
	}, nil
}

// validateAchievementRequest - Validasi input
func (s *AchievementService) validateAchievementRequest(req *SubmitAchievementRequest) error {
	if req.Title == "" {
		return errors.New("title is required")
	}

	if req.AchievementType == "" {
		return errors.New("achievement type is required")
	}

	validTypes := map[string]bool{
		"academic":      true,
		"competition":   true,
		"organization":  true,
		"publication":   true,
		"certification": true,
		"other":         true,
	}

	if !validTypes[req.AchievementType] {
		return errors.New("invalid achievement type")
	}

	return nil
}

// parseDetails - Parse details dari map ke AchievementDetails struct
func (s *AchievementService) parseDetails(achievementType string, detailsMap map[string]interface{}) model.AchievementDetails {
	details := model.AchievementDetails{}

	// Helper function untuk get string pointer
	getStringPtr := func(key string) *string {
		if val, ok := detailsMap[key].(string); ok {
			return &val
		}
		return nil
	}

	// Helper function untuk get float64 pointer
	getFloat64Ptr := func(key string) *float64 {
		if val, ok := detailsMap[key].(float64); ok {
			return &val
		}
		return nil
	}

	// Parse berdasarkan tipe
	switch achievementType {
	case "competition":
		details.CompetitionName = getStringPtr("competitionName")
		details.CompetitionLevel = getStringPtr("competitionLevel")
		details.Rank = getFloat64Ptr("rank")
		details.MedalType = getStringPtr("medalType")

	case "publication":
		details.PublicationType = getStringPtr("publicationType")
		details.PublicationTitle = getStringPtr("publicationTitle")
		details.Publisher = getStringPtr("publisher")
		details.ISSN = getStringPtr("issn")
		
		// Parse authors array
		if authors, ok := detailsMap["authors"].([]interface{}); ok {
			authorStrs := make([]string, 0, len(authors))
			for _, author := range authors {
				if authorStr, ok := author.(string); ok {
					authorStrs = append(authorStrs, authorStr)
				}
			}
			details.Authors = authorStrs
		}

	case "organization":
		details.OrganizationName = getStringPtr("organizationName")
		details.Position = getStringPtr("position")
		
		// Parse period
		if periodMap, ok := detailsMap["period"].(map[string]interface{}); ok {
			period := &model.Period{}
			if start, ok := periodMap["start"].(string); ok {
				if t, err := time.Parse(time.RFC3339, start); err == nil {
					period.Start = t
				}
			}
			if end, ok := periodMap["end"].(string); ok {
				if t, err := time.Parse(time.RFC3339, end); err == nil {
					period.End = t
				}
			}
			details.Period = period
		}

	case "certification":
		details.CertificationName = getStringPtr("certificationName")
		details.IssuedBy = getStringPtr("issuedBy")
		details.CertificationNumber = getStringPtr("certificationNumber")
		
		// Parse validUntil
		if validUntil, ok := detailsMap["validUntil"].(string); ok {
			if t, err := time.Parse(time.RFC3339, validUntil); err == nil {
				details.ValidUntil = &t
			}
		}
	}

	// Parse field umum
	details.EventDate = nil
	if eventDate, ok := detailsMap["eventDate"].(string); ok {
		if t, err := time.Parse(time.RFC3339, eventDate); err == nil {
			details.EventDate = &t
		}
	}

	details.Location = getStringPtr("location")
	details.Organizer = getStringPtr("organizer")
	details.Score = getFloat64Ptr("score")

	return details
}

// GetStudentAchievements - Ambil semua prestasi mahasiswa
func (s *AchievementService) GetStudentAchievements(ctx context.Context, userID uuid.UUID) ([]model.Achievement, error) {
	student, err := s.Repo.GetStudentByUserID(userID)
	if err != nil {
		return nil, errors.New("user is not a student")
	}

	return s.Repo.GetStudentAchievements(ctx, student.ID)
}

// GetAchievementByID - Ambil detail prestasi
func (s *AchievementService) GetAchievementByID(ctx context.Context, mongoID string) (*model.Achievement, error) {
	objectID, err := primitive.ObjectIDFromHex(mongoID)
	if err != nil {
		return nil, errors.New("invalid achievement ID")
	}

	return s.Repo.GetAchievementByID(ctx, objectID)
}

// SubmitForVerificationRequest - DTO untuk submit verifikasi
type SubmitForVerificationRequest struct {
	ReferenceID uuid.UUID `json:"reference_id"`
}

// SubmitForVerificationResponse - DTO untuk response
type SubmitForVerificationResponse struct {
	ReferenceID uuid.UUID  `json:"reference_id"`
	Status      string     `json:"status"`
	SubmittedAt time.Time  `json:"submitted_at"`
	Message     string     `json:"message"`
}

// SubmitForVerification - Flow FR-004: Submit untuk Verifikasi
func (s *AchievementService) SubmitForVerification(ctx context.Context, userID uuid.UUID, referenceID uuid.UUID, notificationService *NotificationService) (*SubmitForVerificationResponse, error) {
	// Validasi: User harus mahasiswa
	student, err := s.Repo.GetStudentByUserID(userID)
	if err != nil {
		return nil, errors.New("user is not a student")
	}

	// 1. Get achievement reference
	reference, err := s.Repo.GetAchievementReferenceByID(referenceID)
	if err != nil {
		return nil, errors.New("achievement reference not found")
	}

	// Validasi: Reference harus milik mahasiswa ini
	if reference.StudentID != student.ID {
		return nil, errors.New("unauthorized: achievement does not belong to you")
	}

	// Precondition: Prestasi berstatus 'draft'
	if reference.Status != "draft" {
		return nil, errors.New("achievement must be in 'draft' status to submit for verification")
	}

	// Get achievement detail dari MongoDB
	achievement, err := s.GetAchievementByID(ctx, reference.MongoAchievementID)
	if err != nil {
		return nil, errors.New("achievement not found in MongoDB")
	}

	// 2. Update status menjadi 'submitted'
	now := time.Now()
	err = s.Repo.UpdateAchievementReferenceStatus(referenceID, "submitted", nil, nil)
	if err != nil {
		return nil, errors.New("failed to update status: " + err.Error())
	}

	// Update submitted_at
	reference.Status = "submitted"
	reference.SubmittedAt = &now

	// 3. Create notification untuk dosen wali
	if notificationService != nil {
		// Get student user info untuk nama
		studentUser, err := s.Repo.GetUserByID(student.UserID)
		if err == nil && studentUser != nil {
			// Get advisor user_id
			advisorInfo, err := s.Repo.GetLecturerByID(student.AdvisorID)
			if err == nil && advisorInfo != nil {
				// Create notification
				_ = notificationService.CreateAchievementSubmittedNotification(
					advisorInfo.UserID,
					studentUser.FullName,
					achievement.Title,
					referenceID,
				)
			}
		}
	}

	// 4. Return updated status
	return &SubmitForVerificationResponse{
		ReferenceID: referenceID,
		Status:      "submitted",
		SubmittedAt: now,
		Message:     "Achievement submitted for verification successfully",
	}, nil
}

// DeleteAchievementRequest - DTO untuk delete achievement
type DeleteAchievementRequest struct {
	ReferenceID uuid.UUID `json:"reference_id"`
}

// DeleteAchievementResponse - DTO untuk response
type DeleteAchievementResponse struct {
	ReferenceID uuid.UUID `json:"reference_id"`
	Message     string    `json:"message"`
	DeletedAt   time.Time `json:"deleted_at"`
}

// DeleteAchievement - Flow FR-005: Hapus Prestasi
func (s *AchievementService) DeleteAchievement(ctx context.Context, userID uuid.UUID, referenceID uuid.UUID) (*DeleteAchievementResponse, error) {
	// Validasi: User harus mahasiswa
	student, err := s.Repo.GetStudentByUserID(userID)
	if err != nil {
		return nil, errors.New("user is not a student")
	}

	// Get achievement reference
	reference, err := s.Repo.GetAchievementReferenceByID(referenceID)
	if err != nil {
		return nil, errors.New("achievement reference not found")
	}

	// Validasi: Reference harus milik mahasiswa ini
	if reference.StudentID != student.ID {
		return nil, errors.New("unauthorized: achievement does not belong to you")
	}

	// Precondition: Prestasi berstatus 'draft'
	if reference.Status != "draft" {
		return nil, errors.New("only draft achievements can be deleted")
	}

	// Validasi: Belum dihapus sebelumnya
	if reference.IsDeleted {
		return nil, errors.New("achievement already deleted")
	}

	// Parse MongoDB ObjectID
	mongoID, err := primitive.ObjectIDFromHex(reference.MongoAchievementID)
	if err != nil {
		return nil, errors.New("invalid mongo achievement ID")
	}

	// 1. Soft delete data di MongoDB
	err = s.Repo.SoftDeleteAchievement(ctx, mongoID)
	if err != nil {
		return nil, errors.New("failed to delete achievement in MongoDB: " + err.Error())
	}

	// 2. Update reference di PostgreSQL (soft delete)
	err = s.Repo.SoftDeleteAchievementReference(referenceID)
	if err != nil {
		return nil, errors.New("failed to delete achievement reference in PostgreSQL: " + err.Error())
	}

	// 3. Return success message
	now := time.Now()
	return &DeleteAchievementResponse{
		ReferenceID: referenceID,
		Message:     "Achievement deleted successfully",
		DeletedAt:   now,
	}, nil
}

// AdvisedStudentAchievement - DTO untuk prestasi mahasiswa bimbingan
type AdvisedStudentAchievement struct {
	Reference   model.AchievementReference `json:"reference"`
	Achievement *model.Achievement         `json:"achievement"`
	Student     *StudentInfo               `json:"student"`
}

// StudentInfo - DTO untuk info mahasiswa
type StudentInfo struct {
	ID           uuid.UUID `json:"id"`
	StudentID    string    `json:"student_id"`
	FullName     string    `json:"full_name"`
	ProgramStudy string    `json:"program_study"`
	AcademicYear string    `json:"academic_year"`
}

// ViewAdvisedStudentsAchievementsResponse - DTO untuk response
type ViewAdvisedStudentsAchievementsResponse struct {
	Achievements []AdvisedStudentAchievement `json:"achievements"`
	Pagination   model.PaginationResponse    `json:"pagination"`
}

// VerifyAchievementRequest - DTO untuk verify achievement
type VerifyAchievementRequest struct {
	ReferenceID uuid.UUID `json:"reference_id"`
	Approved    bool      `json:"approved"` // true = verified, false = rejected
	Note        string    `json:"note,omitempty"`
}

// VerifyAchievementResponse - DTO untuk response
type VerifyAchievementResponse struct {
	ReferenceID uuid.UUID  `json:"reference_id"`
	Status      string     `json:"status"`
	VerifiedBy  uuid.UUID  `json:"verified_by"`
	VerifiedAt  time.Time  `json:"verified_at"`
	Note        *string    `json:"note,omitempty"`
	Message     string     `json:"message"`
}

// VerifyAchievement - Flow FR-007: Verify Prestasi
func (s *AchievementService) VerifyAchievement(ctx context.Context, userID uuid.UUID, referenceID uuid.UUID, approved bool, note string, notificationService *NotificationService) (*VerifyAchievementResponse, error) {
	// Validasi: User harus dosen/lecturer
	lecturer, err := s.Repo.GetLecturerByUserID(userID)
	if err != nil {
		return nil, errors.New("user is not a lecturer")
	}

	// 1. Get achievement reference
	reference, err := s.Repo.GetAchievementReferenceByID(referenceID)
	if err != nil {
		return nil, errors.New("achievement reference not found")
	}

	// Precondition: Prestasi berstatus 'submitted'
	if reference.Status != "submitted" {
		return nil, errors.New("achievement must be in 'submitted' status to verify")
	}

	// Get student info untuk validasi advisor
	student, err := s.Repo.GetStudentByID(reference.StudentID)
	if err != nil {
		return nil, errors.New("student not found")
	}

	// Validasi: Hanya dosen wali yang bisa verify
	if student.AdvisorID != lecturer.ID {
		return nil, errors.New("unauthorized: you are not the advisor of this student")
	}

	// Get achievement detail dari MongoDB
	achievement, err := s.GetAchievementByID(ctx, reference.MongoAchievementID)
	if err != nil {
		return nil, errors.New("achievement not found in MongoDB")
	}

	// 2. Dosen approve/reject prestasi
	var newStatus string
	var rejectionNote *string
	if approved {
		newStatus = "verified"
		rejectionNote = nil
	} else {
		newStatus = "rejected"
		if note != "" {
			rejectionNote = &note
		}
	}

	// 3. Update status menjadi 'verified' atau 'rejected'
	// 4. Set verified_by dan verified_at
	now := time.Now()
	err = s.Repo.UpdateAchievementReferenceStatus(referenceID, newStatus, &lecturer.ID, rejectionNote)
	if err != nil {
		return nil, errors.New("failed to update status: " + err.Error())
	}

	// Create notification untuk mahasiswa
	if notificationService != nil {
		studentUser, err := s.Repo.GetUserByID(student.UserID)
		if err == nil && studentUser != nil {
			lecturerUser, err := s.Repo.GetUserByID(lecturer.UserID)
			if err == nil && lecturerUser != nil {
				if approved {
					_ = notificationService.CreateAchievementVerifiedNotification(
						student.UserID,
						lecturerUser.FullName,
						achievement.Title,
						referenceID,
					)
				} else {
					_ = notificationService.CreateAchievementRejectedNotification(
						student.UserID,
						lecturerUser.FullName,
						achievement.Title,
						referenceID,
						note,
					)
				}
			}
		}
	}

	// 5. Return updated status
	message := "Achievement verified successfully"
	if !approved {
		message = "Achievement rejected"
	}

	return &VerifyAchievementResponse{
		ReferenceID: referenceID,
		Status:      newStatus,
		VerifiedBy:  lecturer.ID,
		VerifiedAt:  now,
		Note:        rejectionNote,
		Message:     message,
	}, nil
}

// ViewAdvisedStudentsAchievements - Flow FR-006: View Prestasi Mahasiswa Bimbingan
func (s *AchievementService) ViewAdvisedStudentsAchievements(ctx context.Context, userID uuid.UUID, pagination model.PaginationRequest) (*ViewAdvisedStudentsAchievementsResponse, error) {
	// Validasi: User harus dosen/lecturer
	lecturer, err := s.Repo.GetLecturerByUserID(userID)
	if err != nil {
		return nil, errors.New("user is not a lecturer")
	}

	// 1. Get list student IDs dari tabel students where advisor_id
	students, err := s.Repo.GetStudentsByAdvisorID(lecturer.ID)
	if err != nil {
		return nil, errors.New("failed to get advised students: " + err.Error())
	}

	// Jika tidak ada mahasiswa bimbingan
	if len(students) == 0 {
		return &ViewAdvisedStudentsAchievementsResponse{
			Achievements: []AdvisedStudentAchievement{},
			Pagination: model.PaginationResponse{
				Page:       pagination.Page,
				PageSize:   pagination.PageSize,
				TotalItems: 0,
				TotalPages: 0,
			},
		}, nil
	}

	// Extract student IDs
	studentIDs := make([]uuid.UUID, len(students))
	studentMap := make(map[uuid.UUID]model.Student)
	for i, student := range students {
		studentIDs[i] = student.ID
		studentMap[student.ID] = student
	}

	// Count total achievements
	totalItems, err := s.Repo.CountAchievementReferencesByStudentIDs(studentIDs)
	if err != nil {
		return nil, errors.New("failed to count achievements: " + err.Error())
	}

	// 2. Get achievements references dengan filter student_ids (dengan pagination)
	references, err := s.Repo.GetAchievementReferencesByStudentIDs(
		studentIDs,
		pagination.GetLimit(),
		pagination.GetOffset(),
	)
	if err != nil {
		return nil, errors.New("failed to get achievement references: " + err.Error())
	}

	// Jika tidak ada prestasi
	if len(references) == 0 {
		return &ViewAdvisedStudentsAchievementsResponse{
			Achievements: []AdvisedStudentAchievement{},
			Pagination: model.PaginationResponse{
				Page:       pagination.Page,
				PageSize:   pagination.PageSize,
				TotalItems: totalItems,
				TotalPages: model.CalculateTotalPages(totalItems, pagination.PageSize),
			},
		}, nil
	}

	// Extract MongoDB IDs
	mongoIDs := make([]primitive.ObjectID, 0, len(references))
	for _, ref := range references {
		mongoID, err := primitive.ObjectIDFromHex(ref.MongoAchievementID)
		if err == nil {
			mongoIDs = append(mongoIDs, mongoID)
		}
	}

	// 3. Fetch detail dari MongoDB
	achievements, err := s.Repo.GetAchievementsByIDs(ctx, mongoIDs)
	if err != nil {
		return nil, errors.New("failed to get achievements from MongoDB: " + err.Error())
	}

	// Create achievement map for quick lookup
	achievementMap := make(map[string]*model.Achievement)
	for i := range achievements {
		achievementMap[achievements[i].ID.Hex()] = &achievements[i]
	}

	// Get user info for students
	userMap := make(map[uuid.UUID]*model.Users)
	for _, student := range students {
		user, err := s.Repo.GetUserByID(student.UserID)
		if err == nil {
			userMap[student.ID] = user
		}
	}

	// Combine data
	result := make([]AdvisedStudentAchievement, 0, len(references))
	for _, ref := range references {
		achievement := achievementMap[ref.MongoAchievementID]
		student := studentMap[ref.StudentID]
		user := userMap[ref.StudentID]

		var studentInfo *StudentInfo
		if user != nil {
			studentInfo = &StudentInfo{
				ID:           student.ID,
				StudentID:    student.StudentID,
				FullName:     user.FullName,
				ProgramStudy: student.ProgramStudy,
				AcademicYear: student.AcademicYear,
			}
		}

		result = append(result, AdvisedStudentAchievement{
			Reference:   ref,
			Achievement: achievement,
			Student:     studentInfo,
		})
	}

	// 4. Return list dengan pagination
	return &ViewAdvisedStudentsAchievementsResponse{
		Achievements: result,
		Pagination: model.PaginationResponse{
			Page:       pagination.Page,
			PageSize:   pagination.PageSize,
			TotalItems: totalItems,
			TotalPages: model.CalculateTotalPages(totalItems, pagination.PageSize),
		},
	}, nil
}
