package usecase

import (
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/pur108/talestoon-be/internal/domain/entity"
	"github.com/pur108/talestoon-be/internal/domain/exception"
	"github.com/pur108/talestoon-be/internal/domain/repository"
	"github.com/pur108/talestoon-be/pkg/utils"
)

type ComicUsecase interface {
	CreateComic(input CreateComicInput) (*entity.Comic, error)
	GetComic(id uuid.UUID) (*entity.Comic, error)
	GetChapter(id uuid.UUID) (*entity.Chapter, error)
	CreateChapter(comicID uuid.UUID, creatorID uuid.UUID, input CreateChapterInput) (*entity.Chapter, error)
	ListComics() ([]entity.Comic, error)
	ListMyComics(creatorID uuid.UUID) ([]entity.Comic, error)
	UpdateComic(id uuid.UUID, creatorID uuid.UUID, input UpdateComicInput) (*entity.Comic, error)
	DeleteComic(id uuid.UUID, creatorID uuid.UUID) error
}

type comicUsecase struct {
	comicRepo repository.ComicRepository
	userRepo  repository.UserRepository
}

func NewComicUsecase(comicRepo repository.ComicRepository, userRepo repository.UserRepository) ComicUsecase {
	return &comicUsecase{comicRepo, userRepo}
}

type CreateComicInput struct {
	CreatorID   uuid.UUID                 `json:"creator_id"`
	Title       entity.MultilingualText   `json:"title"`
	Subtitle    entity.MultilingualText   `json:"subtitle"`
	Description entity.MultilingualText   `json:"description"`
	Author      string                    `json:"author"`
	Genres      []string                  `json:"genres"`
	Tags        []entity.MultilingualText `json:"tags"`
	//ThumbnailURL        string                    `json:"thumbnail_url"`
	CoverImageURL       string             `json:"cover_image_url"`
	BannerImageURL      string             `json:"banner_image_url"`
	Status              entity.ComicStatus `json:"status"`
	Visibility          string             `json:"visibility"`
	NSFW                bool               `json:"nsfw"`
	SchedulePublishAt   *time.Time         `json:"schedule_publish_at"`
	MonetizationEnabled bool               `json:"monetization_enabled"`
	MonetizationType    string             `json:"monetization_type"`
	DefaultUnlockType   string             `json:"default_unlock_type"`
}

type UpdateComicInput struct {
	Title       entity.MultilingualText `json:"title"`
	Subtitle    entity.MultilingualText `json:"subtitle"`
	Description entity.MultilingualText `json:"description"`
	Author      string                  `json:"author"`
	Genres      []string                `json:"genres"`
	//ThumbnailURL        string                  `json:"thumbnail_url"`
	CoverImageURL       string             `json:"cover_image_url"`
	BannerImageURL      string             `json:"banner_image_url"`
	Status              entity.ComicStatus `json:"status"`
	Visibility          string             `json:"visibility"`
	NSFW                bool               `json:"nsfw"`
	MonetizationEnabled bool               `json:"monetization_enabled"`
	MonetizationType    string             `json:"monetization_type"`
	DefaultUnlockType   string             `json:"default_unlock_type"`
}

type CreateChapterInput struct {
	Title         string   `json:"title"`
	ChapterNumber int      `json:"chapter_number"`
	ImageURLs     []string `json:"image_urls"`
	Price         float64  `json:"price"`
}

func (u *comicUsecase) CreateComic(input CreateComicInput) (*entity.Comic, error) {
	comic := &entity.Comic{
		ID:          uuid.New(),
		CreatorID:   input.CreatorID,
		Title:       input.Title,
		Subtitle:    input.Subtitle,
		Description: input.Description,
		Author:      input.Author,
		Genres:      input.Genres,
		CoverImageURL:     input.CoverImageURL,
		BannerImageURL:    input.BannerImageURL,
		Status:            input.Status,
		Visibility:        input.Visibility,
		NSFW:              input.NSFW,
		SchedulePublishAt: input.SchedulePublishAt,
		// MonetizationEnabled: input.MonetizationEnabled,
		// MonetizationType:    input.MonetizationType,
		// DefaultUnlockType:   input.DefaultUnlockType,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	var tags []entity.Tag
	for _, t := range input.Tags {
		tagID := uuid.New()
		slug := utils.SimpleSlug(t.En)

		tags = append(tags, entity.Tag{
			ID:        tagID,
			Slug:      slug,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
			Translations: []entity.TagTranslation{
				{
					ID:       uuid.New(),
					TagID:    tagID,
					Language: "en",
					Name:     t.En,
				},
				{
					ID:       uuid.New(),
					TagID:    tagID,
					Language: "th",
					Name:     t.Th,
				},
			},
		})
	}
	comic.Tags = tags

	if err := u.comicRepo.CreateComic(comic); err != nil {
		return nil, err
	}

	user, err := u.userRepo.FindByID(input.CreatorID)
	if err == nil && user.Role == entity.RoleUser {
		user.Role = entity.RoleCreator
		_ = u.userRepo.Update(user)
	}

	return comic, nil
}

func (u *comicUsecase) CreateChapter(comicID uuid.UUID, creatorID uuid.UUID, input CreateChapterInput) (*entity.Chapter, error) {
	comic, err := u.comicRepo.GetComicByID(comicID)
	if err != nil {
		return nil, err
	}
	if comic.CreatorID != creatorID {
		return nil, exception.ErrUnauthorized
	}

	season, err := u.comicRepo.GetSeasonByComicID(comic.ID, 1)
	if err != nil || season == nil {
		newSeason := &entity.Season{
			ID:           uuid.New(),
			ComicID:      comic.ID,
			SeasonNumber: 1,
			Title:        "Season 1",
		}
		if err := u.comicRepo.CreateSeason(newSeason); err != nil {
			return nil, fmt.Errorf("failed to create season: %w", err)
		}
		season = newSeason
	}

	chapter := &entity.Chapter{
		ID:            uuid.New(),
		SeasonID:      season.ID,
		ChapterNumber: input.ChapterNumber,
		Title:         input.Title,
		Status:        entity.ChapterPublished,
		PublishedAt:   nowPtr(),
		Images:        []entity.ChapterImage{},
	}

	for i, url := range input.ImageURLs {
		chapter.Images = append(chapter.Images, entity.ChapterImage{
			ID:        uuid.New(),
			ChapterID: chapter.ID,
			ImageURL:  url,
			Order:     i + 1,
		})
	}

	if err := u.comicRepo.CreateChapter(chapter); err != nil {
		return nil, err
	}

	return chapter, nil
}

func nowPtr() *time.Time {
	t := time.Now()
	return &t
}

func (u *comicUsecase) GetComic(id uuid.UUID) (*entity.Comic, error) {
	return u.comicRepo.GetComicByID(id)
}

func (u *comicUsecase) GetChapter(id uuid.UUID) (*entity.Chapter, error) {
	return u.comicRepo.GetChapterByID(id)
}

func (u *comicUsecase) ListComics() ([]entity.Comic, error) {
	return u.comicRepo.ListComics()
}

func (u *comicUsecase) ListMyComics(creatorID uuid.UUID) ([]entity.Comic, error) {
	user, err := u.userRepo.FindByID(creatorID)
	if err != nil {
		return nil, err
	}
	fmt.Printf("ListMyComics Usecase: User found: %s. Querying by author: %s\n", user.Username, user.Username)

	return u.comicRepo.ListComicsByCreatorID(user.ID)
}

func (u *comicUsecase) UpdateComic(id uuid.UUID, creatorID uuid.UUID, input UpdateComicInput) (*entity.Comic, error) {
	comic, err := u.comicRepo.GetComicByID(id)
	if err != nil {
		return nil, err
	}

	if comic.CreatorID != creatorID {
		return nil, exception.ErrUnauthorized
	}

	comic.Title = input.Title
	comic.Subtitle = input.Subtitle
	comic.Description = input.Description
	comic.Author = input.Author
	comic.Genres = input.Genres
	//comic.ThumbnailURL = input.ThumbnailURL
	comic.CoverImageURL = input.CoverImageURL
	comic.BannerImageURL = input.BannerImageURL
	comic.Status = input.Status
	comic.Visibility = input.Visibility
	comic.NSFW = input.NSFW
	// comic.MonetizationEnabled = input.MonetizationEnabled
	// comic.MonetizationType = input.MonetizationType
	// comic.DefaultUnlockType = input.DefaultUnlockType
	comic.UpdatedAt = time.Now()

	if err := u.comicRepo.UpdateComic(comic); err != nil {
		return nil, err
	}

	return comic, nil
}

func (u *comicUsecase) DeleteComic(id uuid.UUID, creatorID uuid.UUID) error {
	comic, err := u.comicRepo.GetComicByID(id)
	if err != nil {
		return err
	}

	if comic.CreatorID != creatorID {
		return exception.ErrUnauthorized
	}

	return u.comicRepo.DeleteComic(id)
}
