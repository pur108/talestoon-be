package usecase

import (
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/pur108/webteen-be/internal/domain/entity"
	"github.com/pur108/webteen-be/internal/domain/exception"
	"github.com/pur108/webteen-be/internal/domain/repository"
	"github.com/pur108/webteen-be/pkg/utils"
)

type ComicUsecase interface {
	CreateComic(input CreateComicInput) (*entity.Comic, error)
	GetComic(id uuid.UUID) (*entity.Comic, error)
	GetChapter(id uuid.UUID) (*entity.Chapter, error)
	CreateChapter(comicID uuid.UUID, creatorID uuid.UUID, input CreateChapterInput) (*entity.Chapter, error)
	ListComics() ([]entity.Comic, error)
	ListPendingComics() ([]entity.Comic, error)
	ListMyComics(creatorID uuid.UUID) ([]entity.Comic, error)
	UpdateComic(id uuid.UUID, creatorID uuid.UUID, input UpdateComicInput) (*entity.Comic, error)
	DeleteComic(id uuid.UUID, creatorID uuid.UUID) error
	RequestPublish(id uuid.UUID, creatorID uuid.UUID) error
	ApproveComic(id uuid.UUID) error
	RejectComic(id uuid.UUID, reason string) error
}

type comicUsecase struct {
	comicRepo   repository.ComicRepository
	userRepo    repository.UserRepository
	storageRepo repository.StorageRepository
}

func NewComicUsecase(comicRepo repository.ComicRepository, userRepo repository.UserRepository, storageRepo repository.StorageRepository) ComicUsecase {
	return &comicUsecase{comicRepo, userRepo, storageRepo}
}

type CreateComicInput struct {
	CreatorID           uuid.UUID                       `json:"creator_id"`
	Title               entity.MultilingualText         `json:"title"`
	Subtitle            entity.MultilingualText         `json:"subtitle"`
	Description         entity.MultilingualText         `json:"description"`
	Author              string                          `json:"author"`
	Genres              []string                        `json:"genres"`
	Tags                []entity.MultilingualText       `json:"tags"`
	CoverImageURL       string                          `json:"cover_image_url"`
	BannerImageURL      string                          `json:"banner_image_url"`
	Status              entity.ComicStatus              `json:"status"`
	SerializationStatus entity.ComicSerializationStatus `json:"serialization_status"`
	Visibility          string                          `json:"visibility"`
	NSFW                bool                            `json:"nsfw"`
	SchedulePublishAt   *time.Time                      `json:"schedule_publish_at"`
}

type UpdateComicInput struct {
	Title               entity.MultilingualText         `json:"title"`
	Subtitle            entity.MultilingualText         `json:"subtitle"`
	Description         entity.MultilingualText         `json:"description"`
	Author              string                          `json:"author"`
	Genres              []string                        `json:"genres"`
	CoverImageURL       string                          `json:"cover_image_url"`
	BannerImageURL      string                          `json:"banner_image_url"`
	Status              entity.ComicStatus              `json:"status"`
	SerializationStatus entity.ComicSerializationStatus `json:"serialization_status"`
	Visibility          string                          `json:"visibility"`
	NSFW                bool                            `json:"nsfw"`
}

// ... existing code ...

type CreateChapterInput struct {
	Title         string   `json:"title"`
	SeasonTitle   string   `json:"season_title"`
	ChapterNumber int      `json:"chapter_number"`
	ImageURLs     []string `json:"image_urls"`
	Price         float64  `json:"price"`
}

func (u *comicUsecase) CreateComic(input CreateComicInput) (*entity.Comic, error) {
	comic := &entity.Comic{
		ID:                  uuid.New(),
		CreatorID:           input.CreatorID,
		Title:               input.Title,
		Subtitle:            input.Subtitle,
		Description:         input.Description,
		Author:              input.Author,
		Genres:              input.Genres,
		CoverImageURL:       input.CoverImageURL,
		BannerImageURL:      input.BannerImageURL,
		Status:              input.Status,
		SerializationStatus: input.SerializationStatus,
		Visibility:          input.Visibility,
		NSFW:                input.NSFW,
		SchedulePublishAt:   input.SchedulePublishAt,
		CreatedAt:           time.Now(),
		UpdatedAt:           time.Now(),
	}

	// Default to ongoing if not specified
	if comic.SerializationStatus == "" {
		comic.SerializationStatus = entity.ComicOngoing
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

func (u *comicUsecase) GetComic(id uuid.UUID) (*entity.Comic, error) {
	return u.comicRepo.GetComicByID(id)
}

func (u *comicUsecase) GetChapter(id uuid.UUID) (*entity.Chapter, error) {
	return u.comicRepo.GetChapterByID(id)
}

func (u *comicUsecase) ListComics() ([]entity.Comic, error) {
	return u.comicRepo.ListComics()
}

func (u *comicUsecase) ListPendingComics() ([]entity.Comic, error) {
	return u.comicRepo.ListComicsByStatus(entity.ComicPending)
}

func (u *comicUsecase) ListMyComics(creatorID uuid.UUID) ([]entity.Comic, error) {
	user, err := u.userRepo.FindByID(creatorID)
	if err != nil {
		return nil, err
	}

	return u.comicRepo.ListComicsByCreatorID(user.ID)
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
	comic.CoverImageURL = input.CoverImageURL
	comic.BannerImageURL = input.BannerImageURL
	comic.Status = input.Status
	comic.SerializationStatus = input.SerializationStatus
	comic.Visibility = input.Visibility
	comic.NSFW = input.NSFW
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

func (u *comicUsecase) RequestPublish(id uuid.UUID, creatorID uuid.UUID) error {
	comic, err := u.comicRepo.GetComicByID(id)
	if err != nil {
		return err
	}

	if comic.CreatorID != creatorID {
		return exception.ErrUnauthorized
	}

	if comic.Status != entity.ComicDraft && comic.Status != entity.ComicRejected {
		return fmt.Errorf("only drafts or rejected comics can be submitted for review")
	}

	comic.Status = entity.ComicPending
	comic.UpdatedAt = time.Now()

	return u.comicRepo.UpdateComic(comic)
}

func (u *comicUsecase) ApproveComic(id uuid.UUID) error {
	comic, err := u.comicRepo.GetComicByID(id)
	if err != nil {
		return err
	}

	moveImage := func(url string) (string, error) {
		if url == "" || !strings.Contains(url, "/drafts/") {
			return url, nil
		}

		parts := strings.Split(url, "/media/")
		if len(parts) < 2 {
			return url, nil
		}

		srcPath := parts[1]
		destPath := strings.Replace(srcPath, "drafts/", "public/", 1)

		if err := u.storageRepo.MoveFile("media", srcPath, destPath); err != nil {
			fmt.Printf("Failed to move file %s: %v\n", srcPath, err)
			return url, err
		}

		return strings.Replace(url, "/drafts/", "/public/", 1), nil
	}

	var errMove error
	comic.CoverImageURL, errMove = moveImage(comic.CoverImageURL)
	if errMove != nil {
		return errMove
	}

	comic.BannerImageURL, errMove = moveImage(comic.BannerImageURL)
	if errMove != nil {
		return errMove
	}

	for i := range comic.Seasons {
		for j := range comic.Seasons[i].Chapters {
			for k := range comic.Seasons[i].Chapters[j].Images {
				img := &comic.Seasons[i].Chapters[j].Images[k]
				newURL, err := moveImage(img.ImageURL)
				if err != nil {
					return err
				}
				img.ImageURL = newURL
			}
		}
	}

	now := time.Now()
	comic.Status = entity.ComicPublished
	comic.Visibility = entity.VisibilityPublic
	comic.ApprovedAt = &now

	return u.comicRepo.UpdateComic(comic)
}

func (u *comicUsecase) RejectComic(id uuid.UUID, reason string) error {
	comic, err := u.comicRepo.GetComicByID(id)
	if err != nil {
		return err
	}

	comic.Status = entity.ComicRejected
	comic.RejectionReason = reason
	comic.UpdatedAt = time.Now()

	return u.comicRepo.UpdateComic(comic)
}
