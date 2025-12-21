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

	app.Post("/login", func(c *fiber.Ctx) error {
		var input auth.LoginRequest

		if err := c.BodyParser(&input); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Invalid body request",
			})
		}

		user, err := service.AuthenticateUser(input.Username, input.Password)

		if err != nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Invalid username or password",
			})
		}

		token, err := auth.GenerateJWT(user.ID)

		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "failed to sign token",
			})
		}

		return c.JSON(fiber.Map{
			"token": token,
		})
	})
}
