package server

import (
	"github.com/gofiber/fiber/v2"

	"go-event-stream/internal/database"
)

type FiberServer struct {
	*fiber.App

	db database.Service
}

func New() *FiberServer {
	server := &FiberServer{
		App: fiber.New(fiber.Config{
			ServerHeader: "06-go-event-stream",
			AppName:      "06-go-event-stream",
		}),

		db: database.New(),
	}

	return server
}
