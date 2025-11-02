package handlers

import (
	"golang-live-chat/internal/services"
	"log"

	"github.com/gorilla/websocket"
	"github.com/labstack/echo/v4"
)

var clients = make(map[*websocket.Conn]bool)

type Handler struct {
	Svc services.ChatService
}


func (h *Handler) HandleWebSocket(c echo.Context, upgrader *websocket.Upgrader) error {
	conn, err := upgrader.Upgrade(c.Response(), c.Request(), nil)
	if err != nil {
		log.Println("websocket upgrade error: ", err)
		return err
	}
	defer conn.Close()

	err = h.Svc.AddClient(conn)
	if err != nil {
		log.Printf("failed to connect client: %s", err)
	}
	// defer delete(clients, conn)

	log.Println("client connected")

	for {
		mt, msg, err := conn.ReadMessage()
		if err != nil {
			log.Println("failed to read message: ", err)
			break
		}

		log.Printf("received message: %s", msg)

		for client := range clients {
			if err := client.WriteMessage(mt, msg); err != nil {
				log.Println("write error:", err)
				client.Close()
				delete(clients, client)
			}
		}
	}

	return nil
}
