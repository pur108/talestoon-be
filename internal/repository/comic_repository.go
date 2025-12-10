package repository

import (
	"github.com/google/uuid"
	"github.com/pur108/talestoon-be/internal/domain"
	"gorm.io/gorm"
)

type comicRepository struct {
	db *gorm.DB
}

func NewComicRepository(db *gorm.DB) domain.ComicRepository {
	return &comicRepository{db}
}

func (r *comicRepository) CreateComic(comic *domain.Comic) error {
	return r.db.Create(comic).Error
}

func (r *comicRepository) CreateChapter(chapter *domain.Chapter) error {
	return r.db.Create(chapter).Error
}

func (r *comicRepository) CreateSeason(season *domain.Season) error {
	return r.db.Create(season).Error
}

func (r *comicRepository) GetComicByID(id uuid.UUID) (*domain.Comic, error) {
	var comic domain.Comic
	err := r.db.Preload("Seasons.Chapters").Preload("Tags.Translations").First(&comic, id).Error
	if err != nil {
		return nil, err
	}
	return &comic, nil
}

func (r *comicRepository) GetChapterByID(id uuid.UUID) (*domain.Chapter, error) {
	var chapter domain.Chapter
	err := r.db.Preload("Images.TextLayers.Translations").First(&chapter, id).Error
	if err != nil {
		return nil, err
	}
	return &chapter, nil
}

func (r *comicRepository) GetSeasonByComicID(comicID uuid.UUID, seasonNumber int) (*domain.Season, error) {
	var season domain.Season
	err := r.db.Where("comic_id = ? AND season_number = ?", comicID, seasonNumber).First(&season).Error
	if err != nil {
		return nil, err
	}
	return &season, nil
}

func (r *comicRepository) ListComics() ([]domain.Comic, error) {
	var comics []domain.Comic
	err := r.db.Preload("Tags.Translations").Order("updated_at desc").Limit(20).Find(&comics).Error
	if err != nil {
		return nil, err
	}
	return comics, nil
}

func (r *comicRepository) ListComicsByCreatorID(creatorID uuid.UUID) ([]domain.Comic, error) {
	var comics []domain.Comic
	err := r.db.Preload("Tags.Translations").Where("creator_id = ?", creatorID).Order("updated_at desc").Find(&comics).Error
	if err != nil {
		return nil, err
	}
	return comics, nil
}

func (r *comicRepository) ListComicsByAuthor(author string) ([]domain.Comic, error) {
	var comics []domain.Comic
	err := r.db.Preload("Tags.Translations").Where("author = ?", author).Order("updated_at desc").Find(&comics).Error
	if err != nil {
		return nil, err
	}
	return comics, nil
}

func (r *comicRepository) UpdateComic(comic *domain.Comic) error {
	return r.db.Save(comic).Error
}

func (r *comicRepository) DeleteComic(id uuid.UUID) error {
	return r.db.Delete(&domain.Comic{}, id).Error
}
