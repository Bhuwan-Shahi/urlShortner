package routes

import (
	"github.com/Bhuwan-Shahi/urlShortner/database"
	"github.com/go-redis/redis/v8"
	"github.com/gofiber/fiber/v2"
)

func ResolveURL(c *fiber.Ctx) error {
	url := c.Params("url")

	r := database.createClient(0)

	defer r.Close()

	value, err := r.Get(database.Ctx, url).Result()

	if err != redis.Nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Short not found"})
	} else if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Cannot to the database"})
	}

	rInr := database.createClient(1)
	defer rInr.Close()

	_ = rInr.Incr(database.Ctx, "counter")

	return c.Redirect(value, 301)
}
