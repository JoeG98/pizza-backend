package routes

import (
	"github.com/JoeG98/pizza-backend/internal/auth"
	"github.com/JoeG98/pizza-backend/internal/orders"
	"github.com/gofiber/fiber/v2"
)

func Register(app *fiber.App, orderService *orders.Service, authService *auth.Service) {
	// base health endpoint

	app.Get("/health", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"status": "OK",
		})
	})

	RegisterOrderRoutes(app, orderService)

	RegisterAuthRoutes(app, authService)
}
