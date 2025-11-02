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

	JoinRoom(roomID string, conn *websocket.Conn)
	LeaveRoom(roomID string, conn *websocket.Conn)
	BoradcastMessageToRoom(roomID string, mt int, msg []byte)
}

type chatService struct {
	clients map[*websocket.Conn]bool
	rooms	map[string]map[*websocket.Conn]bool
	mu      sync.Mutex
}

func NewChatService() *chatService {
	return &chatService{
		clients: 	make(map[*websocket.Conn]bool),
		rooms: 		make(map[string]map[*websocket.Conn]bool),
		mu:      	sync.Mutex{},
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

func (c *chatService) JoinRoom(roomID string, conn *websocket.Conn) {
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.rooms[roomID] == nil {
		c.rooms[roomID] = make(map[*websocket.Conn]bool)
	}

	c.rooms[roomID][conn] = true
}

func (c *chatService) LeaveRoom(roomID string, conn *websocket.Conn) {
	c.mu.Lock()
	defer c.mu.Unlock()

	clients, ok := c.rooms[roomID]
	if ok {
		delete(clients, conn)
		if len(clients) == 0 {
			delete(c.rooms, roomID)
		}
	}
}

func (c *chatService) BoradcastMessageToRoom(roomID string, mt int, msg []byte) {
	c.mu.Lock()
	defer c.mu.Unlock()

	for client := range c.rooms[roomID] {
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
