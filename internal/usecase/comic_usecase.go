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
	ListComics(tags []string) ([]entity.Comic, error)
	ListPendingComics() ([]entity.Comic, error)
	ListMyComics(creatorID uuid.UUID) ([]entity.Comic, error)
	ListComicsByAuthor(author string) ([]entity.Comic, error)
	ListTags(filterType string) ([]entity.Tag, error)
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

type ComicTranslationInput struct {
	LanguageCode string `json:"language_code"`
	Title        string `json:"title"`
	Subtitle     string `json:"subtitle"`
	Description  string `json:"description"`
}

type CreateComicInput struct {
	CreatorID           uuid.UUID                       `json:"creator_id"`
	Translations        []ComicTranslationInput         `json:"translations"`
	Author              string                          `json:"author"`
	Tags                []TagTranslationInput           `json:"tags"`
	TagIDs              []uuid.UUID                     `json:"tag_ids"`
	CoverImageURL       string                          `json:"cover_image_url"`
	BannerImageURL      string                          `json:"banner_image_url"`
	Status              entity.ComicStatus              `json:"status"`
	SerializationStatus entity.ComicSerializationStatus `json:"serialization_status"`
	Visibility          string                          `json:"visibility"`
	NSFW                bool                            `json:"nsfw"`
	SchedulePublishAt   *time.Time                      `json:"schedule_publish_at"`
}

type TagTranslationInput struct {
	LanguageCode string `json:"language_code"`
	Name         string `json:"name"`
}

type UpdateComicInput struct {
	Translations        []ComicTranslationInput         `json:"translations"`
	Author              string                          `json:"author"`
	CoverImageURL       string                          `json:"cover_image_url"`
	BannerImageURL      string                          `json:"banner_image_url"`
	Status              entity.ComicStatus              `json:"status"`
	SerializationStatus entity.ComicSerializationStatus `json:"serialization_status"`
	Visibility          string                          `json:"visibility"`
	NSFW                bool                            `json:"nsfw"`
}

type ChapterTranslationInput struct {
	LanguageCode string `json:"language_code"`
	Title        string `json:"title"`
}

type CreateChapterInput struct {
	Translations  []ChapterTranslationInput `json:"translations"`
	ChapterNumber int                       `json:"chapter_number"`
	ThumbnailURL  string                    `json:"thumbnail_url"`
	ImageURLs     []string                  `json:"image_urls"`
}

func (u *comicUsecase) CreateComic(input CreateComicInput) (*entity.Comic, error) {
	comicID := uuid.New()
	comic := &entity.Comic{
		ID:                  comicID,
		CreatorID:           input.CreatorID,
		Author:              input.Author,
		CoverImageURL:       input.CoverImageURL,
		BannerImageURL:      input.BannerImageURL,
		Status:              input.Status,
		SerializationStatus: input.SerializationStatus,
		Visibility:          input.Visibility,
		NSFW:                input.NSFW,
		SchedulePublishAt:   input.SchedulePublishAt,
		CreatedAt:           time.Now(),
		UpdatedAt:           time.Now(),
		Translations:        []entity.ComicTranslation{},
	}

	for _, t := range input.Translations {
		comic.Translations = append(comic.Translations, entity.ComicTranslation{
			ID:               uuid.New(),
			ComicID:          comicID,
			LanguageCode:     t.LanguageCode,
			Title:            t.Title,
			Synopsis:         t.Description,
			AlternativeTitle: t.Subtitle,
			CreatedAt:        time.Now(),
			UpdatedAt:        time.Now(),
		})
	}

	if comic.SerializationStatus == "" {
		comic.SerializationStatus = entity.ComicOngoing
	}

	var tags []entity.Tag
	for _, t := range input.Tags {
		tagID := uuid.New()
		slug := utils.SimpleSlug(t.Name)

		newTag := entity.Tag{
			ID:        tagID,
			Slug:      slug,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
			Type:      "genre",
			Translations: []entity.TagTranslation{
				{
					ID:       uuid.New(),
					TagID:    tagID,
					Language: t.LanguageCode,
					Name:     t.Name,
				},
			},
		}
		tags = append(tags, newTag)
	}

	for _, tagID := range input.TagIDs {
		tags = append(tags, entity.Tag{
			ID: tagID,
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

func (u *comicUsecase) ListComics(tags []string) ([]entity.Comic, error) {
	return u.comicRepo.ListComics(tags)
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

func (u *comicUsecase) ListTags(filterType string) ([]entity.Tag, error) {
	return u.comicRepo.ListTags(filterType)
}

func (u *comicUsecase) ListComicsByAuthor(author string) ([]entity.Comic, error) {
	return u.comicRepo.ListComicsByAuthor(author)
}

func (u *comicUsecase) CreateChapter(comicID uuid.UUID, creatorID uuid.UUID, input CreateChapterInput) (*entity.Chapter, error) {
	comic, err := u.comicRepo.GetComicByID(comicID)
	if err != nil {
		return nil, err
	}
	if comic.CreatorID != creatorID {
		return nil, exception.ErrUnauthorized
	}

	// No Season creation anymore
	chapterID := uuid.New()
	chapter := &entity.Chapter{
		ID:            chapterID,
		ComicID:       comic.ID, // Direct link to Comic
		ChapterNumber: input.ChapterNumber,
		Status:        entity.ChapterPublished,
		ThumbnailURL:  input.ThumbnailURL,
		PublishedAt:   nowPtr(),
		Images:        []entity.ChapterImage{},
		Translations:  []entity.ChapterTranslation{},
	}

	for _, t := range input.Translations {
		chapter.Translations = append(chapter.Translations, entity.ChapterTranslation{
			ID:           uuid.New(),
			ChapterID:    chapterID,
			LanguageCode: t.LanguageCode,
			Title:        t.Title,
		})
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

	// Update base fields
	comic.Author = input.Author
	comic.CoverImageURL = input.CoverImageURL
	comic.BannerImageURL = input.BannerImageURL
	comic.Status = input.Status
	comic.SerializationStatus = input.SerializationStatus
	comic.Visibility = input.Visibility
	comic.NSFW = input.NSFW
	comic.UpdatedAt = time.Now()

	// Update Translations logic (simplified for implementation speed: overwrite if exists in list, or append)
	// Ideally we find En translation and update it.
	updateTranslation := func(lang string, title, synopsis, altTitle string) {
		found := false
		for i := range comic.Translations {
			if comic.Translations[i].LanguageCode == lang {
				comic.Translations[i].Title = title
				comic.Translations[i].Synopsis = synopsis
				comic.Translations[i].AlternativeTitle = altTitle
				comic.Translations[i].UpdatedAt = time.Now()
				found = true
				break
			}
		}
		if !found {
			comic.Translations = append(comic.Translations, entity.ComicTranslation{
				ID:               uuid.New(),
				ComicID:          comic.ID,
				LanguageCode:     lang,
				Title:            title,
				Synopsis:         synopsis,
				AlternativeTitle: altTitle,
				CreatedAt:        time.Now(),
				UpdatedAt:        time.Now(),
			})
		}
	}

	for _, t := range input.Translations {
		updateTranslation(t.LanguageCode, t.Title, t.Description, t.Subtitle)
	}

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

	// Updated Loop for flat Chapters structure
	for i := range comic.Chapters {
		for j := range comic.Chapters[i].Images {
			img := &comic.Chapters[i].Images[j]
			newURL, err := moveImage(img.ImageURL)
			if err != nil {
				return err
			}
			img.ImageURL = newURL
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
