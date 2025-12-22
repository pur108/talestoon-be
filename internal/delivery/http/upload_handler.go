package http

import (
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/pur108/talestoon-be/internal/usecase"
)

type UploadHandler struct {
	uploadUsecase usecase.UploadUsecase
}

func NewUploadHandler(app *fiber.App, uploadUsecase usecase.UploadUsecase) *UploadHandler {
	return &UploadHandler{
		uploadUsecase: uploadUsecase,
	}
}

func (h *UploadHandler) UploadFile(c *fiber.Ctx) error {
	file, err := c.FormFile("file")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "File is required"})
	}

	userID, ok := c.Locals("user_id").(string)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Unauthorized"})
	}

	// Read file content
	f, err := file.Open()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to open file"})
	}
	defer func() {
		_ = f.Close()
	}()


	buffer := make([]byte, file.Size)
	_, err = f.Read(buffer)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to read file"})
	}

	contentType := file.Header.Get("Content-Type")

	resp, err := h.uploadUsecase.UploadFile(userID, file.Filename, buffer, contentType)
	if err != nil {
		if strings.Contains(err.Error(), "invalid file type") || strings.Contains(err.Error(), "file too large") {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
		}
		if strings.Contains(err.Error(), "server storage configuration missing") {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to upload file",
		})
	}

	return c.JSON(resp)
}
