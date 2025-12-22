package http

import (
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/pur108/webteen-be/internal/domain/entity"
	"github.com/pur108/webteen-be/internal/domain/repository"
	"github.com/pur108/webteen-be/internal/usecase"
)

type AdminHandler struct {
	comicUsecase usecase.ComicUsecase
	storageRepo  repository.StorageRepository
}

func NewAdminHandler(app *fiber.App, comicUsecase usecase.ComicUsecase, storageRepo repository.StorageRepository) *AdminHandler {
	return &AdminHandler{
		comicUsecase: comicUsecase,
		storageRepo:  storageRepo,
	}
}

// ListPendingComics Lists all comics with 'pending_review' status
func (h *AdminHandler) ListPendingComics(c *fiber.Ctx) error {
	comics, err := h.comicUsecase.ListPendingComics()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(comics)
}

// ApproveComic
func (h *AdminHandler) ApproveComic(c *fiber.Ctx) error {
	idStr := c.Params("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid comic ID"})
	}

	comic, err := h.comicUsecase.GetComic(id)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Comic not found"})
	}

	if comic.Status != entity.ComicPending {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Comic is not pending review"})
	}
	err = h.comicUsecase.ApproveComic(id)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(fiber.Map{"message": "Comic approved and files moved to public storage"})
}

func (h *AdminHandler) RejectComic(c *fiber.Ctx) error {
	idStr := c.Params("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid comic ID"})
	}

	var req struct {
		Reason string `json:"reason"`
	}
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request"})
	}

	err = h.comicUsecase.RejectComic(id, req.Reason)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(fiber.Map{"message": "Comic rejected"})
}
