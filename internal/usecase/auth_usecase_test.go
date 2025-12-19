package usecase_test

import (
	"errors"
	"testing"

	"github.com/google/uuid"
	"github.com/pur108/talestoon-be/internal/domain/entity"
	"github.com/pur108/talestoon-be/internal/usecase"
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
	type fields struct {
		mockRepo *MockUserRepository
	}
	type args struct {
		username string
		email    string
		password string
		role     entity.UserRole
	}
	tests := []struct {
		name    string
		prepare func(f *fields)
		args    args
		wantErr bool
	}{
		{
			name: "Success",
			prepare: func(f *fields) {
				f.mockRepo.On("FindByEmailOrUsername", "test@example.com").Return(nil, errors.New("not found"))
				f.mockRepo.On("FindByEmailOrUsername", "testuser").Return(nil, errors.New("not found"))
				f.mockRepo.On("Create", mock.AnythingOfType("*entity.User")).Return(nil)
			},
			args: args{
				username: "testuser",
				email:    "test@example.com",
				password: "password123",
				role:     entity.RoleUser,
			},
			wantErr: false,
		},
		{
			name: "UserAlreadyExists",
			prepare: func(f *fields) {
				existingUser := &entity.User{Username: "existinguser", Email: "existing@example.com"}
				f.mockRepo.On("FindByEmailOrUsername", "existing@example.com").Return(existingUser, nil)
			},
			args: args{
				username: "existinguser",
				email:    "existing@example.com",
				password: "password123",
				role:     entity.RoleUser,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := &fields{mockRepo: new(MockUserRepository)}
			if tt.prepare != nil {
				tt.prepare(f)
			}
			u := usecase.NewAuthUsecase(f.mockRepo)
			got, err := u.SignUp(tt.args.username, tt.args.email, tt.args.password, tt.args.role)
			if (err != nil) != tt.wantErr {
				t.Errorf("SignUp() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr {
				assert.NotNil(t, got)
				assert.Equal(t, tt.args.username, got.Username)
			}
			f.mockRepo.AssertExpectations(t)
		})
	}
}

func TestAuthUsecase_Login(t *testing.T) {
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("password123"), bcrypt.DefaultCost)
	user := &entity.User{
		ID:           uuid.New(),
		Username:     "testuser",
		Email:        "test@example.com",
		PasswordHash: string(hashedPassword),
		Role:         entity.RoleUser,
	}

	type fields struct {
		mockRepo *MockUserRepository
	}
	type args struct {
		identifier string
		password   string
	}
	tests := []struct {
		name    string
		prepare func(f *fields)
		args    args
		wantErr bool
	}{
		{
			name: "Success_By_Username",
			prepare: func(f *fields) {
				f.mockRepo.On("FindByEmailOrUsername", "testuser").Return(user, nil)
			},
			args: args{
				identifier: "testuser",
				password:   "password123",
			},
			wantErr: false,
		},
		{
			name: "Success_By_Email",
			prepare: func(f *fields) {
				f.mockRepo.On("FindByEmailOrUsername", "test@example.com").Return(user, nil)
			},
			args: args{
				identifier: "test@example.com",
				password:   "password123",
			},
			wantErr: false,
		},
		{
			name: "InvalidCredentials_UserNotFound",
			prepare: func(f *fields) {
				f.mockRepo.On("FindByEmailOrUsername", "unknown").Return(nil, errors.New("not found"))
			},
			args: args{
				identifier: "unknown",
				password:   "password123",
			},
			wantErr: true,
		},
		{
			name: "InvalidCredentials_WrongPassword",
			prepare: func(f *fields) {
				f.mockRepo.On("FindByEmailOrUsername", "testuser").Return(user, nil)
			},
			args: args{
				identifier: "testuser",
				password:   "wrongpassword",
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := &fields{mockRepo: new(MockUserRepository)}
			if tt.prepare != nil {
				tt.prepare(f)
			}
			u := usecase.NewAuthUsecase(f.mockRepo)
			token, loggedInUser, err := u.Login(tt.args.identifier, tt.args.password)
			if (err != nil) != tt.wantErr {
				t.Errorf("Login() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr {
				assert.NotEmpty(t, token)
				assert.NotNil(t, loggedInUser)
				assert.Equal(t, user.ID, loggedInUser.ID)
			}
			f.mockRepo.AssertExpectations(t)
		})
	}
}
