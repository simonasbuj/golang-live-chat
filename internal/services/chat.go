package services

import (
	"sync"

	"github.com/gorilla/websocket"
)


type ChatService interface {
	AddClient(conn *websocket.Conn) error
}

type chatService struct {
	clients map[*websocket.Conn]bool
	mu sync.Mutex
}

func NewChatService() *chatService {
	return &chatService{
		clients: make(map[*websocket.Conn]bool),
	}
}

func (c *chatService) AddClient(conn *websocket.Conn) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.clients[conn] = true

	return nil
}
