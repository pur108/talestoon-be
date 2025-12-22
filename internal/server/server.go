package server

import (
	"github.com/gofiber/fiber/v2"

	"github.com/pur108/webteen-be/internal/database"
)

type FiberServer struct {
	*fiber.App

	db database.Service
}

func New() *FiberServer {
	server := &FiberServer{
		App: fiber.New(fiber.Config{
			ServerHeader: "webteen",
			AppName:      "webteen",
		}),

		db: database.New(),
	}

	return server
}
