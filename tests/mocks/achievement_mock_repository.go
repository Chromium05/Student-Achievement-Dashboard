package mocks

import (
	"context"
	"errors"
	"fmt"
	"student-report/app/model"
	"student-report/app/repository"
	"sync"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

var _ repository.IAchievementRepository = (*MockAchievementRepository)(nil)

type MockAchievementRepository struct {
	mu                     sync.Mutex
	mongoAchievements      map[string]*model.Achievement      // mongoID -> Achievement
	achievementReferences  map[string]*model.AchievementReference // refID -> Reference
	mongoIDToRefID         map[string]string                 // mongoID -> refID
	studentAchievements    map[string][]string               // studentID -> []refID
	nextRefID              int
}

func NewMockAchievementRepository() *MockAchievementRepository {
	return &MockAchievementRepository{
		mongoAchievements:     make(map[string]*model.Achievement),
		achievementReferences: make(map[string]*model.AchievementReference),
		mongoIDToRefID:        make(map[string]string),
		studentAchievements:   make(map[string][]string),
		nextRefID:             1,
	}
}

// MongoDB Operations
func (m *MockAchievementRepository) CreateAchievementMongo(ctx context.Context, achievement *model.Achievement) (string, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	achievement.ID = primitive.NewObjectID()
	achievement.CreatedAt = time.Now()
	achievement.UpdatedAt = time.Now()
	achievement.IsDeleted = false

	mongoID := achievement.ID.Hex()
	m.mongoAchievements[mongoID] = achievement

	return mongoID, nil
}

func (m *MockAchievementRepository) GetAchievementMongo(ctx context.Context, mongoID string) (*model.Achievement, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	achievement, exists := m.mongoAchievements[mongoID]
	if !exists || achievement.IsDeleted {
		return nil, errors.New("achievement not found")
	}

	return achievement, nil
}

func (m *MockAchievementRepository) UpdateAchievementMongo(ctx context.Context, mongoID string, update *model.Achievement) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	achievement, exists := m.mongoAchievements[mongoID]
	if !exists || achievement.IsDeleted {
		return errors.New("achievement not found")
	}

	achievement.Title = update.Title
	achievement.Description = update.Description
	achievement.Details = update.Details
	achievement.Tags = update.Tags
	achievement.Points = update.Points
	achievement.UpdatedAt = time.Now()

	return nil
}

func (m *MockAchievementRepository) SoftDeleteAchievementMongo(ctx context.Context, mongoID string) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	achievement, exists := m.mongoAchievements[mongoID]
	if !exists {
		return errors.New("achievement not found")
	}

	achievement.IsDeleted = true
	achievement.UpdatedAt = time.Now()

	return nil
}

func (m *MockAchievementRepository) AddAttachmentMongo(ctx context.Context, mongoID string, attachment model.Attachment) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	achievement, exists := m.mongoAchievements[mongoID]
	if !exists || achievement.IsDeleted {
		return errors.New("achievement not found")
	}

	attachment.UploadedAt = time.Now()
	achievement.Attachments = append(achievement.Attachments, attachment)
	achievement.UpdatedAt = time.Now()

	return nil
}

// PostgreSQL Operations
func (m *MockAchievementRepository) CreateAchievementReference(studentID, mongoID string) (string, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	refID := fmt.Sprintf("ref-%d", m.nextRefID)
	m.nextRefID++

	ref := &model.AchievementReference{
		ID:                 refID,
		StudentID:          studentID,
		MongoAchievementID: mongoID,
		Status:             "draft",
		CreatedAt:          time.Now(),
		UpdatedAt:          time.Now(),
	}

	m.achievementReferences[refID] = ref
	m.mongoIDToRefID[mongoID] = refID
	m.studentAchievements[studentID] = append(m.studentAchievements[studentID], refID)

	return refID, nil
}

func (m *MockAchievementRepository) GetAchievementReference(id string) (*model.AchievementReference, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	ref, exists := m.achievementReferences[id]
	if !exists {
		return nil, errors.New("achievement reference not found")
	}

	return ref, nil
}

func (m *MockAchievementRepository) GetAchievementReferenceByMongoID(mongoID string) (*model.AchievementReference, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	refID, exists := m.mongoIDToRefID[mongoID]
	if !exists {
		return nil, errors.New("achievement reference not found")
	}

	return m.achievementReferences[refID], nil
}

func (m *MockAchievementRepository) UpdateAchievementStatus(id, status string, submittedAt *time.Time) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	ref, exists := m.achievementReferences[id]
	if !exists {
		return errors.New("achievement reference not found")
	}

	ref.Status = status
	if submittedAt != nil {
		ref.SubmittedAt = submittedAt
	}
	ref.UpdatedAt = time.Now()

	return nil
}

func (m *MockAchievementRepository) VerifyAchievement(id, verifiedBy string) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	ref, exists := m.achievementReferences[id]
	if !exists {
		return errors.New("achievement reference not found")
	}

	now := time.Now()
	ref.Status = "verified"
	ref.VerifiedAt = &now
	ref.VerifiedBy = &verifiedBy
	ref.UpdatedAt = time.Now()

	return nil
}

func (m *MockAchievementRepository) RejectAchievement(id, note string) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	ref, exists := m.achievementReferences[id]
	if !exists {
		return errors.New("achievement reference not found")
	}

	ref.Status = "rejected"
	ref.RejectionNote = &note
	ref.UpdatedAt = time.Now()

	return nil
}

func (m *MockAchievementRepository) GetAchievementsByStudentID(studentID string) ([]model.AchievementReference, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	refIDs := m.studentAchievements[studentID]
	var refs []model.AchievementReference

	for _, refID := range refIDs {
		if ref, exists := m.achievementReferences[refID]; exists {
			refs = append(refs, *ref)
		}
	}

	return refs, nil
}

func (m *MockAchievementRepository) GetAchievementReferencesByStudentIDs(studentIDs []string) ([]model.AchievementReference, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	var refs []model.AchievementReference
	for _, studentID := range studentIDs {
		refIDs := m.studentAchievements[studentID]
		for _, refID := range refIDs {
			if ref, exists := m.achievementReferences[refID]; exists {
				refs = append(refs, *ref)
			}
		}
	}

	return refs, nil
}

func (m *MockAchievementRepository) GetAllAchievementReferences() ([]model.AchievementReference, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	var refs []model.AchievementReference
	for _, ref := range m.achievementReferences {
		refs = append(refs, *ref)
	}

	return refs, nil
}

func (m *MockAchievementRepository) GetAchievementsWithFilter(ctx context.Context, filter model.AchievementFilter) ([]model.Achievement, int64, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	var results []model.Achievement
	for _, achievement := range m.mongoAchievements {
		if achievement.IsDeleted {
			continue
		}

		// Apply filters
		if filter.AchievementType != nil && achievement.AchievementType != *filter.AchievementType {
			continue
		}

		if filter.StudentID != nil && achievement.StudentID != *filter.StudentID {
			continue
		}

		results = append(results, *achievement)
	}

	return results, int64(len(results)), nil
}

func (m *MockAchievementRepository) GetAchievementStatistics(ctx context.Context, studentIDs []string) (*model.AchievementStatistics, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	stats := &model.AchievementStatistics{
		ByType:   make(map[string]int),
		ByStatus: make(map[string]int),
	}

	// Count by type
	for _, achievement := range m.mongoAchievements {
		if !achievement.IsDeleted {
			stats.TotalAchievements++
			stats.ByType[achievement.AchievementType]++
			stats.TotalPoints += achievement.Points
		}
	}

	// Count by status
	for _, ref := range m.achievementReferences {
		stats.ByStatus[ref.Status]++
		if ref.Status == "submitted" {
			stats.PendingVerification++
		}
		if ref.Status == "verified" {
			stats.RecentVerified++
		}
	}

	return stats, nil
}

func (m *MockAchievementRepository) GetTopStudents(ctx context.Context, limit int) ([]model.StudentAchievementCount, error) {
	return []model.StudentAchievementCount{}, nil
}

func (m *MockAchievementRepository) GetStudentInfo(studentID string) (string, string) {
	return "Test Student", "123456"
}

func (m *MockAchievementRepository) GetAchievementsByStudentIDs(ctx context.Context, studentIDs []string) ([]model.Achievement, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	var results []model.Achievement
	for _, achievement := range m.mongoAchievements {
		if achievement.IsDeleted {
			continue
		}
		for _, studentID := range studentIDs {
			if achievement.StudentID == studentID {
				results = append(results, *achievement)
				break
			}
		}
	}

	return results, nil
}