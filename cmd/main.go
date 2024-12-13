package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

type Message struct {
	Role     string `json:"role"`
	Username string `json:"username"`
	Content  string `json:"content"`
	BlogID   string `json:"blogId"`
}

type ChatbotRequest struct {
	IDPost   string `json:"idPost"`
	Question string `json:"question"`
	Token    string `json:"token"`
}

type ChatbotResponse struct {
	Answer string `json:"answer"`
}

var clients = make(map[*websocket.Conn]bool)
var broadcast = make(chan Message)

const chatbotToken = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VySUQiOiIxIiwic2Vzc2lvbklEIjoiOTg4NzExZTYtNzJlZi00MGYwLTk4NmQtODA0OWVkMDJjZWMwIiwiaWF0IjoxNzMzNjQ4OTgyLCJleHAiOjE3MzQyNTM3ODJ9.6Aw3zklu4G_iaX8ID3TZQwzKSuMfQtjAojmmAzZFLvk"

func main() {
	http.HandleFunc("/ws", handleConnections)
	go handleMessages()

	fmt.Println("Server started on :8080")
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		panic("Error starting server: " + err.Error())
	}
}

func handleConnections(w http.ResponseWriter, r *http.Request) {
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		fmt.Println("Error upgrading connection:", err)
		return
	}
	defer ws.Close()

	clients[ws] = true

	for {
		var msg Message
		err := ws.ReadJSON(&msg)
		if err != nil {
			fmt.Println("Error reading JSON from client:", err)
			delete(clients, ws)
			break
		}

		if msg.Role == "user" && msg.Content != "" {
			broadcast <- msg
			fmt.Println("Received:", msg)

			answer, err := getAnswerFromApi(msg.BlogID, msg.Content)
			if err != nil {
				fmt.Println("Error calling AI API:", err)
				answer = "Error retrieving answer"
			}

			fmt.Println(answer)
			broadcast <- Message{
				Role:     "bot",
				Username: "bot",
				Content:  answer,
				BlogID:   msg.BlogID,
			}
		}
	}
}

func getAnswerFromApi(blogID, content string) (string, error) {
	apiURL := "https://b696-103-99-246-49.ngrok-free.app/chatbot2"
	if apiURL == "" {
		return "", fmt.Errorf("AI API URL is not set")
	}

	requestBody := ChatbotRequest{
		IDPost:   blogID,
		Question: content,
		Token:    chatbotToken,
	}

	body, err := json.Marshal(requestBody)
	if err != nil {
		return "", fmt.Errorf("error marshaling request body: %v", err)
	}

	resp, err := http.Post(apiURL, "application/json", bytes.NewBuffer(body))
	if err != nil {
		return "", fmt.Errorf("error making request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("API returned status code: %d", resp.StatusCode)
	}

	bodyBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("error reading response: %v", err)
	}

	var apiResp ChatbotResponse
	err = json.Unmarshal(bodyBytes, &apiResp)
	if err != nil {
		return "", fmt.Errorf("error unmarshaling response: %v", err)
	}

	return apiResp.Answer, nil
}

func handleMessages() {
	for {
		msg := <-broadcast
		for client := range clients {
			err := client.WriteJSON(msg)
			if err != nil {
				fmt.Println("Error writing JSON:", err)
				client.Close()
				delete(clients, client)
			}
		}
	}
}
