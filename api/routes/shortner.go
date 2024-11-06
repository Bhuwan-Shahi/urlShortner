package routes

import (
	"time"

	"github.com/Bhuwan-Shahi/urlShortner/helpers"

	"github.com/gofiber/fiber/v2"
)

type request struct {
	URL         string        `json:"url"`
	cusotmShort string        `json:"short"`
	expiry      time.Duration `json:"expiry"`
}

type resonse struct {
	URL            string        `json:"url"`
	cusotmShort    string        `json:"short"`
	expiry         time.Duration `json:"expiry"`
	XRateRemaining int           `json:"rate_Limit"`
	XRateLimitRest time.Duration `json:"rate_limit_res"`
}

func shortenURL(c *fiber.Ctx) error {
	body := new(request)

	if err := c.BodyParser(body); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "cannot parse json"})
	}

	//implementing ratelimiting

	//check if the input sent by user is actual URL
	if !govalidator.IsURL(body.URL) {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid URL"})

	}

	//check for domain error
	if !helpers.RemoveDomainError(body.URL) {
		return c.Status(fiber.StatusServiceUnavailable).JSON(fiber.Map{"error": "You can't hack our system"})
	}
	//enforse https, SSL

	body.URL = helpers.EnforceHTTP(body.URL)
}
