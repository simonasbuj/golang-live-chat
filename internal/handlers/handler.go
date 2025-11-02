package handlers

import (
	"fmt"
	"golang-live-chat/internal/services"
	"log"

	"github.com/gorilla/websocket"
	"github.com/labstack/echo/v4"
)

type Handler struct {
	Svc services.ChatService
}

func (h *Handler) HandleWebSocket(c echo.Context, upgrader *websocket.Upgrader) error {
	conn, err := upgrader.Upgrade(c.Response(), c.Request(), nil)
	if err != nil {
		log.Println("websocket upgrade error: ", err)

		return fmt.Errorf("failed to upgrade websocket: %w", err)
	}

	defer func() {
		err = conn.Close()
		if err != nil {
			log.Printf("failed to close connection: %s", err)
		}
	}()

	h.Svc.AddClient(conn)
	defer h.Svc.RemoveClient(conn)

	log.Println("client connected")

	for {
		mt, msg, err := conn.ReadMessage()
		if err != nil {
			log.Println("failed to read message: ", err)

			break
		}

		log.Printf("received message: %s", msg)
		h.Svc.BroadcastMessage(mt, msg)
	}

	return nil
}
