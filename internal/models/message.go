package models

type Message struct {
	Role    string `json:"role"`
	Room    string `json:"room"`
	Content string `json:"content"`
}
