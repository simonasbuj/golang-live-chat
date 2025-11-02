package config

type AppConfig struct {
	ChatAddress		string 	`yaml:"chat_addres" env:"CHAT_ADDRESS" env-default:"localhost:7072"`
}