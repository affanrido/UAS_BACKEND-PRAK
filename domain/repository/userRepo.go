package repository

import (
	model "UAS_BACKEND/domain/Model"
	"database/sql"
	"errors"
	"time"

	"github.com/google/uuid"
)

type UserRepository struct {
	DB *sql.DB
}

func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{DB: db}
}

// CreateUser - Create new user
func (r *UserRepository) CreateUser(user *model.Users) error {
	query := `
		INSERT INTO users (id, username, email, password_hash, full_name, role_id, is_active, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
	`

	now := time.Now()
	user.CreatedAt = now
	user.UpdatedAt = now

	_, err := r.DB.Exec(query,
		user.ID,
		user.Username,
		user.Email,
		user.PasswordHash,
		user.FullName,
		user.RoleID,
		user.IsActive,
		user.CreatedAt,
		user.UpdatedAt,
	)

	return err
}

// UpdateUser - Update existing user
func (r *UserRepository) UpdateUser(user *model.Users) error {
	query := `
		UPDATE users
		SET username = $1, email = $2, password_hash = $3, full_name = $4, 
		    role_id = $5, is_active = $6, updated_at = $7
		WHERE id = $8
	`

	user.UpdatedAt = time.Now()

	_, err := r.DB.Exec(query,
		user.Username,
		user.Email,
		user.PasswordHash,
		user.FullName,
		user.RoleID,
		user.IsActive,
		user.UpdatedAt,
		user.ID,
	)

	return err
}

// DeleteUser - Delete user
func (r *UserRepository) DeleteUser(userID uuid.UUID) error {
	query := `DELETE FROM users WHERE id = $1`
	_, err := r.DB.Exec(query, userID)
	return err
}

// GetUserByID - Get user by ID
func (r *UserRepository) GetUserByID(userID uuid.UUID) (*model.Users, error) {
	query := `
		SELECT id, username, email, password_hash, full_name, role_id, is_active, created_at, updated_at
		FROM users
		WHERE id = $1
	`

	var user model.Users
	err := r.DB.QueryRow(query, userID).Scan(
		&user.ID,
		&user.Username,
		&user.Email,
		&user.PasswordHash,
		&user.FullName,
		&user.RoleID,
		&user.IsActive,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.New("user not found")
		}
		return nil, err
	}

	return &user, nil
}

// GetUserByUsername - Get user by username
func (r *UserRepository) GetUserByUsername(username string) (*model.Users, error) {
	query := `
		SELECT id, username, email, password_hash, full_name, role_id, is_active, created_at, updated_at
		FROM users
		WHERE username = $1
	`

	var user model.Users
	err := r.DB.QueryRow(query, username).Scan(
		&user.ID,
		&user.Username,
		&user.Email,
		&user.PasswordHash,
		&user.FullName,
		&user.RoleID,
		&user.IsActive,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.New("user not found")
		}
		return nil, err
	}

	return &user, nil
}

// GetUserByEmail - Get user by email
func (r *UserRepository) GetUserByEmail(email string) (*model.Users, error) {
	query := `
		SELECT id, username, email, password_hash, full_name, role_id, is_active, created_at, updated_at
		FROM users
		WHERE email = $1
	`

	var user model.Users
	err := r.DB.QueryRow(query, email).Scan(
		&user.ID,
		&user.Username,
		&user.Email,
		&user.PasswordHash,
		&user.FullName,
		&user.RoleID,
		&user.IsActive,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.New("user not found")
		}
		return nil, err
	}

	return &user, nil
}

// GetAllUsers - Get all users with pagination
func (r *UserRepository) GetAllUsers(limit, offset int) ([]model.Users, int, error) {
	// Get total count
	var total int
	countQuery := `SELECT COUNT(*) FROM users`
	err := r.DB.QueryRow(countQuery).Scan(&total)
	if err != nil {
		return nil, 0, err
	}

	// Get users with pagination
	query := `
		SELECT id, username, email, password_hash, full_name, role_id, is_active, created_at, updated_at
		FROM users
		ORDER BY created_at DESC
		LIMIT $1 OFFSET $2
	`

	rows, err := r.DB.Query(query, limit, offset)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var users []model.Users
	for rows.Next() {
		var user model.Users
		err := rows.Scan(
			&user.ID,
			&user.Username,
			&user.Email,
			&user.PasswordHash,
			&user.FullName,
			&user.RoleID,
			&user.IsActive,
			&user.CreatedAt,
			&user.UpdatedAt,
		)
		if err != nil {
			return nil, 0, err
		}
		users = append(users, user)
	}

	return users, total, nil
}

// GetRoleByID - Get role by ID
func (r *UserRepository) GetRoleByID(roleID uuid.UUID) (*model.Roles, error) {
	query := `
		SELECT id, name, description, created_at
		FROM roles
		WHERE id = $1
	`

	var role model.Roles
	err := r.DB.QueryRow(query, roleID).Scan(
		&role.ID,
		&role.Name,
		&role.Description,
		&role.CreatedAt,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.New("role not found")
		}
		return nil, err
	}

	return &role, nil
}

// CreateStudent - Create student profile
func (r *UserRepository) CreateStudent(student *model.Student) error {
	query := `
		INSERT INTO students (id, user_id, student_id, program_study, academic_year, advisor_id, created_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
	`

	now := time.Now()
	student.CreatedAt = now

	_, err := r.DB.Exec(query,
		student.ID,
		student.UserID,
		student.StudentID,
		student.ProgramStudy,
		student.AcademicYear,
		student.AdvisorID,
		student.CreatedAt,
	)

	return err
}

// UpdateStudent - Update student profile
func (r *UserRepository) UpdateStudent(student *model.Student) error {
	query := `
		UPDATE students
		SET student_id = $1, program_study = $2, academic_year = $3, advisor_id = $4
		WHERE id = $5
	`

	_, err := r.DB.Exec(query,
		student.StudentID,
		student.ProgramStudy,
		student.AcademicYear,
		student.AdvisorID,
		student.ID,
	)

	return err
}

// GetStudentByUserID - Get student by user ID
func (r *UserRepository) GetStudentByUserID(userID uuid.UUID) (*model.Student, error) {
	query := `
		SELECT id, user_id, student_id, program_study, academic_year, advisor_id, created_at
		FROM students
		WHERE user_id = $1
	`

	var student model.Student
	err := r.DB.QueryRow(query, userID).Scan(
		&student.ID,
		&student.UserID,
		&student.StudentID,
		&student.ProgramStudy,
		&student.AcademicYear,
		&student.AdvisorID,
		&student.CreatedAt,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil // Not an error, just no student profile
		}
		return nil, err
	}

	return &student, nil
}

// CreateLecturer - Create lecturer profile
func (r *UserRepository) CreateLecturer(lecturer *model.Lecturer) error {
	query := `
		INSERT INTO lecturers (id, user_id, lecturer_id, department, created_at)
		VALUES ($1, $2, $3, $4, $5)
	`

	now := time.Now()
	lecturer.CreatedAt = now

	_, err := r.DB.Exec(query,
		lecturer.ID,
		lecturer.UserID,
		lecturer.LecturerID,
		lecturer.Department,
		lecturer.CreatedAt,
	)

	return err
}

// UpdateLecturer - Update lecturer profile
func (r *UserRepository) UpdateLecturer(lecturer *model.Lecturer) error {
	query := `
		UPDATE lecturers
		SET lecturer_id = $1, department = $2
		WHERE id = $3
	`

	_, err := r.DB.Exec(query,
		lecturer.LecturerID,
		lecturer.Department,
		lecturer.ID,
	)

	return err
}

// GetLecturerByUserID - Get lecturer by user ID
func (r *UserRepository) GetLecturerByUserID(userID uuid.UUID) (*model.Lecturer, error) {
	query := `
		SELECT id, user_id, lecturer_id, department, created_at
		FROM lecturers
		WHERE user_id = $1
	`

	var lecturer model.Lecturer
	err := r.DB.QueryRow(query, userID).Scan(
		&lecturer.ID,
		&lecturer.UserID,
		&lecturer.LecturerID,
		&lecturer.Department,
		&lecturer.CreatedAt,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil // Not an error, just no lecturer profile
		}
		return nil, err
	}

	return &lecturer, nil
}

// GetLecturerByID - Get lecturer by ID
func (r *UserRepository) GetLecturerByID(lecturerID uuid.UUID) (*model.Lecturer, error) {
	query := `
		SELECT id, user_id, lecturer_id, department, created_at
		FROM lecturers
		WHERE id = $1
	`

	var lecturer model.Lecturer
	err := r.DB.QueryRow(query, lecturerID).Scan(
		&lecturer.ID,
		&lecturer.UserID,
		&lecturer.LecturerID,
		&lecturer.Department,
		&lecturer.CreatedAt,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.New("lecturer not found")
		}
		return nil, err
	}

	return &lecturer, nil
}
