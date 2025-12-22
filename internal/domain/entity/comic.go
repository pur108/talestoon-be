package entity

import (
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/lib/pq"
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

type MultilingualText struct {
	En string `json:"en"`
	Th string `json:"th"`
}

func (m MultilingualText) Value() (interface{}, error) {
	return json.Marshal(m)
}

func (m *MultilingualText) Scan(value interface{}) error {
	bytes, ok := value.([]byte)
	if !ok {
		return errors.New(fmt.Sprint("Failed to unmarshal JSONB value:", value))
	}

	return json.Unmarshal(bytes, m)
}

type Comic struct {
	ID                  uuid.UUID                `gorm:"type:uuid;primary_key;" json:"id"`
	CreatorID           uuid.UUID                `gorm:"type:uuid;not null;index" json:"creator_id"`
	Title               MultilingualText         `gorm:"type:jsonb;serializer:json" json:"title"`
	Subtitle            MultilingualText         `gorm:"type:jsonb;serializer:json" json:"subtitle"`
	Description         MultilingualText         `gorm:"type:jsonb;serializer:json" json:"description"`
	Author              string                   `gorm:"index" json:"author"`
	Genres              pq.StringArray           `gorm:"type:text[]" json:"genres"`
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
	Seasons             []Season                 `json:"seasons,omitempty"`
}

type Tag struct {
	ID           uuid.UUID        `gorm:"type:uuid;primary_key;" json:"id"`
	Slug         string           `gorm:"uniqueIndex;not null" json:"slug"`
	Translations []TagTranslation `json:"translations"`
	CreatedAt    time.Time        `json:"created_at"`
	UpdatedAt    time.Time        `json:"updated_at"`
}

type TagTranslation struct {
	ID       uuid.UUID `gorm:"type:uuid;primary_key;" json:"id"`
	TagID    uuid.UUID `gorm:"type:uuid;not null;index" json:"tag_id"`
	Language string    `gorm:"not null;index" json:"language"`
	Name     string    `gorm:"not null" json:"name"`
}

type Season struct {
	ID           uuid.UUID `gorm:"type:uuid;primary_key;" json:"id"`
	ComicID      uuid.UUID `gorm:"type:uuid;not null;index" json:"comic_id"`
	SeasonNumber int       `gorm:"not null" json:"season_number"`
	Title        string    `json:"title"`
	Chapters     []Chapter `json:"chapters,omitempty"`
}

type Chapter struct {
	ID            uuid.UUID      `gorm:"type:uuid;primary_key;" json:"id"`
	SeasonID      uuid.UUID      `gorm:"type:uuid;not null;index" json:"season_id"`
	ChapterNumber int            `gorm:"not null" json:"chapter_number"`
	Title         string         `json:"title"`
	Status        ChapterStatus  `gorm:"default:'draft'" json:"status"`
	PublishedAt   *time.Time     `json:"published_at"`
	Images        []ChapterImage `json:"images,omitempty"`
}

type ChapterImage struct {
	ID        uuid.UUID `gorm:"type:uuid;primary_key;" json:"id"`
	ChapterID uuid.UUID `gorm:"type:uuid;not null;index" json:"chapter_id"`
	ImageURL  string    `gorm:"not null" json:"image_url"`
	Order     int       `gorm:"not null" json:"order"`
}
