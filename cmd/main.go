package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/your_project/internal/socket"
)

func main() {
	server, err := socket.NewServer()
	if err != nil {
		log.Fatalf("Failed to create Socket.IO server: %v", err)
	}

	http.Handle("/socket.io/", server)

	port := ":8080"
	fmt.Printf("Server is running on http://localhost%s\n", port)
	log.Fatal(http.ListenAndServe(port, nil))
}
