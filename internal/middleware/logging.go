package middleware

import (
	"log"
	"time"

	"github.com/JoeG98/pizza-backend/internal/models"
	"github.com/gofiber/fiber/v2"
)

func AccessLog(c *fiber.Ctx) error {
	start := time.Now()

	// process request

	err := c.Next()

	duration := time.Since(start)

	user, _ := c.Locals("user").(models.User)

	log.Printf(
		"[ACCESS] user=%s method=%s path=%s status=%d duration=%s",
		user.ID,
		c.Method(),
		c.OriginalURL(),
		c.Response().StatusCode(),
		duration,
	)

	return err

}
