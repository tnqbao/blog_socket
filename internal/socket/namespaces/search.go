package namespaces

import (
	"fmt"

	socketio "github.com/googollee/go-socket.io"
)

func RegisterSearchNamespace(server *socketio.Server) {
	server.OnNamespace("/search", func(c socketio.Conn) error {
		fmt.Printf("User connected to /search: %s\n", c.ID())

		c.On("query", func(query string) {
			fmt.Printf("Search query from %s: %s\n", c.ID(), query)

			results := []string{"Result 1", "Result 2", "Result 3"}
			c.Emit("results", results)
		})

		c.OnDisconnect(func(reason string) {
			fmt.Printf("User disconnected from /search: %s\n", c.ID())
		})

		return nil
	})
}
