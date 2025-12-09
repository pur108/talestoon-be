package usecase

import (
	"errors"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/pur108/talestoon-be.git/internal/domain"
	"golang.org/x/crypto/bcrypt"
)

type AuthUsecase interface {
	SignUp(username, email, password string, role domain.UserRole) (*domain.User, error)
	Login(identifier, password string) (string, *domain.User, error)
}

type authUsecase struct {
	userRepo domain.UserRepository
}

func NewAuthUsecase(userRepo domain.UserRepository) AuthUsecase {
	return &authUsecase{userRepo}
}

func (u *authUsecase) SignUp(username, email, password string, role domain.UserRole) (*domain.User, error) {
	existingUser, _ := u.userRepo.FindByEmailOrUsername(email)
	if existingUser != nil {
		return nil, errors.New("email or username already exists")
	}

	existingUserByUsername, _ := u.userRepo.FindByEmailOrUsername(username)
	if existingUserByUsername != nil {
		return nil, errors.New("username already exists")
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	user := &domain.User{
		ID:           uuid.New(),
		Username:     username,
		Email:        email,
		PasswordHash: string(hashedPassword),
		Role:         role,
	}

	if err := u.userRepo.Create(user); err != nil {
		return nil, err
	}

	return user, nil
}

func (u *authUsecase) Login(identifier, password string) (string, *domain.User, error) {
	user, err := u.userRepo.FindByEmailOrUsername(identifier)
	if err != nil {
		return "", nil, errors.New("invalid credentials")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password)); err != nil {
		return "", nil, errors.New("invalid credentials")
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": user.ID.String(),
		"role":    user.Role,
		"exp":     time.Now().Add(time.Hour * 72).Unix(),
	})

	t, err := token.SignedString([]byte(os.Getenv("JWT_SECRET")))
	if err != nil {
		return "", nil, err
	}

	return t, user, nil
}
