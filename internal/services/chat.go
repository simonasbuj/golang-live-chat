package services

import (
	"log"
	"sync"

	"github.com/gorilla/websocket"
)

type ChatService interface {
	AddClient(conn *websocket.Conn)
	RemoveClient(conn *websocket.Conn)
	BroadcastMessage(mt int, msg []byte)
}

type chatService struct {
	clients map[*websocket.Conn]bool
	mu      sync.Mutex
}

func NewChatService() *chatService {
	return &chatService{
		clients: make(map[*websocket.Conn]bool),
		mu:      sync.Mutex{},
	}
}

func (c *chatService) AddClient(conn *websocket.Conn) {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.clients[conn] = true
}

func (c *chatService) RemoveClient(conn *websocket.Conn) {
	c.mu.Lock()
	defer c.mu.Unlock()

	delete(c.clients, conn)
}

func (c *chatService) BroadcastMessage(mt int, msg []byte) {
	c.mu.Lock()
	defer c.mu.Unlock()

	for client := range c.clients {
		err := client.WriteMessage(mt, msg)
		if err != nil {
			log.Println("failed to write message to a client:", err)

			closeErr := client.Close()
			if closeErr != nil {
				log.Println("failed to close client")
			}

			delete(c.clients, client)
		}
	}
}
