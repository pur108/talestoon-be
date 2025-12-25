package usecase

import (
	"time"

	"github.com/google/uuid"
	"github.com/pur108/webteen-be/internal/domain/entity"
	"github.com/pur108/webteen-be/internal/domain/exception"
	"github.com/pur108/webteen-be/internal/domain/repository"
	"github.com/pur108/webteen-be/pkg/utils"
)

type LibraryUsecase interface {
	GetMyLibrary(userID uuid.UUID) (*entity.LibraryFolder, error)
	AddToLibrary(userID, comicID uuid.UUID) error
	RemoveFromLibrary(userID, comicID uuid.UUID) error
	CheckInLibrary(userID, comicID uuid.UUID) (bool, error)
	CreateFolder(userID uuid.UUID, name string) (*entity.LibraryFolder, error)
	ListFolders(userID uuid.UUID) ([]entity.LibraryFolder, error)
	GetFolder(userID, folderID uuid.UUID) (*entity.LibraryFolder, error)
	DeleteFolder(userID, folderID uuid.UUID) error
}

type libraryUsecase struct {
	libraryRepo repository.LibraryRepository
	comicRepo   repository.ComicRepository
}

func NewLibraryUsecase(libraryRepo repository.LibraryRepository, comicRepo repository.ComicRepository) LibraryUsecase {
	return &libraryUsecase{libraryRepo, comicRepo}
}

func (u *libraryUsecase) GetMyLibrary(userID uuid.UUID) (*entity.LibraryFolder, error) {
	folder, err := u.libraryRepo.GetDefaultFolder(userID)
	if err == nil {
		return folder, nil
	}

	newFolder := &entity.LibraryFolder{
		ID:        uuid.New(),
		UserID:    userID,
		Name:      "My Library",
		Slug:      utils.SimpleSlug("my-library"),
		IsDefault: true,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	if err := u.libraryRepo.CreateFolder(newFolder); err != nil {
		return nil, err
	}

	return newFolder, nil
}

func (u *libraryUsecase) AddToLibrary(userID, comicID uuid.UUID) error {
	folder, err := u.GetMyLibrary(userID)
	if err != nil {
		return err
	}

	_, err = u.comicRepo.GetComicByID(comicID)
	if err != nil {
		return exception.ErrNotFound
	}
	existing, _ := u.libraryRepo.GetItem(folder.ID, comicID)
	if existing != nil {
		return nil
	}

	item := &entity.LibraryFolderItem{
		ID:       uuid.New(),
		FolderID: folder.ID,
		ComicID:  comicID,
		AddedAt:  time.Now(),
	}

	return u.libraryRepo.AddItem(item)
}

func (u *libraryUsecase) RemoveFromLibrary(userID, comicID uuid.UUID) error {
	folder, err := u.GetMyLibrary(userID)
	if err != nil {
		return err
	}
	return u.libraryRepo.RemoveItem(folder.ID, comicID)
}

func (u *libraryUsecase) CheckInLibrary(userID, comicID uuid.UUID) (bool, error) {
	folder, err := u.GetMyLibrary(userID)
	if err != nil {
		return false, err
	}
	item, err := u.libraryRepo.GetItem(folder.ID, comicID)
	if err != nil {
		return false, nil
	}
	return item != nil, nil
}

func (u *libraryUsecase) CreateFolder(userID uuid.UUID, name string) (*entity.LibraryFolder, error) {
	folder := &entity.LibraryFolder{
		ID:        uuid.New(),
		UserID:    userID,
		Name:      name,
		Slug:      utils.SimpleSlug(name),
		IsDefault: false,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	if err := u.libraryRepo.CreateFolder(folder); err != nil {
		return nil, err
	}
	return folder, nil
}

func (u *libraryUsecase) ListFolders(userID uuid.UUID) ([]entity.LibraryFolder, error) {
	_, err := u.GetMyLibrary(userID)
	if err != nil {
		return nil, err
	}
	return u.libraryRepo.ListFolders(userID)
}

func (u *libraryUsecase) DeleteFolder(userID, folderID uuid.UUID) error {
	folder, err := u.libraryRepo.GetFolder(folderID)
	if err != nil {
		return err
	}
	if folder.UserID != userID {
		return exception.ErrForbidden
	}
	if folder.IsDefault {
		return exception.ErrForbidden
	}
	return u.libraryRepo.DeleteFolder(folderID)
}

func (u *libraryUsecase) UpdateFolder(userID, folderID uuid.UUID, name string, isPublic bool) (*entity.LibraryFolder, error) {
	folder, err := u.libraryRepo.GetFolder(folderID)
	if err != nil {
		return nil, err
	}
	if folder.UserID != userID {
		return nil, exception.ErrForbidden
	}
	if folder.IsDefault && name != "" && name != folder.Name {
		return nil, exception.ErrForbidden
	}

	if name != "" {
		folder.Name = name
		folder.Slug = utils.SimpleSlug(name)
	}
	folder.IsPublic = isPublic
	folder.UpdatedAt = time.Now()
	return folder, nil
}

func (u *libraryUsecase) GetFolder(userID, folderID uuid.UUID) (*entity.LibraryFolder, error) {
	folder, err := u.libraryRepo.GetFolder(folderID)
	if err != nil {
		return nil, err
	}
	if folder.UserID != userID && !folder.IsPublic {
		return nil, exception.ErrForbidden
	}
	return folder, nil
}
