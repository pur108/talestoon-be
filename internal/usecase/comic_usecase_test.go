package usecase

import (
	"errors"
	"testing"

	"github.com/google/uuid"
	"github.com/pur108/talestoon-be/internal/domain/entity"
	"github.com/pur108/talestoon-be/internal/domain/exception"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockComicRepository implements repository.ComicRepository for testing
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

func (m *MockComicRepository) CreateSeason(season *entity.Season) error {
	args := m.Called(season)
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

func (m *MockComicRepository) GetSeasonByComicID(comicID uuid.UUID, seasonNumber int) (*entity.Season, error) {
	args := m.Called(comicID, seasonNumber)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.Season), args.Error(1)
}

func (m *MockComicRepository) ListComics() ([]entity.Comic, error) {
	args := m.Called()
	return args.Get(0).([]entity.Comic), args.Error(1)
}

func (m *MockComicRepository) ListComicsByCreatorID(creatorID uuid.UUID) ([]entity.Comic, error) {
	args := m.Called(creatorID)
	return args.Get(0).([]entity.Comic), args.Error(1)
}

func (m *MockComicRepository) ListComicsByAuthor(author string) ([]entity.Comic, error) {
	args := m.Called(author)
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

func TestComicUsecase_CreateComic(t *testing.T) {
	creatorID := uuid.New()

	type fields struct {
		mockComicRepo *MockComicRepository
		mockUserRepo  *MockUserRepository
	}
	type args struct {
		input CreateComicInput
	}
	tests := []struct {
		name    string
		prepare func(f *fields)
		args    args
		wantErr bool
	}{
		{
			name: "Success_And_Promote_User_To_Creator",
			prepare: func(f *fields) {
				user := &entity.User{
					ID:   creatorID,
					Role: entity.RoleUser,
				}
				f.mockComicRepo.On("CreateComic", mock.AnythingOfType("*entity.Comic")).Return(nil)
				f.mockUserRepo.On("FindByID", creatorID).Return(user, nil)
				f.mockUserRepo.On("Update", mock.MatchedBy(func(u *entity.User) bool {
					return u.ID == creatorID && u.Role == entity.RoleCreator
				})).Return(nil)
			},
			args: args{
				input: CreateComicInput{
					CreatorID: creatorID,
					Title:     entity.MultilingualText{En: "Title", Th: "Title TH"},
					Status:    entity.ComicPublished,
					Tags:      []entity.MultilingualText{{En: "Action", Th: "Action TH"}},
				},
			},
			wantErr: false,
		},
		{
			name: "Failure_RepoError",
			prepare: func(f *fields) {
				f.mockComicRepo.On("CreateComic", mock.AnythingOfType("*entity.Comic")).Return(errors.New("db error"))
			},
			args: args{
				input: CreateComicInput{
					CreatorID: uuid.New(),
				},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := &fields{
				mockComicRepo: new(MockComicRepository),
				mockUserRepo:  new(MockUserRepository),
			}
			if tt.prepare != nil {
				tt.prepare(f)
			}
			u := NewComicUsecase(f.mockComicRepo, f.mockUserRepo)
			got, err := u.CreateComic(tt.args.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("CreateComic() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr {
				assert.NotNil(t, got)
				assert.Equal(t, tt.args.input.CreatorID, got.CreatorID)
			}
			f.mockComicRepo.AssertExpectations(t)
			f.mockUserRepo.AssertExpectations(t)
		})
	}
}

func TestComicUsecase_UpdateComic(t *testing.T) {
	comicID := uuid.New()
	creatorID := uuid.New()
	otherID := uuid.New()

	type fields struct {
		mockComicRepo *MockComicRepository
		mockUserRepo  *MockUserRepository
	}
	type args struct {
		id        uuid.UUID
		creatorID uuid.UUID
		input     UpdateComicInput
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
				existingComic := &entity.Comic{
					ID:        comicID,
					CreatorID: creatorID,
					Title:     entity.MultilingualText{En: "Old", Th: "Old"},
				}
				f.mockComicRepo.On("GetComicByID", comicID).Return(existingComic, nil)
				f.mockComicRepo.On("UpdateComic", mock.MatchedBy(func(c *entity.Comic) bool {
					return c.ID == comicID && c.Title.En == "New"
				})).Return(nil)
			},
			args: args{
				id:        comicID,
				creatorID: creatorID,
				input: UpdateComicInput{
					Title: entity.MultilingualText{En: "New", Th: "New"},
				},
			},
			wantErr: false,
		},
		{
			name: "Unauthorized",
			prepare: func(f *fields) {
				existingComic := &entity.Comic{
					ID:        comicID,
					CreatorID: creatorID,
				}
				f.mockComicRepo.On("GetComicByID", comicID).Return(existingComic, nil)
			},
			args: args{
				id:        comicID,
				creatorID: otherID,
				input:     UpdateComicInput{},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := &fields{
				mockComicRepo: new(MockComicRepository),
				mockUserRepo:  new(MockUserRepository),
			}
			if tt.prepare != nil {
				tt.prepare(f)
			}
			u := NewComicUsecase(f.mockComicRepo, f.mockUserRepo)
			got, err := u.UpdateComic(tt.args.id, tt.args.creatorID, tt.args.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("UpdateComic() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.wantErr {
				if tt.name == "Unauthorized" {
					assert.ErrorIs(t, err, exception.ErrUnauthorized)
				}
			} else {
				assert.NotNil(t, got)
				assert.Equal(t, "New", got.Title.En)
			}
			f.mockComicRepo.AssertExpectations(t)
			f.mockUserRepo.AssertExpectations(t)
		})
	}
}

func TestComicUsecase_DeleteComic(t *testing.T) {
	comicID := uuid.New()
	creatorID := uuid.New()

	type fields struct {
		mockComicRepo *MockComicRepository
		mockUserRepo  *MockUserRepository
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
				existingComic := &entity.Comic{ID: comicID, CreatorID: creatorID}
				f.mockComicRepo.On("GetComicByID", comicID).Return(existingComic, nil)
				f.mockComicRepo.On("DeleteComic", comicID).Return(nil)
			},
			args: args{
				id:        comicID,
				creatorID: creatorID,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := &fields{
				mockComicRepo: new(MockComicRepository),
				mockUserRepo:  new(MockUserRepository),
			}
			if tt.prepare != nil {
				tt.prepare(f)
			}
			u := NewComicUsecase(f.mockComicRepo, f.mockUserRepo)
			err := u.DeleteComic(tt.args.id, tt.args.creatorID)
			if (err != nil) != tt.wantErr {
				t.Errorf("DeleteComic() error = %v, wantErr %v", err, tt.wantErr)
			}
			f.mockComicRepo.AssertExpectations(t)
			f.mockUserRepo.AssertExpectations(t)
		})
	}
}

func TestComicUsecase_ListMyComics(t *testing.T) {
	creatorID := uuid.New()

	type fields struct {
		mockComicRepo *MockComicRepository
		mockUserRepo  *MockUserRepository
	}
	type args struct {
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
				user := &entity.User{ID: creatorID, Username: "creator"}
				comics := []entity.Comic{{ID: uuid.New(), CreatorID: creatorID}}
				f.mockUserRepo.On("FindByID", creatorID).Return(user, nil)
				f.mockComicRepo.On("ListComicsByCreatorID", creatorID).Return(comics, nil)
			},
			args: args{
				creatorID: creatorID,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := &fields{
				mockComicRepo: new(MockComicRepository),
				mockUserRepo:  new(MockUserRepository),
			}
			if tt.prepare != nil {
				tt.prepare(f)
			}
			u := NewComicUsecase(f.mockComicRepo, f.mockUserRepo)
			got, err := u.ListMyComics(tt.args.creatorID)
			if (err != nil) != tt.wantErr {
				t.Errorf("ListMyComics() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr {
				assert.Len(t, got, 1)
			}
			f.mockComicRepo.AssertExpectations(t)
			f.mockUserRepo.AssertExpectations(t)
		})
	}
}

func TestComicUsecase_CreateChapter(t *testing.T) {
	comicID := uuid.New()
	creatorID := uuid.New()

	type fields struct {
		mockComicRepo *MockComicRepository
		mockUserRepo  *MockUserRepository
	}
	type args struct {
		comicID   uuid.UUID
		creatorID uuid.UUID
		input     CreateChapterInput
	}
	tests := []struct {
		name    string
		prepare func(f *fields)
		args    args
		wantErr bool
	}{
		{
			name: "Success_NewSeason",
			prepare: func(f *fields) {
				comic := &entity.Comic{ID: comicID, CreatorID: creatorID}
				f.mockComicRepo.On("GetComicByID", comicID).Return(comic, nil)
				f.mockComicRepo.On("GetSeasonByComicID", comicID, 1).Return(nil, errors.New("not found"))
				f.mockComicRepo.On("CreateSeason", mock.AnythingOfType("*entity.Season")).Return(nil)
				f.mockComicRepo.On("CreateChapter", mock.MatchedBy(func(c *entity.Chapter) bool {
					return c.ChapterNumber == 1 && len(c.Images) == 1
				})).Return(nil)
			},
			args: args{
				comicID:   comicID,
				creatorID: creatorID,
				input: CreateChapterInput{
					ChapterNumber: 1,
					Title:         "Ch1",
					ImageURLs:     []string{"http://img.com/1.jpg"},
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := &fields{
				mockComicRepo: new(MockComicRepository),
				mockUserRepo:  new(MockUserRepository),
			}
			if tt.prepare != nil {
				tt.prepare(f)
			}
			u := NewComicUsecase(f.mockComicRepo, f.mockUserRepo)
			got, err := u.CreateChapter(tt.args.comicID, tt.args.creatorID, tt.args.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("CreateChapter() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr {
				assert.NotNil(t, got)
				assert.Equal(t, 1, got.ChapterNumber)
			}
			f.mockComicRepo.AssertExpectations(t)
			f.mockUserRepo.AssertExpectations(t)
		})
	}
}
