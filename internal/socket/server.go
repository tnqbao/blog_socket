package socket

import (
	"fmt"
	"log"

	socketio "github.com/googollee/go-socket.io"
	"github.com/tnqbao/blog_socket/internal/socket/namespaces"
)

func NewSocketServer() (*socketio.Server, error) {
	server := socketio.NewServer(nil)
	if server == nil {
		return nil, fmt.Errorf("failed to create Socket.IO server")
	}
	namespaces.RegisterChatNamespace(server)

	server.OnDisconnect("/", func(c socketio.Conn, reason string) {
		log.Println("Client disconnected:", reason)
	})

	go server.Serve()
	return server, nil
}
