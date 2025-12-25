package sse

import (
	"bufio"
	"encoding/json"

	"github.com/gofiber/fiber/v2"
)

func (h *Hub) HandleOrderByID() fiber.Handler {
	return func(c *fiber.Ctx) error {

		// 1️⃣ Authenticate user
		_, err := authenticateUser(c)
		if err != nil {
			return c.Status(fiber.StatusUnauthorized).SendString("unauthorized")
		}

		// 2️⃣ Get order ID from URL
		orderID := c.Params("id")
		if orderID == "" {
			return c.Status(fiber.StatusBadRequest).SendString("missing order id")
		}

		// 3️⃣ SSE headers
		c.Set("Content-Type", "text/event-stream")
		c.Set("Cache-Control", "no-cache")
		c.Set("Connection", "keep-alive")
		c.Set("Access-Control-Allow-Origin", "*")
		c.Set("Access-Control-Allow-Private-Network", "true")

		client := make(Client)
		h.Register <- client

		defer func() {
			h.Unregister <- client
		}()

		c.Context().SetBodyStreamWriter(func(w *bufio.Writer) {

			// initial frame
			w.WriteString(": connected\n\n")
			w.Flush()

			for msg := range client {
				var event Event
				if err := json.Unmarshal([]byte(msg), &event); err != nil {
					continue
				}

				// only forward relevant order events
				order, ok := event.Data.(map[string]interface{})
				if !ok {
					continue
				}

				if order["id"] == orderID {
					w.WriteString("data: " + msg + "\n\n")
					w.Flush()
				}
			}
		})

		return nil
	}
}
