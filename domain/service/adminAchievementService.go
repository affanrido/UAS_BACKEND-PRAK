package service

import (
	model "UAS_BACKEND/domain/Model"
	"UAS_BACKEND/domain/repository"
	"context"
	"errors"

	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type AdminAchievementService struct {
	Repo *repository.AchievementRepository
}

func NewAdminAchievementService(repo *repository.AchievementRepository) *AdminAchievementService {
	return &AdminAchievementService{Repo: repo}
}



// AdminAchievementResponse - DTO untuk response
type AdminAchievementResponse struct {
	Reference   model.AchievementReference `json:"reference"`
	Achievement *model.Achievement         `json:"achievement"`
	Student     *AdminStudentInfo          `json:"student"`
	Advisor     *AdminAdvisorInfo          `json:"advisor,omitempty"`
}

// AdminStudentInfo - DTO untuk info mahasiswa
type AdminStudentInfo struct {
	ID           uuid.UUID `json:"id"`
	UserID       uuid.UUID `json:"user_id"`
	StudentID    string    `json:"student_id"`
	FullName     string    `json:"full_name"`
	Email        string    `json:"email"`
	ProgramStudy string    `json:"program_study"`
	AcademicYear string    `json:"academic_year"`
}

// AdminAdvisorInfo - DTO untuk info dosen wali
type AdminAdvisorInfo struct {
	ID         uuid.UUID `json:"id"`
	UserID     uuid.UUID `json:"user_id"`
	LecturerID string    `json:"lecturer_id"`
	FullName   string    `json:"full_name"`
	Email      string    `json:"email"`
	Department string    `json:"department"`
}

// ViewAllAchievementsRequest - DTO untuk request
type ViewAllAchievementsRequest struct {
	Page   int                           `json:"page"`
	Size   int                           `json:"size"`
	Filter *model.AdminAchievementFilter `json:"filter,omitempty"`
	Sort   *model.AdminAchievementSort   `json:"sort,omitempty"`
}

// ViewAllAchievementsResponse - DTO untuk response
type ViewAllAchievementsResponse struct {
	Achievements []AdminAchievementResponse   `json:"achievements"`
	Pagination   model.PaginationResponse     `json:"pagination"`
	Summary      *model.AchievementSummary    `json:"summary"`
}



// ViewAllAchievements - Flow FR-010: View All Achievements
func (s *AdminAchievementService) ViewAllAchievements(ctx context.Context, req *ViewAllAchievementsRequest) (*ViewAllAchievementsResponse, error) {
	// Set default values
	if req.Page < 1 {
		req.Page = 1
	}
	if req.Size < 1 {
		req.Size = 10
	}
	if req.Size > 100 {
		req.Size = 100
	}

	// Set default sort
	if req.Sort == nil {
		req.Sort = &model.AdminAchievementSort{
			Field: "created_at",
			Order: "desc",
		}
	}

	offset := (req.Page - 1) * req.Size

	// 1. Get all achievement references dengan filter dan pagination
	references, totalItems, err := s.Repo.GetAllAchievementReferencesAdmin(req.Size, offset, req.Filter, req.Sort)
	if err != nil {
		return nil, errors.New("failed to get achievement references: " + err.Error())
	}

	// Jika tidak ada data
	if len(references) == 0 {
		summary, _ := s.getAchievementSummary(req.Filter)
		return &ViewAllAchievementsResponse{
			Achievements: []AdminAchievementResponse{},
			Pagination: model.PaginationResponse{
				Page:       req.Page,
				PageSize:   req.Size,
				TotalItems: totalItems,
				TotalPages: model.CalculateTotalPages(totalItems, req.Size),
			},
			Summary: summary,
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

	// 2. Fetch details dari MongoDB
	achievements, err := s.Repo.GetAchievementsByIDs(ctx, mongoIDs)
	if err != nil {
		return nil, errors.New("failed to get achievements from MongoDB: " + err.Error())
	}

	// Create achievement map for quick lookup
	achievementMap := make(map[string]*model.Achievement)
	for i := range achievements {
		achievementMap[achievements[i].ID.Hex()] = &achievements[i]
	}

	// Get student and advisor info
	studentMap := make(map[uuid.UUID]*AdminStudentInfo)
	advisorMap := make(map[uuid.UUID]*AdminAdvisorInfo)

	for _, ref := range references {
		// Get student info if not cached
		if _, exists := studentMap[ref.StudentID]; !exists {
			studentInfo, err := s.getStudentInfo(ref.StudentID)
			if err == nil {
				studentMap[ref.StudentID] = studentInfo
				
				// Get advisor info if not cached
				if studentInfo != nil && studentInfo.ID != uuid.Nil {
					student, err := s.Repo.GetStudentByID(ref.StudentID)
					if err == nil && student != nil {
						if _, advisorExists := advisorMap[student.AdvisorID]; !advisorExists {
							advisorInfo, err := s.getAdvisorInfo(student.AdvisorID)
							if err == nil {
								advisorMap[student.AdvisorID] = advisorInfo
							}
						}
					}
				}
			}
		}
	}

	// Combine data
	result := make([]AdminAchievementResponse, 0, len(references))
	for _, ref := range references {
		achievement := achievementMap[ref.MongoAchievementID]
		student := studentMap[ref.StudentID]
		
		var advisor *AdminAdvisorInfo
		if student != nil {
			studentRecord, err := s.Repo.GetStudentByID(ref.StudentID)
			if err == nil && studentRecord != nil {
				advisor = advisorMap[studentRecord.AdvisorID]
			}
		}

		result = append(result, AdminAchievementResponse{
			Reference:   ref,
			Achievement: achievement,
			Student:     student,
			Advisor:     advisor,
		})
	}

	// Get summary statistics
	summary, _ := s.getAchievementSummary(req.Filter)

	// 4. Return dengan pagination
	return &ViewAllAchievementsResponse{
		Achievements: result,
		Pagination: model.PaginationResponse{
			Page:       req.Page,
			PageSize:   req.Size,
			TotalItems: totalItems,
			TotalPages: model.CalculateTotalPages(totalItems, req.Size),
		},
		Summary: summary,
	}, nil
}

// getStudentInfo - Get student information
func (s *AdminAchievementService) getStudentInfo(studentID uuid.UUID) (*AdminStudentInfo, error) {
	student, err := s.Repo.GetStudentByID(studentID)
	if err != nil {
		return nil, err
	}

	user, err := s.Repo.GetUserByID(student.UserID)
	if err != nil {
		return nil, err
	}

	return &AdminStudentInfo{
		ID:           student.ID,
		UserID:       student.UserID,
		StudentID:    student.StudentID,
		FullName:     user.FullName,
		Email:        user.Email,
		ProgramStudy: student.ProgramStudy,
		AcademicYear: student.AcademicYear,
	}, nil
}

// getAdvisorInfo - Get advisor information
func (s *AdminAchievementService) getAdvisorInfo(advisorID uuid.UUID) (*AdminAdvisorInfo, error) {
	lecturer, err := s.Repo.GetLecturerByID(advisorID)
	if err != nil {
		return nil, err
	}

	user, err := s.Repo.GetUserByID(lecturer.UserID)
	if err != nil {
		return nil, err
	}

	return &AdminAdvisorInfo{
		ID:         lecturer.ID,
		UserID:     lecturer.UserID,
		LecturerID: lecturer.LecturerID,
		FullName:   user.FullName,
		Email:      user.Email,
		Department: lecturer.Department,
	}, nil
}

// getAchievementSummary - Get achievement statistics
func (s *AdminAchievementService) getAchievementSummary(filter *model.AdminAchievementFilter) (*model.AchievementSummary, error) {
	summary, err := s.Repo.GetAchievementSummaryAdmin(filter)
	if err != nil {
		return &model.AchievementSummary{}, nil // Return empty summary on error
	}
	return summary, nil
}

// GetAchievementByReferenceID - Get achievement detail by reference ID
func (s *AdminAchievementService) GetAchievementByReferenceID(ctx context.Context, referenceID uuid.UUID) (*AdminAchievementResponse, error) {
	// Get reference
	reference, err := s.Repo.GetAchievementReferenceByID(referenceID)
	if err != nil {
		return nil, errors.New("achievement reference not found")
	}

	// Get achievement from MongoDB
	mongoID, err := primitive.ObjectIDFromHex(reference.MongoAchievementID)
	if err != nil {
		return nil, errors.New("invalid mongo achievement ID")
	}

	achievement, err := s.Repo.GetAchievementByID(ctx, mongoID)
	if err != nil {
		return nil, errors.New("achievement not found in MongoDB")
	}

	// Get student info
	student, err := s.getStudentInfo(reference.StudentID)
	if err != nil {
		return nil, errors.New("student not found")
	}

	// Get advisor info
	var advisor *AdminAdvisorInfo
	studentRecord, err := s.Repo.GetStudentByID(reference.StudentID)
	if err == nil && studentRecord != nil {
		advisor, _ = s.getAdvisorInfo(studentRecord.AdvisorID)
	}

	return &AdminAchievementResponse{
		Reference:   *reference,
		Achievement: achievement,
		Student:     student,
		Advisor:     advisor,
	}, nil
}