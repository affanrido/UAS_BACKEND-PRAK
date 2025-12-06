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
