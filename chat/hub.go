package chat

// Hub maintains the set of active clients and broadcasts messages to the
// clients.
type Hub struct {
	// Registered clients.
	clients     map[*Client]bool
	rooms       map[*Room]bool
	activeRooms map[*Room]bool

	// Inbound messages from the clients.
	broadcast       chan []byte
	broadcastToRoom chan RoomMessage

	// Register requests from the clients.
	registerClient chan *Client

	// Unregister requests from clients.
	unregisterClient chan *Client

	registerRoom chan *Room

	registerToRoom   chan RegisterToRoomInfo
	unregisterToRoom chan RegisterToRoomInfo
	// Unregister requests from clients.
	unregisterRoom chan *Room
}

func NewHub() *Hub {
	return &Hub{
		broadcast:        make(chan []byte),
		broadcastToRoom:  make(chan RoomMessage),
		registerClient:   make(chan *Client),
		unregisterClient: make(chan *Client),
		registerRoom:     make(chan *Room),
		unregisterRoom:   make(chan *Room),
		registerToRoom:   make(chan RegisterToRoomInfo),
		clients:          make(map[*Client]bool),
		rooms:            make(map[*Room]bool),
		activeRooms:      make(map[*Room]bool),
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
		case roomMessage := <-h.broadcastToRoom:

			for client := range roomMessage.room.clients {
				select {
				case client.send <- roomMessage.message:
				default:
					close(client.send)
					delete(roomMessage.room.clients, client)
				}
			}
		case info := <-h.registerToRoom:
			info.room.clients[info.client] = true
		case info := <-h.unregisterToRoom:
			if _, ok := info.room.clients[info.client]; ok {
				delete(info.room.clients, info.client)
			}
			if len(info.room.clients) == 0 {
				h.unregisterRoom <- info.room
			}
		case message := <-h.broadcast:
			for client := range h.clients {
				select {
				case client.send <- message:
				default:
					close(client.send)
					delete(h.clients, client)
				}
			}
		}
	}
}
