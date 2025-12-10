package server

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"

	"github.com/pur108/talestoon-be/internal/delivery/http"
	"github.com/pur108/talestoon-be/internal/repository"
	"github.com/pur108/talestoon-be/internal/usecase"
)

func (s *FiberServer) RegisterFiberRoutes() {
	s.App.Use(logger.New())
	s.App.Use(cors.New(cors.Config{
		AllowOrigins:     "*",
		AllowMethods:     "GET,POST,PUT,DELETE,OPTIONS,PATCH",
		AllowHeaders:     "Accept,Authorization,Content-Type",
		AllowCredentials: false, // credentials require explicit origins
		MaxAge:           300,
	}))

	s.App.Get("/", s.HelloWorldHandler)
	s.App.Get("/health", s.healthHandler)

	db := s.db.GetDB()

	// user routes
	userRepo := repository.NewUserRepository(db)
	userUsecase := usecase.NewUserUsecase(userRepo)
	http.NewUserHandler(s.App, userUsecase)

	// auth routes
	authUsecase := usecase.NewAuthUsecase(userRepo)
	http.NewAuthHandler(s.App, authUsecase)
}

func (s *FiberServer) HelloWorldHandler(c *fiber.Ctx) error {
	resp := fiber.Map{
		"message": "Hello World",
	}

	return c.JSON(resp)
}

func (s *FiberServer) healthHandler(c *fiber.Ctx) error {
	return c.JSON(s.db.Health())
}
