package usecase

import (
	"errors"
	"strings"
	"time"

	"github.com/pur108/talestoon-be/internal/domain/entity"
	"github.com/pur108/talestoon-be/internal/domain/repository"

	"github.com/google/uuid"
)

type LibraryUsecase interface {
	AddToLibrary(userID, comicID uuid.UUID) error
	RemoveFromLibrary(userID, comicID uuid.UUID) error
	GetUserLibrary(userID uuid.UUID) ([]entity.LibraryEntry, error)
	CreateFolder(userID uuid.UUID, name string, description string, isPublic bool) (*entity.LibraryFolder, error)
	GetFolder(folderID uuid.UUID) (*entity.LibraryFolder, error)
	GetFolderBySlug(slug string) (*entity.LibraryFolder, error)
	GetUserFolders(userID uuid.UUID) ([]entity.LibraryFolder, error)
	AddToFolder(userID, folderID, comicID uuid.UUID) error
	RemoveFromFolder(userID, folderID, comicID uuid.UUID) error
	DeleteFolder(userID, folderID uuid.UUID) error
}

type libraryUsecase struct {
	libraryRepo repository.LibraryRepository
	comicRepo   repository.ComicRepository
	userRepo    repository.UserRepository
}

func NewLibraryUsecase(libraryRepo repository.LibraryRepository, comicRepo repository.ComicRepository, userRepo repository.UserRepository) LibraryUsecase {
	return &libraryUsecase{
		libraryRepo: libraryRepo,
		comicRepo:   comicRepo,
		userRepo:    userRepo,
	}
}

func (u *libraryUsecase) AddToLibrary(userID, comicID uuid.UUID) error {
	exists, err := u.libraryRepo.IsInLibrary(userID, comicID)
	if err != nil {
		return err
	}
	if exists {
		return errors.New("comic already in library")
	}

	entry := &entity.LibraryEntry{
		ID:        uuid.New(),
		UserID:    userID,
		ComicID:   comicID,
		CreatedAt: time.Now(),
	}

	return u.libraryRepo.AddToLibrary(entry)
}

func (u *libraryUsecase) RemoveFromLibrary(userID, comicID uuid.UUID) error {
	return u.libraryRepo.RemoveFromLibrary(userID, comicID)
}

func (u *libraryUsecase) GetUserLibrary(userID uuid.UUID) ([]entity.LibraryEntry, error) {
	return u.libraryRepo.GetUserLibrary(userID)
}

func (u *libraryUsecase) CreateFolder(userID uuid.UUID, name string, description string, isPublic bool) (*entity.LibraryFolder, error) {
	slug := strings.ToLower(strings.ReplaceAll(name, " ", "-")) + "-" + uuid.New().String()[:8]

	folder := &entity.LibraryFolder{
		ID:          uuid.New(),
		UserID:      userID,
		Name:        name,
		Description: description,
		IsPublic:    isPublic,
		Slug:        slug,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	if err := u.libraryRepo.CreateFolder(folder); err != nil {
		return nil, err
	}

	return folder, nil
}

func (u *libraryUsecase) GetFolder(folderID uuid.UUID) (*entity.LibraryFolder, error) {
	return u.libraryRepo.GetFolderByID(folderID)
}

func (u *libraryUsecase) GetFolderBySlug(slug string) (*entity.LibraryFolder, error) {
	return u.libraryRepo.GetFolderBySlug(slug)
}

func (u *libraryUsecase) GetUserFolders(userID uuid.UUID) ([]entity.LibraryFolder, error) {
	return u.libraryRepo.GetUserFolders(userID)
}

func (u *libraryUsecase) AddToFolder(userID, folderID, comicID uuid.UUID) error {
	folder, err := u.libraryRepo.GetFolderByID(folderID)
	if err != nil {
		return err
	}
	if folder == nil {
		return errors.New("folder not found")
	}
	if folder.UserID != userID {
		return errors.New("unauthorized")
	}

	item := &entity.LibraryFolderItem{
		ID:       uuid.New(),
		FolderID: folderID,
		ComicID:  comicID,
		Order:    len(folder.Items),
		AddedAt:  time.Now(),
	}

	return u.libraryRepo.AddToFolder(item)
}

func (u *libraryUsecase) RemoveFromFolder(userID, folderID, comicID uuid.UUID) error {
	folder, err := u.libraryRepo.GetFolderByID(folderID)
	if err != nil {
		return err
	}
	if folder == nil {
		return errors.New("folder not found")
	}
	if folder.UserID != userID {
		return errors.New("unauthorized")
	}

	return u.libraryRepo.RemoveFromFolder(folderID, comicID)
}

func (u *libraryUsecase) DeleteFolder(userID, folderID uuid.UUID) error {
	folder, err := u.libraryRepo.GetFolderByID(folderID)
	if err != nil {
		return err
	}
	if folder == nil {
		return errors.New("folder not found")
	}
	if folder.UserID != userID {
		return errors.New("unauthorized")
	}

	return u.libraryRepo.DeleteFolder(folderID)
}
