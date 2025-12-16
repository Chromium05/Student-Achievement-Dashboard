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

// PostgreSQL Operations
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
