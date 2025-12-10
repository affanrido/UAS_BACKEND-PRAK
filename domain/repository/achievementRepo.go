package repository

import (
	model "UAS_BACKEND/domain/Model"
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type AchievementRepository struct {
	PostgresDB *sql.DB
	MongoDB    *mongo.Database
}

func NewAchievementRepository(postgresDB *sql.DB, mongoDB *mongo.Database) *AchievementRepository {
	return &AchievementRepository{
		PostgresDB: postgresDB,
		MongoDB:    mongoDB,
	}
}

// CreateAchievement - Simpan achievement ke MongoDB
func (r *AchievementRepository) CreateAchievement(ctx context.Context, achievement *model.Achievement) (primitive.ObjectID, error) {
	collection := r.MongoDB.Collection("achievements")

	// Set timestamps
	now := time.Now()
	achievement.CreatedAt = now
	achievement.UpdatedAt = now

	result, err := collection.InsertOne(ctx, achievement)
	if err != nil {
		return primitive.NilObjectID, err
	}

	objectID, ok := result.InsertedID.(primitive.ObjectID)
	if !ok {
		return primitive.NilObjectID, errors.New("failed to get inserted ID")
	}

	return objectID, nil
}

// CreateAchievementReference - Simpan reference ke PostgreSQL
func (r *AchievementRepository) CreateAchievementReference(ref *model.AchievementReference) error {
	query := `
		INSERT INTO achievement_references 
		(id, student_id, mongo_achievement_id, status, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6)
	`

	now := time.Now()
	ref.CreatedAt = now
	ref.UpdatedAt = now

	_, err := r.PostgresDB.Exec(query,
		ref.ID,
		ref.StudentID,
		ref.MongoAchievementID,
		ref.Status,
		ref.CreatedAt,
		ref.UpdatedAt,
	)

	return err
}

// GetAchievementByID - Ambil achievement dari MongoDB
func (r *AchievementRepository) GetAchievementByID(ctx context.Context, id primitive.ObjectID) (*model.Achievement, error) {
	collection := r.MongoDB.Collection("achievements")

	var achievement model.Achievement
	err := collection.FindOne(ctx, bson.M{"_id": id}).Decode(&achievement)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, errors.New("achievement not found")
		}
		return nil, err
	}

	return &achievement, nil
}

// GetAchievementReferenceByID - Ambil reference dari PostgreSQL
func (r *AchievementRepository) GetAchievementReferenceByID(id uuid.UUID) (*model.AchievementReference, error) {
	query := `
		SELECT id, student_id, mongo_achievement_id, status, 
		       submitted_at, verified_at, verified_by, rejection_note,
		       created_at, updated_at
		FROM achievement_references
		WHERE id = $1
	`

	var ref model.AchievementReference
	err := r.PostgresDB.QueryRow(query, id).Scan(
		&ref.ID,
		&ref.StudentID,
		&ref.MongoAchievementID,
		&ref.Status,
		&ref.SubmittedAt,
		&ref.VerifiedAt,
		&ref.VerifiedBy,
		&ref.RejectionNote,
		&ref.CreatedAt,
		&ref.UpdatedAt,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.New("achievement reference not found")
		}
		return nil, err
	}

	return &ref, nil
}

// GetStudentAchievements - Ambil semua achievements mahasiswa
func (r *AchievementRepository) GetStudentAchievements(ctx context.Context, studentID uuid.UUID) ([]model.Achievement, error) {
	collection := r.MongoDB.Collection("achievements")

	cursor, err := collection.Find(ctx, bson.M{"studentId": studentID})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var achievements []model.Achievement
	if err := cursor.All(ctx, &achievements); err != nil {
		return nil, err
	}

	return achievements, nil
}

// GetStudentAchievementReferences - Ambil semua references mahasiswa
func (r *AchievementRepository) GetStudentAchievementReferences(studentID uuid.UUID) ([]model.AchievementReference, error) {
	query := `
		SELECT id, student_id, mongo_achievement_id, status, 
		       submitted_at, verified_at, verified_by, rejection_note,
		       created_at, updated_at
		FROM achievement_references
		WHERE student_id = $1
		ORDER BY created_at DESC
	`

	rows, err := r.PostgresDB.Query(query, studentID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var references []model.AchievementReference
	for rows.Next() {
		var ref model.AchievementReference
		err := rows.Scan(
			&ref.ID,
			&ref.StudentID,
			&ref.MongoAchievementID,
			&ref.Status,
			&ref.SubmittedAt,
			&ref.VerifiedAt,
			&ref.VerifiedBy,
			&ref.RejectionNote,
			&ref.CreatedAt,
			&ref.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		references = append(references, ref)
	}

	return references, nil
}

// UpdateAchievement - Update achievement di MongoDB
func (r *AchievementRepository) UpdateAchievement(ctx context.Context, id primitive.ObjectID, achievement *model.Achievement) error {
	collection := r.MongoDB.Collection("achievements")

	achievement.UpdatedAt = time.Now()

	update := bson.M{
		"$set": achievement,
	}

	result, err := collection.UpdateOne(ctx, bson.M{"_id": id}, update)
	if err != nil {
		return err
	}

	if result.MatchedCount == 0 {
		return errors.New("achievement not found")
	}

	return nil
}

// UpdateAchievementReferenceStatus - Update status reference di PostgreSQL
func (r *AchievementRepository) UpdateAchievementReferenceStatus(id uuid.UUID, status string, verifiedBy *uuid.UUID, rejectionNote *string) error {
	query := `
		UPDATE achievement_references
		SET status = $1, 
		    verified_at = $2,
		    verified_by = $3,
		    rejection_note = $4,
		    updated_at = $5
		WHERE id = $6
	`

	now := time.Now()
	var verifiedAt *time.Time
	if status == "verified" || status == "rejected" {
		verifiedAt = &now
	}

	_, err := r.PostgresDB.Exec(query, status, verifiedAt, verifiedBy, rejectionNote, now, id)
	return err
}

// GetStudentByUserID - Ambil data student berdasarkan user_id
func (r *AchievementRepository) GetStudentByUserID(userID uuid.UUID) (*model.Student, error) {
	query := `
		SELECT id, user_id, student_id, program_study, academic_year, advisor_id, created_at
		FROM students
		WHERE user_id = $1
	`

	var student model.Student
	err := r.PostgresDB.QueryRow(query, userID).Scan(
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
			return nil, errors.New("student not found")
		}
		return nil, err
	}

	return &student, nil
}

// GetUserByID - Ambil data user berdasarkan ID
func (r *AchievementRepository) GetUserByID(userID uuid.UUID) (*model.Users, error) {
	query := `
		SELECT id, username, email, password_hash, full_name, role_id, is_active, created_at, updated_at
		FROM users
		WHERE id = $1
	`

	var user model.Users
	err := r.PostgresDB.QueryRow(query, userID).Scan(
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

// GetLecturerByID - Ambil data lecturer berdasarkan ID
func (r *AchievementRepository) GetLecturerByID(lecturerID uuid.UUID) (*model.Lecturer, error) {
	query := `
		SELECT id, user_id, lecturer_id, department, created_at
		FROM lecturers
		WHERE id = $1
	`

	var lecturer model.Lecturer
	err := r.PostgresDB.QueryRow(query, lecturerID).Scan(
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

// SoftDeleteAchievement - Soft delete achievement di MongoDB
func (r *AchievementRepository) SoftDeleteAchievement(ctx context.Context, id primitive.ObjectID) error {
	collection := r.MongoDB.Collection("achievements")

	now := time.Now()
	update := bson.M{
		"$set": bson.M{
			"isDeleted": true,
			"deletedAt": now,
			"updatedAt": now,
		},
	}

	result, err := collection.UpdateOne(ctx, bson.M{"_id": id}, update)
	if err != nil {
		return err
	}

	if result.MatchedCount == 0 {
		return errors.New("achievement not found")
	}

	return nil
}

// SoftDeleteAchievementReference - Soft delete reference di PostgreSQL
func (r *AchievementRepository) SoftDeleteAchievementReference(id uuid.UUID) error {
	query := `
		UPDATE achievement_references
		SET is_deleted = true,
		    deleted_at = $1,
		    updated_at = $2
		WHERE id = $3
	`

	now := time.Now()
	_, err := r.PostgresDB.Exec(query, now, now, id)
	return err
}

// HardDeleteAchievement - Hard delete achievement di MongoDB (untuk cleanup)
func (r *AchievementRepository) HardDeleteAchievement(ctx context.Context, id primitive.ObjectID) error {
	collection := r.MongoDB.Collection("achievements")

	result, err := collection.DeleteOne(ctx, bson.M{"_id": id})
	if err != nil {
		return err
	}

	if result.DeletedCount == 0 {
		return errors.New("achievement not found")
	}

	return nil
}

// HardDeleteAchievementReference - Hard delete reference di PostgreSQL (untuk cleanup)
func (r *AchievementRepository) HardDeleteAchievementReference(id uuid.UUID) error {
	query := `DELETE FROM achievement_references WHERE id = $1`
	_, err := r.PostgresDB.Exec(query, id)
	return err
}

// GetStudentsByAdvisorID - Ambil semua mahasiswa bimbingan berdasarkan advisor_id
func (r *AchievementRepository) GetStudentsByAdvisorID(advisorID uuid.UUID) ([]model.Student, error) {
	query := `
		SELECT id, user_id, student_id, program_study, academic_year, advisor_id, created_at
		FROM students
		WHERE advisor_id = $1
		ORDER BY created_at DESC
	`

	rows, err := r.PostgresDB.Query(query, advisorID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var students []model.Student
	for rows.Next() {
		var student model.Student
		err := rows.Scan(
			&student.ID,
			&student.UserID,
			&student.StudentID,
			&student.ProgramStudy,
			&student.AcademicYear,
			&student.AdvisorID,
			&student.CreatedAt,
		)
		if err != nil {
			return nil, err
		}
		students = append(students, student)
	}

	return students, nil
}

// GetAchievementReferencesByStudentIDs - Ambil achievement references berdasarkan student_ids dengan pagination
func (r *AchievementRepository) GetAchievementReferencesByStudentIDs(studentIDs []uuid.UUID, limit, offset int) ([]model.AchievementReference, error) {
	if len(studentIDs) == 0 {
		return []model.AchievementReference{}, nil
	}

	query := `
		SELECT id, student_id, mongo_achievement_id, status, 
		       submitted_at, verified_at, verified_by, rejection_note,
		       is_deleted, deleted_at, created_at, updated_at
		FROM achievement_references
		WHERE student_id = ANY($1) AND is_deleted = false
		ORDER BY created_at DESC
		LIMIT $2 OFFSET $3
	`

	rows, err := r.PostgresDB.Query(query, studentIDs, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var references []model.AchievementReference
	for rows.Next() {
		var ref model.AchievementReference
		err := rows.Scan(
			&ref.ID,
			&ref.StudentID,
			&ref.MongoAchievementID,
			&ref.Status,
			&ref.SubmittedAt,
			&ref.VerifiedAt,
			&ref.VerifiedBy,
			&ref.RejectionNote,
			&ref.IsDeleted,
			&ref.DeletedAt,
			&ref.CreatedAt,
			&ref.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		references = append(references, ref)
	}

	return references, nil
}

// CountAchievementReferencesByStudentIDs - Hitung total achievement references berdasarkan student_ids
func (r *AchievementRepository) CountAchievementReferencesByStudentIDs(studentIDs []uuid.UUID) (int, error) {
	if len(studentIDs) == 0 {
		return 0, nil
	}

	query := `
		SELECT COUNT(*)
		FROM achievement_references
		WHERE student_id = ANY($1) AND is_deleted = false
	`

	var count int
	err := r.PostgresDB.QueryRow(query, studentIDs).Scan(&count)
	return count, err
}

// GetAchievementsByIDs - Ambil multiple achievements dari MongoDB berdasarkan IDs
func (r *AchievementRepository) GetAchievementsByIDs(ctx context.Context, ids []primitive.ObjectID) ([]model.Achievement, error) {
	if len(ids) == 0 {
		return []model.Achievement{}, nil
	}

	collection := r.MongoDB.Collection("achievements")

	filter := bson.M{
		"_id": bson.M{"$in": ids},
		"isDeleted": false,
	}

	cursor, err := collection.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var achievements []model.Achievement
	if err := cursor.All(ctx, &achievements); err != nil {
		return nil, err
	}

	return achievements, nil
}

// GetLecturerByUserID - Ambil data lecturer berdasarkan user_id
func (r *AchievementRepository) GetLecturerByUserID(userID uuid.UUID) (*model.Lecturer, error) {
	query := `
		SELECT id, user_id, lecturer_id, department, created_at
		FROM lecturers
		WHERE user_id = $1
	`

	var lecturer model.Lecturer
	err := r.PostgresDB.QueryRow(query, userID).Scan(
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

// GetAllAchievementReferencesAdmin - Get all achievement references for admin with filter and pagination
func (r *AchievementRepository) GetAllAchievementReferencesAdmin(limit, offset int, filter *model.AdminAchievementFilter, sort *model.AdminAchievementSort) ([]model.AchievementReference, int, error) {
	// Build WHERE clause
	whereClause := "WHERE ar.is_deleted = false"
	args := []interface{}{}
	argIndex := 1

	if filter != nil {
		if filter.Status != "" {
			whereClause += " AND ar.status = $" + fmt.Sprintf("%d", argIndex)
			args = append(args, filter.Status)
			argIndex++
		}
		if filter.StudentID != "" {
			whereClause += " AND s.student_id = $" + fmt.Sprintf("%d", argIndex)
			args = append(args, filter.StudentID)
			argIndex++
		}
		if filter.ProgramStudy != "" {
			whereClause += " AND s.program_study ILIKE $" + fmt.Sprintf("%d", argIndex)
			args = append(args, "%"+filter.ProgramStudy+"%")
			argIndex++
		}
		if filter.AdvisorID != "" {
			whereClause += " AND l.lecturer_id = $" + fmt.Sprintf("%d", argIndex)
			args = append(args, filter.AdvisorID)
			argIndex++
		}
	}

	// Count total
	countQuery := `
		SELECT COUNT(*)
		FROM achievement_references ar
		JOIN students s ON ar.student_id = s.id
		JOIN lecturers l ON s.advisor_id = l.id
		` + whereClause

	var total int
	err := r.PostgresDB.QueryRow(countQuery, args...).Scan(&total)
	if err != nil {
		return nil, 0, err
	}

	// Build ORDER BY clause
	orderClause := "ORDER BY ar.created_at DESC"
	if sort != nil {
		switch sort.Field {
		case "created_at":
			orderClause = "ORDER BY ar.created_at"
		case "updated_at":
			orderClause = "ORDER BY ar.updated_at"
		case "status":
			orderClause = "ORDER BY ar.status"
		case "student_name":
			orderClause = "ORDER BY u.full_name"
		}
		
		if sort.Order == "asc" {
			orderClause += " ASC"
		} else {
			orderClause += " DESC"
		}
	}

	// Get data with pagination
	query := `
		SELECT ar.id, ar.student_id, ar.mongo_achievement_id, ar.status, 
		       ar.submitted_at, ar.verified_at, ar.verified_by, ar.rejection_note,
		       ar.is_deleted, ar.deleted_at, ar.created_at, ar.updated_at
		FROM achievement_references ar
		JOIN students s ON ar.student_id = s.id
		JOIN users u ON s.user_id = u.id
		JOIN lecturers l ON s.advisor_id = l.id
		` + whereClause + " " + orderClause + `
		LIMIT $` + fmt.Sprintf("%d", argIndex) + ` OFFSET $` + fmt.Sprintf("%d", argIndex+1)

	args = append(args, limit, offset)

	rows, err := r.PostgresDB.Query(query, args...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var references []model.AchievementReference
	for rows.Next() {
		var ref model.AchievementReference
		err := rows.Scan(
			&ref.ID,
			&ref.StudentID,
			&ref.MongoAchievementID,
			&ref.Status,
			&ref.SubmittedAt,
			&ref.VerifiedAt,
			&ref.VerifiedBy,
			&ref.RejectionNote,
			&ref.IsDeleted,
			&ref.DeletedAt,
			&ref.CreatedAt,
			&ref.UpdatedAt,
		)
		if err != nil {
			return nil, 0, err
		}
		references = append(references, ref)
	}

	return references, total, nil
}

// GetAchievementSummaryAdmin - Get achievement statistics for admin
func (r *AchievementRepository) GetAchievementSummaryAdmin(filter *model.AdminAchievementFilter) (*model.AchievementSummary, error) {
	// Build WHERE clause
	whereClause := "WHERE ar.is_deleted = false"
	args := []interface{}{}
	argIndex := 1

	if filter != nil {
		if filter.StudentID != "" {
			whereClause += " AND s.student_id = $" + fmt.Sprintf("%d", argIndex)
			args = append(args, filter.StudentID)
			argIndex++
		}
		if filter.ProgramStudy != "" {
			whereClause += " AND s.program_study ILIKE $" + fmt.Sprintf("%d", argIndex)
			args = append(args, "%"+filter.ProgramStudy+"%")
			argIndex++
		}
		if filter.AdvisorID != "" {
			whereClause += " AND l.lecturer_id = $" + fmt.Sprintf("%d", argIndex)
			args = append(args, filter.AdvisorID)
			argIndex++
		}
	}

	query := `
		SELECT 
			COUNT(*) as total,
			COUNT(CASE WHEN ar.status = 'draft' THEN 1 END) as draft,
			COUNT(CASE WHEN ar.status = 'submitted' THEN 1 END) as submitted,
			COUNT(CASE WHEN ar.status = 'verified' THEN 1 END) as verified,
			COUNT(CASE WHEN ar.status = 'rejected' THEN 1 END) as rejected
		FROM achievement_references ar
		JOIN students s ON ar.student_id = s.id
		JOIN lecturers l ON s.advisor_id = l.id
		` + whereClause

	var summary model.AchievementSummary
	err := r.PostgresDB.QueryRow(query, args...).Scan(
		&summary.Total,
		&summary.Draft,
		&summary.Submitted,
		&summary.Verified,
		&summary.Rejected,
	)

	return &summary, err
}

// GetStudentByID - Ambil data student berdasarkan ID
func (r *AchievementRepository) GetStudentByID(studentID uuid.UUID) (*model.Student, error) {
	query := `
		SELECT id, user_id, student_id, program_study, academic_year, advisor_id, created_at
		FROM students
		WHERE id = $1
	`

	var student model.Student
	err := r.PostgresDB.QueryRow(query, studentID).Scan(
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
			return nil, errors.New("student not found")
		}
		return nil, err
	}

	return &student, nil
}
