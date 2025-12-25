package entity

import (
	"time"

	"github.com/google/uuid"
)

type LibraryFolder struct {
	ID        uuid.UUID           `gorm:"type:uuid;primary_key;" json:"id"`
	UserID    uuid.UUID           `gorm:"type:uuid;not null;index" json:"user_id"`
	Name      string              `gorm:"not null" json:"name"`
	Slug      string              `gorm:"index" json:"slug"`
	IsPublic  bool                `gorm:"default:false" json:"is_public"`
	IsDefault bool                `gorm:"default:false" json:"is_default"`
	Items     []LibraryFolderItem `gorm:"foreignKey:FolderID;constraint:OnDelete:CASCADE;" json:"items,omitempty"`
	CreatedAt time.Time           `json:"created_at"`
	UpdatedAt time.Time           `json:"updated_at"`
}

type LibraryFolderItem struct {
	ID       uuid.UUID `gorm:"type:uuid;primary_key;" json:"id"`
	FolderID uuid.UUID `gorm:"type:uuid;not null;index" json:"folder_id"`
	ComicID  uuid.UUID `gorm:"type:uuid;not null;index" json:"comic_id"`
	Comic    Comic     `json:"comic,omitempty"`
	Order    int       `gorm:"default:0" json:"order"`
	AddedAt  time.Time `json:"added_at"`
}
