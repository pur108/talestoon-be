package usecase

import (
	"errors"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/pur108/talestoon-be/internal/domain/entity"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestUserUsecase_GetProfile(t *testing.T) {
	userId := uuid.New()
	user := &entity.User{
		ID:        userId,
		Username:  "testuser",
		Email:     "test@example.com",
		Role:      entity.RoleUser,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	type fields struct {
		mockRepo *MockUserRepository
	}
	type args struct {
		id uuid.UUID
	}
	tests := []struct {
		name    string
		prepare func(f *fields)
		args    args
		want    *entity.User
		wantErr bool
	}{
		{
			name: "Success",
			prepare: func(f *fields) {
				f.mockRepo.On("FindByID", userId).Return(user, nil)
			},
			args:    args{id: userId},
			want:    user,
			wantErr: false,
		},
		{
			name: "UserNotFound",
			prepare: func(f *fields) {
				f.mockRepo.On("FindByID", userId).Return(nil, errors.New("not found"))
			},
			args:    args{id: userId},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := &fields{mockRepo: new(MockUserRepository)}
			if tt.prepare != nil {
				tt.prepare(f)
			}

			u := NewUserUsecase(f.mockRepo)
			got, err := u.GetProfile(tt.args.id)

			if (err != nil) != tt.wantErr {
				t.Errorf("GetProfile() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr {
				assert.Equal(t, tt.want, got)
			}
			f.mockRepo.AssertExpectations(t)
		})
	}
}

func TestUserUsecase_BecomeCreator(t *testing.T) {
	userId := uuid.New()

	type fields struct {
		mockRepo *MockUserRepository
	}
	tests := []struct {
		name    string
		prepare func(f *fields)
		id      uuid.UUID
		wantErr bool
	}{
		{
			name: "Success_UserToCreator",
			prepare: func(f *fields) {
				user := &entity.User{ID: userId, Role: entity.RoleUser}
				f.mockRepo.On("FindByID", userId).Return(user, nil)
				f.mockRepo.On("Update", mock.MatchedBy(func(u *entity.User) bool {
					return u.ID == userId && u.Role == entity.RoleCreator
				})).Return(nil)
			},
			id:      userId,
			wantErr: false,
		},
		{
			name: "AlreadyCreator",
			prepare: func(f *fields) {
				user := &entity.User{ID: userId, Role: entity.RoleCreator}
				f.mockRepo.On("FindByID", userId).Return(user, nil)
			},
			id:      userId,
			wantErr: true,
		},
		{
			name: "AlreadyAdmin",
			prepare: func(f *fields) {
				user := &entity.User{ID: userId, Role: entity.RoleAdmin}
				f.mockRepo.On("FindByID", userId).Return(user, nil)
			},
			id:      userId,
			wantErr: true,
		},
		{
			name: "UserNotFound",
			prepare: func(f *fields) {
				f.mockRepo.On("FindByID", userId).Return(nil, errors.New("not found"))
			},
			id:      userId,
			wantErr: true,
		},
		{
			name: "UpdateError",
			prepare: func(f *fields) {
				user := &entity.User{ID: userId, Role: entity.RoleUser}
				f.mockRepo.On("FindByID", userId).Return(user, nil)
				f.mockRepo.On("Update", mock.MatchedBy(func(u *entity.User) bool {
					return u.Role == entity.RoleCreator
				})).Return(errors.New("db error"))
			},
			id:      userId,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := &fields{mockRepo: new(MockUserRepository)}
			if tt.prepare != nil {
				tt.prepare(f)
			}

			u := NewUserUsecase(f.mockRepo)
			err := u.BecomeCreator(tt.id)

			if (err != nil) != tt.wantErr {
				t.Errorf("BecomeCreator() error = %v, wantErr %v", err, tt.wantErr)
			}
			f.mockRepo.AssertExpectations(t)
		})
	}
}
