package sse

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
	for {
		select {

		// add new client
		case client := <-h.Register:
			h.Clients[client] = true

		// remove disconnected client
		case client := <-h.Unregister:
			if _, ok := h.Clients[client]; ok {
				delete(h.Clients, client)
			}

		// incoming SSE message
		case msg := <-h.Broadcast:
			for client := range h.Clients {
				select {
				case client <- msg:
				default:
					// prevent blocked writer
					delete(h.Clients, client)
				}
			}
		}
	}
}
