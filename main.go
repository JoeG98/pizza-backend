package main

import (
	"github.com/JoeG98/pizza-backend/internal/auth"
	database "github.com/JoeG98/pizza-backend/internal/database"
	"github.com/JoeG98/pizza-backend/internal/orders"
	"github.com/JoeG98/pizza-backend/internal/routes"
	"github.com/JoeG98/pizza-backend/internal/sse"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors" // Add this import
)

func main() {
	auth.LoadKeys()

	app := fiber.New(fiber.Config{
		StreamRequestBody:     true,
		DisableStartupMessage: false,
		// Disable compression for SSE
		EnableTrustedProxyCheck: false,
	})

	app.Use(cors.New(cors.Config{
		AllowOrigins: "*", // Allow all origins (use specific origins in production)
		AllowHeaders: "Origin, Content-Type, Accept, Authorization",
		AllowMethods: "GET, POST, PUT, DELETE, OPTIONS",
	}))

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
