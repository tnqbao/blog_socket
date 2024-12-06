package namespaces

import (
	"fmt"

	socketio "github.com/googollee/go-socket.io"
	"github.com/tnqbao/blog_socket/internal/models"
)

func RegisterChatNamespace(server *socketio.Server) {
	server.OnNamespace("/chat", func(c socketio.Conn) error {
		fmt.Printf("User connected to /chat: %s\n", c.ID())

		c.On("join_room", func(room string) {
			c.Join(room)
			fmt.Printf("User %s joined room %s\n", c.ID(), room)
		})

		c.On("message", func(msg models.Message) {
			fmt.Printf("Message from %s: %s\n", c.ID(), msg.Content)

			server.BroadcastToRoom("/chat", msg.Room, "message", msg)
		})

		c.OnDisconnect(func(reason string) {
			fmt.Printf("User disconnected from /chat: %s\n", c.ID())
		})

		return nil
	})
}
