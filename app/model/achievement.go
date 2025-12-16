package model

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// MongoDB Achievement Model
type Achievement struct {
	ID              primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	StudentID       string             `bson:"student_id" json:"student_id"`
	AchievementType string             `bson:"achievementType" json:"achievementType"` // academic, competition, organization, publication, certification, other
	Title           string             `bson:"title" json:"title"`
	Description     string             `bson:"description" json:"description"`
	Details         AchievementDetails `bson:"details" json:"details"`
	Attachments     []Attachment       `bson:"attachments" json:"attachments"`
	Tags            []string           `bson:"tags" json:"tags"`
	Points          int                `bson:"points" json:"points"`
	IsDeleted       bool               `bson:"isDeleted" json:"-"`
	CreatedAt       time.Time          `bson:"createdAt" json:"createdAt"`
	UpdatedAt       time.Time          `bson:"updatedAt" json:"updatedAt"`
}

// Dynamic achievement details based on type
type AchievementDetails struct {
	// Competition fields
	CompetitionName  string    `bson:"competitionName,omitempty" json:"competitionName,omitempty"`
	CompetitionLevel string    `bson:"competitionLevel,omitempty" json:"competitionLevel,omitempty"` // international, national, regional, local
	Rank             int       `bson:"rank,omitempty" json:"rank,omitempty"`
	MedalType        string    `bson:"medalType,omitempty" json:"medalType,omitempty"`
	
	// Publication fields
	PublicationType string   `bson:"publicationType,omitempty" json:"publicationType,omitempty"` // journal, conference, book
	PublicationTitle string  `bson:"publicationTitle,omitempty" json:"publicationTitle,omitempty"`
	Authors          []string `bson:"authors,omitempty" json:"authors,omitempty"`
	Publisher        string   `bson:"publisher,omitempty" json:"publisher,omitempty"`
	ISSN             string   `bson:"issn,omitempty" json:"issn,omitempty"`
	
	// Organization fields
	OrganizationName string        `bson:"organizationName,omitempty" json:"organizationName,omitempty"`
	Position         string        `bson:"position,omitempty" json:"position,omitempty"`
	Period           *PeriodDetail `bson:"period,omitempty" json:"period,omitempty"`
	
	// Certification fields
	CertificationName   string     `bson:"certificationName,omitempty" json:"certificationName,omitempty"`
	IssuedBy            string     `bson:"issuedBy,omitempty" json:"issuedBy,omitempty"`
	CertificationNumber string     `bson:"certificationNumber,omitempty" json:"certificationNumber,omitempty"`
	ValidUntil          *time.Time `bson:"validUntil,omitempty" json:"validUntil,omitempty"`
	
	// Common fields
	EventDate    *time.Time             `bson:"eventDate,omitempty" json:"eventDate,omitempty"`
	Location     string                 `bson:"location,omitempty" json:"location,omitempty"`
	Organizer    string                 `bson:"organizer,omitempty" json:"organizer,omitempty"`
	Score        float64                `bson:"score,omitempty" json:"score,omitempty"`
	CustomFields map[string]interface{} `bson:"customFields,omitempty" json:"customFields,omitempty"`
}

type PeriodDetail struct {
	Start time.Time `bson:"start" json:"start"`
	End   time.Time `bson:"end" json:"end"`
}

type Attachment struct {
	FileName   string    `bson:"fileName" json:"fileName"`
	FileURL    string    `bson:"fileUrl" json:"fileUrl"`
	FileType   string    `bson:"fileType" json:"fileType"`
	UploadedAt time.Time `bson:"uploadedAt" json:"uploadedAt"`
}

// PostgreSQL Achievement Reference Model
type AchievementReference struct {
	ID                 string     `json:"id"`
	StudentID          string     `json:"student_id"`
	MongoAchievementID string     `json:"mongoAchievementId"`
	Status             string     `json:"status"` // draft, submitted, verified, rejected
	SubmittedAt        *time.Time `json:"submittedAt"`
	VerifiedAt         *time.Time `json:"verifiedAt"`
	VerifiedBy         *string    `json:"verifiedBy"`
	RejectionNote      *string    `json:"rejectionNote"`
	CreatedAt          time.Time  `json:"createdAt"`
	UpdatedAt          time.Time  `json:"updatedAt"`
}

// Combined Achievement Response (MongoDB + PostgreSQL)
type AchievementResponse struct {
	Achievement
	Status        string     `json:"status"`
	SubmittedAt   *time.Time `json:"submittedAt,omitempty"`
	VerifiedAt    *time.Time `json:"verifiedAt,omitempty"`
	VerifiedBy    *string    `json:"verifiedBy,omitempty"`
	RejectionNote *string    `json:"rejectionNote,omitempty"`
}

// Request DTOs
type CreateAchievementRequest struct {
	AchievementType string             `json:"achievementType" validate:"required"`
	Title           string             `json:"title" validate:"required"`
	Description     string             `json:"description"`
	Details         AchievementDetails `json:"details"`
	Tags            []string           `json:"tags"`
	Points          int                `json:"points"`
}

type UpdateAchievementRequest struct {
	Title       string             `json:"title"`
	Description string             `json:"description"`
	Details     AchievementDetails `json:"details"`
	Tags        []string           `json:"tags"`
	Points      int                `json:"points"`
}

type VerifyAchievementRequest struct {
	Action string  `json:"action" validate:"required"` // verify or reject
	Note   *string `json:"note"`                       // required for reject
}

type UploadAttachmentRequest struct {
	FileName string `json:"fileName" validate:"required"`
	FileURL  string `json:"fileUrl" validate:"required"`
	FileType string `json:"fileType" validate:"required"`
}
