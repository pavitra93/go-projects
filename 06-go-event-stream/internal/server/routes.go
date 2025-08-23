package server

import (
	"bufio"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/valyala/fasthttp"
	"time"
)

func (s *FiberServer) RegisterFiberRoutes() {
	// Apply CORS middleware
	s.App.Use(cors.New(cors.Config{
		AllowOrigins:     "*",
		AllowMethods:     "GET,POST,PUT,DELETE,OPTIONS,PATCH",
		AllowHeaders:     "Accept,Authorization,Content-Type",
		AllowCredentials: false, // credentials require explicit origins
		MaxAge:           300,
	}))

	s.App.Get("/", s.HelloWorldHandler)
	s.App.Get("/events", s.EventsStreamHandler)
	s.App.Get("/health", s.healthHandler)

}

func (s *FiberServer) EventsStreamHandler(c *fiber.Ctx) error {
	c.Set("Content-Type", "text/event-stream")
	events := []string{
		"event1",
		"event2",
		"event3",
		"event4",
		"event5",
		"event6",
		"event7",
		"event8",
		"event9",
		"event10",
		"event11",
		"event12",
		"event13",
		"event14",
		"event15",
		"event16",
		"event17",
	}

	c.Status(fiber.StatusOK)
	c.Context().SetBodyStreamWriter(fasthttp.StreamWriter(func(w *bufio.Writer) {
		for _, e := range events {
			if _, err := fmt.Fprintf(w, "data: %s\n\n", e); err != nil {
				return // client likely disconnected
			}
			if err := w.Flush(); err != nil {
				return
			}
			time.Sleep(1 * time.Second) // simulate streaming
		}

	}))

	return nil

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
