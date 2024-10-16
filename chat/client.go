package chat

import (
	"bytes"
	auth_views "chatapp/auth"
	db "chatapp/db"
	"encoding/json"
	"errors"
	"fmt"
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

	rooms       map[*Room]bool
	activeRooms map[*Room]bool

	// The websocket connection.
	conn *websocket.Conn

	// Buffered channel of outbound messages.
	send        chan []byte
	sendActions chan string

	dbUser db.User
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

		c.dbUser = user

		roomMap := make(map[uint]*Room)
		for room := range c.hub.rooms {
			roomMap[room.dbRoomID] = room
		}

		for _, roomToCheck := range user.Chats {
			if existingRoom, exists := roomMap[roomToCheck.ID]; exists {
				c.rooms[existingRoom] = true

				existingRoom.RegisterClient(c)
			} else {
				newRoom := NewRoom(c.hub, roomToCheck.ID, roomToCheck.SecondaryID)

				c.hub.registerRoom <- newRoom

				c.rooms[newRoom] = true
				newRoom.RegisterClient(c)
			}
		}

		for i, chat := range c.dbUser.Chats {
			// Preload the users for the current chat
			if err := db.DB.Preload("Users").Find(&c.dbUser.Chats[i], chat.ID).Error; err != nil {
				// Handle error if loading users fails for a chat
				response := map[string]string{
					"action":  "ERROR_LOADING_CHAT_USERS",
					"message": "Error loading users for chat",
				}
				jsonData, _ := json.Marshal(response)
				c.send <- jsonData
				return
			}
		}

		initialData := struct {
			Action string `json:"action"`
			Data   struct {
				Chats []db.Chat `json:"chats"`
			} `json:"data"`
		}{
			Action: "INITIAL_DATA",
			Data: struct {
				Chats []db.Chat `json:"chats"`
			}{
				Chats: user.Chats,
			},
		}

		connectionRequests := struct {
			Action string `json:"action"`
			Data   struct {
				ConnectionRequests []db.ConnectionRequest `json:"connectionRequests"`
			} `json:"data"`
		}{
			Action: "CONNECTION_REQUESTS",
			Data: struct {
				ConnectionRequests []db.ConnectionRequest `json:"connectionRequests"`
			}{
				ConnectionRequests: user.ReceivedRequests,
			},
		}

		go c.writePump()
		go c.readPump()
		jsonData, err := json.Marshal(initialData)
		connectionRequestsJsonData, err := json.Marshal(connectionRequests)
		c.send <- jsonData
		c.send <- connectionRequestsJsonData
	}
}

// type IncomingData struct {
//     Action string      `json:"ACTION"`
//     Data   json.RawMessage `json:"DATA"` // Use RawMessage for lazy decoding
// }

type Message struct {
	Action string          `json:"action"`
	Data   json.RawMessage `json:"data"` // Generic to allow varying structures
}

// readPump pumps messages from the websocket connection to the hub.
//
// The application runs readPump in a per-connection goroutine. The application
// ensures that there is at most one reader on a connection by executing all
// reads from this goroutine.
func (c *Client) readPump() {
	defer func() {
		c.hub.unregisterClient <- c

		for room := range c.rooms {
			delete(c.rooms, room)
			room.UnregisterClient(c)
		}

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
			log.Printf("%s", err)
			break
		}
		message = bytes.TrimSpace(bytes.Replace(message, newline, space, -1))

		// Define an instance of the Message struct
		var msg Message
		// Unmarshal the JSON data into the msg struct
		err = json.Unmarshal(message, &msg)

		if err != nil {
			log.Printf("error unmarshaling JSON: %v", err)
			continue
		}

		switch msg.Action {
		case "CREATECHAT":
			CreateChat(&msg, c)

		case "MESSAGE":
			HandleMessage(&msg, c)

		case "CHECK_IF_USER_EXIST":
			CheckIfUserExist(&msg, c)

		}
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
			// n := len(c.send)
			// for i := 0; i < n; i++ {
			// 	msg := <-c.send
			// 	w.Write(newline)
			// 	w.Write(msg)
			// }

			if err := w.Close(); err != nil {
				fmt.Println(err)
				return
			}
		case <-ticker.C:
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := c.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		case action := <-c.sendActions:
			switch action {
			case "SEND_CONNECTION_REQUESTS":
				db.DB.Preload("ReceivedRequests").
					Preload("ReceivedRequests.SendBy").
					Preload("ReceivedRequests.Chat").
					Preload("ReceivedRequests").
					Find(&c.dbUser, c.dbUser.ID)

				connectionRequests := struct {
					Action string `json:"action"`
					Data   struct {
						ConnectionRequests []db.ConnectionRequest `json:"connectionRequests"`
					} `json:"data"`
				}{
					Action: "CONNECTION_REQUESTS",
					Data: struct {
						ConnectionRequests []db.ConnectionRequest `json:"connectionRequests"`
					}{
						ConnectionRequests: c.dbUser.ReceivedRequests,
					},
				}
				connectionRequestsJsonData, _ := json.Marshal(connectionRequests)
				c.send <- connectionRequestsJsonData
			}
			fmt.Println(action)
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
	client := &Client{
		hub:         hub,
		conn:        conn,
		send:        make(chan []byte, 256),
		sendActions: make(chan string, 256),
		rooms:       make(map[*Room]bool),
	}
	client.hub.registerClient <- client

	go client.waitForAuth()
}

func getUserAndRoomInfo(userSecondaryId string, user *db.User) error {
	if err := db.DB.Where("secondary_id = ?", userSecondaryId).First(user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return err // No user found, return the error or handle it
		}
		return err // Other errors like DB connection failure, etc.
	}

	// If user is found, preload Chats
	if err := db.DB.Preload("Chats", func(db *gorm.DB) *gorm.DB {
		return db.Order("id desc")
	}).First(&user, user.ID).Error; err != nil {
		return err
	}

	err := db.DB.Preload("ReceivedRequests").
		Preload("ReceivedRequests.SendBy").
		Preload("ReceivedRequests.Chat").
		Preload("ReceivedRequests").
		Find(&user, user.ID)

	if err != nil {
		return err.Error
	}

	return nil
}
