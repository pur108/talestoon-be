package http

import (
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/pur108/talestoon-be/internal/usecase"
)

type LibraryHandler struct {
	libraryUsecase usecase.LibraryUsecase
}

func NewLibraryHandler(app *fiber.App, libraryUsecase usecase.LibraryUsecase) *LibraryHandler {
	return &LibraryHandler{libraryUsecase}
}

func (h *LibraryHandler) AddToLibrary(c *fiber.Ctx) error {
	userIDStr := c.Locals("user_id").(string)
	userID, _ := uuid.Parse(userIDStr)

	type reqBody struct {
		ComicID string `json:"comic_id"`
	}
	var req reqBody
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request body"})
	}

	comicID, err := uuid.Parse(req.ComicID)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid comic ID"})
	}

	if err := h.libraryUsecase.AddToLibrary(userID, comicID); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(fiber.Map{"message": "Added to library"})
}

func (h *LibraryHandler) RemoveFromLibrary(c *fiber.Ctx) error {
	userIDStr := c.Locals("user_id").(string)
	userID, _ := uuid.Parse(userIDStr)
	comicIDStr := c.Params("comic_id")
	comicID, err := uuid.Parse(comicIDStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid comic ID"})
	}

	if err := h.libraryUsecase.RemoveFromLibrary(userID, comicID); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(fiber.Map{"message": "Removed from library"})
}

func (h *LibraryHandler) GetUserLibrary(c *fiber.Ctx) error {
	userIDStr := c.Locals("user_id").(string)
	userID, _ := uuid.Parse(userIDStr)

	entries, err := h.libraryUsecase.GetUserLibrary(userID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(entries)
}

func (h *LibraryHandler) CreateFolder(c *fiber.Ctx) error {
	userIDStr := c.Locals("user_id").(string)
	userID, _ := uuid.Parse(userIDStr)

	type reqBody struct {
		Name        string `json:"name"`
		Description string `json:"description"`
		IsPublic    bool   `json:"is_public"`
	}
	var req reqBody
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request body"})
	}

	folder, err := h.libraryUsecase.CreateFolder(userID, req.Name, req.Description, req.IsPublic)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(folder)
}

func (h *LibraryHandler) DeleteFolder(c *fiber.Ctx) error {
	userIDStr := c.Locals("user_id").(string)
	userID, _ := uuid.Parse(userIDStr)
	folderIDStr := c.Params("id")
	folderID, err := uuid.Parse(folderIDStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid folder ID"})
	}

	if err := h.libraryUsecase.DeleteFolder(userID, folderID); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(fiber.Map{"message": "Folder deleted"})
}

func (h *LibraryHandler) GetUserFolders(c *fiber.Ctx) error {
	userIDStr := c.Locals("user_id").(string)
	userID, _ := uuid.Parse(userIDStr)

	folders, err := h.libraryUsecase.GetUserFolders(userID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(folders)
}

func (h *LibraryHandler) GetFolder(c *fiber.Ctx) error {
	folderIDStr := c.Params("id")
	folderID, err := uuid.Parse(folderIDStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid folder ID"})
	}

	folder, err := h.libraryUsecase.GetFolder(folderID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	if folder == nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Folder not found"})
	}

	// Check ownership if private or generally for edit rights info?
	// For this endpoint we assume authenticated user wants to view their own or a public one?
	// The requirement was "folders... like a playlist".
	// If it is my folder, I see it. If public, I see it via public Link.
	// This endpoint is under protected group, so it might be for owner.
	userIDStr := c.Locals("user_id").(string)
	userID, _ := uuid.Parse(userIDStr)

	if folder.UserID != userID && !folder.IsPublic {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"error": "Unauthorized"})
	}

	return c.JSON(folder)
}

func (h *LibraryHandler) GetPublicFolder(c *fiber.Ctx) error {
	slug := c.Params("slug")
	folder, err := h.libraryUsecase.GetFolderBySlug(slug)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	if folder == nil || !folder.IsPublic {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Folder not found"})
	}

	return c.JSON(folder)
}

func (h *LibraryHandler) AddToFolder(c *fiber.Ctx) error {
	userIDStr := c.Locals("user_id").(string)
	userID, _ := uuid.Parse(userIDStr)
	folderIDStr := c.Params("id")
	folderID, err := uuid.Parse(folderIDStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid folder ID"})
	}

	type reqBody struct {
		ComicID string `json:"comic_id"`
	}
	var req reqBody
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request body"})
	}
	comicID, err := uuid.Parse(req.ComicID)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid comic ID"})
	}

	if err := h.libraryUsecase.AddToFolder(userID, folderID, comicID); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(fiber.Map{"message": "Added to folder"})
}

func (h *LibraryHandler) RemoveFromFolder(c *fiber.Ctx) error {
	userIDStr := c.Locals("user_id").(string)
	userID, _ := uuid.Parse(userIDStr)
	folderIDStr := c.Params("id")
	folderID, err := uuid.Parse(folderIDStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid folder ID"})
	}

	comicIDStr := c.Params("comic_id")
	comicID, err := uuid.Parse(comicIDStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid comic ID"})
	}

	if err := h.libraryUsecase.RemoveFromFolder(userID, folderID, comicID); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(fiber.Map{"message": "Removed from folder"})
}
