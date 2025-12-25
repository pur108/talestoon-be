package http

import (
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/pur108/webteen-be/internal/domain/entity"
	"github.com/pur108/webteen-be/internal/usecase"
)

type LibraryHandler struct {
	libraryUsecase usecase.LibraryUsecase
}

func NewLibraryHandler(app *fiber.App, libraryUsecase usecase.LibraryUsecase) *LibraryHandler {
	return &LibraryHandler{libraryUsecase}
}

func (h *LibraryHandler) GetUserLibrary(c *fiber.Ctx) error {
	userID, err := getUserID(c)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Unauthorized"})
	}

	folder, err := h.libraryUsecase.GetMyLibrary(userID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	type LibraryEntryResponse struct {
		ID        string      `json:"id"`
		Comic     interface{} `json:"comic"`
		CreatedAt string      `json:"created_at"`
	}

	response := make([]LibraryEntryResponse, 0)
	for _, item := range folder.Items {
		response = append(response, LibraryEntryResponse{
			ID:        item.ID.String(),
			Comic:     item.Comic,
			CreatedAt: item.AddedAt.Format("2006-01-02T15:04:05Z07:00"),
		})
	}

	return c.JSON(response)
}

func (h *LibraryHandler) AddToLibrary(c *fiber.Ctx) error {
	userID, err := getUserID(c)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Unauthorized"})
	}

	var req struct {
		ComicID string `json:"comic_id"`
	}
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
	userID, err := getUserID(c)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Unauthorized"})
	}

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

func (h *LibraryHandler) CheckInLibrary(c *fiber.Ctx) error {
	userID, err := getUserID(c)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Unauthorized"})
	}

	comicIDStr := c.Params("comic_id")
	comicID, err := uuid.Parse(comicIDStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid comic ID"})
	}

	inLibrary, err := h.libraryUsecase.CheckInLibrary(userID, comicID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(fiber.Map{"in_library": inLibrary})
}

func (h *LibraryHandler) ListFolders(c *fiber.Ctx) error {
	userID, err := getUserID(c)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Unauthorized"})
	}

	folders, err := h.libraryUsecase.ListFolders(userID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	if folders == nil {
		folders = []entity.LibraryFolder{}
	}

	return c.JSON(folders)
}

func (h *LibraryHandler) CreateFolder(c *fiber.Ctx) error {
	userID, err := getUserID(c)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Unauthorized"})
	}

	var req struct {
		Name string `json:"name"`
	}
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request body"})
	}

	if req.Name == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Folder name is required"})
	}

	folder, err := h.libraryUsecase.CreateFolder(userID, req.Name)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(folder)
}

func (h *LibraryHandler) DeleteFolder(c *fiber.Ctx) error {
	userID, err := getUserID(c)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Unauthorized"})
	}

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

func (h *LibraryHandler) GetFolder(c *fiber.Ctx) error {
	userID, err := getUserID(c)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Unauthorized"})
	}

	folderIDStr := c.Params("id")
	folderID, err := uuid.Parse(folderIDStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid folder ID"})
	}

	folder, err := h.libraryUsecase.GetFolder(userID, folderID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(folder)
}

func getUserID(c *fiber.Ctx) (uuid.UUID, error) {
	userIDStr := c.Locals("user_id")
	if userIDStr == nil {
		return uuid.Nil, fiber.NewError(fiber.StatusUnauthorized, "Unauthorized")
	}
	return uuid.Parse(userIDStr.(string))
}
