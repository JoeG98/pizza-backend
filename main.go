package main

import (
	database "github.com/JoeG98/pizza-backend/internal/database"
	"github.com/JoeG98/pizza-backend/internal/orders"
	"github.com/JoeG98/pizza-backend/internal/routes"
	"github.com/gofiber/fiber/v2"
)

func main() {
	app := fiber.New()

	// Connect DB
	db := database.Connect()

	orderService := orders.OrderService(db)

	routes.Register(app, orderService)

	app.Listen(":3000")
}
