package main

import (
	"github.com/JoeG98/pizza-backend/internal/auth"
	database "github.com/JoeG98/pizza-backend/internal/database"
	"github.com/JoeG98/pizza-backend/internal/orders"
	"github.com/JoeG98/pizza-backend/internal/routes"
	"github.com/JoeG98/pizza-backend/internal/sse"
	"github.com/gofiber/fiber/v2"
)

func main() {
	auth.LoadKeys()

	app := fiber.New()

	// Connect DB
	db := database.Connect()

	auth.DB = db

	hub := sse.NewHub()
	go hub.Run()

	orderService := orders.OrderService(db, hub)
	authService := auth.AuthService(db)

	routes.Register(app, orderService, authService, hub)

	app.Listen(":3000")
}
