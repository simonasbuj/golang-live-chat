package main

import (
	"fmt"
	"golang-live-chat/config"
	"golang-live-chat/internal/handler"
	"golang-live-chat/internal/services"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gorilla/websocket"
	"github.com/labstack/echo/v4"

	"github.com/ilyakaznacheev/cleanenv"
)

func main() {
	environment := os.Getenv("CHAT_ENVIRONMENT")
	cfgPath := fmt.Sprintf("config/%s.yml", environment)

	var cfg config.AppConfig

	err := cleanenv.ReadConfig(cfgPath, &cfg)
	if err != nil {
		log.Panic("failed to load config: %w", err)
	}

	log.Printf("AppConfig => %+v\n", cfg)

	upgrader := websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			return true // allow any origin for now
		},
		HandshakeTimeout:  time.Duration(cfg.ChatHandshakeTimeout) * time.Second,
		ReadBufferSize:    cfg.ChatReadBufferSize,
		WriteBufferSize:   cfg.ChatWriteBufferSize,
		WriteBufferPool:   nil,
		Subprotocols:      nil,
		Error:             nil,
		EnableCompression: false,
	}

	chatSvc := services.NewChatService()
	handler := handler.New(chatSvc, &upgrader)

	e := echo.New()

	e.GET("/global-chat", handler.HandleGlobalChat)
	e.GET("/chat", handler.HandleRoomChat)

	err = e.Start(cfg.ChatAddress)
	if err != nil {
		log.Fatal(err)
	}
}
