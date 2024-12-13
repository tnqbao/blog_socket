package namespaces

import (
	"log"
	"strings"

	socketio "github.com/googollee/go-socket.io"
)

func RegisterChatNamespace(server *socketio.Server) {
	server.OnConnect("/chat", func(c socketio.Conn) error {
		log.Println("Connected to /chat:", c.ID())
		return nil
	})

	server.OnEvent("/chat", "join_room", func(c socketio.Conn, room string) {
		if room != "" {
			c.Join(room)
			log.Println("Client joined room:", room)
			c.Emit("message", "Welcome to the "+room+" room!")
		} else {
			c.Emit("message", "Room name is invalid!")
		}
	})

	server.OnEvent("/chat", "message", func(c socketio.Conn, msg string) {
		log.Println("Received message:", msg)

		room := getRoom(c)

		if msg != "" {
			server.BroadcastToRoom("/chat", room, "message", msg)

			botReply := generateBotReply(msg)

			server.BroadcastToRoom("/chat", room, "message", botReply)
		} else {
			c.Emit("message", "Please send a valid message.")
		}
	})

	server.OnDisconnect("/chat", func(c socketio.Conn, reason string) {
		log.Println("Disconnected from /chat:", reason)
	})
}

func getRoom(c socketio.Conn) string {
	for _, room := range c.Rooms() {
		if room != "" {
			return room
		}
	}
	return "default"
}

func generateBotReply(userMessage string) string {
	log.Printf("Generating bot response for: %s\n", userMessage)

	if strings.Contains(strings.ToLower(userMessage), "hello") {
		return "Hello! How can I assist you today?"
	}

	return "I'm sorry, I didn't understand that. Could you please rephrase?"
}
