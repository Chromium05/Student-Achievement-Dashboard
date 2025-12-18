package repository

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"student-report/app/model"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type AchievementRepository struct {
	mongoDB *mongo.Database
	sqlDB   *sql.DB
}

var _ IAchievementRepository = (*AchievementRepository)(nil)

// IAchievementRepository defines the interface for achievement repository operations
type IAchievementRepository interface {
	// MongoDB Operations
	CreateAchievementMongo(ctx context.Context, achievement *model.Achievement) (string, error)
	GetAchievementMongo(ctx context.Context, mongoID string) (*model.Achievement, error)
	UpdateAchievementMongo(ctx context.Context, mongoID string, update *model.Achievement) error
	SoftDeleteAchievementMongo(ctx context.Context, mongoID string) error
	AddAttachmentMongo(ctx context.Context, mongoID string, attachment model.Attachment) error
	GetAchievementsByStudentIDs(ctx context.Context, studentIDs []string) ([]model.Achievement, error)
	GetAchievementsWithFilter(ctx context.Context, filter model.AchievementFilter) ([]model.Achievement, int64, error)
	GetAchievementStatistics(ctx context.Context, studentIDs []string) (*model.AchievementStatistics, error)
	GetTopStudents(ctx context.Context, limit int) ([]model.StudentAchievementCount, error)
	GetStudentInfo(studentID string) (string, string)

	// PostgreSQL Operations
	CreateAchievementReference(studentID, mongoID string) (string, error)
	GetAchievementReference(id string) (*model.AchievementReference, error)
	GetAchievementReferenceByMongoID(mongoID string) (*model.AchievementReference, error)
	UpdateAchievementStatus(id, status string, submittedAt *time.Time) error
	VerifyAchievement(id, verifiedBy string) error
	RejectAchievement(id, note string) error
	GetAchievementsByStudentID(studentID string) ([]model.AchievementReference, error)
	GetAchievementReferencesByStudentIDs(studentIDs []string) ([]model.AchievementReference, error)
	GetAllAchievementReferences() ([]model.AchievementReference, error)
}

func NewAchievementRepository(mongoDB *mongo.Database, sqlDB *sql.DB) *AchievementRepository {
	return &AchievementRepository{
		mongoDB: mongoDB,
		sqlDB:   sqlDB,
	}
}

// MongoDB Operations
func (r *AchievementRepository) CreateAchievementMongo(ctx context.Context, achievement *model.Achievement) (string, error) {
	collection := r.mongoDB.Collection("achievements")
	achievement.CreatedAt = time.Now()
	achievement.UpdatedAt = time.Now()
	achievement.IsDeleted = false

	result, err := collection.InsertOne(ctx, achievement)
	if err != nil {
		return "", err
	}

	return result.InsertedID.(primitive.ObjectID).Hex(), nil
}

func (r *AchievementRepository) GetAchievementMongo(ctx context.Context, mongoID string) (*model.Achievement, error) {
	collection := r.mongoDB.Collection("achievements")
	objectID, err := primitive.ObjectIDFromHex(mongoID)
	if err != nil {
		return nil, err
	}

	var achievement model.Achievement
	err = collection.FindOne(ctx, bson.M{"_id": objectID, "isDeleted": false}).Decode(&achievement)
	if err != nil {
		return nil, err
	}

	return &achievement, nil
}

func (r *AchievementRepository) UpdateAchievementMongo(ctx context.Context, mongoID string, update *model.Achievement) error {
	collection := r.mongoDB.Collection("achievements")
	objectID, err := primitive.ObjectIDFromHex(mongoID)
	if err != nil {
		return err
	}

	update.UpdatedAt = time.Now()
	updateDoc := bson.M{
		"$set": bson.M{
			"title":       update.Title,
			"description": update.Description,
			"details":     update.Details,
			"tags":        update.Tags,
			"points":      update.Points,
			"updatedAt":   update.UpdatedAt,
		},
	}

	_, err = collection.UpdateOne(ctx, bson.M{"_id": objectID}, updateDoc)
	return err
}

func (r *AchievementRepository) SoftDeleteAchievementMongo(ctx context.Context, mongoID string) error {
	collection := r.mongoDB.Collection("achievements")
	objectID, err := primitive.ObjectIDFromHex(mongoID)
	if err != nil {
		return err
	}

	updateDoc := bson.M{
		"$set": bson.M{
			"isDeleted": true,
			"updatedAt": time.Now(),
		},
	}

	_, err = collection.UpdateOne(ctx, bson.M{"_id": objectID}, updateDoc)
	return err
}

func (r *AchievementRepository) AddAttachmentMongo(ctx context.Context, mongoID string, attachment model.Attachment) error {
	collection := r.mongoDB.Collection("achievements")
	objectID, err := primitive.ObjectIDFromHex(mongoID)
	if err != nil {
		return err
	}

	attachment.UploadedAt = time.Now()
	updateDoc := bson.M{
		"$push": bson.M{"attachments": attachment},
		"$set":  bson.M{"updatedAt": time.Now()},
	}

	_, err = collection.UpdateOne(ctx, bson.M{"_id": objectID}, updateDoc)
	return err
}

func (r *AchievementRepository) GetAchievementsByStudentIDs(ctx context.Context, studentIDs []string) ([]model.Achievement, error) {
	collection := r.mongoDB.Collection("achievements")
	
	filter := bson.M{
		"studentId": bson.M{"$in": studentIDs},
		"isDeleted": false,
	}
	
	cursor, err := collection.Find(ctx, filter, options.Find().SetSort(bson.M{"createdAt": -1}))
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var achievements []model.Achievement
	if err = cursor.All(ctx, &achievements); err != nil {
		return nil, err
	}

	return achievements, nil
}

// New function for filtered query with sorting and pagination
func (r *AchievementRepository) GetAchievementsWithFilter(ctx context.Context, filter model.AchievementFilter) ([]model.Achievement, int64, error) {
	collection := r.mongoDB.Collection("achievements")
	
	// Build MongoDB filter
	mongoFilter := bson.M{"isDeleted": false}
	
	if filter.AchievementType != nil && *filter.AchievementType != "" {
		mongoFilter["achievementType"] = *filter.AchievementType
	}
	
	if filter.StudentID != nil && *filter.StudentID != "" {
		mongoFilter["studentId"] = *filter.StudentID
	}
	
	if filter.DateFrom != nil && *filter.DateFrom != "" {
		dateFrom, _ := time.Parse("2006-01-02", *filter.DateFrom)
		mongoFilter["createdAt"] = bson.M{"$gte": dateFrom}
	}
	
	if filter.DateTo != nil && *filter.DateTo != "" {
		dateTo, _ := time.Parse("2006-01-02", *filter.DateTo)
		if mongoFilter["createdAt"] != nil {
			mongoFilter["createdAt"].(bson.M)["$lte"] = dateTo
		} else {
			mongoFilter["createdAt"] = bson.M{"$lte": dateTo}
		}
	}
	
	// Count total
	total, err := collection.CountDocuments(ctx, mongoFilter)
	if err != nil {
		return nil, 0, err
	}
	
	// Build sort
	sortField := "createdAt"
	if filter.SortBy != nil && *filter.SortBy != "" {
		sortField = *filter.SortBy
	}
	
	sortOrder := -1 // desc by default
	if filter.SortOrder != nil && *filter.SortOrder == "asc" {
		sortOrder = 1
	}
	
	// Pagination
	page := 1
	if filter.Page > 0 {
		page = filter.Page
	}
	
	limit := 10
	if filter.Limit > 0 {
		limit = filter.Limit
	}
	
	skip := (page - 1) * limit
	
	opts := options.Find().
		SetSort(bson.M{sortField: sortOrder}).
		SetSkip(int64(skip)).
		SetLimit(int64(limit))
	
	cursor, err := collection.Find(ctx, mongoFilter, opts)
	if err != nil {
		return nil, 0, err
	}
	defer cursor.Close(ctx)
	
	var achievements []model.Achievement
	if err = cursor.All(ctx, &achievements); err != nil {
		return nil, 0, err
	}
	
	return achievements, total, nil
}

// New function for statistics aggregation methods
func (r *AchievementRepository) GetAchievementStatistics(ctx context.Context, studentIDs []string) (*model.AchievementStatistics, error) {
	collection := r.mongoDB.Collection("achievements")
	
	// Build filter - if studentIDs provided, filter by them
	matchFilter := bson.M{"isDeleted": false}
	if len(studentIDs) > 0 {
		matchFilter["studentId"] = bson.M{"$in": studentIDs}
	}
	
	// Aggregation pipeline for statistics
	pipeline := mongo.Pipeline{
		{{Key: "$match", Value: matchFilter}},
		{{Key: "$facet", Value: bson.M{
			"totalCount": bson.A{
				bson.M{"$count": "count"},
			},
			"byType": bson.A{
				bson.M{"$group": bson.M{
					"_id":   "$achievementType",
					"count": bson.M{"$sum": 1},
				}},
			},
			"byCompetitionLevel": bson.A{
				bson.M{"$match": bson.M{"achievementType": "competition"}},
				bson.M{"$group": bson.M{
					"_id":   "$details.competitionLevel",
					"count": bson.M{"$sum": 1},
				}},
			},
			"totalPoints": bson.A{
				bson.M{"$group": bson.M{
					"_id":         nil,
					"totalPoints": bson.M{"$sum": "$points"},
				}},
			},
		}}},
	}
	
	cursor, err := collection.Aggregate(ctx, pipeline)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)
	
	var results []bson.M
	if err = cursor.All(ctx, &results); err != nil {
		return nil, err
	}
	
	stats := &model.AchievementStatistics{
		ByType:             make(map[string]int),
		ByStatus:           make(map[string]int),
		ByCompetitionLevel: make(map[string]int),
		TopStudents:        []model.StudentAchievementCount{},
	}
	
	if len(results) > 0 {
		result := results[0]
		
		// Total count
		if totalCount, ok := result["totalCount"].(bson.A); ok && len(totalCount) > 0 {
			if countDoc, ok := totalCount[0].(bson.M); ok {
				if count, ok := countDoc["count"].(int32); ok {
					stats.TotalAchievements = int(count)
				}
			}
		}
		
		// By type
		if byType, ok := result["byType"].(bson.A); ok {
			for _, item := range byType {
				if doc, ok := item.(bson.M); ok {
					typeStr := doc["_id"].(string)
					count := int(doc["count"].(int32))
					stats.ByType[typeStr] = count
				}
			}
		}
		
		// By competition level
		if byLevel, ok := result["byCompetitionLevel"].(bson.A); ok {
			for _, item := range byLevel {
				if doc, ok := item.(bson.M); ok {
					if level, ok := doc["_id"].(string); ok {
						count := int(doc["count"].(int32))
						stats.ByCompetitionLevel[level] = count
					}
				}
			}
		}
		
		// Total points
		if totalPoints, ok := result["totalPoints"].(bson.A); ok && len(totalPoints) > 0 {
			if pointsDoc, ok := totalPoints[0].(bson.M); ok {
				if points, ok := pointsDoc["totalPoints"].(int32); ok {
					stats.TotalPoints = int(points)
				} else if points, ok := pointsDoc["totalPoints"].(int64); ok {
					stats.TotalPoints = int(points)
				}
			}
		}
	}
	
	// Get status statistics from PostgreSQL
	statusStats, err := r.getStatusStatistics(studentIDs)
	if err == nil {
		stats.ByStatus = statusStats
		stats.PendingVerification = statusStats["submitted"]
		stats.RecentVerified = statusStats["verified"]
	}
	
	return stats, nil
}

func (r *AchievementRepository) getStatusStatistics(studentIDs []string) (map[string]int, error) {
	stats := make(map[string]int)
	
	var query string
	var args []interface{}
	
	if len(studentIDs) > 0 {
		placeholders := ""
		args = make([]interface{}, len(studentIDs))
		for i, id := range studentIDs {
			if i > 0 {
				placeholders += ", "
			}
			placeholders += fmt.Sprintf("$%d", i+1)
			args[i] = id
		}
		query = fmt.Sprintf(`
			SELECT status, COUNT(*) as count
			FROM achievement_references
			WHERE student_id IN (%s)
			GROUP BY status
		`, placeholders)
	} else {
		query = `
			SELECT status, COUNT(*) as count
			FROM achievement_references
			GROUP BY status
		`
	}
	
	rows, err := r.sqlDB.Query(query, args...)
	if err != nil {
		return stats, err
	}
	defer rows.Close()
	
	for rows.Next() {
		var status string
		var count int
		if err := rows.Scan(&status, &count); err == nil {
			stats[status] = count
		}
	}
	
	return stats, nil
}

// New function for getting top students by achievement count
func (r *AchievementRepository) GetTopStudents(ctx context.Context, limit int) ([]model.StudentAchievementCount, error) {
	collection := r.mongoDB.Collection("achievements")
	
	pipeline := mongo.Pipeline{
		{{Key: "$match", Value: bson.M{"isDeleted": false}}},
		{{Key: "$group", Value: bson.M{
			"_id":         "$studentId",
			"count":       bson.M{"$sum": 1},
			"totalPoints": bson.M{"$sum": "$points"},
		}}},
		{{Key: "$sort", Value: bson.M{"count": -1}}},
		{{Key: "$limit", Value: limit}},
	}
	
	cursor, err := collection.Aggregate(ctx, pipeline)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)
	
	var topStudents []model.StudentAchievementCount
	for cursor.Next(ctx) {
		var result struct {
			ID          string `bson:"_id"`
			Count       int32  `bson:"count"`
			TotalPoints int32  `bson:"totalPoints"`
		}
		
		if err := cursor.Decode(&result); err != nil {
			continue
		}
		
		// Get student details from PostgreSQL
		studentName, nim := r.getStudentInfo(result.ID)
		
		topStudents = append(topStudents, model.StudentAchievementCount{
			StudentID:   result.ID,
			StudentName: studentName,
			NIM:         nim,
			Count:       int(result.Count),
			TotalPoints: int(result.TotalPoints),
		})
	}
	
	return topStudents, nil
}

func (r *AchievementRepository) getStudentInfo(studentID string) (string, string) {
	query := `
		SELECT u.full_name, s.student_id
		FROM students s
		JOIN users u ON s.user_id = u.id
		WHERE s.id = $1
	`
	
	var fullName, nim string
	err := r.sqlDB.QueryRow(query, studentID).Scan(&fullName, &nim)
	if err != nil {
		return "Unknown", "Unknown"
	}
	
	return fullName, nim
}

func (r *AchievementRepository) GetStudentInfo(studentID string) (string, string) {
	return r.getStudentInfo(studentID)
}

func (r *AchievementRepository) GetAchievementReference(id string) (*model.AchievementReference, error) {
	query := `
		SELECT id, student_id, mongo_achievement_id, status, submitted_at, verified_at, 
		       verified_by, rejection_note, created_at, updated_at
		FROM achievement_references
		WHERE id = $1
	`
	var ref model.AchievementReference
	err := r.sqlDB.QueryRow(query, id).Scan(
		&ref.ID, &ref.StudentID, &ref.MongoAchievementID, &ref.Status,
		&ref.SubmittedAt, &ref.VerifiedAt, &ref.VerifiedBy, &ref.RejectionNote,
		&ref.CreatedAt, &ref.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	return &ref, nil
}

func (r *AchievementRepository) GetAchievementReferenceByMongoID(mongoID string) (*model.AchievementReference, error) {
	query := `
		SELECT id, student_id, mongo_achievement_id, status, submitted_at, verified_at, 
		       verified_by, rejection_note, created_at, updated_at
		FROM achievement_references
		WHERE mongo_achievement_id = $1
	`
	var ref model.AchievementReference
	err := r.sqlDB.QueryRow(query, mongoID).Scan(
		&ref.ID, &ref.StudentID, &ref.MongoAchievementID, &ref.Status,
		&ref.SubmittedAt, &ref.VerifiedAt, &ref.VerifiedBy, &ref.RejectionNote,
		&ref.CreatedAt, &ref.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	return &ref, nil
}

func (r *AchievementRepository) UpdateAchievementStatus(id, status string, submittedAt *time.Time) error {
	query := `
		UPDATE achievement_references
		SET status = $1, submitted_at = $2, updated_at = NOW()
		WHERE id = $3
	`
	_, err := r.sqlDB.Exec(query, status, submittedAt, id)
	return err
}

func (r *AchievementRepository) VerifyAchievement(id, verifiedBy string) error {
	now := time.Now()
	query := `
		UPDATE achievement_references
		SET status = 'verified', verified_at = $1, verified_by = $2, updated_at = NOW()
		WHERE id = $3
	`
	_, err := r.sqlDB.Exec(query, now, verifiedBy, id)
	return err
}

func (r *AchievementRepository) RejectAchievement(id, note string) error {
	query := `
		UPDATE achievement_references
		SET status = 'rejected', rejection_note = $1, updated_at = NOW()
		WHERE id = $2
	`
	_, err := r.sqlDB.Exec(query, note, id)
	return err
}

func (r *AchievementRepository) GetAchievementsByStudentID(studentID string) ([]model.AchievementReference, error) {
	query := `
		SELECT id, student_id, mongo_achievement_id, status, submitted_at, verified_at,
		       verified_by, rejection_note, created_at, updated_at
		FROM achievement_references
		WHERE student_id = $1
		ORDER BY created_at DESC
	`
	rows, err := r.sqlDB.Query(query, studentID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var refs []model.AchievementReference
	for rows.Next() {
		var ref model.AchievementReference
		err := rows.Scan(
			&ref.ID, &ref.StudentID, &ref.MongoAchievementID, &ref.Status,
			&ref.SubmittedAt, &ref.VerifiedAt, &ref.VerifiedBy, &ref.RejectionNote,
			&ref.CreatedAt, &ref.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		refs = append(refs, ref)
	}
	return refs, nil
}

func (r *AchievementRepository) GetAchievementReferencesByStudentIDs(studentIDs []string) ([]model.AchievementReference, error) {
	if len(studentIDs) == 0 {
		return []model.AchievementReference{}, nil
	}

	// Build IN clause
	placeholders := ""
	args := make([]interface{}, len(studentIDs))
	for i, id := range studentIDs {
		if i > 0 {
			placeholders += ", "
		}
		placeholders += fmt.Sprintf("$%d", i+1)
		args[i] = id
	}

	query := fmt.Sprintf(`
		SELECT id, student_id, mongo_achievement_id, status, submitted_at, verified_at,
		       verified_by, rejection_note, created_at, updated_at
		FROM achievement_references
		WHERE student_id IN (%s)
		ORDER BY created_at DESC
	`, placeholders)

	rows, err := r.sqlDB.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var refs []model.AchievementReference
	for rows.Next() {
		var ref model.AchievementReference
		err := rows.Scan(
			&ref.ID, &ref.StudentID, &ref.MongoAchievementID, &ref.Status,
			&ref.SubmittedAt, &ref.VerifiedAt, &ref.VerifiedBy, &ref.RejectionNote,
			&ref.CreatedAt, &ref.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		refs = append(refs, ref)
	}
	return refs, nil
}

func (r *AchievementRepository) GetAllAchievementReferences() ([]model.AchievementReference, error) {
	query := `
		SELECT id, student_id, mongo_achievement_id, status, submitted_at, verified_at,
		       verified_by, rejection_note, created_at, updated_at
		FROM achievement_references
		ORDER BY created_at DESC
	`
	rows, err := r.sqlDB.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var refs []model.AchievementReference
	for rows.Next() {
		var ref model.AchievementReference
		err := rows.Scan(
			&ref.ID, &ref.StudentID, &ref.MongoAchievementID, &ref.Status,
			&ref.SubmittedAt, &ref.VerifiedAt, &ref.VerifiedBy, &ref.RejectionNote,
			&ref.CreatedAt, &ref.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		refs = append(refs, ref)
	}
	return refs, nil
}

func (r *AchievementRepository) CreateAchievementReference(studentID, mongoID string) (string, error) {
	query := `
		INSERT INTO achievement_references (id, student_id, mongo_achievement_id, status, created_at, updated_at)
		VALUES (gen_random_uuid(), $1, $2, 'draft', NOW(), NOW())
		RETURNING id
	`
	var id string
	err := r.sqlDB.QueryRow(query, studentID, mongoID).Scan(&id)
	return id, err
}
