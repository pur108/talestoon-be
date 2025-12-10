package usecase

import (
	"errors"

	"github.com/google/uuid"
	"github.com/pur108/talestoon-be/internal/domain"
)

type UserUsecase interface {
	GetProfile(id uuid.UUID) (*domain.User, error)
	BecomeCreator(id uuid.UUID) error
}

type userUsecase struct {
	userRepo domain.UserRepository
}

func NewUserUsecase(userRepo domain.UserRepository) UserUsecase {
	return &userUsecase{userRepo}
}

func (u *userUsecase) GetProfile(id uuid.UUID) (*domain.User, error) {
	return u.userRepo.FindByID(id)
}

func (u *userUsecase) BecomeCreator(id uuid.UUID) error {
	user, err := u.userRepo.FindByID(id)
	if err != nil {
		return err
	}

	if user.Role == domain.RoleCreator || user.Role == domain.RoleAdmin {
		return errors.New("user is already a creator or admin")
	}

	user.Role = domain.RoleCreator
	return u.userRepo.Update(user)
}
