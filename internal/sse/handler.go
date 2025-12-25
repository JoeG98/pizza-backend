// sse/handler.go
package sse

import (
	"bufio"
	"fmt"
	"time"

	"github.com/gofiber/fiber/v2"
)

func (h *Hub) Handle() fiber.Handler {
	return func(c *fiber.Ctx) error {
		fmt.Println("🚨 SSE handler HIT")

		// if _, err := authenticateAdmin(c); err != nil {
		// 	return c.Status(fiber.StatusUnauthorized).SendString("unauthorized")
		// }

		// Set SSE headers BEFORE SetBodyStreamWriter
		c.Set("Content-Type", "text/event-stream")
		c.Set("Cache-Control", "no-cache")
		c.Set("Connection", "keep-alive")
		c.Set("Access-Control-Allow-Origin", "*")
		c.Set("X-Accel-Buffering", "no")

		// Disable response compression
		c.Context().Response.Header.Del("Content-Encoding")

		// Create client channel with buffer
		client := make(Client, 10)
		h.Register <- client

		// Cleanup flag
		disconnected := make(chan struct{})

		// Set stream writer
		c.Context().SetBodyStreamWriter(func(w *bufio.Writer) {
			fmt.Println("✅ Stream writer started")

			defer func() {
				fmt.Println("🔌 Stream writer closing")
				h.Unregister <- client
				close(disconnected)
			}()

			// Send a comment first to establish connection
			fmt.Fprintf(w, ": connected\n\n")
			w.Flush()

			// Then send welcome data event
			welcomeMsg := `{"type":"connected","data":{"message":"Connected to order stream"}}`
			fmt.Fprintf(w, "data: %s\n\n", welcomeMsg)

			// Flush multiple times to force through any buffers
			w.Flush()
			w.Flush()

			fmt.Println("✅ Welcome message sent and flushed")

			// Keep-alive ticker - reduced to 5 seconds for faster connection establishment
			ticker := time.NewTicker(5 * time.Second)
			defer ticker.Stop()

			// Read messages from client channel
			for {
				select {
				case msg, ok := <-client:
					if !ok {
						fmt.Println("❌ Client channel closed")
						return
					}

					fmt.Printf("📤 Sending message: %s\n", msg)

					// Write SSE message
					fmt.Fprintf(w, "data: %s\n\n", msg)

					if err := w.Flush(); err != nil {
						fmt.Printf("❌ Flush failed: %v\n", err)
						return
					}

				case <-ticker.C:
					// Send keep-alive comment with timestamp
					fmt.Fprintf(w, ": keepalive %d\n\n", time.Now().Unix())
					if err := w.Flush(); err != nil {
						fmt.Println("❌ Keep-alive flush failed, client disconnected")
						return
					}
					fmt.Println("💓 Keepalive sent")
				}
			}
		})

		// Wait for disconnection before returning
		<-disconnected
		return nil
	}
}
