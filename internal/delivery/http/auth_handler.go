package http

import (
	"github.com/gofiber/fiber/v2"
	"github.com/pur108/talestoon-be.git/internal/domain"
	"github.com/pur108/talestoon-be.git/internal/usecase"
)

type AuthHandler struct {
	authUsecase usecase.AuthUsecase
}

func NewAuthHandler(app *fiber.App, authUsecase usecase.AuthUsecase) {
	handler := &AuthHandler{authUsecase}
	group := app.Group("/api/auth")
	group.Post("/signup", handler.SignUp)
	group.Post("/login", handler.Login)
}

func (h *AuthHandler) SignUp(c *fiber.Ctx) error {
	type Request struct {
		Username string          `json:"username"`
		Email    string          `json:"email"`
		Password string          `json:"password"`
		Role     domain.UserRole `json:"role"`
	}

	var req Request
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request"})
	}

	if req.Username == "" || req.Email == "" || req.Password == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Username, email, and password are required"})
	}

	if req.Role == "" {
		req.Role = domain.RoleUser
	}

	user, err := h.authUsecase.SignUp(req.Username, req.Email, req.Password, req.Role)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(fiber.StatusCreated).JSON(user)
}

func (h *AuthHandler) Login(c *fiber.Ctx) error {
	type Request struct {
		Identifier string `json:"identifier"`
		Password   string `json:"password"`
	}

	var req Request
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request"})
	}

	token, user, err := h.authUsecase.Login(req.Identifier, req.Password)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(fiber.Map{
		"token": token,
		"user":  user,
	})
}
