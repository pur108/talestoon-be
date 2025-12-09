package domain

import (
	"time"

	"github.com/google/uuid"
)

type UserRole string

const (
	RoleUser    UserRole = "user"
	RoleCreator UserRole = "creator"
	RoleAdmin   UserRole = "admin"
)

type User struct {
	ID           uuid.UUID `gorm:"type:uuid;primary_key;" json:"id"`
	Username     string    `gorm:"unique;not null" json:"username"`
	Email        string    `gorm:"unique;not null" json:"email"`
	PasswordHash string    `gorm:"not null" json:"-"`
	Role         UserRole  `gorm:"default:'user'" json:"role"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

type UserRepository interface {
	Create(user *User) error
	Update(user *User) error
	FindByEmail(email string) (*User, error)
	FindByEmailOrUsername(identifier string) (*User, error)
	FindByID(id uuid.UUID) (*User, error)
}
