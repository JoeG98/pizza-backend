package routes

import (
	"strings"
	"time"

	"github.com/JoeG98/pizza-backend/internal/auth"
	"github.com/JoeG98/pizza-backend/internal/models"
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

		if input.Role == "" {
			input.Role = "customer"
		}

		if input.Role != "customer" && input.Role != "admin" {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "role must be 'customer' or 'admin'",
			})
		}

		err := service.CreateUser(input.Username, input.Password, input.Role)

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

		refreshToken, err := auth.CreateRefreshToken(user.ID)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "failed to create refresh token"})
		}

		return c.JSON(fiber.Map{
			"token":         token,
			"refresh_token": refreshToken,
		})
	})

	app.Post("/refresh", func(c *fiber.Ctx) error {

		// Extract refresh token from Authorization header
		authHeader := c.Get("Authorization")

		if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "missing or invalid refresh token header",
			})
		}

		refreshToken := strings.TrimPrefix(authHeader, "Bearer ")

		// validate token exists in DB
		var rt models.RefreshToken
		if err := auth.DB.DB.Where("token = ?", refreshToken).First(&rt).Error; err != nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "invalid refresh token",
			})
		}

		// check expiration
		if time.Now().After(rt.ExpiresAt) {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "refresh token expired",
			})
		}

		// load user
		var user models.User
		if err := auth.DB.DB.First(&user, "id = ?", rt.UserID).Error; err != nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "user not found",
			})
		}

		// success → generate NEW access token
		newToken, err := auth.GenerateJWT(user.ID)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "failed to create access token",
			})
		}

		return c.JSON(fiber.Map{
			"token": newToken,
		})
	})

}
