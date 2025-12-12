package http

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/pur108/talestoon-be/internal/domain/exception"
	"github.com/pur108/talestoon-be/internal/usecase"
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

	if req.Title == "" || req.ChapterNumber == 0 || len(req.ImageURLs) == 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Title, chapter number, and at least one image are required"})
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

	if req.Title.En == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "English title is required"})
	}
	if req.Title.En == "" && req.Title.Th == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Title must have at least one language"})
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
	comics, err := h.comicUsecase.ListComics()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to fetch comics"})
	}

	return c.JSON(comics)
}
