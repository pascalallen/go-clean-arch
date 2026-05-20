package websocket

import (
	"bytes"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/oklog/ulid/v2"
)

const (
	writeWait      = 10 * time.Second
	pongWait       = 60 * time.Second
	pingPeriod     = (pongWait * 9) / 10
	maxMessageSize = 512
)

var (
	newline = []byte{'\n'}
	space   = []byte{' '}
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

// Client is a middleman between the websocket connection and the hub.
type Client struct {
	hub     *Hub
	conn    *websocket.Conn
	send    chan []byte
	groupID ulid.ULID
}

func (c *Client) readPump() {
	defer func() {
		c.hub.unregister <- &UnregisterRequest{Client: c, GroupID: c.groupID}
		c.conn.Close()
	}()
	c.conn.SetReadLimit(maxMessageSize)
	c.conn.SetReadDeadline(time.Now().Add(pongWait))
	c.conn.SetPongHandler(func(string) error {
		c.conn.SetReadDeadline(time.Now().Add(pongWait))
		return nil
	})
	for {
		_, message, err := c.conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("error: %v", err)
			}
			break
		}
		message = bytes.TrimSpace(bytes.Replace(message, newline, space, -1))
		c.hub.broadcast <- &Message{GroupID: c.groupID, Content: message}
	}
}

func (c *Client) writePump() {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		c.conn.Close()
	}()
	for {
		select {
		case message, ok := <-c.send:
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if !ok {
				c.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}
			if err := c.conn.WriteMessage(websocket.TextMessage, message); err != nil {
				return
			}
		case <-ticker.C:
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := c.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}

// ServeWs handles WebSocket requests and binds clients to ULID groups.
func ServeWs(hub *Hub, groupId ulid.ULID, c *gin.Context) {
	isWsPath := strings.HasSuffix(strings.TrimSuffix(c.Request.URL.Path, "/"), "/ws")
	hasWsHeader := strings.EqualFold(c.GetHeader("Upgrade"), "websocket") ||
		strings.Contains(strings.ToLower(c.GetHeader("Connection")), "upgrade") ||
		c.GetHeader("Sec-WebSocket-Key") != ""

	if c.Request.Method == http.MethodGet && (isWsPath || hasWsHeader) {
		c.Request.Header.Set("Connection", "Upgrade")
		c.Request.Header.Set("Upgrade", "websocket")
		if c.GetHeader("Sec-WebSocket-Version") == "" {
			c.Request.Header.Set("Sec-WebSocket-Version", "13")
		}
	}

	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Printf("WebSocket upgrade error: %v\n", err)
		return
	}
	client := &Client{
		hub:     hub,
		conn:    conn,
		send:    make(chan []byte, 256),
		groupID: groupId,
	}
	hub.register <- &RegisterRequest{Client: client, GroupID: groupId}

	go client.writePump()
	go client.readPump()
}
