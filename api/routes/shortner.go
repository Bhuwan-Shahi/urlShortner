package routes

import (
	"os"
	"strconv"
	"time"

	"github.com/Bhuwan-Shahi/urlShortner/database"
	"github.com/Bhuwan-Shahi/urlShortner/helpers"
	"github.com/asaskevich/govalidator"
	"github.com/go-redis/redis/v8"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type Request struct {
	URL         string        `json:"url"`
	CustomShort string        `json:"short"`
	Expiry      time.Duration `json:"expiry"`
}

type Response struct {
	URL             string        `json:"url"`
	CustomShort     string        `json:"short"`
	Expiry          time.Duration `json:"expiry"`
	XRateRemaining  int           `json:"rate_limit"`
	XRateLimitReset time.Duration `json:"rate_limit_reset"`
}

func ShortenURL(c *fiber.Ctx) error {
	body := new(Request)
	if err := c.BodyParser(body); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "cannot parse JSON",
		})
	}

	r2, err := database.CreateClient(1)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "failed to connect to rate limit database",
		})
	}
	defer r2.Close()

	// Get current user's rate limit
	val, err := r2.Get(database.Ctx, c.IP()).Result()
	if err == redis.Nil {
		// Set initial quota from environment variable
		quota, _ := strconv.Atoi(os.Getenv("API_QUOTA"))
		err = r2.Set(database.Ctx, c.IP(), quota, 30*60*time.Second).Err()
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "failed to set rate limit",
			})
		}
	}

	// Check remaining quota
	valInt, _ := strconv.Atoi(val)
	if valInt <= 0 {
		limit, _ := r2.TTL(database.Ctx, c.IP()).Result()
		return c.Status(fiber.StatusServiceUnavailable).JSON(fiber.Map{
			"error":            "rate limit exceeded",
			"rate_limit_reset": limit / time.Minute,
		})
	}

	// Validate URL
	if !govalidator.IsURL(body.URL) {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "invalid URL",
		})
	}

	// Check for domain security
	if !helpers.RemoveDomainError(body.URL) {
		return c.Status(fiber.StatusServiceUnavailable).JSON(fiber.Map{
			"error": "domain not allowed",
		})
	}

	// Enforce HTTPS
	body.URL = helpers.EnforceHTTP(body.URL)

	// Create or use custom short URL
	var id string
	if body.CustomShort == "" {
		id = uuid.New().String()[:6]
	} else {
		id = body.CustomShort
	}

	r, err := database.CreateClient(0)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "failed to connect to database",
		})
	}
	defer r.Close()

	// Check if custom short URL already exists
	val, _ = r.Get(database.Ctx, id).Result()
	if val != "" {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"error": "URL custom short already in use",
		})
	}

	// Set expiry time
	if body.Expiry == 0 {
		body.Expiry = 24 // Default to 24 hours
	}

	// Store URL in Redis
	err = r.Set(database.Ctx, id, body.URL, body.Expiry*3600*time.Second).Err()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "failed to store URL in database",
		})
	}

	// Decrement rate limit
	r2.Decr(database.Ctx, c.IP())

	// Prepare response
	resp := Response{
		URL:             body.URL,
		CustomShort:     "",
		Expiry:          body.Expiry,
		XRateRemaining:  valInt - 1,
		XRateLimitReset: 30,
	}

	// Add domain to short URL
	resp.CustomShort = os.Getenv("DOMAIN") + "/" + id

	return c.Status(fiber.StatusOK).JSON(resp)
}
