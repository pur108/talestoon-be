package entity

import (
	"time"

	"github.com/google/uuid"
)

type Tag struct {
	ID           uuid.UUID        `gorm:"type:uuid;primary_key;" json:"id"`
	Slug         string           `gorm:"uniqueIndex;not null" json:"slug"`
	Type         string           `gorm:"index;not null;default:'genre'" json:"type"`
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
