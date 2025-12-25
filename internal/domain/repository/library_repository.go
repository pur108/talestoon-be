package repository

import (
	"github.com/google/uuid"
	"github.com/pur108/webteen-be/internal/domain/entity"
)

type LibraryRepository interface {
	GetDefaultFolder(userID uuid.UUID) (*entity.LibraryFolder, error)
	CreateFolder(folder *entity.LibraryFolder) error
	GetFolder(id uuid.UUID) (*entity.LibraryFolder, error)
	GetFolderBySlug(userID uuid.UUID, slug string) (*entity.LibraryFolder, error)
	DeleteFolder(id uuid.UUID) error
	ListFolders(userID uuid.UUID) ([]entity.LibraryFolder, error)

	AddItem(item *entity.LibraryFolderItem) error
	RemoveItem(folderID, comicID uuid.UUID) error
	GetItem(folderID, comicID uuid.UUID) (*entity.LibraryFolderItem, error)
}
