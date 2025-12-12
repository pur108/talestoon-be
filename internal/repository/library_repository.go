package repository

import (
	"errors"

	"github.com/pur108/talestoon-be/internal/domain/entity"
	"github.com/pur108/talestoon-be/internal/domain/repository"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type LibraryRepository struct {
	db *gorm.DB
}

func NewLibraryRepository(db *gorm.DB) repository.LibraryRepository {
	return &LibraryRepository{db: db}
}

func (r *LibraryRepository) AddToLibrary(entry *entity.LibraryEntry) error {
	return r.db.Create(entry).Error
}

func (r *LibraryRepository) RemoveFromLibrary(userID, comicID uuid.UUID) error {
	return r.db.Where("user_id = ? AND comic_id = ?", userID, comicID).Delete(&entity.LibraryEntry{}).Error
}

func (r *LibraryRepository) IsInLibrary(userID, comicID uuid.UUID) (bool, error) {
	var count int64
	err := r.db.Model(&entity.LibraryEntry{}).Where("user_id = ? AND comic_id = ?", userID, comicID).Count(&count).Error
	return count > 0, err
}

func (r *LibraryRepository) GetUserLibrary(userID uuid.UUID) ([]entity.LibraryEntry, error) {
	var entries []entity.LibraryEntry
	err := r.db.Preload("Comic").Where("user_id = ?", userID).Order("created_at desc").Find(&entries).Error
	return entries, err
}

func (r *LibraryRepository) CreateFolder(folder *entity.LibraryFolder) error {
	return r.db.Create(folder).Error
}

func (r *LibraryRepository) UpdateFolder(folder *entity.LibraryFolder) error {
	return r.db.Save(folder).Error
}

func (r *LibraryRepository) DeleteFolder(folderID uuid.UUID) error {
	return r.db.Delete(&entity.LibraryFolder{}, folderID).Error
}

func (r *LibraryRepository) GetFolderByID(folderID uuid.UUID) (*entity.LibraryFolder, error) {
	var folder entity.LibraryFolder
	err := r.db.Preload("Items.Comic").First(&folder, folderID).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &folder, err
}

func (r *LibraryRepository) GetFolderBySlug(slug string) (*entity.LibraryFolder, error) {
	var folder entity.LibraryFolder
	err := r.db.Preload("Items.Comic").Where("slug = ?", slug).First(&folder).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &folder, err
}

func (r *LibraryRepository) GetUserFolders(userID uuid.UUID) ([]entity.LibraryFolder, error) {
	var folders []entity.LibraryFolder
	err := r.db.Where("user_id = ?", userID).Order("updated_at desc").Find(&folders).Error
	return folders, err
}

func (r *LibraryRepository) AddToFolder(item *entity.LibraryFolderItem) error {
	return r.db.Create(item).Error
}

func (r *LibraryRepository) RemoveFromFolder(folderID, comicID uuid.UUID) error {
	return r.db.Where("folder_id = ? AND comic_id = ?", folderID, comicID).Delete(&entity.LibraryFolderItem{}).Error
}

func (r *LibraryRepository) GetFolderItems(folderID uuid.UUID) ([]entity.LibraryFolderItem, error) {
	var items []entity.LibraryFolderItem
	err := r.db.Preload("Comic").Where("folder_id = ?", folderID).Order("\"order\" asc").Find(&items).Error
	return items, err
}
