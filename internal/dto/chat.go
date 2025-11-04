package dto

import "time"

type Message struct {
	Type    string `json:"type"`
	User    string `json:"user"`
	Content string `json:"content"`
	Time    int64  `json:"time"`
}

func NewMessage(msgType, user, content string) Message {
	return Message{
		Type:    msgType,
		User:    user,
		Content: content,
		Time:    time.Now().Unix(),
	}
}

type ErrorResponse struct {
	Type  string `json:"type"`
	Error string `json:"error"`
}
