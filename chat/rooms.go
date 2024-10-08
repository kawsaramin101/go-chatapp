package chat

import (
	"log"
	"sync"
)

type Room struct {
	hub *Hub

	clients map[*Client]bool

	// Inbound messages from the clients.
	broadcast chan []byte

	// Register requests from the clients.
	register chan *Client

	// Unregister requests from clients.
	unregister chan *Client

	dbRoomID          uint
	dbRoomSecondaryID string

	done chan struct{}
	once sync.Once
}

func (r *Room) RunRoom() {
	defer func() {
		if err := recover(); err != nil {
			log.Printf("Panic in RunRoom: %v", err)
		}
	}()

	for {
		select {
		case client := <-r.register:
			r.clients[client] = true
		case client := <-r.unregister:
			if _, ok := r.clients[client]; ok {
				delete(r.clients, client)
				// close(client.send)
			}
		case message := <-r.broadcast:
			for client := range r.clients {
				select {
				case client.send <- message:
				default:
					close(client.send)
					delete(r.clients, client)
				}
			}
		case <-r.done:
			return

		}

	}
}

func (r *Room) Stop() {

	close(r.done)
	close(r.broadcast)
	close(r.register)
	close(r.unregister)

}
