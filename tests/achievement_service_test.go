package service_test

import (
	"context"
	"student-report/app/model"
	"student-report/app/service"
	"student-report/tests/mocks"
	"testing"
)

// Test Create Achievement
func TestAchievementService_Create(t *testing.T) {
	ctx := context.Background()

	tests := []struct {
		name      string
		studentID string
		input     model.CreateAchievementRequest
		wantErr   bool
		errMsg    string
	}{
		{
			name:      "Valid achievement creation",
			studentID: "student-1",
			input: model.CreateAchievementRequest{
				AchievementType: "competition",
				Title:           "First Place in Programming Contest",
				Description:     "Won first place in national programming contest",
				Points:          100,
				Tags:            []string{"programming", "competition"},
			},
			wantErr: false,
		},
		{
			name:      "Invalid - empty achievement type",
			studentID: "student-2",
			input: model.CreateAchievementRequest{
				AchievementType: "",
				Title:           "Test Achievement",
				Description:     "Test description",
				Points:          50,
			},
			wantErr: true,
		},
		{
			name:      "Invalid - empty title",
			studentID: "student-3",
			input: model.CreateAchievementRequest{
				AchievementType: "academic",
				Title:           "",
				Description:     "Test description",
				Points:          50,
			},
			wantErr: true,
		},
		{
			name:      "Valid publication achievement",
			studentID: "student-4",
			input: model.CreateAchievementRequest{
				AchievementType: "publication",
				Title:           "Research Paper Published",
				Description:     "Published paper in international journal",
				Points:          150,
				Tags:            []string{"research", "publication"},
				Details: model.AchievementDetails{
					PublicationType:  "journal",
					PublicationTitle: "Advances in AI",
					Publisher:        "IEEE",
				},
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := mocks.NewMockAchievementRepository()
			svc := service.NewAchievementService(mockRepo)

			result, err := svc.CreateAchievement(ctx, tt.studentID, tt.input)

			if tt.wantErr {
				if err == nil {
					t.Errorf("expected error but got nil")
				}
			} else {
				if err != nil {
					t.Errorf("unexpected error: %v", err)
				}
				if result == nil {
					t.Errorf("expected result but got nil")
				}
				if result != nil {
					if result.StudentID != tt.studentID {
						t.Errorf("expected studentID %s, got %s", tt.studentID, result.StudentID)
					}
					if result.Title != tt.input.Title {
						t.Errorf("expected title %s, got %s", tt.input.Title, result.Title)
					}
					if result.Status != "draft" {
						t.Errorf("expected status draft, got %s", result.Status)
					}
				}
			}
		})
	}
}

// Test Submit for Verification
func TestAchievementService_Submit(t *testing.T) {
	ctx := context.Background()

	tests := []struct {
		name       string
		setupData  *model.CreateAchievementRequest
		studentID  string
		submitID   string
		wantErr    bool
		wantStatus string
	}{
		{
			name: "Successfully submit draft achievement",
			setupData: &model.CreateAchievementRequest{
				AchievementType: "competition",
				Title:           "Test Achievement",
				Description:     "Test description",
				Points:          100,
			},
			studentID:  "student-1",
			wantErr:    false,
			wantStatus: "submitted",
		},
		{
			name: "Cannot submit already submitted achievement",
			setupData: &model.CreateAchievementRequest{
				AchievementType: "academic",
				Title:           "Already Submitted",
				Description:     "This is already submitted",
				Points:          50,
			},
			studentID: "student-2",
			wantErr:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := mocks.NewMockAchievementRepository()
			svc := service.NewAchievementService(mockRepo)

			// Setup: Create achievement
			created, err := svc.CreateAchievement(ctx, tt.studentID, *tt.setupData)
			if err != nil {
				t.Fatalf("failed to setup test data: %v", err)
			}

			refID := created.ID.Hex()

			// For second test case, submit once first
			if tt.name == "Cannot submit already submitted achievement" {
				_, err = svc.SubmitForVerification(ctx, refID, tt.studentID)
				if err != nil {
					t.Fatalf("failed to setup submitted achievement: %v", err)
				}
			}

			// Test: Submit for verification
			result, err := svc.SubmitForVerification(ctx, refID, tt.studentID)

			if tt.wantErr {
				if err == nil {
					t.Errorf("expected error but got nil")
				}
			} else {
				if err != nil {
					t.Errorf("unexpected error: %v", err)
				}
				if result == nil {
					t.Errorf("expected result but got nil")
				}
				if result != nil && result.Status != tt.wantStatus {
					t.Errorf("expected status %s, got %s", tt.wantStatus, result.Status)
				}
			}
		})
	}
}

// Test Verify Achievement
func TestAchievementService_Verify(t *testing.T) {
	ctx := context.Background()

	tests := []struct {
		name        string
		setupData   *model.CreateAchievementRequest
		studentID   string
		lecturerID  string
		verifyReq   model.VerifyAchievementRequest
		wantErr     bool
		wantStatus  string
	}{
		{
			name: "Successfully verify submitted achievement",
			setupData: &model.CreateAchievementRequest{
				AchievementType: "competition",
				Title:           "Programming Contest Winner",
				Description:     "Won programming contest",
				Points:          100,
			},
			studentID:  "student-1",
			lecturerID: "lecturer-1",
			verifyReq: model.VerifyAchievementRequest{
				Action: "verify",
			},
			wantErr:    false,
			wantStatus: "verified",
		},
		{
			name: "Successfully reject with note",
			setupData: &model.CreateAchievementRequest{
				AchievementType: "academic",
				Title:           "Invalid Achievement",
				Description:     "This will be rejected",
				Points:          50,
			},
			studentID:  "student-2",
			lecturerID: "lecturer-1",
			verifyReq: model.VerifyAchievementRequest{
				Action: "reject",
				Note:   stringPtr("Invalid evidence provided"),
			},
			wantErr:    false,
			wantStatus: "rejected",
		},
		{
			name: "Reject without note should fail",
			setupData: &model.CreateAchievementRequest{
				AchievementType: "competition",
				Title:           "Test Achievement",
				Description:     "Test",
				Points:          50,
			},
			studentID:  "student-3",
			lecturerID: "lecturer-1",
			verifyReq: model.VerifyAchievementRequest{
				Action: "reject",
				Note:   nil,
			},
			wantErr: true,
		},
		{
			name: "Invalid action should fail",
			setupData: &model.CreateAchievementRequest{
				AchievementType: "competition",
				Title:           "Test Achievement",
				Description:     "Test",
				Points:          50,
			},
			studentID:  "student-4",
			lecturerID: "lecturer-1",
			verifyReq: model.VerifyAchievementRequest{
				Action: "invalid_action",
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := mocks.NewMockAchievementRepository()
			svc := service.NewAchievementService(mockRepo)

			// Setup: Create and submit achievement
			created, err := svc.CreateAchievement(ctx, tt.studentID, *tt.setupData)
			if err != nil {
				t.Fatalf("failed to create achievement: %v", err)
			}

			refID := created.ID.Hex()

			// Submit for verification
			_, err = svc.SubmitForVerification(ctx, refID, tt.studentID)
			if err != nil {
				t.Fatalf("failed to submit achievement: %v", err)
			}

			// Test: Verify or reject
			result, err := svc.VerifyAchievement(ctx, refID, tt.lecturerID, tt.verifyReq)

			if tt.wantErr {
				if err == nil {
					t.Errorf("expected error but got nil")
				}
			} else {
				if err != nil {
					t.Errorf("unexpected error: %v", err)
				}
				if result == nil {
					t.Errorf("expected result but got nil")
				}
				if result != nil && result.Status != tt.wantStatus {
					t.Errorf("expected status %s, got %s", tt.wantStatus, result.Status)
				}
				if tt.verifyReq.Action == "verify" && result.VerifiedBy == nil {
					t.Errorf("expected verifiedBy to be set")
				}
				if tt.verifyReq.Action == "reject" && result.RejectionNote == nil {
					t.Errorf("expected rejectionNote to be set")
				}
			}
		})
	}
}

// Test Upload Attachment
func TestAchievementService_UploadAttachment(t *testing.T) {
	ctx := context.Background()

	tests := []struct {
		name       string
		setupData  *model.CreateAchievementRequest
		studentID  string
		uploadReq  model.UploadAttachmentRequest
		wantErr    bool
	}{
		{
			name: "Successfully upload attachment",
			setupData: &model.CreateAchievementRequest{
				AchievementType: "competition",
				Title:           "Test Achievement",
				Description:     "Test",
				Points:          100,
			},
			studentID: "student-1",
			uploadReq: model.UploadAttachmentRequest{
				FileName: "certificate.pdf",
				FileURL:  "https://example.com/certificate.pdf",
				FileType: "application/pdf",
			},
			wantErr: false,
		},
		{
			name: "Upload to non-existent achievement",
			setupData: &model.CreateAchievementRequest{
				AchievementType: "competition",
				Title:           "Test Achievement",
				Description:     "Test",
				Points:          100,
			},
			studentID: "student-2",
			uploadReq: model.UploadAttachmentRequest{
				FileName: "test.pdf",
				FileURL:  "https://example.com/test.pdf",
				FileType: "application/pdf",
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := mocks.NewMockAchievementRepository()
			svc := service.NewAchievementService(mockRepo)

			// Setup: Create achievement
			created, err := svc.CreateAchievement(ctx, tt.studentID, *tt.setupData)
			if err != nil {
				t.Fatalf("failed to create achievement: %v", err)
			}

			refID := created.ID.Hex()

			// For error test case, use wrong refID
			if tt.name == "Upload to non-existent achievement" {
				refID = "non-existent-ref-id"
			}

			// Test: Upload attachment
			err = svc.UploadAttachment(ctx, refID, tt.studentID, tt.uploadReq)

			if tt.wantErr {
				if err == nil {
					t.Errorf("expected error but got nil")
				}
			} else {
				if err != nil {
					t.Errorf("unexpected error: %v", err)
				}

				// Verify attachment was added
				result, _ := svc.GetAchievementByID(ctx, refID)
				if result != nil && len(result.Attachments) == 0 {
					t.Errorf("expected attachment to be added")
				}
			}
		})
	}
}

// Test Update Achievement
func TestAchievementService_Update(t *testing.T) {
	ctx := context.Background()

	tests := []struct {
		name       string
		setupData  *model.CreateAchievementRequest
		updateData model.UpdateAchievementRequest
		studentID  string
		wantErr    bool
	}{
		{
			name: "Successfully update draft achievement",
			setupData: &model.CreateAchievementRequest{
				AchievementType: "competition",
				Title:           "Original Title",
				Description:     "Original description",
				Points:          100,
			},
			updateData: model.UpdateAchievementRequest{
				Title:       "Updated Title",
				Description: "Updated description",
				Points:      150,
			},
			studentID: "student-1",
			wantErr:   false,
		},
		{
			name: "Cannot update submitted achievement",
			setupData: &model.CreateAchievementRequest{
				AchievementType: "academic",
				Title:           "Submitted Achievement",
				Description:     "This is submitted",
				Points:          50,
			},
			updateData: model.UpdateAchievementRequest{
				Title:       "Try to Update",
				Description: "Should fail",
				Points:      60,
			},
			studentID: "student-2",
			wantErr:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := mocks.NewMockAchievementRepository()
			svc := service.NewAchievementService(mockRepo)

			// Setup: Create achievement
			created, err := svc.CreateAchievement(ctx, tt.studentID, *tt.setupData)
			if err != nil {
				t.Fatalf("failed to create achievement: %v", err)
			}

			refID := created.ID.Hex()

			// For second test, submit first
			if tt.name == "Cannot update submitted achievement" {
				_, err = svc.SubmitForVerification(ctx, refID, tt.studentID)
				if err != nil {
					t.Fatalf("failed to submit achievement: %v", err)
				}
			}

			// Test: Update achievement
			result, err := svc.UpdateAchievement(ctx, refID, tt.studentID, tt.updateData)

			if tt.wantErr {
				if err == nil {
					t.Errorf("expected error but got nil")
				}
			} else {
				if err != nil {
					t.Errorf("unexpected error: %v", err)
				}
				if result == nil {
					t.Errorf("expected result but got nil")
				}
				if result != nil {
					if result.Title != tt.updateData.Title {
						t.Errorf("expected title %s, got %s", tt.updateData.Title, result.Title)
					}
					if result.Description != tt.updateData.Description {
						t.Errorf("expected description %s, got %s", tt.updateData.Description, result.Description)
					}
				}
			}
		})
	}
}

// Test Delete Achievement
func TestAchievementService_Delete(t *testing.T) {
	ctx := context.Background()

	tests := []struct {
		name      string
		setupData *model.CreateAchievementRequest
		studentID string
		wantErr   bool
	}{
		{
			name: "Successfully delete draft achievement",
			setupData: &model.CreateAchievementRequest{
				AchievementType: "competition",
				Title:           "To Be Deleted",
				Description:     "This will be deleted",
				Points:          100,
			},
			studentID: "student-1",
			wantErr:   false,
		},
		{
			name: "Cannot delete submitted achievement",
			setupData: &model.CreateAchievementRequest{
				AchievementType: "academic",
				Title:           "Submitted Achievement",
				Description:     "Cannot delete",
				Points:          50,
			},
			studentID: "student-2",
			wantErr:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := mocks.NewMockAchievementRepository()
			svc := service.NewAchievementService(mockRepo)

			// Setup: Create achievement
			created, err := svc.CreateAchievement(ctx, tt.studentID, *tt.setupData)
			if err != nil {
				t.Fatalf("failed to create achievement: %v", err)
			}

			refID := created.ID.Hex()

			// For second test, submit first
			if tt.name == "Cannot delete submitted achievement" {
				_, err = svc.SubmitForVerification(ctx, refID, tt.studentID)
				if err != nil {
					t.Fatalf("failed to submit achievement: %v", err)
				}
			}

			// Test: Delete achievement
			err = svc.DeleteAchievement(ctx, refID, tt.studentID)

			if tt.wantErr {
				if err == nil {
					t.Errorf("expected error but got nil")
				}
			} else {
				if err != nil {
					t.Errorf("unexpected error: %v", err)
				}
			}
		})
	}
}

// Test Get Achievement by Student ID
func TestAchievementService_GetByStudentID(t *testing.T) {
	ctx := context.Background()

	tests := []struct {
		name          string
		setupData     []*model.CreateAchievementRequest
		studentID     string
		expectedCount int
		wantErr       bool
	}{
		{
			name: "Get multiple achievements for student",
			setupData: []*model.CreateAchievementRequest{
				{
					AchievementType: "competition",
					Title:           "Achievement 1",
					Description:     "First achievement",
					Points:          100,
				},
				{
					AchievementType: "academic",
					Title:           "Achievement 2",
					Description:     "Second achievement",
					Points:          50,
				},
				{
					AchievementType: "publication",
					Title:           "Achievement 3",
					Description:     "Third achievement",
					Points:          150,
				},
			},
			studentID:     "student-1",
			expectedCount: 3,
			wantErr:       false,
		},
		{
			name:          "Get achievements for student with no achievements",
			setupData:     []*model.CreateAchievementRequest{},
			studentID:     "student-2",
			expectedCount: 0,
			wantErr:       false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := mocks.NewMockAchievementRepository()
			svc := service.NewAchievementService(mockRepo)

			// Setup: Create achievements
			for _, data := range tt.setupData {
				_, err := svc.CreateAchievement(ctx, tt.studentID, *data)
				if err != nil {
					t.Fatalf("failed to create achievement: %v", err)
				}
			}

			// Test: Get achievements by student ID
			results, err := svc.GetAchievementsByStudentID(ctx, tt.studentID)

			if tt.wantErr {
				if err == nil {
					t.Errorf("expected error but got nil")
				}
			} else {
				if err != nil {
					t.Errorf("unexpected error: %v", err)
				}
				if len(results) != tt.expectedCount {
					t.Errorf("expected %d achievements, got %d", tt.expectedCount, len(results))
				}
			}
		})
	}
}

// Helper function
func stringPtr(s string) *string {
	return &s
}