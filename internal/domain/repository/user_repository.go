package repository

import (
	"github.com/google/uuid"
	"github.com/pur108/talestoon-be/internal/domain/entity"
)

type UserRepository interface {
	Create(user *entity.User) error
	Update(user *entity.User) error
	FindByEmailOrUsername(identifier string) (*entity.User, error)
	FindByID(id uuid.UUID) (*entity.User, error)
}
