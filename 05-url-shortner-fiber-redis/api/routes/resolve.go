package routes

import (
	"errors"
	"github.com/gofiber/fiber/v2"
	"github.com/pavitra93/05-url-shortner-fiber-redis/database"
	"github.com/redis/go-redis/v9"
)

func ResolveURL(c *fiber.Ctx) error {
	url := c.Params("url")
	r := database.CreateClient(0)
	defer r.Close()

	value, err := r.Get(database.Ctx, url).Result()
	if errors.Is(err, redis.Nil) {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Short not found",
		})
	} else if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Something went wrong, Cannot connect to DB",
		})
	}

	rInr := database.CreateClient(1)
	defer rInr.Close()
	rInr.Incr(database.Ctx, "counter")

	return c.Redirect(value, 302)
}
