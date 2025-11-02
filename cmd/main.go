package main

import (
	"fmt"
	"golang-live-chat/config"
	"golang-live-chat/internal/handlers"
	"golang-live-chat/internal/services"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/websocket"
	"github.com/labstack/echo/v4"

	"github.com/ilyakaznacheev/cleanenv"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true // allow any origin for now
	},
}

func main() {
	environment := os.Getenv("CHAT_ENVIRONMENT")
	cfgPath := fmt.Sprintf("config/%s.yml", environment)
	var cfg config.AppConfig

	err := cleanenv.ReadConfig(cfgPath, &cfg)
	if err != nil {
		log.Panic("failed to load config: %w", err)
	}
	log.Printf("AppConfig => %+v\n", cfg)

	chatSvc := services.NewChatService()
	handler := handlers.Handler{Svc: chatSvc}

	e := echo.New()

	e.GET("/ws", func(c echo.Context) error {
		return handler.HandleWebSocket(c, &upgrader)
	})

	err = e.Start(cfg.ChatAddress)
	if err != nil {
		log.Fatal(err)
	}
}
