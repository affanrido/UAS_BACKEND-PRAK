package repository

import (
	model "UAS_BACKEND/domain/Model"
	"context"
	"database/sql"
	"errors"
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
