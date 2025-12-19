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
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "No file uploaded"})
	}

	url, err := h.uploadUsecase.UploadFile(file)
	if err != nil {
		if strings.Contains(err.Error(), "invalid file type") || strings.Contains(err.Error(), "invalid bucket") {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
		}
		if strings.Contains(err.Error(), "server storage configuration missing") {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "File upload failed"})
	}

	return c.JSON(fiber.Map{
		"url": url,
	})
}
