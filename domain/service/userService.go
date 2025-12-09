package service

import (
	model "UAS_BACKEND/domain/Model"
	"UAS_BACKEND/domain/repository"
	"errors"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type UserService struct {
	Repo *repository.UserRepository
}

func NewUserService(repo *repository.UserRepository) *UserService {
	return &UserService{Repo: repo}
}

// CreateUserRequest - DTO untuk create user
type CreateUserRequest struct {
	Username     string    `json:"username"`
	Email        string    `json:"email"`
	Password     string    `json:"password"`
	FullName     string    `json:"full_name"`
	RoleID       uuid.UUID `json:"role_id"`
	IsActive     bool      `json:"is_active"`
	StudentData  *StudentProfileRequest  `json:"student_data,omitempty"`
	LecturerData *LecturerProfileRequest `json:"lecturer_data,omitempty"`
}

// UpdateUserRequest - DTO untuk update user
type UpdateUserRequest struct {
	Username string    `json:"username,omitempty"`
	Email    string    `json:"email,omitempty"`
	Password string    `json:"password,omitempty"`
	FullName string    `json:"full_name,omitempty"`
	RoleID   uuid.UUID `json:"role_id,omitempty"`
	IsActive *bool     `json:"is_active,omitempty"`
}

// StudentProfileRequest - DTO untuk student profile
type StudentProfileRequest struct {
	StudentID    string    `json:"student_id"`
	ProgramStudy string    `json:"program_study"`
	AcademicYear string    `json:"academic_year"`
	AdvisorID    uuid.UUID `json:"advisor_id"`
}

// LecturerProfileRequest - DTO untuk lecturer profile
type LecturerProfileRequest struct {
	LecturerID string `json:"lecturer_id"`
	Department string `json:"department"`
}

// UserResponse - DTO untuk response
type UserResponse struct {
	User     model.Users              `json:"user"`
	Student  *model.Student           `json:"student,omitempty"`
	Lecturer *model.Lecturer          `json:"lecturer,omitempty"`
	Role     *model.Roles             `json:"role,omitempty"`
}

// CreateUser - Flow FR-009: Create user
func (s *UserService) CreateUser(req *CreateUserRequest) (*UserResponse, error) {
	// Validasi input
	if err := s.validateCreateUserRequest(req); err != nil {
		return nil, err
	}

	// Check username sudah ada
	existingUser, _ := s.Repo.GetUserByUsername(req.Username)
	if existingUser != nil {
		return nil, errors.New("username already exists")
	}

	// Check email sudah ada
	existingUser, _ = s.Repo.GetUserByEmail(req.Email)
	if existingUser != nil {
		return nil, errors.New("email already exists")
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, errors.New("failed to hash password")
	}

	// 1. Create user
	user := &model.Users{
		ID:           uuid.New(),
		Username:     req.Username,
		Email:        req.Email,
		PasswordHash: string(hashedPassword),
		FullName:     req.FullName,
		RoleID:       req.RoleID,
		IsActive:     req.IsActive,
	}

	err = s.Repo.CreateUser(user)
	if err != nil {
		return nil, errors.New("failed to create user: " + err.Error())
	}

	// 2. Assign role (already set in user creation)
	role, _ := s.Repo.GetRoleByID(user.RoleID)

	var student *model.Student
	var lecturer *model.Lecturer

	// 3. Set student/lecturer profile
	if req.StudentData != nil {
		student = &model.Student{
			ID:           uuid.New(),
			UserID:       user.ID,
			StudentID:    req.StudentData.StudentID,
			ProgramStudy: req.StudentData.ProgramStudy,
			AcademicYear: req.StudentData.AcademicYear,
			AdvisorID:    req.StudentData.AdvisorID,
		}
		err = s.Repo.CreateStudent(student)
		if err != nil {
			return nil, errors.New("failed to create student profile: " + err.Error())
		}
	}

	if req.LecturerData != nil {
		lecturer = &model.Lecturer{
			ID:         uuid.New(),
			UserID:     user.ID,
			LecturerID: req.LecturerData.LecturerID,
			Department: req.LecturerData.Department,
		}
		err = s.Repo.CreateLecturer(lecturer)
		if err != nil {
			return nil, errors.New("failed to create lecturer profile: " + err.Error())
		}
	}

	return &UserResponse{
		User:     *user,
		Student:  student,
		Lecturer: lecturer,
		Role:     role,
	}, nil
}

// UpdateUser - Flow FR-009: Update user
func (s *UserService) UpdateUser(userID uuid.UUID, req *UpdateUserRequest) (*UserResponse, error) {
	// Get existing user
	user, err := s.Repo.GetUserByID(userID)
	if err != nil {
		return nil, errors.New("user not found")
	}

	// Update fields if provided
	if req.Username != "" {
		// Check username conflict
		existingUser, _ := s.Repo.GetUserByUsername(req.Username)
		if existingUser != nil && existingUser.ID != userID {
			return nil, errors.New("username already exists")
		}
		user.Username = req.Username
	}

	if req.Email != "" {
		// Check email conflict
		existingUser, _ := s.Repo.GetUserByEmail(req.Email)
		if existingUser != nil && existingUser.ID != userID {
			return nil, errors.New("email already exists")
		}
		user.Email = req.Email
	}

	if req.Password != "" {
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
		if err != nil {
			return nil, errors.New("failed to hash password")
		}
		user.PasswordHash = string(hashedPassword)
	}

	if req.FullName != "" {
		user.FullName = req.FullName
	}

	if req.RoleID != uuid.Nil {
		user.RoleID = req.RoleID
	}

	if req.IsActive != nil {
		user.IsActive = *req.IsActive
	}

	// Update user
	err = s.Repo.UpdateUser(user)
	if err != nil {
		return nil, errors.New("failed to update user: " + err.Error())
	}

	// Get role
	role, _ := s.Repo.GetRoleByID(user.RoleID)

	// Get student/lecturer profile if exists
	student, _ := s.Repo.GetStudentByUserID(userID)
	lecturer, _ := s.Repo.GetLecturerByUserID(userID)

	return &UserResponse{
		User:     *user,
		Student:  student,
		Lecturer: lecturer,
		Role:     role,
	}, nil
}

// DeleteUser - Flow FR-009: Delete user
func (s *UserService) DeleteUser(userID uuid.UUID) error {
	// Check user exists
	_, err := s.Repo.GetUserByID(userID)
	if err != nil {
		return errors.New("user not found")
	}

	// Delete user (cascade will delete student/lecturer profiles)
	err = s.Repo.DeleteUser(userID)
	if err != nil {
		return errors.New("failed to delete user: " + err.Error())
	}

	return nil
}

// AssignRole - Flow FR-009: Assign role to user
func (s *UserService) AssignRole(userID uuid.UUID, roleID uuid.UUID) (*UserResponse, error) {
	// Get user
	user, err := s.Repo.GetUserByID(userID)
	if err != nil {
		return nil, errors.New("user not found")
	}

	// Check role exists
	role, err := s.Repo.GetRoleByID(roleID)
	if err != nil {
		return nil, errors.New("role not found")
	}

	// Update role
	user.RoleID = roleID
	err = s.Repo.UpdateUser(user)
	if err != nil {
		return nil, errors.New("failed to assign role: " + err.Error())
	}

	// Get student/lecturer profile if exists
	student, _ := s.Repo.GetStudentByUserID(userID)
	lecturer, _ := s.Repo.GetLecturerByUserID(userID)

	return &UserResponse{
		User:     *user,
		Student:  student,
		Lecturer: lecturer,
		Role:     role,
	}, nil
}

// SetStudentProfile - Flow FR-009: Set student profile
func (s *UserService) SetStudentProfile(userID uuid.UUID, req *StudentProfileRequest) (*model.Student, error) {
	// Check user exists
	_, err := s.Repo.GetUserByID(userID)
	if err != nil {
		return nil, errors.New("user not found")
	}

	// Check if student profile already exists
	existingStudent, _ := s.Repo.GetStudentByUserID(userID)
	if existingStudent != nil {
		// Update existing
		existingStudent.StudentID = req.StudentID
		existingStudent.ProgramStudy = req.ProgramStudy
		existingStudent.AcademicYear = req.AcademicYear
		existingStudent.AdvisorID = req.AdvisorID

		err = s.Repo.UpdateStudent(existingStudent)
		if err != nil {
			return nil, errors.New("failed to update student profile: " + err.Error())
		}
		return existingStudent, nil
	}

	// Create new student profile
	student := &model.Student{
		ID:           uuid.New(),
		UserID:       userID,
		StudentID:    req.StudentID,
		ProgramStudy: req.ProgramStudy,
		AcademicYear: req.AcademicYear,
		AdvisorID:    req.AdvisorID,
	}

	err = s.Repo.CreateStudent(student)
	if err != nil {
		return nil, errors.New("failed to create student profile: " + err.Error())
	}

	return student, nil
}

// SetLecturerProfile - Flow FR-009: Set lecturer profile
func (s *UserService) SetLecturerProfile(userID uuid.UUID, req *LecturerProfileRequest) (*model.Lecturer, error) {
	// Check user exists
	_, err := s.Repo.GetUserByID(userID)
	if err != nil {
		return nil, errors.New("user not found")
	}

	// Check if lecturer profile already exists
	existingLecturer, _ := s.Repo.GetLecturerByUserID(userID)
	if existingLecturer != nil {
		// Update existing
		existingLecturer.LecturerID = req.LecturerID
		existingLecturer.Department = req.Department

		err = s.Repo.UpdateLecturer(existingLecturer)
		if err != nil {
			return nil, errors.New("failed to update lecturer profile: " + err.Error())
		}
		return existingLecturer, nil
	}

	// Create new lecturer profile
	lecturer := &model.Lecturer{
		ID:         uuid.New(),
		UserID:     userID,
		LecturerID: req.LecturerID,
		Department: req.Department,
	}

	err = s.Repo.CreateLecturer(lecturer)
	if err != nil {
		return nil, errors.New("failed to create lecturer profile: " + err.Error())
	}

	return lecturer, nil
}

// SetAdvisor - Flow FR-009: Set advisor untuk mahasiswa
func (s *UserService) SetAdvisor(studentID uuid.UUID, advisorID uuid.UUID) (*model.Student, error) {
	// Get student
	student, err := s.Repo.GetStudentByUserID(studentID)
	if err != nil {
		return nil, errors.New("student not found")
	}

	// Check advisor exists
	_, err = s.Repo.GetLecturerByID(advisorID)
	if err != nil {
		return nil, errors.New("advisor not found")
	}

	// Update advisor
	student.AdvisorID = advisorID
	err = s.Repo.UpdateStudent(student)
	if err != nil {
		return nil, errors.New("failed to set advisor: " + err.Error())
	}

	return student, nil
}

// GetAllUsers - Get all users with pagination
func (s *UserService) GetAllUsers(limit, offset int) ([]UserResponse, int, error) {
	users, total, err := s.Repo.GetAllUsers(limit, offset)
	if err != nil {
		return nil, 0, err
	}

	responses := make([]UserResponse, len(users))
	for i, user := range users {
		role, _ := s.Repo.GetRoleByID(user.RoleID)
		student, _ := s.Repo.GetStudentByUserID(user.ID)
		lecturer, _ := s.Repo.GetLecturerByUserID(user.ID)

		responses[i] = UserResponse{
			User:     user,
			Student:  student,
			Lecturer: lecturer,
			Role:     role,
		}
	}

	return responses, total, nil
}

// GetUserByID - Get user by ID
func (s *UserService) GetUserByID(userID uuid.UUID) (*UserResponse, error) {
	user, err := s.Repo.GetUserByID(userID)
	if err != nil {
		return nil, errors.New("user not found")
	}

	role, _ := s.Repo.GetRoleByID(user.RoleID)
	student, _ := s.Repo.GetStudentByUserID(user.ID)
	lecturer, _ := s.Repo.GetLecturerByUserID(user.ID)

	return &UserResponse{
		User:     *user,
		Student:  student,
		Lecturer: lecturer,
		Role:     role,
	}, nil
}

// validateCreateUserRequest - Validasi input
func (s *UserService) validateCreateUserRequest(req *CreateUserRequest) error {
	if req.Username == "" {
		return errors.New("username is required")
	}
	if req.Email == "" {
		return errors.New("email is required")
	}
	if req.Password == "" {
		return errors.New("password is required")
	}
	if len(req.Password) < 6 {
		return errors.New("password must be at least 6 characters")
	}
	if req.FullName == "" {
		return errors.New("full name is required")
	}
	if req.RoleID == uuid.Nil {
		return errors.New("role ID is required")
	}
	return nil
}
