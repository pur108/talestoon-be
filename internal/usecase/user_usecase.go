package usecase

import (
	"errors"

	"github.com/google/uuid"
	"github.com/pur108/talestoon-be/internal/domain/entity"
	"github.com/pur108/talestoon-be/internal/domain/repository"
)

type UserUsecase interface {
	GetProfile(id uuid.UUID) (*entity.User, error)
	BecomeCreator(id uuid.UUID) error
}

type userUsecase struct {
	userRepo repository.UserRepository
}

func NewUserUsecase(userRepo repository.UserRepository) UserUsecase {
	return &userUsecase{userRepo}
}

func (u *userUsecase) GetProfile(id uuid.UUID) (*entity.User, error) {
	return u.userRepo.FindByID(id)
}

func (u *userUsecase) BecomeCreator(id uuid.UUID) error {
	user, err := u.userRepo.FindByID(id)
	if err != nil {
		return err
	}

	if user.Role == entity.RoleCreator || user.Role == entity.RoleAdmin {
		return errors.New("user is already a creator or admin")
	}

	user.Role = entity.RoleCreator
	return u.userRepo.Update(user)
}
