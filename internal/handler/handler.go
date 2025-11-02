package handler

import (
	"fmt"
	"golang-live-chat/internal/services"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
	"github.com/labstack/echo/v4"
)

type Handler struct {
	svc      services.ChatService
	upgrader *websocket.Upgrader
}

func New(svc services.ChatService, upgrader *websocket.Upgrader) *Handler {
	return &Handler{
		svc:      svc,
		upgrader: upgrader,
	}
}

func (h *Handler) HandleGlobalChat(c echo.Context) error {
	conn, err := h.upgrader.Upgrade(c.Response(), c.Request(), nil)
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

	h.svc.AddClient(conn)
	defer h.svc.RemoveClient(conn)

	log.Println("client connected")

	for {
		mt, msg, err := conn.ReadMessage()
		if err != nil {
			log.Println("failed to read message: ", err)

			break
		}

		log.Printf("received message: %s", msg)
		h.svc.BroadcastMessage(mt, msg)
	}

	return nil
}

func (h *Handler) HandleRoomChat(c echo.Context) error {
	roomID := c.QueryParam("room")
	if roomID == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{ //nolint:wrapcheck
			"error": "missing required query parameter: room",
		})
	}

	conn, err := h.upgrader.Upgrade(c.Response(), c.Request(), nil)
	if err != nil {
		log.Println("websocket upgrade error: ", err)

		return fmt.Errorf("failed to upgrade websocket: %w", err)
	}

	defer func() { _ = conn.Close() }()

	h.svc.JoinRoom(roomID, conn)
	defer h.svc.LeaveRoom(roomID, conn)

	log.Printf("client joined room: %s", roomID)

	for {
		mt, msg, err := conn.ReadMessage()
		if err != nil {
			if websocket.IsCloseError(err, websocket.CloseGoingAway, websocket.CloseNormalClosure) {
				log.Println("client disconnected")
			} else {
				log.Println("failed to read message:", err)
			}

			break
		}

		h.svc.BoradcastMessageToRoom(roomID, mt, msg)
	}

	return nil
}
