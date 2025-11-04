package handler

import (
	"encoding/json"
	"fmt"
	"golang-live-chat/internal/dto"
	"golang-live-chat/internal/services"
	"log"
	"net/http"
	"time"

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
	roomID, user, err := h.validateRoomAndUser(c)
	if err != nil {
		return fmt.Errorf("failed to validate request to join: %w", err)
	}

	conn, err := h.upgrader.Upgrade(c.Response(), c.Request(), nil)
	if err != nil {
		log.Println("websocket upgrade error: ", err)

		return fmt.Errorf("failed to upgrade websocket: %w", err)
	}

	defer func() { _ = conn.Close() }()

	h.svc.JoinRoom(roomID, conn)
	defer h.svc.LeaveRoom(roomID, conn)

	log.Printf("user %s joined room: %s", user, roomID)

	h.svc.BoradcastMessageToRoom(
		roomID,
		websocket.TextMessage,
		dto.NewMessage("joined", user, "joined the room"),
	)

	for {
		mt, msg, err := conn.ReadMessage()
		if err != nil {
			if websocket.IsCloseError(err, websocket.CloseGoingAway, websocket.CloseNormalClosure) {
				log.Println("client disconnected")
				h.svc.BoradcastMessageToRoom(
					roomID,
					websocket.TextMessage,
					dto.NewMessage("disconnected", user, "disconnected the room"),
				)
			} else {
				log.Println("failed to read message:", err)
			}

			break
		}

		var msgDto dto.Message

		err = json.Unmarshal(msg, &msgDto)
		if err != nil {
			log.Println("failed to unmarshal the message: %w", err)
			h.svc.WriteErrorToClient(conn, "failed to unmarshal the message")

			continue
		}

		msgDto.Type = "message"
		msgDto.User = user
		msgDto.Time = time.Now().Unix()

		h.svc.BoradcastMessageToRoom(roomID, mt, msgDto)
	}

	return nil
}

func (h *Handler) validateRoomAndUser(c echo.Context) (string, string, error) {
	roomID := c.QueryParam("room")
	if roomID == "" {
		return "", "", c.JSON( //nolint:wrapcheck
			http.StatusBadRequest,
			map[string]string{"error": "missing required query parameter: room"},
		)
	}

	user := c.Request().Header.Get("Token")
	if user == "" {
		user = c.QueryParam("token")
	}

	if user == "" {
		return "", "", c.JSON( //nolint:wrapcheck
			http.StatusUnauthorized,
			map[string]string{"error": "bad token"},
		)
	}

	return roomID, user, nil
}
