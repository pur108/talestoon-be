package repository

import (
	"github.com/google/uuid"
	"github.com/pur108/webteen-be/internal/domain/entity"
	"github.com/pur108/webteen-be/internal/domain/repository"
	"gorm.io/gorm"
)

type libraryRepository struct {
	db *gorm.DB
}

func NewLibraryRepository(db *gorm.DB) repository.LibraryRepository {
	return &libraryRepository{db}
}

func (r *libraryRepository) GetDefaultFolder(userID uuid.UUID) (*entity.LibraryFolder, error) {
	var folder entity.LibraryFolder
	err := r.db.Preload("Items").Preload("Items.Comic").
		Where("user_id = ? AND is_default = ?", userID, true).First(&folder).Error
	if err != nil {
		return nil, err
	}
	return &folder, nil
}

func (r *libraryRepository) CreateFolder(folder *entity.LibraryFolder) error {
	return r.db.Create(folder).Error
}

func (r *libraryRepository) GetFolder(id uuid.UUID) (*entity.LibraryFolder, error) {
	var folder entity.LibraryFolder
	err := r.db.Preload("Items").Preload("Items.Comic").First(&folder, id).Error
	if err != nil {
		return nil, err
	}
	return &folder, nil
}

func (r *libraryRepository) GetFolderBySlug(userID uuid.UUID, slug string) (*entity.LibraryFolder, error) {
	var folder entity.LibraryFolder
	err := r.db.Preload("Items").Preload("Items.Comic").
		Where("user_id = ? AND slug = ?", userID, slug).First(&folder).Error
	if err != nil {
		return nil, err
	}
	return &folder, nil
}

func (r *libraryRepository) DeleteFolder(id uuid.UUID) error {
	return r.db.Delete(&entity.LibraryFolder{}, id).Error
}

func (r *libraryRepository) ListFolders(userID uuid.UUID) ([]entity.LibraryFolder, error) {
	var folders []entity.LibraryFolder
	err := r.db.Where("user_id = ?", userID).Order("is_default DESC, created_at ASC").Find(&folders).Error
	return folders, err
}

func (r *libraryRepository) AddItem(item *entity.LibraryFolderItem) error {
	return r.db.Create(item).Error
}

func (r *libraryRepository) RemoveItem(folderID, comicID uuid.UUID) error {
	return r.db.Where("folder_id = ? AND comic_id = ?", folderID, comicID).Delete(&entity.LibraryFolderItem{}).Error
}

func (r *libraryRepository) GetItem(folderID, comicID uuid.UUID) (*entity.LibraryFolderItem, error) {
	var item entity.LibraryFolderItem
	err := r.db.Where("folder_id = ? AND comic_id = ?", folderID, comicID).First(&item).Error
	if err != nil {
		return nil, err
	}
	return &item, nil
}
