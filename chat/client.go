package chat

import (
	"bytes"
	auth_views "chatapp/auth/views"
	db "chatapp/db"
	"errors"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
	"gorm.io/gorm"
)

const (
	// Time allowed to write a message to the peer.
	writeWait = 10 * time.Second

	// Time allowed to read the next pong message from the peer.
	pongWait = 60 * time.Second

	// Send pings to peer with this period. Must be less than pongWait.
	pingPeriod = (pongWait * 9) / 10

	// Maximum message size allowed from peer.
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
	hub *Hub

	rooms map[*Room]bool

	// The websocket connection.
	conn *websocket.Conn

	// Buffered channel of outbound messages.
	send chan []byte

	dbUserID          uint
	dbUserUsername    string
	dbUserSecondaryID string
}

type Room struct {
	hub *Hub

	clients map[*Client]bool

	dbRoomID          uint
	dbRoomSecondaryID string
}

func (c *Client) waitForAuth() {

	c.conn.SetReadLimit(maxMessageSize)
	c.conn.SetReadDeadline(time.Now().Add(pongWait))
	c.conn.SetPongHandler(func(string) error { c.conn.SetReadDeadline(time.Now().Add(pongWait)); return nil })

	for {
		_, authToken, err := c.conn.ReadMessage()
		if err != nil {
			log.Printf("Error reading auth token: %v", err)
			c.conn.Close()
			return
		}

		authToken = bytes.TrimSpace(bytes.Replace(authToken, newline, space, -1))

		_, userSecondaryId, err := auth_views.ValidateToken(string(authToken))

		if err != nil {
			log.Printf("Error validating auth token: %v", err)
			c.conn.Close()
			return
		}

		var user db.User

		err = getUserAndRoomInfo(userSecondaryId, &user)

		if err != nil {
			log.Printf("Error finding the user info: %v", err)
			c.conn.Close()
			return
		}

		c.dbUserID = user.ID
		c.dbUserUsername = user.Username
		c.dbUserSecondaryID = user.SecondaryID

		roomMap := make(map[uint]*Room)
		for room := range c.hub.rooms {
			roomMap[room.dbRoomID] = room
		}

		for _, roomToCheck := range user.Chats {
			if existingRoom, exists := roomMap[roomToCheck.ID]; !exists {
				c.rooms[existingRoom] = true
				existingRoom.clients[c] = true

			} else {
				newRoom := Room{hub: c.hub, clients: make(map[*Client]bool), dbRoomID: roomToCheck.ID, dbRoomSecondaryID: roomToCheck.SecondaryID}
				newRoom.clients[c] = true
				c.rooms[&newRoom] = true
			}
		}

		go c.writePump()
		go c.readPump()

		return

	}
}

type Message struct {
	Action string `json:"action"`
	Value  string `json:"value"`
}

// readPump pumps messages from the websocket connection to the hub.
//
// The application runs readPump in a per-connection goroutine. The application
// ensures that there is at most one reader on a connection by executing all
// reads from this goroutine.
func (c *Client) readPump() {
	defer func() {
		c.hub.unregister <- c
		c.conn.Close()
	}()
	c.conn.SetReadLimit(maxMessageSize)
	c.conn.SetReadDeadline(time.Now().Add(pongWait))
	c.conn.SetPongHandler(func(string) error { c.conn.SetReadDeadline(time.Now().Add(pongWait)); return nil })

	for {
		_, message, err := c.conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("error: %v", err)
			}
			break
		}
		message = bytes.TrimSpace(bytes.Replace(message, newline, space, -1))
		log.Printf("Received: %s", message)
		c.hub.broadcast <- message
	}
}

// writePump pumps messages from the hub to the websocket connection.
//
// A goroutine running writePump is started for each connection. The
// application ensures that there is at most one writer to a connection by
// executing all writes from this goroutine.
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
				// The hub closed the channel.
				c.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			w, err := c.conn.NextWriter(websocket.TextMessage)
			if err != nil {
				return
			}
			w.Write(message)

			// Add queued chat messages to the current websocket message.
			n := len(c.send)
			for i := 0; i < n; i++ {
				msg := <-c.send
				w.Write(newline)
				w.Write(msg)
				log.Printf("Sent: %s", msg)
			}

			if err := w.Close(); err != nil {
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

// ServeWs handles websocket requests from the peer.
func ServeWs(hub *Hub, w http.ResponseWriter, r *http.Request) {

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}
	client := &Client{hub: hub, conn: conn, send: make(chan []byte, 256)}
	client.hub.register <- client

	// Allow collection of memory referenced by the caller by doing all work in
	// new goroutines.

	go client.waitForAuth()

}

func getUserAndRoomInfo(userSecondaryId string, user *db.User) error {
	if err := db.DB.Where("secondary_id = ?", userSecondaryId).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return err // No user found, return the error or handle it
		}
		return err // Other errors like DB connection failure, etc.
	}

	// If user is found, preload Chats
	if err := db.DB.Preload("Chats").First(&user, user.ID).Error; err != nil {
		return err
	}

	return nil
}
