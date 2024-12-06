package socket

import (
	socketio "github.com/googollee/go-socket.io"
	"github.com/tnqbao/blog_socket/internal/socket/namespaces"
)

func NewServer() (*socketio.Server, error) {
	server := socketio.NewServer(nil)

	namespaces.RegisterChatNamespace(server)
	namespaces.RegisterSearchNamespace(server)

	server.OnConnect("/", func(s socketio.Conn) error {
		s.SetContext("")
		return nil
	})

	server.OnDisconnect("/", func(s socketio.Conn, reason string) {
		s.Close()
	})

	return server, nil
}
