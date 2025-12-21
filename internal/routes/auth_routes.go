package routes

import (
	"github.com/JoeG98/pizza-backend/internal/auth"
	"github.com/gofiber/fiber/v2"
)

func RegisterAuthRoutes(app *fiber.App, service *auth.Service) {
	app.Post("/signup", func(c *fiber.Ctx) error {
		var input auth.SignupRequest

		if err := c.BodyParser(&input); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "invalid body request",
			})
		}

		if input.Username == "" || input.Password == "" {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "username and password are required!",
			})
		}

		err := service.CreateUser(input.Username, input.Password)

		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": err.Error(),
			})
		}

		return c.Status(fiber.StatusCreated).JSON(fiber.Map{
			"message": "user created",
		})
	})
}
