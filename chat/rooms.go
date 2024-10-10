package chat

import (
	"fmt"
	"log"
	"sync"
)

type Room struct {
	hub *Hub

	clients map[*Client]bool

	// Inbound messages from the clients.
	broadcast chan []byte

	// Register requests from the clients.
	registerClient chan *Client

	// Unregister requests from clients.
	unregisterClient chan *Client

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
		case client := <-r.registerClient:
			r.clients[client] = true
		case client := <-r.unregisterClient:
			if _, ok := r.clients[client]; ok {
				delete(r.clients, client)
			}
			if len(r.clients) == 0 {
				r.Stop()
				r.hub.unregisterRoom <- r
			}

		case message := <-r.broadcast:
			for client := range r.clients {
				fmt.Println("send to ", client.dbUser.Username)
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
	close(r.registerClient)
	close(r.unregisterClient)
}

func NewRoom(hub *Hub, dbRoomID uint, dbRoomSecondaryID string) *Room {
	return &Room{
		hub:               hub,
		dbRoomID:          dbRoomID,
		dbRoomSecondaryID: dbRoomSecondaryID,
		broadcast:         make(chan []byte),
		registerClient:    make(chan *Client),
		unregisterClient:  make(chan *Client),
		clients:           make(map[*Client]bool),
		done:              make(chan struct{}),
	}
}
