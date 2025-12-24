package http

import (
	"fmt"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/pur108/webteen-be/internal/domain/exception"
	"github.com/pur108/webteen-be/internal/usecase"
)

type ComicHandler struct {
	comicUsecase usecase.ComicUsecase
}

func NewComicHandler(app *fiber.App, comicUsecase usecase.ComicUsecase) *ComicHandler {
	return &ComicHandler{comicUsecase}
}

func (h *ComicHandler) CreateChapter(c *fiber.Ctx) error {
	comicIDStr := c.Params("id")
	comicID, err := uuid.Parse(comicIDStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid comic ID"})
	}

	userIDStr := c.Locals("user_id").(string)
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid user ID"})
	}

	var req usecase.CreateChapterInput
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request"})
	}

	if len(req.Translations) == 0 || req.ChapterNumber == 0 || len(req.ImageURLs) == 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Translations (at least one), chapter number, and at least one image are required"})
	}

	// Validate title presence in translations
	hasTitle := false
	for _, t := range req.Translations {
		if t.Title != "" {
			hasTitle = true
			break
		}
	}
	if !hasTitle {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "At least one translation must have a title"})
	}

	chapter, err := h.comicUsecase.CreateChapter(comicID, userID, req)
	if err != nil {
		if err == exception.ErrUnauthorized {
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"error": "Unauthorized"})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(fiber.StatusCreated).JSON(chapter)
}

func (h *ComicHandler) CreateComic(c *fiber.Ctx) error {
	var req usecase.CreateComicInput
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request"})
	}

	userIDStr := c.Locals("user_id").(string)
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid user ID"})
	}
	req.CreatorID = userID

	if len(req.Translations) == 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "At least one translation is required"})
	}

	hasTitle := false
	for _, t := range req.Translations {
		if t.Title != "" {
			hasTitle = true
			break
		}
	}

	if !hasTitle {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Title is required in at least one language"})
	}

	comic, err := h.comicUsecase.CreateComic(req)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(fiber.StatusCreated).JSON(comic)
}

func (h *ComicHandler) ListMyComics(c *fiber.Ctx) error {
	fmt.Println("ListMyComics called")
	userIDStr := c.Locals("user_id").(string)
	fmt.Printf("UserID from locals: %s\n", userIDStr)
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		fmt.Printf("Error parsing UUID: %v\n", err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid user ID"})
	}

	comics, err := h.comicUsecase.ListMyComics(userID)
	if err != nil {
		fmt.Printf("Error fetching comics: %v\n", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to fetch comics"})
	}
	fmt.Printf("Found %d comics for user\n", len(comics))

	return c.JSON(comics)
}

func (h *ComicHandler) UpdateComic(c *fiber.Ctx) error {
	idStr := c.Params("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid comic ID"})
	}

	userIDStr := c.Locals("user_id").(string)
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid user ID"})
	}

	var req usecase.UpdateComicInput
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request"})
	}

	comic, err := h.comicUsecase.UpdateComic(id, userID, req)
	if err != nil {
		if err == exception.ErrUnauthorized {
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"error": "Unauthorized"})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(comic)
}

func (h *ComicHandler) DeleteComic(c *fiber.Ctx) error {
	idStr := c.Params("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid comic ID"})
	}

	userIDStr := c.Locals("user_id").(string)
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid user ID"})
	}

	err = h.comicUsecase.DeleteComic(id, userID)
	if err != nil {
		if err == exception.ErrUnauthorized {
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"error": "Unauthorized"})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to delete comic"})
	}

	return c.SendStatus(fiber.StatusNoContent)
}

func (h *ComicHandler) RequestPublish(c *fiber.Ctx) error {
	idStr := c.Params("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid comic ID"})
	}

	userIDStr := c.Locals("user_id").(string)
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid user ID"})
	}

	err = h.comicUsecase.RequestPublish(id, userID)
	if err != nil {
		if err == exception.ErrUnauthorized {
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"error": "Unauthorized"})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(fiber.Map{"message": "Comic submitted for review"})
}

func (h *ComicHandler) GetComic(c *fiber.Ctx) error {
	idStr := c.Params("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid comic ID"})
	}

	comic, err := h.comicUsecase.GetComic(id)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Comic not found"})
	}

	return c.JSON(comic)
}

func (h *ComicHandler) GetChapter(c *fiber.Ctx) error {
	idStr := c.Params("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid chapter ID"})
	}

	chapter, err := h.comicUsecase.GetChapter(id)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Chapter not found"})
	}

	return c.JSON(chapter)
}

func (h *ComicHandler) ListComics(c *fiber.Ctx) error {
	tagsParam := c.Query("tags")
	var tags []string
	if tagsParam != "" {
		tags = strings.Split(tagsParam, ",")
	}

	comics, err := h.comicUsecase.ListComics(tags)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to fetch comics"})
	}

	return c.JSON(comics)
}

func (h *ComicHandler) ListTags(c *fiber.Ctx) error {
	filterType := c.Query("type")
	tags, err := h.comicUsecase.ListTags(filterType)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to fetch tags"})
	}

	return c.JSON(tags)
}
