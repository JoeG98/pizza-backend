package main

import (
	"fmt"
	"time"

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

	orderService := orders.OrderService(db)
	authService := auth.AuthService(db)

	hub := sse.NewHub()
	go hub.Run()

	routes.Register(app, orderService, authService, hub)

	go func() {
		for {
			time.Sleep(2 * time.Second)
			fmt.Println("🔥 broadcasting SSE tick")
			hub.Broadcast <- "tick from SSE"
		}
	}()

	app.Listen(":3000")
}
