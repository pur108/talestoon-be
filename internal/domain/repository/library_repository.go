package repository

import (
	"github.com/google/uuid"
	"github.com/pur108/talestoon-be/internal/domain/entity"
)

type LibraryRepository interface {
	// Library Entry
	AddToLibrary(entry *entity.LibraryEntry) error
	RemoveFromLibrary(userID, comicID uuid.UUID) error
	IsInLibrary(userID, comicID uuid.UUID) (bool, error)
	GetUserLibrary(userID uuid.UUID) ([]entity.LibraryEntry, error)

	// Folder Management
	CreateFolder(folder *entity.LibraryFolder) error
	UpdateFolder(folder *entity.LibraryFolder) error
	DeleteFolder(folderID uuid.UUID) error
	GetFolderByID(folderID uuid.UUID) (*entity.LibraryFolder, error)
	GetFolderBySlug(slug string) (*entity.LibraryFolder, error)
	GetUserFolders(userID uuid.UUID) ([]entity.LibraryFolder, error)

	// Folder Items
	AddToFolder(item *entity.LibraryFolderItem) error
	RemoveFromFolder(folderID, comicID uuid.UUID) error
	GetFolderItems(folderID uuid.UUID) ([]entity.LibraryFolderItem, error)
}
