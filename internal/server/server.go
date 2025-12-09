package server

import (
	"github.com/gofiber/fiber/v2"

	"github.com/pur108/talestoon-be.git/internal/database"
)

type FiberServer struct {
	*fiber.App

	db database.Service
}

func New() *FiberServer {
	server := &FiberServer{
		App: fiber.New(fiber.Config{
			ServerHeader: "talestoon",
			AppName:      "talestoon",
		}),

		db: database.New(),
	}

	return server
}
