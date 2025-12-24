package repository

import (
	"log"

	"github.com/google/uuid"
	"github.com/pur108/webteen-be/internal/domain/entity"
	"github.com/pur108/webteen-be/internal/domain/repository"
	"gorm.io/gorm"
)

type comicRepository struct {
	db *gorm.DB
}

func NewComicRepository(db *gorm.DB) repository.ComicRepository {
	return &comicRepository{db}
}

func (r *comicRepository) CreateComic(comic *entity.Comic) error {
	return r.db.Create(comic).Error
}

func (r *comicRepository) CreateChapter(chapter *entity.Chapter) error {
	return r.db.Create(chapter).Error
}

func (r *comicRepository) GetComicByID(id uuid.UUID) (*entity.Comic, error) {
	var comic entity.Comic
	// Preload Translations, Tags, and Chapters (flat)
	// Also preload Chapter Translations
	err := r.db.
		Preload("Translations").
		Preload("Tags.Translations").
		Preload("Chapters", func(db *gorm.DB) *gorm.DB {
			return db.Order("chapter_number ASC")
		}).
		Preload("Chapters.Translations").
		Preload("Chapters.Images").
		First(&comic, id).Error

	if err != nil {
		return nil, err
	}
	return &comic, nil
}

func (r *comicRepository) GetChapterByID(id uuid.UUID) (*entity.Chapter, error) {
	var chapter entity.Chapter
	err := r.db.
		Preload("Translations").
		Preload("Images", func(db *gorm.DB) *gorm.DB {
			return db.Order("chapter_images.\"order\" ASC")
		}).First(&chapter, id).Error
	if err != nil {
		return nil, err
	}

	log.Println("chapter: ", chapter)

	return &chapter, nil
}

func (r *comicRepository) ListComics(tags []string) ([]entity.Comic, error) {
	var comics []entity.Comic
	query := r.db.Preload("Translations").Preload("Tags.Translations").Where("status = ?", entity.ComicPublished)

	if len(tags) > 0 {
		query = query.Joins("JOIN comic_tags ON comic_tags.comic_id = comics.id").
			Joins("JOIN tags ON tags.id = comic_tags.tag_id").
			Where("tags.slug IN ?", tags).
			Group("comics.id")
	}

	err := query.Order("updated_at desc").Limit(20).Find(&comics).Error
	if err != nil {
		return nil, err
	}
	return comics, nil
}

func (r *comicRepository) ListComicsByStatus(status entity.ComicStatus) ([]entity.Comic, error) {
	var comics []entity.Comic
	err := r.db.Preload("Translations").Preload("Tags.Translations").Where("status = ?", status).Order("updated_at desc").Find(&comics).Error
	if err != nil {
		return nil, err
	}
	return comics, nil
}

func (r *comicRepository) ListComicsByCreatorID(creatorID uuid.UUID) ([]entity.Comic, error) {
	var comics []entity.Comic
	err := r.db.Preload("Translations").Preload("Tags.Translations").Where("creator_id = ?", creatorID).Order("updated_at desc").Find(&comics).Error
	if err != nil {
		return nil, err
	}
	return comics, nil
}

func (r *comicRepository) ListComicsByAuthor(author string) ([]entity.Comic, error) {
	var comics []entity.Comic
	err := r.db.Preload("Translations").Preload("Tags.Translations").Where("author = ?", author).Order("updated_at desc").Find(&comics).Error
	if err != nil {
		return nil, err
	}
	return comics, nil
}

func (r *comicRepository) UpdateComic(comic *entity.Comic) error {
	return r.db.Session(&gorm.Session{FullSaveAssociations: true}).Save(comic).Error
}

func (r *comicRepository) DeleteComic(id uuid.UUID) error {
	return r.db.Delete(&entity.Comic{}, id).Error
}

func (r *comicRepository) ListTags(filterType string) ([]entity.Tag, error) {
	var tags []entity.Tag
	query := r.db.Preload("Translations")

	if filterType != "" {
		query = query.Where("type = ?", filterType)
	}

	err := query.Find(&tags).Error
	if err != nil {
		return nil, err
	}
	return tags, nil
}
