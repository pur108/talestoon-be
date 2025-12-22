package repository

import (
	"github.com/google/uuid"
	"github.com/pur108/talestoon-be/internal/domain/entity"
)

type ComicRepository interface {
	CreateComic(comic *entity.Comic) error
	CreateChapter(chapter *entity.Chapter) error
	CreateSeason(season *entity.Season) error
	GetComicByID(id uuid.UUID) (*entity.Comic, error)
	GetChapterByID(id uuid.UUID) (*entity.Chapter, error)
	GetSeasonByComicID(comicID uuid.UUID, seasonNumber int) (*entity.Season, error)
	ListComics() ([]entity.Comic, error)
	ListComicsByStatus(status entity.ComicStatus) ([]entity.Comic, error)
	ListComicsByCreatorID(creatorID uuid.UUID) ([]entity.Comic, error)
	ListComicsByAuthor(author string) ([]entity.Comic, error)
	UpdateComic(comic *entity.Comic) error
	DeleteComic(id uuid.UUID) error
}
