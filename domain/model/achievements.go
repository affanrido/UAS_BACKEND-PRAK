package model

import (
	"time"

	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Achievement - Collection achievements (MongoDB)
type Achievement struct {
	ID             primitive.ObjectID `json:"_id" bson:"_id,omitempty"`
	StudentID      uuid.UUID          `json:"studentId" bson:"studentId"`
	AchievementType string            `json:"achievementType" bson:"achievementType"` // 'academic', 'competition', 'organization', 'publication', 'certification', 'other'
	Title          string             `json:"title" bson:"title"`
	Description    string             `json:"description" bson:"description"`
	Details        AchievementDetails `json:"details" bson:"details"`
	CustomFields   map[string]interface{} `json:"customFields,omitempty" bson:"customFields,omitempty"`
	Attachments    []Attachment       `json:"attachments" bson:"attachments"`
	Tags           []string           `json:"tags" bson:"tags"`
	Points         float64            `json:"points" bson:"points"`
	IsDeleted      bool               `json:"isDeleted" bson:"isDeleted"`           // Soft delete flag
	DeletedAt      *time.Time         `json:"deletedAt,omitempty" bson:"deletedAt,omitempty"` // Soft delete timestamp
	CreatedAt      time.Time          `json:"createdAt" bson:"createdAt"`
	UpdatedAt      time.Time          `json:"updatedAt" bson:"updatedAt"`
}

// AchievementDetails - Field dinamis berdasarkan tipe prestasi
type AchievementDetails struct {
	// Untuk competition
	CompetitionName  *string  `json:"competitionName,omitempty" bson:"competitionName,omitempty"`
	CompetitionLevel *string  `json:"competitionLevel,omitempty" bson:"competitionLevel,omitempty"` // 'international', 'national', 'regional', 'local'
	Rank             *float64 `json:"rank,omitempty" bson:"rank,omitempty"`
	MedalType        *string  `json:"medalType,omitempty" bson:"medalType,omitempty"`

	// Untuk publication
	PublicationType *string  `json:"publicationType,omitempty" bson:"publicationType,omitempty"` // 'journal', 'conference', 'book'
	PublicationTitle *string `json:"publicationTitle,omitempty" bson:"publicationTitle,omitempty"`
	Authors          []string `json:"authors,omitempty" bson:"authors,omitempty"`
	Publisher        *string  `json:"publisher,omitempty" bson:"publisher,omitempty"`
	ISSN             *string  `json:"issn,omitempty" bson:"issn,omitempty"`

	// Untuk organization
	OrganizationName *string `json:"organizationName,omitempty" bson:"organizationName,omitempty"`
	Position         *string `json:"position,omitempty" bson:"position,omitempty"`
	Period           *Period `json:"period,omitempty" bson:"period,omitempty"`

	// Untuk certification
	CertificationName   *string     `json:"certificationName,omitempty" bson:"certificationName,omitempty"`
	IssuedBy            *string     `json:"issuedBy,omitempty" bson:"issuedBy,omitempty"`
	CertificationNumber *string     `json:"certificationNumber,omitempty" bson:"certificationNumber,omitempty"`
	ValidUntil          *time.Time  `json:"validUntil,omitempty" bson:"validUntil,omitempty"`

	// Field umum yang bisa ada
	EventDate *time.Time `json:"eventDate,omitempty" bson:"eventDate,omitempty"`
	Location  *string    `json:"location,omitempty" bson:"location,omitempty"`
	Organizer *string    `json:"organizer,omitempty" bson:"organizer,omitempty"`
	Score     *float64   `json:"score,omitempty" bson:"score,omitempty"`
}

// Period - Untuk organization period
type Period struct {
	Start time.Time `json:"start" bson:"start"`
	End   time.Time `json:"end" bson:"end"`
}

// Attachment - Untuk attachments array
type Attachment struct {
	FileName   string    `json:"fileName" bson:"fileName"`
	FileURL    string    `json:"fileUrl" bson:"fileUrl"`
	FileType   string    `json:"fileType" bson:"fileType"`
	UploadedAt time.Time `json:"uploadedAt" bson:"uploadedAt"`
}
