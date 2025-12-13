package usecase

import (
	"errors"
	"testing"

	"github.com/google/uuid"
	"github.com/pur108/talestoon-be/internal/domain/entity"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"golang.org/x/crypto/bcrypt"
)

type MockUserRepository struct {
	mock.Mock
}

func (m *MockUserRepository) Create(user *entity.User) error {
	args := m.Called(user)
	return args.Error(0)
}

func (m *MockUserRepository) Update(user *entity.User) error {
	args := m.Called(user)
	return args.Error(0)
}

func (m *MockUserRepository) FindByEmailOrUsername(identifier string) (*entity.User, error) {
	args := m.Called(identifier)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.User), args.Error(1)
}

func (m *MockUserRepository) FindByID(id uuid.UUID) (*entity.User, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.User), args.Error(1)
}

func TestAuthUsecase_SignUp(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		mockRepo := new(MockUserRepository)
		authUsecase := NewAuthUsecase(mockRepo)

		username := "testuser"
		email := "test@example.com"
		password := "password123"
		role := entity.RoleUser

		mockRepo.On("FindByEmailOrUsername", email).Return(nil, errors.New("not found"))
		mockRepo.On("FindByEmailOrUsername", username).Return(nil, errors.New("not found"))
 
		mockRepo.On("Create", mock.AnythingOfType("*entity.User")).Return(nil)

		user, err := authUsecase.SignUp(username, email, password, role)

		assert.NoError(t, err)
		assert.NotNil(t, user)
		assert.Equal(t, username, user.Username)
		assert.Equal(t, email, user.Email)
		assert.Equal(t, role, user.Role)
		assert.NotEmpty(t, user.PasswordHash)

		mockRepo.AssertExpectations(t)
	})

	t.Run("UserAlreadyExists", func(t *testing.T) {
		mockRepo := new(MockUserRepository)
		authUsecase := NewAuthUsecase(mockRepo)

		username := "existinguser"
		email := "existing@example.com"
		password := "password123"
		role := entity.RoleUser

		existingUser := &entity.User{
			Username: username,
			Email:    email,
		}

		mockRepo.On("FindByEmailOrUsername", email).Return(existingUser, nil)

		user, err := authUsecase.SignUp(username, email, password, role)

		assert.Error(t, err)
		assert.Nil(t, user)
		assert.Equal(t, "email or username already exists", err.Error())

		mockRepo.AssertNotCalled(t, "Create", mock.Anything)
	})
}

func TestAuthUsecase_Login(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		mockRepo := new(MockUserRepository)
		authUsecase := NewAuthUsecase(mockRepo)

		username := "testuser"
		password := "password123"
		hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)

		user := &entity.User{
			ID:           uuid.New(),
			Username:     username,
			PasswordHash: string(hashedPassword),
			Role:         entity.RoleUser,
		}

		mockRepo.On("FindByEmailOrUsername", username).Return(user, nil)

		token, loggedInUser, err := authUsecase.Login(username, password)

		assert.NoError(t, err)
		assert.NotNil(t, loggedInUser)
		assert.NotEmpty(t, token)
		assert.Equal(t, user, loggedInUser)

		mockRepo.AssertExpectations(t)
	})

	t.Run("InvalidCredentials_UserNotFound", func(t *testing.T) {
		mockRepo := new(MockUserRepository)
		authUsecase := NewAuthUsecase(mockRepo)

		username := "unknown"
		password := "pass"

		mockRepo.On("FindByEmailOrUsername", username).Return(nil, errors.New("not found"))

		token, loggedInUser, err := authUsecase.Login(username, password)

		assert.Error(t, err)
		assert.Nil(t, loggedInUser)
		assert.Empty(t, token)
		assert.Equal(t, "invalid credentials", err.Error())
	})

	t.Run("InvalidCredentials_WrongPassword", func(t *testing.T) {
		mockRepo := new(MockUserRepository)
		authUsecase := NewAuthUsecase(mockRepo)

		username := "testuser"
		password := "wrongpassword"
		hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("correctpassword"), bcrypt.DefaultCost)

		user := &entity.User{
			Username:     username,
			PasswordHash: string(hashedPassword),
		}

		mockRepo.On("FindByEmailOrUsername", username).Return(user, nil)

		token, loggedInUser, err := authUsecase.Login(username, password)

		assert.Error(t, err)
		assert.Nil(t, loggedInUser)
		assert.Empty(t, token)
		assert.Equal(t, "invalid credentials", err.Error())
	})
}
