package websocket

import (
	"sync"

	"github.com/oklog/ulid/v2"
)

// Hub manages active clients, grouped by ULIDs.
type Hub struct {
	clients    map[ulid.ULID]map[*Client]bool
	broadcast  chan *Message
	register   chan *RegisterRequest
	unregister chan *UnregisterRequest
	mu         sync.Mutex
}

// Message represents a message with its associated ULID group.
type Message struct {
	GroupID ulid.ULID
	Content []byte
}

// RegisterRequest represents a new client connection and its ULID group.
type RegisterRequest struct {
	Client  *Client
	GroupID ulid.ULID
}

// UnregisterRequest for removing a client from a ULID group.
type UnregisterRequest struct {
	Client  *Client
	GroupID ulid.ULID
}

func NewHub() *Hub {
	return &Hub{
		clients:    make(map[ulid.ULID]map[*Client]bool),
		broadcast:  make(chan *Message),
		register:   make(chan *RegisterRequest),
		unregister: make(chan *UnregisterRequest),
	}
}

// Run handles registering, unregistering, and broadcasting.
func (h *Hub) Run() {
	for {
		select {
		case req := <-h.register:
			h.mu.Lock()
			if _, exists := h.clients[req.GroupID]; !exists {
				h.clients[req.GroupID] = make(map[*Client]bool)
			}
			h.clients[req.GroupID][req.Client] = true
			h.mu.Unlock()

		case req := <-h.unregister:
			h.mu.Lock()
			if group, exists := h.clients[req.GroupID]; exists {
				if _, found := group[req.Client]; found {
					delete(group, req.Client)
					close(req.Client.send)
					if len(group) == 0 {
						delete(h.clients, req.GroupID)
					}
				}
			}
			h.mu.Unlock()

		case message := <-h.broadcast:
			h.mu.Lock()
			if group, exists := h.clients[message.GroupID]; exists {
				for client := range group {
					select {
					case client.send <- message.Content:
					default:
						close(client.send)
						delete(group, client)
					}
				}
			}
			h.mu.Unlock()
		}
	}
}

// Broadcast sends a message to the given ULID group.
func (h *Hub) Broadcast(message *Message) {
	h.broadcast <- message
}
