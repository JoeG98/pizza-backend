package routes

import (
	"github.com/JoeG98/pizza-backend/internal/auth"
	"github.com/JoeG98/pizza-backend/internal/orders"
	"github.com/JoeG98/pizza-backend/internal/sse"
	"github.com/gofiber/fiber/v2"
)

func Register(app *fiber.App, orderService *orders.Service, authService *auth.Service, hub *sse.Hub) {
	// base health endpoint

	app.Get("/api/health", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"status": "OK",
		})
	})

	RegisterOrderRoutes(app, orderService)

	RegisterAuthRoutes(app, authService)
	app.Get("/api/events/orders", hub.Handle())
	app.Get("/api/events/orders/:id", hub.HandleOrderByID())
}
