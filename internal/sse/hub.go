// sse/hub.go
package sse

import "fmt"

// Client channel type
type Client chan string

// Hub manages all SSE clients
type Hub struct {
	Clients    map[Client]bool
	Register   chan Client
	Unregister chan Client
	Broadcast  chan string
}

// NewHub creates a Hub
func NewHub() *Hub {
	return &Hub{
		Clients:    make(map[Client]bool),
		Register:   make(chan Client),
		Unregister: make(chan Client),
		Broadcast:  make(chan string),
	}
}

// Run processes client registration, removal, and messages
func (h *Hub) Run() {
	fmt.Println("🚀 Hub started")
	for {
		select {
		// Add new client
		case client := <-h.Register:
			h.Clients[client] = true
			fmt.Printf("✅ Client registered. Total clients: %d\n", len(h.Clients))

		// Remove disconnected client
		case client := <-h.Unregister:
			if _, ok := h.Clients[client]; ok {
				delete(h.Clients, client)
				close(client) // Close the channel
				fmt.Printf("❌ Client unregistered. Total clients: %d\n", len(h.Clients))
			}

		// Broadcast message to all clients
		case msg := <-h.Broadcast:
			fmt.Printf("📢 Broadcasting to %d clients: %s\n", len(h.Clients), msg)

			for client := range h.Clients {
				select {
				case client <- msg:
					// Message sent successfully
				default:
					// Client channel is full or blocked, remove it
					fmt.Println("⚠️ Client channel blocked, removing")
					delete(h.Clients, client)
					close(client)
				}
			}
		}
	}
}
