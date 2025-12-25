package usecase_test

import (
	"errors"
	"testing"

	"github.com/google/uuid"
	"github.com/pur108/webteen-be/internal/domain/entity"
	"github.com/pur108/webteen-be/internal/usecase"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockComicRepository struct {
	mock.Mock
}

func (m *MockComicRepository) CreateComic(comic *entity.Comic) error {
	args := m.Called(comic)
	return args.Error(0)
}

func (m *MockComicRepository) CreateChapter(chapter *entity.Chapter) error {
	args := m.Called(chapter)
	return args.Error(0)
}

func (m *MockComicRepository) GetComicByID(id uuid.UUID) (*entity.Comic, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.Comic), args.Error(1)
}

func (m *MockComicRepository) GetChapterByID(id uuid.UUID) (*entity.Chapter, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.Chapter), args.Error(1)
}

func (m *MockComicRepository) ListComics(tags []string) ([]entity.Comic, error) {
	args := m.Called(tags)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]entity.Comic), args.Error(1)
}

func (m *MockComicRepository) ListComicsByStatus(status entity.ComicStatus) ([]entity.Comic, error) {
	args := m.Called(status)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]entity.Comic), args.Error(1)
}

func (m *MockComicRepository) ListComicsByCreatorID(creatorID uuid.UUID) ([]entity.Comic, error) {
	args := m.Called(creatorID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]entity.Comic), args.Error(1)
}

func (m *MockComicRepository) ListComicsByAuthor(author string) ([]entity.Comic, error) {
	args := m.Called(author)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]entity.Comic), args.Error(1)
}

func (m *MockComicRepository) UpdateComic(comic *entity.Comic) error {
	args := m.Called(comic)
	return args.Error(0)
}

func (m *MockComicRepository) DeleteComic(id uuid.UUID) error {
	args := m.Called(id)
	return args.Error(0)
}

func (m *MockComicRepository) ListTags(filterType string) ([]entity.Tag, error) {
	args := m.Called(filterType)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]entity.Tag), args.Error(1)
}

type MockStorageRepository struct {
	mock.Mock
}

func (m *MockStorageRepository) UploadFile(bucketName string, filePath string, data []byte, contentType string) (string, error) {
	args := m.Called(bucketName, filePath, data, contentType)
	return args.String(0), args.Error(1)
}

func (m *MockStorageRepository) MoveFile(bucketName string, srcPath string, destPath string) error {
	args := m.Called(bucketName, srcPath, destPath)
	return args.Error(0)
}

func TestComicUsecase_CreateComic(t *testing.T) {
	userUUID := uuid.New()
	type fields struct {
		comicRepo   *MockComicRepository
		userRepo    *MockUserRepository
		storageRepo *MockStorageRepository
	}
	type args struct {
		input usecase.CreateComicInput
	}
	tests := []struct {
		name    string
		prepare func(f *fields)
		args    args
		want    *entity.Comic
		wantErr bool
	}{
		{
			name: "Success",
			prepare: func(f *fields) {
				f.comicRepo.On("CreateComic", mock.AnythingOfType("*entity.Comic")).Return(nil)
				f.userRepo.On("FindByID", userUUID).Return(&entity.User{ID: userUUID, Role: entity.RoleUser}, nil)
				f.userRepo.On("Update", mock.AnythingOfType("*entity.User")).Return(nil)
			},
			args: args{
				input: usecase.CreateComicInput{
					CreatorID: userUUID,
					Translations: []usecase.ComicTranslationInput{
						{LanguageCode: "en", Title: "Title En", Description: "Desc En"},
						{LanguageCode: "th", Title: "Title Th", Description: "Desc Th"},
					},
					Status: entity.ComicDraft,
				},
			},
			wantErr: false,
		},
		{
			name: "Error_CreateComic_Failed",
			prepare: func(f *fields) {
				f.comicRepo.On("CreateComic", mock.AnythingOfType("*entity.Comic")).Return(errors.New("db error"))
			},
			args: args{
				input: usecase.CreateComicInput{
					CreatorID: userUUID,
					Translations: []usecase.ComicTranslationInput{
						{LanguageCode: "en", Title: "Title En"},
					},
				},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := &fields{
				comicRepo:   new(MockComicRepository),
				userRepo:    new(MockUserRepository),
				storageRepo: new(MockStorageRepository),
			}
			if tt.prepare != nil {
				tt.prepare(f)
			}
			u := usecase.NewComicUsecase(f.comicRepo, f.userRepo, f.storageRepo)
			got, err := u.CreateComic(tt.args.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("CreateComic() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr {
				assert.NotNil(t, got)
				assert.NotEmpty(t, got.Translations)
				assert.Equal(t, tt.args.input.Translations[0].Title, got.Translations[0].Title)
			}
			f.comicRepo.AssertExpectations(t)
			f.userRepo.AssertExpectations(t)
			f.storageRepo.AssertExpectations(t)
		})
	}
}

func TestComicUsecase_GetComic(t *testing.T) {
	comicID := uuid.New()
	comic := &entity.Comic{
		ID:           comicID,
		Translations: []entity.ComicTranslation{{Title: "Test Comic", LanguageCode: "en"}},
	}

	type fields struct {
		comicRepo *MockComicRepository
	}
	type args struct {
		id uuid.UUID
	}
	tests := []struct {
		name    string
		prepare func(f *fields)
		args    args
		want    *entity.Comic
		wantErr bool
	}{
		{
			name: "Success",
			prepare: func(f *fields) {
				f.comicRepo.On("GetComicByID", comicID).Return(comic, nil)
			},
			args:    args{id: comicID},
			want:    comic,
			wantErr: false,
		},
		{
			name: "NotFound",
			prepare: func(f *fields) {
				f.comicRepo.On("GetComicByID", comicID).Return(nil, errors.New("not found"))
			},
			args:    args{id: comicID},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := &fields{comicRepo: new(MockComicRepository)}
			if tt.prepare != nil {
				tt.prepare(f)
			}
			u := usecase.NewComicUsecase(f.comicRepo, nil, nil)
			got, err := u.GetComic(tt.args.id)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetComic() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr {
				assert.Equal(t, tt.want, got)
			}
			f.comicRepo.AssertExpectations(t)
		})
	}
}

func TestComicUsecase_ListComics(t *testing.T) {
	comics := []entity.Comic{
		{ID: uuid.New(), Translations: []entity.ComicTranslation{{Title: "Comic 1"}}},
		{ID: uuid.New(), Translations: []entity.ComicTranslation{{Title: "Comic 2"}}},
	}

	type fields struct {
		comicRepo *MockComicRepository
	}
	tests := []struct {
		name    string
		prepare func(f *fields)
		want    []entity.Comic
		wantErr bool
	}{
		{
			name: "Success",
			prepare: func(f *fields) {
				f.comicRepo.On("ListComics", []string(nil)).Return(comics, nil)
			},
			want:    comics,
			wantErr: false,
		},
		{
			name: "Error",
			prepare: func(f *fields) {
				f.comicRepo.On("ListComics", []string(nil)).Return(nil, errors.New("db error"))
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := &fields{comicRepo: new(MockComicRepository)}
			if tt.prepare != nil {
				tt.prepare(f)
			}
			u := usecase.NewComicUsecase(f.comicRepo, nil, nil)
			got, err := u.ListComics(nil)
			if (err != nil) != tt.wantErr {
				t.Errorf("ListComics() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr {
				assert.Equal(t, tt.want, got)
			}
			f.comicRepo.AssertExpectations(t)
		})
	}
}

func TestComicUsecase_UpdateComic(t *testing.T) {
	comicID := uuid.New()
	creatorID := uuid.New()
	comic := &entity.Comic{ID: comicID, CreatorID: creatorID, Translations: []entity.ComicTranslation{{Title: "Old Title", LanguageCode: "en"}}}

	type fields struct {
		comicRepo *MockComicRepository
	}
	type args struct {
		id        uuid.UUID
		creatorID uuid.UUID
		input     usecase.UpdateComicInput
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
				f.comicRepo.On("GetComicByID", comicID).Return(comic, nil)
				f.comicRepo.On("UpdateComic", mock.AnythingOfType("*entity.Comic")).Return(nil)
			},
			args: args{
				id:        comicID,
				creatorID: creatorID,
				input: usecase.UpdateComicInput{
					Translations: []usecase.ComicTranslationInput{
						{LanguageCode: "en", Title: "New Title"},
					},
				},
			},
			wantErr: false,
		},
		{
			name: "Unauthorized",
			prepare: func(f *fields) {
				f.comicRepo.On("GetComicByID", comicID).Return(comic, nil)
			},
			args: args{
				id:        comicID,
				creatorID: uuid.New(), // Different creator ID
				input:     usecase.UpdateComicInput{},
			},
			wantErr: true,
		},
		{
			name: "NotFound",
			prepare: func(f *fields) {
				f.comicRepo.On("GetComicByID", comicID).Return(nil, errors.New("not found"))
			},
			args: args{
				id:        comicID,
				creatorID: creatorID,
				input:     usecase.UpdateComicInput{},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := &fields{comicRepo: new(MockComicRepository)}
			if tt.prepare != nil {
				tt.prepare(f)
			}
			u := usecase.NewComicUsecase(f.comicRepo, nil, nil)
			_, err := u.UpdateComic(tt.args.id, tt.args.creatorID, tt.args.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("UpdateComic() error = %v, wantErr %v", err, tt.wantErr)
			}
			f.comicRepo.AssertExpectations(t)
		})
	}
}

func TestComicUsecase_DeleteComic(t *testing.T) {
	comicID := uuid.New()
	creatorID := uuid.New()
	comic := &entity.Comic{ID: comicID, CreatorID: creatorID}

	type fields struct {
		comicRepo *MockComicRepository
	}
	type args struct {
		id        uuid.UUID
		creatorID uuid.UUID
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
				f.comicRepo.On("GetComicByID", comicID).Return(comic, nil)
				f.comicRepo.On("DeleteComic", comicID).Return(nil)
			},
			args: args{
				id:        comicID,
				creatorID: creatorID,
			},
			wantErr: false,
		},
		{
			name: "Unauthorized",
			prepare: func(f *fields) {
				f.comicRepo.On("GetComicByID", comicID).Return(comic, nil)
			},
			args: args{
				id:        comicID,
				creatorID: uuid.New(),
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := &fields{comicRepo: new(MockComicRepository)}
			if tt.prepare != nil {
				tt.prepare(f)
			}
			u := usecase.NewComicUsecase(f.comicRepo, nil, nil)
			err := u.DeleteComic(tt.args.id, tt.args.creatorID)
			if (err != nil) != tt.wantErr {
				t.Errorf("DeleteComic() error = %v, wantErr %v", err, tt.wantErr)
			}
			f.comicRepo.AssertExpectations(t)
		})
	}
}

func TestComicUsecase_ApproveComic(t *testing.T) {
	comicID := uuid.New()

	type fields struct {
		comicRepo   *MockComicRepository
		storageRepo *MockStorageRepository
	}
	tests := []struct {
		name    string
		prepare func(f *fields)
		wantErr bool
	}{
		{
			name: "Success",
			prepare: func(f *fields) {
				comic := &entity.Comic{
					ID:             comicID,
					Status:         entity.ComicPending,
					CoverImageURL:  "https://example.com/media/drafts/cover.jpg",
					BannerImageURL: "https://example.com/media/drafts/banner.jpg",
					Chapters:       []entity.Chapter{},
				}
				f.comicRepo.On("GetComicByID", comicID).Return(comic, nil)
				f.storageRepo.On("MoveFile", "media", "drafts/cover.jpg", "public/cover.jpg").Return(nil)
				f.storageRepo.On("MoveFile", "media", "drafts/banner.jpg", "public/banner.jpg").Return(nil)
				f.comicRepo.On("UpdateComic", mock.AnythingOfType("*entity.Comic")).Return(nil)
			},
			wantErr: false,
		},
		{
			name: "Error_MoveFile_Failed",
			prepare: func(f *fields) {
				comic := &entity.Comic{
					ID:             comicID,
					Status:         entity.ComicPending,
					CoverImageURL:  "https://example.com/media/drafts/cover.jpg",
					BannerImageURL: "https://example.com/media/drafts/banner.jpg",
					Chapters:       []entity.Chapter{},
				}
				f.comicRepo.On("GetComicByID", comicID).Return(comic, nil)
				f.storageRepo.On("MoveFile", "media", "drafts/cover.jpg", "public/cover.jpg").Return(errors.New("move failed"))
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := &fields{
				comicRepo:   new(MockComicRepository),
				storageRepo: new(MockStorageRepository),
			}
			if tt.prepare != nil {
				tt.prepare(f)
			}
			u := usecase.NewComicUsecase(f.comicRepo, nil, f.storageRepo)
			err := u.ApproveComic(comicID)
			if (err != nil) != tt.wantErr {
				t.Errorf("ApproveComic() error = %v, wantErr %v", err, tt.wantErr)
			}
			f.comicRepo.AssertExpectations(t)
			f.storageRepo.AssertExpectations(t)
		})
	}
}
