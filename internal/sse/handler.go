package sse

import (
	"bufio"
	"fmt"

	"github.com/gofiber/fiber/v2"
)

func (h *Hub) Handle() fiber.Handler {
	return func(c *fiber.Ctx) error {

		fmt.Println("🚨 SSE handler HIT")

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

			// 🔥 CRITICAL: initial SSE frame
			w.WriteString(": connected\n\n")
			w.Flush()

			for msg := range client {
				data := "data: " + msg + "\n\n"

				if _, err := w.WriteString(data); err != nil {
					fmt.Println("write failed:", err)
					return
				}

				if err := w.Flush(); err != nil {
					fmt.Println("flush failed:", err)
					return
				}
			}
		})

		return nil
	}
}
