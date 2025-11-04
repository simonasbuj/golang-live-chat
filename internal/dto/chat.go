package dto

type Message struct {
	Type    string `json:"type"`
	User    string `json:"user"`
	Content string `json:"content"`
	Time    int64  `json:"time"`
}

type ErrorResponse struct {
	Type  string `json:"type"`
	Error string `json:"error"`
}
