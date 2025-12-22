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

func (r *comicRepository) CreateSeason(season *entity.Season) error {
	return r.db.Create(season).Error
}

func (r *comicRepository) GetComicByID(id uuid.UUID) (*entity.Comic, error) {
	var comic entity.Comic
	err := r.db.Preload("Seasons.Chapters.Images").Preload("Seasons.Chapters").Preload("Tags.Translations").First(&comic, id).Error
	if err != nil {
		return nil, err
	}
	return &comic, nil
}

func (r *comicRepository) GetChapterByID(id uuid.UUID) (*entity.Chapter, error) {
	var chapter entity.Chapter
	err := r.db.Preload("Images", func(db *gorm.DB) *gorm.DB {
		return db.Order("chapter_images.\"order\" ASC")
	}).First(&chapter, id).Error
	if err != nil {
		return nil, err
	}

	log.Println("chapter: ", chapter)

	return &chapter, nil
}

func (r *comicRepository) GetSeasonByComicID(comicID uuid.UUID, seasonNumber int) (*entity.Season, error) {
	var season entity.Season
	err := r.db.Where("comic_id = ? AND season_number = ?", comicID, seasonNumber).First(&season).Error
	if err != nil {
		return nil, err
	}
	return &season, nil
}

func (r *comicRepository) ListComics() ([]entity.Comic, error) {
	var comics []entity.Comic
	err := r.db.Preload("Tags.Translations").Where("status = ?", entity.ComicPublished).Order("updated_at desc").Limit(20).Find(&comics).Error
	if err != nil {
		return nil, err
	}
	return comics, nil
}

func (r *comicRepository) ListComicsByStatus(status entity.ComicStatus) ([]entity.Comic, error) {
	var comics []entity.Comic
	err := r.db.Preload("Tags.Translations").Where("status = ?", status).Order("updated_at desc").Find(&comics).Error
	if err != nil {
		return nil, err
	}
	return comics, nil
}

func (r *comicRepository) ListComicsByCreatorID(creatorID uuid.UUID) ([]entity.Comic, error) {
	var comics []entity.Comic
	err := r.db.Preload("Tags.Translations").Where("creator_id = ?", creatorID).Order("updated_at desc").Find(&comics).Error
	if err != nil {
		return nil, err
	}
	return comics, nil
}

func (r *comicRepository) ListComicsByAuthor(author string) ([]entity.Comic, error) {
	var comics []entity.Comic
	err := r.db.Preload("Tags.Translations").Where("author = ?", author).Order("updated_at desc").Find(&comics).Error
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
