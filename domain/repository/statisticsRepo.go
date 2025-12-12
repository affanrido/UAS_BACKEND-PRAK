package repository

import (
	model "UAS_BACKEND/domain/Model"
	"UAS_BACKEND/domain/service"
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type StatisticsRepository struct {
	PostgresDB *sql.DB
	MongoDB    *mongo.Database
}

func NewStatisticsRepository(postgresDB *sql.DB, mongoDB *mongo.Database) *StatisticsRepository {
	return &StatisticsRepository{
		PostgresDB: postgresDB,
		MongoDB:    mongoDB,
	}
}

// GetAchievementTypeStats - Get statistics per achievement type
func (r *StatisticsRepository) GetAchievementTypeStats(req *service.StatisticsRequest) ([]service.AchievementTypeStats, error) {
	// Build WHERE clause
	whereClause := "WHERE ar.is_deleted = false AND ar.status = 'verified'"
	args := []interface{}{}
	argIndex := 1

	if req.UserID != nil {
		whereClause += " AND s.user_id = $" + fmt.Sprintf("%d", argIndex)
		args = append(args, *req.UserID)
		argIndex++
	}

	if req.AdvisorID != nil {
		whereClause += " AND s.advisor_id = $" + fmt.Sprintf("%d", argIndex)
		args = append(args, *req.AdvisorID)
		argIndex++
	}

	if req.StartDate != nil {
		whereClause += " AND ar.created_at >= $" + fmt.Sprintf("%d", argIndex)
		args = append(args, *req.StartDate)
		argIndex++
	}

	if req.EndDate != nil {
		whereClause += " AND ar.created_at <= $" + fmt.Sprintf("%d", argIndex)
		args = append(args, *req.EndDate)
		argIndex++
	}

	// Get achievement type from MongoDB and count from PostgreSQL
	query := `
		SELECT 
			COALESCE(type_counts.achievement_type, 'unknown') as type,
			COALESCE(type_counts.count, 0) as count,
			COALESCE(type_counts.total, 0) as total
		FROM (
			SELECT 
				'academic' as achievement_type,
				0 as count,
				0 as total
			UNION ALL SELECT 'competition', 0, 0
			UNION ALL SELECT 'organization', 0, 0
			UNION ALL SELECT 'publication', 0, 0
			UNION ALL SELECT 'certification', 0, 0
			UNION ALL SELECT 'other', 0, 0
		) all_types
		LEFT JOIN (
			SELECT 
				ar.mongo_achievement_id,
				COUNT(*) as count,
				COUNT(*) as total
			FROM achievement_references ar
			JOIN students s ON ar.student_id = s.id
			` + whereClause + `
			GROUP BY ar.mongo_achievement_id
		) type_counts ON true
		ORDER BY type
	`

	rows, err := r.PostgresDB.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	// Get MongoDB achievement types
	ctx := context.Background()
	collection := r.MongoDB.Collection("achievements")

	// Build MongoDB filter
	mongoFilter := bson.M{"isDeleted": false}
	if req.StartDate != nil || req.EndDate != nil {
		dateFilter := bson.M{}
		if req.StartDate != nil {
			dateFilter["$gte"] = *req.StartDate
		}
		if req.EndDate != nil {
			dateFilter["$lte"] = *req.EndDate
		}
		mongoFilter["createdAt"] = dateFilter
	}

	// Aggregate by achievement type
	pipeline := []bson.M{
		{"$match": mongoFilter},
		{"$group": bson.M{
			"_id":   "$achievementType",
			"count": bson.M{"$sum": 1},
		}},
	}

	cursor, err := collection.Aggregate(ctx, pipeline)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	typeMap := make(map[string]int)
	for cursor.Next(ctx) {
		var result struct {
			ID    string `bson:"_id"`
			Count int    `bson:"count"`
		}
		if err := cursor.Decode(&result); err != nil {
			continue
		}
		typeMap[result.ID] = result.Count
	}

	// Build response
	var stats []service.AchievementTypeStats
	types := []string{"academic", "competition", "organization", "publication", "certification", "other"}
	
	for _, t := range types {
		count := typeMap[t]
		stats = append(stats, service.AchievementTypeStats{
			Type:  t,
			Count: count,
			Total: count,
		})
	}

	return stats, nil
}

// GetAchievementPeriodStats - Get statistics per period
func (r *StatisticsRepository) GetAchievementPeriodStats(req *service.StatisticsRequest) ([]service.AchievementPeriodStats, error) {
	// Build WHERE clause
	whereClause := "WHERE ar.is_deleted = false AND ar.status = 'verified'"
	args := []interface{}{}
	argIndex := 1

	if req.UserID != nil {
		whereClause += " AND s.user_id = $" + fmt.Sprintf("%d", argIndex)
		args = append(args, *req.UserID)
		argIndex++
	}

	if req.AdvisorID != nil {
		whereClause += " AND s.advisor_id = $" + fmt.Sprintf("%d", argIndex)
		args = append(args, *req.AdvisorID)
		argIndex++
	}

	if req.StartDate != nil {
		whereClause += " AND ar.created_at >= $" + fmt.Sprintf("%d", argIndex)
		args = append(args, *req.StartDate)
		argIndex++
	}

	if req.EndDate != nil {
		whereClause += " AND ar.created_at <= $" + fmt.Sprintf("%d", argIndex)
		args = append(args, *req.EndDate)
		argIndex++
	}

	query := `
		SELECT 
			TO_CHAR(ar.created_at, 'YYYY-MM') as period,
			COUNT(*) as count,
			COUNT(*) as total
		FROM achievement_references ar
		JOIN students s ON ar.student_id = s.id
		` + whereClause + `
		GROUP BY TO_CHAR(ar.created_at, 'YYYY-MM')
		ORDER BY period DESC
		LIMIT 12
	`

	rows, err := r.PostgresDB.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var stats []service.AchievementPeriodStats
	for rows.Next() {
		var stat service.AchievementPeriodStats
		err := rows.Scan(&stat.Period, &stat.Count, &stat.Total)
		if err != nil {
			return nil, err
		}
		stats = append(stats, stat)
	}

	return stats, nil
}

// GetTopStudentStats - Get top students statistics
func (r *StatisticsRepository) GetTopStudentStats(req *service.StatisticsRequest) ([]service.TopStudentStats, error) {
	// Build WHERE clause
	whereClause := "WHERE ar.is_deleted = false AND ar.status = 'verified'"
	args := []interface{}{}
	argIndex := 1

	if req.AdvisorID != nil {
		whereClause += " AND s.advisor_id = $" + fmt.Sprintf("%d", argIndex)
		args = append(args, *req.AdvisorID)
		argIndex++
	}

	if req.StartDate != nil {
		whereClause += " AND ar.created_at >= $" + fmt.Sprintf("%d", argIndex)
		args = append(args, *req.StartDate)
		argIndex++
	}

	if req.EndDate != nil {
		whereClause += " AND ar.created_at <= $" + fmt.Sprintf("%d", argIndex)
		args = append(args, *req.EndDate)
		argIndex++
	}

	query := `
		SELECT 
			s.student_id,
			u.full_name,
			s.program_study,
			s.academic_year,
			COUNT(*) as total_count,
			COALESCE(SUM(points.total_points), 0) as total_points
		FROM achievement_references ar
		JOIN students s ON ar.student_id = s.id
		JOIN users u ON s.user_id = u.id
		LEFT JOIN (
			SELECT 
				ar2.id,
				COALESCE(mongo_points.points, 0) as total_points
			FROM achievement_references ar2
			LEFT JOIN (
				SELECT 1 as dummy, 0 as points
			) mongo_points ON true
		) points ON ar.id = points.id
		` + whereClause + `
		GROUP BY s.student_id, u.full_name, s.program_study, s.academic_year
		ORDER BY total_count DESC, total_points DESC
		LIMIT 10
	`

	rows, err := r.PostgresDB.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var stats []service.TopStudentStats
	for rows.Next() {
		var stat service.TopStudentStats
		err := rows.Scan(
			&stat.StudentID,
			&stat.StudentName,
			&stat.ProgramStudy,
			&stat.AcademicYear,
			&stat.TotalCount,
			&stat.TotalPoints,
		)
		if err != nil {
			return nil, err
		}
		stats = append(stats, stat)
	}

	// Get actual points from MongoDB
	for i := range stats {
		points, err := r.getStudentPointsFromMongo(stats[i].StudentID)
		if err == nil {
			stats[i].TotalPoints = points
		}
	}

	return stats, nil
}

// GetCompetitionLevelStats - Get competition level distribution
func (r *StatisticsRepository) GetCompetitionLevelStats(ctx context.Context, req *service.StatisticsRequest) ([]service.CompetitionLevelStats, error) {
	collection := r.MongoDB.Collection("achievements")

	// Build MongoDB filter
	mongoFilter := bson.M{
		"isDeleted":       false,
		"achievementType": "competition",
	}

	if req.StartDate != nil || req.EndDate != nil {
		dateFilter := bson.M{}
		if req.StartDate != nil {
			dateFilter["$gte"] = *req.StartDate
		}
		if req.EndDate != nil {
			dateFilter["$lte"] = *req.EndDate
		}
		mongoFilter["createdAt"] = dateFilter
	}

	// If filtering by user or advisor, get student IDs first
	if req.UserID != nil || req.AdvisorID != nil {
		studentIDs, err := r.getFilteredStudentIDs(req)
		if err != nil {
			return nil, err
		}
		if len(studentIDs) == 0 {
			return []service.CompetitionLevelStats{}, nil
		}
		mongoFilter["studentId"] = bson.M{"$in": studentIDs}
	}

	// Aggregate by competition level
	pipeline := []bson.M{
		{"$match": mongoFilter},
		{"$group": bson.M{
			"_id":   "$details.competitionLevel",
			"count": bson.M{"$sum": 1},
		}},
	}

	cursor, err := collection.Aggregate(ctx, pipeline)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	levelMap := make(map[string]int)
	for cursor.Next(ctx) {
		var result struct {
			ID    *string `bson:"_id"`
			Count int     `bson:"count"`
		}
		if err := cursor.Decode(&result); err != nil {
			continue
		}
		level := "unknown"
		if result.ID != nil {
			level = *result.ID
		}
		levelMap[level] = result.Count
	}

	// Build response
	var stats []service.CompetitionLevelStats
	levels := []string{"international", "national", "regional", "local", "unknown"}
	
	for _, level := range levels {
		count := levelMap[level]
		if count > 0 {
			stats = append(stats, service.CompetitionLevelStats{
				Level: level,
				Count: count,
				Total: count,
			})
		}
	}

	return stats, nil
}

// GetStatisticsSummary - Get summary statistics
func (r *StatisticsRepository) GetStatisticsSummary(req *service.StatisticsRequest) (*service.StatisticsSummary, error) {
	// Build WHERE clause
	whereClause := "WHERE ar.is_deleted = false"
	args := []interface{}{}
	argIndex := 1

	if req.UserID != nil {
		whereClause += " AND s.user_id = $" + fmt.Sprintf("%d", argIndex)
		args = append(args, *req.UserID)
		argIndex++
	}

	if req.AdvisorID != nil {
		whereClause += " AND s.advisor_id = $" + fmt.Sprintf("%d", argIndex)
		args = append(args, *req.AdvisorID)
		argIndex++
	}

	if req.StartDate != nil {
		whereClause += " AND ar.created_at >= $" + fmt.Sprintf("%d", argIndex)
		args = append(args, *req.StartDate)
		argIndex++
	}

	if req.EndDate != nil {
		whereClause += " AND ar.created_at <= $" + fmt.Sprintf("%d", argIndex)
		args = append(args, *req.EndDate)
		argIndex++
	}

	query := `
		SELECT 
			COUNT(*) as total_achievements,
			COUNT(CASE WHEN ar.status = 'verified' THEN 1 END) as verified_count,
			COUNT(CASE WHEN ar.status = 'submitted' THEN 1 END) as pending_count
		FROM achievement_references ar
		JOIN students s ON ar.student_id = s.id
		` + whereClause

	var summary service.StatisticsSummary
	err := r.PostgresDB.QueryRow(query, args...).Scan(
		&summary.TotalAchievements,
		&summary.VerifiedCount,
		&summary.PendingCount,
	)
	if err != nil {
		return nil, err
	}

	// Get total points from MongoDB (simplified)
	summary.TotalPoints = float64(summary.VerifiedCount * 50) // Placeholder calculation
	if summary.TotalAchievements > 0 {
		summary.AveragePoints = summary.TotalPoints / float64(summary.TotalAchievements)
	}

	return &summary, nil
}

// GetAchievementTrends - Get achievement trends over time
func (r *StatisticsRepository) GetAchievementTrends(req *service.StatisticsRequest, months int) ([]service.TrendData, error) {
	// Build WHERE clause
	whereClause := "WHERE ar.is_deleted = false AND ar.created_at >= NOW() - INTERVAL '" + fmt.Sprintf("%d", months) + " months'"
	args := []interface{}{}
	argIndex := 1

	if req.UserID != nil {
		whereClause += " AND s.user_id = $" + fmt.Sprintf("%d", argIndex)
		args = append(args, *req.UserID)
		argIndex++
	}

	if req.AdvisorID != nil {
		whereClause += " AND s.advisor_id = $" + fmt.Sprintf("%d", argIndex)
		args = append(args, *req.AdvisorID)
		argIndex++
	}

	query := `
		SELECT 
			TO_CHAR(ar.created_at, 'YYYY-MM') as month,
			COUNT(*) as count,
			COUNT(CASE WHEN ar.status = 'verified' THEN 1 END) as verified,
			COUNT(CASE WHEN ar.status = 'submitted' THEN 1 END) as submitted
		FROM achievement_references ar
		JOIN students s ON ar.student_id = s.id
		` + whereClause + `
		GROUP BY TO_CHAR(ar.created_at, 'YYYY-MM')
		ORDER BY month DESC
	`

	rows, err := r.PostgresDB.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var trends []service.TrendData
	for rows.Next() {
		var trend service.TrendData
		err := rows.Scan(&trend.Month, &trend.Count, &trend.Verified, &trend.Submitted)
		if err != nil {
			return nil, err
		}
		trend.Points = float64(trend.Verified * 50) // Placeholder calculation
		trends = append(trends, trend)
	}

	return trends, nil
}

// Helper methods

// getFilteredStudentIDs - Get student IDs based on filter
func (r *StatisticsRepository) getFilteredStudentIDs(req *service.StatisticsRequest) ([]uuid.UUID, error) {
	whereClause := "WHERE 1=1"
	args := []interface{}{}
	argIndex := 1

	if req.UserID != nil {
		whereClause += " AND user_id = $" + fmt.Sprintf("%d", argIndex)
		args = append(args, *req.UserID)
		argIndex++
	}

	if req.AdvisorID != nil {
		whereClause += " AND advisor_id = $" + fmt.Sprintf("%d", argIndex)
		args = append(args, *req.AdvisorID)
		argIndex++
	}

	query := `SELECT id FROM students ` + whereClause

	rows, err := r.PostgresDB.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var studentIDs []uuid.UUID
	for rows.Next() {
		var id uuid.UUID
		if err := rows.Scan(&id); err != nil {
			continue
		}
		studentIDs = append(studentIDs, id)
	}

	return studentIDs, nil
}

// getStudentPointsFromMongo - Get total points for student from MongoDB
func (r *StatisticsRepository) getStudentPointsFromMongo(studentID string) (float64, error) {
	// Get student UUID first
	query := `SELECT id FROM students WHERE student_id = $1`
	var studentUUID uuid.UUID
	err := r.PostgresDB.QueryRow(query, studentID).Scan(&studentUUID)
	if err != nil {
		return 0, err
	}

	// Get points from MongoDB
	ctx := context.Background()
	collection := r.MongoDB.Collection("achievements")

	pipeline := []bson.M{
		{"$match": bson.M{
			"studentId":  studentUUID,
			"isDeleted": false,
		}},
		{"$group": bson.M{
			"_id":         nil,
			"totalPoints": bson.M{"$sum": "$points"},
		}},
	}

	cursor, err := collection.Aggregate(ctx, pipeline)
	if err != nil {
		return 0, err
	}
	defer cursor.Close(ctx)

	if cursor.Next(ctx) {
		var result struct {
			TotalPoints float64 `bson:"totalPoints"`
		}
		if err := cursor.Decode(&result); err != nil {
			return 0, err
		}
		return result.TotalPoints, nil
	}

	return 0, nil
}

// GetStudentByUserID - Get student by user ID
func (r *StatisticsRepository) GetStudentByUserID(userID uuid.UUID) (*model.Student, error) {
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
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("student not found")
		}
		return nil, err
	}

	return &student, nil
}

// GetLecturerByUserID - Get lecturer by user ID
func (r *StatisticsRepository) GetLecturerByUserID(userID uuid.UUID) (*model.Lecturer, error) {
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
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("lecturer not found")
		}
		return nil, err
	}

	return &lecturer, nil
}