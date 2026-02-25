package ws

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"sync"
	"time"

	"nhooyr.io/websocket"
)

// Message represents a WebSocket message
type Message struct {
	Type    string          `json:"type"`
	Channel string          `json:"channel,omitempty"`
	Data    json.RawMessage `json:"data,omitempty"`
	From    int64           `json:"from,omitempty"`
}

// Client represents a connected WebSocket client
type Client struct {
	UserID int64
	Conn   *websocket.Conn
	Hub    *Hub
	Send   chan []byte
}

// Hub maintains all active clients and broadcasts messages
type Hub struct {
	clients    map[int64]*Client
	broadcast  chan []byte
	register   chan *Client
	unregister chan *Client
	mu         sync.RWMutex
}

func NewHub() *Hub {
	return &Hub{
		clients:    make(map[int64]*Client),
		broadcast:  make(chan []byte, 256),
		register:   make(chan *Client),
		unregister: make(chan *Client),
	}
}

func (h *Hub) Run(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			return
		case client := <-h.register:
			h.mu.Lock()
			h.clients[client.UserID] = client
			h.mu.Unlock()
			log.Printf("[WS] Client connected: user=%d (total=%d)", client.UserID, len(h.clients))
		case client := <-h.unregister:
			h.mu.Lock()
			if _, ok := h.clients[client.UserID]; ok {
				delete(h.clients, client.UserID)
				close(client.Send)
			}
			h.mu.Unlock()
			log.Printf("[WS] Client disconnected: user=%d (total=%d)", client.UserID, len(h.clients))
		case msg := <-h.broadcast:
			h.mu.RLock()
			for _, client := range h.clients {
				select {
				case client.Send <- msg:
				default:
					// drop message if client buffer full
				}
			}
			h.mu.RUnlock()
		}
	}
}

// SendToUser sends a message to a specific user
func (h *Hub) SendToUser(userID int64, msg Message) {
	h.mu.RLock()
	client, ok := h.clients[userID]
	h.mu.RUnlock()
	if !ok {
		return
	}

	data, err := json.Marshal(msg)
	if err != nil {
		return
	}

	select {
	case client.Send <- data:
	default:
	}
}

// Broadcast sends a message to all connected clients
func (h *Hub) Broadcast(msg Message) {
	data, err := json.Marshal(msg)
	if err != nil {
		return
	}
	h.broadcast <- data
}

// HandleWS handles WebSocket upgrade and connection
func (h *Hub) HandleWS(w http.ResponseWriter, r *http.Request, userID int64) {
	conn, err := websocket.Accept(w, r, &websocket.AcceptOptions{
		InsecureSkipVerify: true, // allow all origins in dev
	})
	if err != nil {
		log.Printf("[WS] Accept error: %v", err)
		return
	}

	client := &Client{
		UserID: userID,
		Conn:   conn,
		Hub:    h,
		Send:   make(chan []byte, 256),
	}

	h.register <- client

	go client.writePump(r.Context())
	go client.readPump(r.Context())
}

func (c *Client) readPump(ctx context.Context) {
	defer func() {
		c.Hub.unregister <- c
		c.Conn.Close(websocket.StatusNormalClosure, "")
	}()

	for {
		_, data, err := c.Conn.Read(ctx)
		if err != nil {
			return
		}

		var msg Message
		if err := json.Unmarshal(data, &msg); err != nil {
			continue
		}
		msg.From = c.UserID

		// Route message based on type
		switch msg.Type {
		case "broadcast":
			c.Hub.Broadcast(msg)
		case "ping":
			resp, _ := json.Marshal(Message{Type: "pong"})
			select {
			case c.Send <- resp:
			default:
			}
		default:
			// Handle other message types as needed
		}
	}
}

func (c *Client) writePump(ctx context.Context) {
	ticker := time.NewTicker(30 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case msg, ok := <-c.Send:
			if !ok {
				c.Conn.Close(websocket.StatusNormalClosure, "")
				return
			}
			if err := c.Conn.Write(ctx, websocket.MessageText, msg); err != nil {
				return
			}
		case <-ticker.C:
			if err := c.Conn.Ping(ctx); err != nil {
				return
			}
		}
	}
}
