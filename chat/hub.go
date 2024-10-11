package chat

import "fmt"

// Hub maintains the set of active clients and broadcasts messages to the
// clients.
type Hub struct {
	// Registered clients.
	clients map[*Client]bool
	rooms   map[*Room]bool

	// Register requests from the clients.
	registerClient chan *Client
	// Unregister requests from the clients.
	unregisterClient chan *Client

	// Register a room
	registerRoom chan *Room
	// Unregister a room
	unregisterRoom chan *Room
}

func NewHub() *Hub {
	return &Hub{
		clients: make(map[*Client]bool),
		rooms:   make(map[*Room]bool),

		registerClient:   make(chan *Client),
		unregisterClient: make(chan *Client),

		registerRoom:   make(chan *Room),
		unregisterRoom: make(chan *Room),
	}
}

type RoomMessage struct {
	room    *Room
	message []byte
}

type RegisterToRoomInfo struct {
	room   *Room
	client *Client
}

func (h *Hub) Run() {
	for {
		select {
		case client := <-h.registerClient:
			h.clients[client] = true
		case client := <-h.unregisterClient:
			if _, ok := h.clients[client]; ok {
				delete(h.clients, client)
				close(client.send)
			}
		case room := <-h.registerRoom:
			h.rooms[room] = true
		case room := <-h.unregisterRoom:
			if _, ok := h.rooms[room]; ok {
				delete(h.rooms, room)
				fmt.Println("Room removed")
			}

			// case message := <-h.broadcast:
			// 	for client := range h.clients {
			// 		select {
			// 		case client.send <- message:
			// 		default:
			// 			close(client.send)
			// 			delete(h.clients, client)
			// 		}
			// 	}
		}
	}
}
