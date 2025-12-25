package routes

import (
	"github.com/JoeG98/pizza-backend/internal/auth"
	"github.com/JoeG98/pizza-backend/internal/middleware"
	"github.com/JoeG98/pizza-backend/internal/orders"
	"github.com/gofiber/fiber/v2"
)

func RegisterOrderRoutes(app *fiber.App, service *orders.Service) {
	app.Post("/api/orders", func(c *fiber.Ctx) error {
		var input orders.CreateOrderRequest

		// Parse json

		if err := c.BodyParser(&input); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Invalid request body",
			})
		}

		// Create Order

		order, err := service.CreateOrder(input)

		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": err.Error(),
			})
		}

		// return success
		return c.Status(fiber.StatusCreated).JSON(order)
	})

	app.Get("/api/orders/:id", func(c *fiber.Ctx) error {
		id := c.Params("id")

		order, err := service.GetOrder(id)

		if err != nil {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": "Order not found",
			})
		}

		return c.JSON(order)
	})

	app.Get("/api/orders", func(c *fiber.Ctx) error {
		orders, err := service.GetAllOrders()

		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": err.Error(),
			})
		}

		return c.JSON(orders)
	})

	app.Patch("/api/orders/:id/status", auth.JWTMiddleware, middleware.AccessLog, auth.RequireAdmin, func(c *fiber.Ctx) error {
		id := c.Params("id")

		var input orders.UpdateOrderStatusRequest

		// parse JSON
		if err := c.BodyParser(&input); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "invalid request body",
			})
		}

		// update order
		order, err := service.UpdateOrderStatus(id, input.Status)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": err.Error(),
			})
		}

		return c.JSON(order)
	})

	app.Delete("/api/orders/:id", auth.JWTMiddleware, middleware.AccessLog, auth.RequireAdmin, func(c *fiber.Ctx) error {
		id := c.Params("id")

		err := service.DeleteOrder(id)
		if err != nil {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": "Order Not Found",
			})
		}

		return c.SendStatus(fiber.StatusNoContent)
	})
}
