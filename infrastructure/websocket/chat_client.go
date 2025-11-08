package websocket

import (
	"encoding/json"
	"log"
	"time"

	"github.com/gofiber/websocket/v2"
)

const (
	// Time allowed to write a message to the peer
	writeWait = 10 * time.Second

	// Time allowed to read the next pong message from the peer
	pongWait = 60 * time.Second

	// Send pings to peer with this period (must be less than pongWait)
	pingPeriod = (pongWait * 9) / 10

	// Maximum message size allowed from peer
	maxMessageSize = 512 * 1024 // 512 KB
)

// ReadPump pumps messages from the WebSocket connection to the hub
func (c *ChatClient) ReadPump() {
	defer func() {
		c.Hub.unregister <- c
		c.Conn.Close()
	}()

	c.Conn.SetReadDeadline(time.Now().Add(pongWait))
	c.Conn.SetPongHandler(func(string) error {
		c.Conn.SetReadDeadline(time.Now().Add(pongWait))
		return nil
	})

	for {
		_, messageData, err := c.Conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("WebSocket error: %v", err)
			}
			break
		}

		// Parse message
		var message ChatMessage
		if err := json.Unmarshal(messageData, &message); err != nil {
			log.Printf("Failed to unmarshal message: %v", err)
			c.sendError("invalid_message", "Invalid message format")
			continue
		}

		// Update online status in Redis on any activity
		_ = c.Hub.redisService.SetUserOnline(c.Hub.ctx, c.UserID)

		// Route message to appropriate handler
		c.Hub.routeMessage(c, &message)
	}
}

// WritePump pumps messages from the hub to the WebSocket connection
func (c *ChatClient) WritePump() {
	// Check if connection is valid
	if c.Conn == nil {
		log.Printf("❌ WritePump: Connection is nil for user %s", c.UserID)
		if c.Ready != nil {
			close(c.Ready)
		}
		return
	}

	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		if c.Conn != nil {
			c.Conn.Close()
		}
	}()

	// Signal that WritePump is ready to receive messages
	if c.Ready != nil {
		close(c.Ready)
	}

	for {
		select {
		case message, ok := <-c.Send:
			if c.Conn == nil {
				log.Printf("❌ WritePump: Connection became nil for user %s", c.UserID)
				return
			}

			c.Conn.SetWriteDeadline(time.Now().Add(writeWait))
			if !ok {
				// The hub closed the channel
				c.Conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			// Send message
			if err := c.Conn.WriteMessage(websocket.TextMessage, message); err != nil {
				log.Printf("Failed to write message: %v", err)
				return
			}

		case <-ticker.C:
			if c.Conn == nil {
				log.Printf("❌ WritePump: Connection became nil for user %s (ping)", c.UserID)
				return
			}

			// Send ping
			c.Conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := c.Conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				log.Printf("Failed to send ping: %v", err)
				return
			}

			// Update online status in Redis on ping
			_ = c.Hub.redisService.SetUserOnline(c.Hub.ctx, c.UserID)
		}
	}
}

// sendError sends an error message to the client
func (c *ChatClient) sendError(code string, message string) {
	errorMsg := &ChatMessage{
		Type: "error",
		Error: &ChatError{
			Code:    code,
			Message: message,
		},
	}

	data, _ := json.Marshal(errorMsg)
	select {
	case c.Send <- data:
	default:
		// Send buffer full, skip error message
	}
}
