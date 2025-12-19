package usecase_test

import (
	"errors"
	"strings"
	"testing"

	"github.com/google/uuid"
	"github.com/pur108/talestoon-be/internal/domain/entity"
	"github.com/pur108/talestoon-be/internal/usecase"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockLibraryRepository struct {
	mock.Mock
}

func (m *MockLibraryRepository) AddToLibrary(entry *entity.LibraryEntry) error {
	args := m.Called(entry)
	return args.Error(0)
}

func (m *MockLibraryRepository) RemoveFromLibrary(userID, comicID uuid.UUID) error {
	args := m.Called(userID, comicID)
	return args.Error(0)
}

func (m *MockLibraryRepository) IsInLibrary(userID, comicID uuid.UUID) (bool, error) {
	args := m.Called(userID, comicID)
	return args.Bool(0), args.Error(1)
}

func (m *MockLibraryRepository) GetUserLibrary(userID uuid.UUID) ([]entity.LibraryEntry, error) {
	args := m.Called(userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]entity.LibraryEntry), args.Error(1)
}

func (m *MockLibraryRepository) CreateFolder(folder *entity.LibraryFolder) error {
	args := m.Called(folder)
	return args.Error(0)
}

func (m *MockLibraryRepository) UpdateFolder(folder *entity.LibraryFolder) error {
	args := m.Called(folder)
	return args.Error(0)
}

func (m *MockLibraryRepository) DeleteFolder(folderID uuid.UUID) error {
	args := m.Called(folderID)
	return args.Error(0)
}

func (m *MockLibraryRepository) GetFolderByID(folderID uuid.UUID) (*entity.LibraryFolder, error) {
	args := m.Called(folderID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.LibraryFolder), args.Error(1)
}

func (m *MockLibraryRepository) GetFolderBySlug(slug string) (*entity.LibraryFolder, error) {
	args := m.Called(slug)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.LibraryFolder), args.Error(1)
}

func (m *MockLibraryRepository) GetUserFolders(userID uuid.UUID) ([]entity.LibraryFolder, error) {
	args := m.Called(userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]entity.LibraryFolder), args.Error(1)
}

func (m *MockLibraryRepository) AddToFolder(item *entity.LibraryFolderItem) error {
	args := m.Called(item)
	return args.Error(0)
}

func (m *MockLibraryRepository) RemoveFromFolder(folderID, comicID uuid.UUID) error {
	args := m.Called(folderID, comicID)
	return args.Error(0)
}

func (m *MockLibraryRepository) GetFolderItems(folderID uuid.UUID) ([]entity.LibraryFolderItem, error) {
	args := m.Called(folderID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]entity.LibraryFolderItem), args.Error(1)
}

func TestLibraryUsecase_AddToLibrary(t *testing.T) {
	type fields struct {
		mockRepo     *MockLibraryRepository
		mockUserRepo *MockUserRepository
		// ComicRepo is null for this test as it's not used in AddToLibrary directly
	}
	type args struct {
		userID  uuid.UUID
		comicID uuid.UUID
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
				f.mockRepo.On("IsInLibrary", mock.Anything, mock.Anything).Return(false, nil)
				f.mockRepo.On("AddToLibrary", mock.AnythingOfType("*entity.LibraryEntry")).Return(nil)
			},
			args: args{
				userID:  uuid.New(),
				comicID: uuid.New(),
			},
			wantErr: false,
		},
		{
			name: "AlreadyInLibrary",
			prepare: func(f *fields) {
				f.mockRepo.On("IsInLibrary", mock.Anything, mock.Anything).Return(true, nil)
			},
			args: args{
				userID:  uuid.New(),
				comicID: uuid.New(),
			},
			wantErr: true,
		},
		{
			name: "RepoError_IsInLibrary",
			prepare: func(f *fields) {
				f.mockRepo.On("IsInLibrary", mock.Anything, mock.Anything).Return(false, errors.New("db error"))
			},
			args: args{
				userID:  uuid.New(),
				comicID: uuid.New(),
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := &fields{
				mockRepo:     new(MockLibraryRepository),
				mockUserRepo: new(MockUserRepository),
			}
			if tt.prepare != nil {
				tt.prepare(f)
			}
			u := usecase.NewLibraryUsecase(f.mockRepo, nil, f.mockUserRepo)
			if err := u.AddToLibrary(tt.args.userID, tt.args.comicID); (err != nil) != tt.wantErr {
				t.Errorf("AddToLibrary() error = %v, wantErr %v", err, tt.wantErr)
			}
			f.mockRepo.AssertExpectations(t)
		})
	}
}

func TestLibraryUsecase_CreateFolder(t *testing.T) {
	type fields struct {
		mockRepo *MockLibraryRepository
	}
	type args struct {
		userID      uuid.UUID
		name        string
		description string
		isPublic    bool
	}
	tests := []struct {
		name    string
		prepare func(f *fields)
		args    args
		want    *entity.LibraryFolder
		wantErr bool
	}{
		{
			name: "Success",
			prepare: func(f *fields) {
				f.mockRepo.On("CreateFolder", mock.MatchedBy(func(folder *entity.LibraryFolder) bool {
					return folder.Name == "My Folder" && strings.Contains(folder.Slug, "my-folder")
				})).Return(nil)
			},
			args: args{
				userID:      uuid.New(),
				name:        "My Folder",
				description: "Desc",
				isPublic:    true,
			},
			wantErr: false,
		},
		{
			name: "RepoError",
			prepare: func(f *fields) {
				f.mockRepo.On("CreateFolder", mock.Anything).Return(errors.New("db error"))
			},
			args:    args{userID: uuid.New(), name: "F"},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := &fields{mockRepo: new(MockLibraryRepository)}
			if tt.prepare != nil {
				tt.prepare(f)
			}
			u := usecase.NewLibraryUsecase(f.mockRepo, nil, nil)
			got, err := u.CreateFolder(tt.args.userID, tt.args.name, tt.args.description, tt.args.isPublic)
			if (err != nil) != tt.wantErr {
				t.Errorf("CreateFolder() error = %v, wantErr %v", err, tt.wantErr)
			}
			if !tt.wantErr {
				assert.NotNil(t, got)
				assert.Equal(t, tt.args.name, got.Name)
			}
			f.mockRepo.AssertExpectations(t)
		})
	}
}

func TestLibraryUsecase_AddToFolder(t *testing.T) {
	userID := uuid.New()
	otherUserID := uuid.New()
	folderID := uuid.New()

	type fields struct {
		mockRepo *MockLibraryRepository
	}
	tests := []struct {
		name    string
		prepare func(f *fields)
		userID  uuid.UUID
		wantErr bool
	}{
		{
			name: "Success",
			prepare: func(f *fields) {
				folder := &entity.LibraryFolder{ID: folderID, UserID: userID}
				f.mockRepo.On("GetFolderByID", folderID).Return(folder, nil)
				f.mockRepo.On("AddToFolder", mock.MatchedBy(func(item *entity.LibraryFolderItem) bool {
					return item.FolderID == folderID
				})).Return(nil)
			},
			userID:  userID,
			wantErr: false,
		},
		{
			name: "FolderNotFound",
			prepare: func(f *fields) {
				f.mockRepo.On("GetFolderByID", folderID).Return(nil, nil)
			},
			userID:  userID,
			wantErr: true,
		},
		{
			name: "Unauthorized",
			prepare: func(f *fields) {
				folder := &entity.LibraryFolder{ID: folderID, UserID: otherUserID}
				f.mockRepo.On("GetFolderByID", folderID).Return(folder, nil)
			},
			userID:  userID,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := &fields{mockRepo: new(MockLibraryRepository)}
			if tt.prepare != nil {
				tt.prepare(f)
			}
			u := usecase.NewLibraryUsecase(f.mockRepo, nil, nil)
			if err := u.AddToFolder(tt.userID, folderID, uuid.New()); (err != nil) != tt.wantErr {
				t.Errorf("AddToFolder() error = %v, wantErr %v", err, tt.wantErr)
			}
			f.mockRepo.AssertExpectations(t)
		})
	}
}

func TestLibraryUsecase_DeleteFolder(t *testing.T) {
	userID := uuid.New()
	folderID := uuid.New()

	type fields struct {
		mockRepo *MockLibraryRepository
	}
	tests := []struct {
		name    string
		prepare func(f *fields)
		userID  uuid.UUID
		wantErr bool
	}{
		{
			name: "Success",
			prepare: func(f *fields) {
				folder := &entity.LibraryFolder{ID: folderID, UserID: userID}
				f.mockRepo.On("GetFolderByID", folderID).Return(folder, nil)
				f.mockRepo.On("DeleteFolder", folderID).Return(nil)
			},
			userID:  userID,
			wantErr: false,
		},
		{
			name: "Unauthorized",
			prepare: func(f *fields) {
				folder := &entity.LibraryFolder{ID: folderID, UserID: uuid.New()}
				f.mockRepo.On("GetFolderByID", folderID).Return(folder, nil)
			},
			userID:  userID,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := &fields{mockRepo: new(MockLibraryRepository)}
			if tt.prepare != nil {
				tt.prepare(f)
			}
			u := usecase.NewLibraryUsecase(f.mockRepo, nil, nil)
			err := u.DeleteFolder(tt.userID, folderID)
			if (err != nil) != tt.wantErr {
				t.Errorf("DeleteFolder() error = %v, wantErr %v", err, tt.wantErr)
			}
			f.mockRepo.AssertExpectations(t)
		})
	}
}
