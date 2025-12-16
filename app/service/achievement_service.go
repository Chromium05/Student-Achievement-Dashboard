package service

import (
	"context"
	"errors"
	"fmt"
	"time"

	"student-report/app/model"
	"student-report/app/repository"
)

type AchievementService struct {
	repo *repository.AchievementRepository
}

func NewAchievementService(repo *repository.AchievementRepository) *AchievementService {
	return &AchievementService{repo: repo}
}

// FR-003: Create Achievement (Student)
func (s *AchievementService) CreateAchievement(ctx context.Context, studentID string, req model.CreateAchievementRequest) (*model.AchievementResponse, error) {
	// Create achievement in MongoDB
	achievement := &model.Achievement{
		StudentID:       studentID,
		AchievementType: req.AchievementType,
		Title:           req.Title,
		Description:     req.Description,
		Details:         req.Details,
		Tags:            req.Tags,
		Points:          req.Points,
		Attachments:     []model.Attachment{},
	}

	mongoID, err := s.repo.CreateAchievementMongo(ctx, achievement)
	if err != nil {
		return nil, fmt.Errorf("failed to create achievement in MongoDB: %w", err)
	}

	// Create reference in PostgreSQL
	refID, err := s.repo.CreateAchievementReference(studentID, mongoID)
	if err != nil {
		return nil, fmt.Errorf("failed to create achievement reference: %w", err)
	}

	// Get the created achievement
	createdAchievement, err := s.repo.GetAchievementMongo(ctx, mongoID)
	if err != nil {
		return nil, err
	}

	ref, err := s.repo.GetAchievementReference(refID)
	if err != nil {
		return nil, err
	}

	return s.combineAchievementResponse(createdAchievement, ref), nil
}

// Get Achievement by ID
func (s *AchievementService) GetAchievementByID(ctx context.Context, refID string) (*model.AchievementResponse, error) {
	ref, err := s.repo.GetAchievementReference(refID)
	if err != nil {
		return nil, fmt.Errorf("achievement reference not found: %w", err)
	}

	achievement, err := s.repo.GetAchievementMongo(ctx, ref.MongoAchievementID)
	if err != nil {
		return nil, fmt.Errorf("achievement not found in MongoDB: %w", err)
	}

	return s.combineAchievementResponse(achievement, ref), nil
}

// Update Achievement (Student - only draft status)
func (s *AchievementService) UpdateAchievement(ctx context.Context, refID, studentID string, req model.UpdateAchievementRequest) (*model.AchievementResponse, error) {
	ref, err := s.repo.GetAchievementReference(refID)
	if err != nil {
		return nil, fmt.Errorf("achievement not found: %w", err)
	}

	// Check ownership
	if ref.StudentID != studentID {
		return nil, errors.New("unauthorized: you can only update your own achievements")
	}

	// Check status - only draft can be updated
	if ref.Status != "draft" {
		return nil, errors.New("can only update achievements in draft status")
	}

	// Update in MongoDB
	update := &model.Achievement{
		Title:       req.Title,
		Description: req.Description,
		Details:     req.Details,
		Tags:        req.Tags,
		Points:      req.Points,
	}

	err = s.repo.UpdateAchievementMongo(ctx, ref.MongoAchievementID, update)
	if err != nil {
		return nil, fmt.Errorf("failed to update achievement: %w", err)
	}

	return s.GetAchievementByID(ctx, refID)
}

// FR-004: Submit for Verification
func (s *AchievementService) SubmitForVerification(ctx context.Context, refID, studentID string) (*model.AchievementResponse, error) {
	ref, err := s.repo.GetAchievementReference(refID)
	if err != nil {
		return nil, fmt.Errorf("achievement not found: %w", err)
	}

	// Check ownership
	if ref.StudentID != studentID {
		return nil, errors.New("unauthorized: you can only submit your own achievements")
	}

	// Check status - only draft can be submitted
	if ref.Status != "draft" {
		return nil, errors.New("can only submit achievements in draft status")
	}

	// Update status to submitted
	now := time.Now()
	err = s.repo.UpdateAchievementStatus(refID, "submitted", &now)
	if err != nil {
		return nil, fmt.Errorf("failed to submit achievement: %w", err)
	}

	return s.GetAchievementByID(ctx, refID)
}

// FR-005: Delete Achievement (Student - only draft)
func (s *AchievementService) DeleteAchievement(ctx context.Context, refID, studentID string) error {
	ref, err := s.repo.GetAchievementReference(refID)
	if err != nil {
		return fmt.Errorf("achievement not found: %w", err)
	}

	// Check ownership
	if ref.StudentID != studentID {
		return errors.New("unauthorized: you can only delete your own achievements")
	}

	// Check status - only draft can be deleted
	if ref.Status != "draft" {
		return errors.New("can only delete achievements in draft status")
	}

	// Soft delete in MongoDB
	err = s.repo.SoftDeleteAchievementMongo(ctx, ref.MongoAchievementID)
	if err != nil {
		return fmt.Errorf("failed to delete achievement: %w", err)
	}

	// Update status in PostgreSQL
	err = s.repo.UpdateAchievementStatus(refID, "deleted", nil)
	if err != nil {
		return fmt.Errorf("failed to update reference status: %w", err)
	}

	return nil
}

// FR-007: Verify Achievement (Lecturer)
func (s *AchievementService) VerifyAchievement(ctx context.Context, refID, lecturerID string, req model.VerifyAchievementRequest) (*model.AchievementResponse, error) {
	ref, err := s.repo.GetAchievementReference(refID)
	if err != nil {
		return nil, fmt.Errorf("achievement not found: %w", err)
	}

	// Check status - only submitted can be verified/rejected
	if ref.Status != "submitted" {
		return nil, errors.New("can only verify/reject achievements in submitted status")
	}

	if req.Action == "verify" {
		err = s.repo.VerifyAchievement(refID, lecturerID)
		if err != nil {
			return nil, fmt.Errorf("failed to verify achievement: %w", err)
		}
	} else if req.Action == "reject" {
		if req.Note == nil || *req.Note == "" {
			return nil, errors.New("rejection note is required")
		}
		err = s.repo.RejectAchievement(refID, *req.Note)
		if err != nil {
			return nil, fmt.Errorf("failed to reject achievement: %w", err)
		}
	} else {
		return nil, errors.New("invalid action: must be 'verify' or 'reject'")
	}

	return s.GetAchievementByID(ctx, refID)
}

// Get Achievements by Student ID
func (s *AchievementService) GetAchievementsByStudentID(ctx context.Context, studentID string) ([]model.AchievementResponse, error) {
	refs, err := s.repo.GetAchievementsByStudentID(studentID)
	if err != nil {
		return nil, err
	}

	return s.combineMultipleAchievements(ctx, refs)
}

// FR-006: Get Achievements for Lecturer's Advisees
func (s *AchievementService) GetAchievementsForAdvisees(ctx context.Context, studentIDs []string) ([]model.AchievementResponse, error) {
	refs, err := s.repo.GetAchievementReferencesByStudentIDs(studentIDs)
	if err != nil {
		return nil, err
	}

	return s.combineMultipleAchievements(ctx, refs)
}

// FR-010: Get All Achievements (Admin)
func (s *AchievementService) GetAllAchievements(ctx context.Context) ([]model.AchievementResponse, error) {
	refs, err := s.repo.GetAllAchievementReferences()
	if err != nil {
		return nil, err
	}

	return s.combineMultipleAchievements(ctx, refs)
}

// Upload Attachment
func (s *AchievementService) UploadAttachment(ctx context.Context, refID, studentID string, req model.UploadAttachmentRequest) error {
	ref, err := s.repo.GetAchievementReference(refID)
	if err != nil {
		return fmt.Errorf("achievement not found: %w", err)
	}

	// Check ownership
	if ref.StudentID != studentID {
		return errors.New("unauthorized: you can only upload attachments to your own achievements")
	}

	attachment := model.Attachment{
		FileName: req.FileName,
		FileURL:  req.FileURL,
		FileType: req.FileType,
	}

	return s.repo.AddAttachmentMongo(ctx, ref.MongoAchievementID, attachment)
}

// Helper functions
func (s *AchievementService) combineAchievementResponse(achievement *model.Achievement, ref *model.AchievementReference) *model.AchievementResponse {
	return &model.AchievementResponse{
		Achievement:   *achievement,
		Status:        ref.Status,
		SubmittedAt:   ref.SubmittedAt,
		VerifiedAt:    ref.VerifiedAt,
		VerifiedBy:    ref.VerifiedBy,
		RejectionNote: ref.RejectionNote,
	}
}

func (s *AchievementService) combineMultipleAchievements(ctx context.Context, refs []model.AchievementReference) ([]model.AchievementResponse, error) {
	var results []model.AchievementResponse

	for _, ref := range refs {
		achievement, err := s.repo.GetAchievementMongo(ctx, ref.MongoAchievementID)
		if err != nil {
			// Skip if MongoDB document not found
			continue
		}
		results = append(results, *s.combineAchievementResponse(achievement, &ref))
	}

	return results, nil
}
