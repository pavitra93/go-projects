package routes

import (
	"github.com/asaskevich/govalidator"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/pavitra93/05-url-shortner-fiber-redis/database"
	"github.com/pavitra93/05-url-shortner-fiber-redis/helpers"
	"github.com/redis/go-redis/v9"
	"os"
	"strconv"
	"time"
)

type request struct {
	URL         string        `json:"url"`
	CustomShort string        `json:"custom_short"`
	Expiry      time.Duration `json:"expiry"`
}

type response struct {
	URL             string        `json:"url"`
	CustomShort     string        `json:"custom_short"`
	Expiry          time.Duration `json:"expiry"`
	XRateRemaining  int           `json:"x_rate_remaining"`
	XRateLimitReset time.Duration `json:"x_rate_limit_rest"`
}

func ShortenURL(c *fiber.Ctx) error {
	body := new(request)
	err := c.BodyParser(body)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Cannot parse request",
		})
	}
	// Implement rate limiting
	rdb0 := database.CreateClient(0)
	rdb1 := database.CreateClient(1)
	defer rdb1.Close()
	defer rdb0.Close()

	value, err := rdb1.Get(database.Ctx, c.IP()).Result()
	if err == redis.Nil {
		_ = rdb1.Set(database.Ctx, c.IP(), 10, 30*60*time.Second).Err()
	} else {
		valInt, _ := strconv.Atoi(value)
		if valInt <= 0 {
			limit, _ := rdb1.TTL(database.Ctx, c.IP()).Result()
			return c.Status(fiber.StatusServiceUnavailable).JSON(fiber.Map{
				"error":            "Rate Limit execeeded",
				"rate_limit_reset": limit / time.Nanosecond / time.Minute,
			})
		}
	}

	// Validate URL
	if !govalidator.IsURL(body.URL) {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid URL",
		})
	}

	// Remove same domain checks
	if !helpers.RemoveDomainError(body.URL) {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Domain Error",
		})
	}

	// Enforce HTTP
	body.URL = helpers.EnforceHTTP(body.URL)

	// Shorten url logic
	var id string
	if body.CustomShort == "" {
		id = uuid.New().String()[:6]
	} else {
		id = body.CustomShort
	}

	value, _ = rdb1.Get(database.Ctx, id).Result()
	if value != "" {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"error": "Your custom short is already used",
		})
	}

	if body.Expiry == 0 {
		body.Expiry = 24
	}

	err = rdb0.Set(database.Ctx, id, body.URL, body.Expiry*time.Minute).Err()

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Unable to connect to server",
		})
	}

	// decrement the rate limit by 1
	rdb1.Decr(database.Ctx, c.IP())

	resp := response{
		URL:             body.URL,
		CustomShort:     "",
		Expiry:          body.Expiry,
		XRateLimitReset: 10,
		XRateRemaining:  30,
	}

	value, _ = rdb1.Get(database.Ctx, c.IP()).Result()
	resp.XRateRemaining, _ = strconv.Atoi(value)

	ttl, _ := rdb1.TTL(database.Ctx, c.IP()).Result()
	resp.XRateLimitReset = ttl / time.Nanosecond / time.Minute

	resp.CustomShort = os.Getenv("DOMAIN") + "/" + id

	return c.Status(fiber.StatusOK).JSON(resp)
}
