package entity

import (
	"time"

	"github.com/google/uuid"
)

type ComicStatus string
type ComicSerializationStatus string
type ChapterStatus string

const (
	ComicDraft     ComicStatus = "draft"
	ComicPending   ComicStatus = "pending_review"
	ComicPublished ComicStatus = "published"
	ComicRejected  ComicStatus = "rejected"

	ComicOngoing   ComicSerializationStatus = "ongoing"
	ComicHiatus    ComicSerializationStatus = "hiatus"
	ComicCompleted ComicSerializationStatus = "completed"

	VisibilityPublic   = "public"
	VisibilityPrivate  = "private"
	VisibilityUnlisted = "unlisted"

	ChapterDraft     ChapterStatus = "draft"
	ChapterPublished ChapterStatus = "published"
	ChapterScheduled ChapterStatus = "scheduled"
)

type Comic struct {
	ID                  uuid.UUID                `gorm:"type:uuid;primary_key;" json:"id"`
	CreatorID           uuid.UUID                `gorm:"type:uuid;not null;index" json:"creator_id"`
	Author              string                   `gorm:"index" json:"author"`
	Tags                []Tag                    `gorm:"many2many:comic_tags;" json:"tags"`
	CoverImageURL       string                   `json:"cover_image_url"`
	BannerImageURL      string                   `json:"banner_image_url"`
	Status              ComicStatus              `gorm:"default:'draft'" json:"status"`
	SerializationStatus ComicSerializationStatus `gorm:"default:'ongoing'" json:"serialization_status"`
	Visibility          string                   `gorm:"default:'public'" json:"visibility"`
	NSFW                bool                     `gorm:"default:false" json:"nsfw"`
	SchedulePublishAt   *time.Time               `json:"schedule_publish_at"`
	ApprovedAt          *time.Time               `json:"approved_at"`
	RejectionReason     string                   `json:"rejection_reason"`
	CreatedAt           time.Time                `json:"created_at"`
	UpdatedAt           time.Time                `gorm:"index" json:"updated_at"`
	Translations        []ComicTranslation       `json:"translations"`
	Chapters            []Chapter                `json:"chapters,omitempty"`
}

type ComicTranslation struct {
	ID               uuid.UUID `gorm:"type:uuid;primary_key;" json:"id"`
	ComicID          uuid.UUID `gorm:"type:uuid;not null;index;uniqueIndex:idx_comic_lang" json:"comic_id"`
	LanguageCode     string    `gorm:"not null;index;uniqueIndex:idx_comic_lang" json:"language_code"`
	Title            string    `gorm:"not null" json:"title"`
	Synopsis         string    `gorm:"type:text" json:"synopsis"`
	AlternativeTitle string    `json:"alternative_title"`
	CreatedAt        time.Time `json:"created_at"`
	UpdatedAt        time.Time `json:"updated_at"`
}

type Chapter struct {
	ID            uuid.UUID            `gorm:"type:uuid;primary_key;" json:"id"`
	ComicID       uuid.UUID            `gorm:"type:uuid;not null;index" json:"comic_id"`
	ChapterNumber int                  `gorm:"not null" json:"chapter_number"`
	Status        ChapterStatus        `gorm:"default:'draft'" json:"status"`
	ThumbnailURL  string               `json:"thumbnail_url"`
	PublishedAt   *time.Time           `json:"published_at"`
	Images        []ChapterImage       `json:"images,omitempty"`
	Translations  []ChapterTranslation `json:"translations"`
}

type ChapterTranslation struct {
	ID           uuid.UUID `gorm:"type:uuid;primary_key;" json:"id"`
	ChapterID    uuid.UUID `gorm:"type:uuid;not null;index;uniqueIndex:idx_chapter_lang" json:"chapter_id"`
	LanguageCode string    `gorm:"not null;index;uniqueIndex:idx_chapter_lang" json:"language_code"`
	Title        string    `json:"title"`
}

type ChapterImage struct {
	ID        uuid.UUID `gorm:"type:uuid;primary_key;" json:"id"`
	ChapterID uuid.UUID `gorm:"type:uuid;not null;index" json:"chapter_id"`
	ImageURL  string    `gorm:"not null" json:"image_url"`
	Order     int       `gorm:"not null" json:"order"`
}
