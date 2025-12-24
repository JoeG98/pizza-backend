package routes

import (
	"github.com/JoeG98/pizza-backend/internal/auth"
	"github.com/JoeG98/pizza-backend/internal/orders"
	"github.com/JoeG98/pizza-backend/internal/sse"
	"github.com/gofiber/fiber/v2"
)

func Register(app *fiber.App, orderService *orders.Service, authService *auth.Service, hub *sse.Hub) {
	// base health endpoint

	app.Get("/health", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"status": "OK",
		})
	})

	app.Get("/events/orders", hub.Handle())

	// app.Get("/event/orders", func(c *fiber.Ctx) error {
	// 	return c.JSON(fiber.Map{
	// 		"status": "OK Good",
	// 	})
	// })

	RegisterOrderRoutes(app, orderService)

	RegisterAuthRoutes(app, authService)
}
