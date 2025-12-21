package main

import (
	"github.com/JoeG98/pizza-backend/internal/auth"
	database "github.com/JoeG98/pizza-backend/internal/database"
	"github.com/JoeG98/pizza-backend/internal/orders"
	"github.com/JoeG98/pizza-backend/internal/routes"
	"github.com/gofiber/fiber/v2"
)

func main() {
	auth.LoadKeys()

	app := fiber.New()

	// Connect DB
	db := database.Connect()

	orderService := orders.OrderService(db)
	authService := auth.AuthService(db)

	routes.Register(app, orderService, authService)

	app.Listen(":3000")
}
