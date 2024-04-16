package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

// UserToken represents a user token in the database.
type UserToken struct {
	UserID string
	Token  string
}

// userTokens is a simple in-memory map for demonstration purposes.
var userTokens = make(map[string]UserToken)

// generateItchOAuthURL generates an OAuth URL for itch.io.
func generateItchOAuthURL(clientID string) string {
	return fmt.Sprintf("https://itch.io/user/oauth?client_id=%s&response_type=token", clientID)
}

// redirectToItchOAuth sends the OAuth URL to the client.
func redirectToItchOAuth(conn *websocket.Conn, clientID string) {
	oauthURL := generateItchOAuthURL(clientID)
	err := conn.WriteMessage(websocket.TextMessage, []byte(oauthURL))
	if err != nil {
		log.Println(err)
	}
}

// saveUserToken saves the user token in the database.
func saveUserToken(userID, token string) {
	userTokens[userID] = UserToken{UserID: userID, Token: token}
}

// handler handles WebSocket connections.
func handler(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}
	defer conn.Close()

	// Example client ID. Replace with your actual client ID.
	clientID := "your_client_id_here"

	// Redirect the user to the OAuth URL.
	redirectToItchOAuth(conn, clientID)

	// Simulate receiving the token from the client.
	// In a real application, you would parse the token from the redirect URL.
	userID := "example_user_id"
	token := "example_token"
	saveUserToken(userID, token)

	// Send a confirmation message to the client.
	err = conn.WriteMessage(websocket.TextMessage, []byte("Token saved successfully."))
	if err != nil {
		log.Println(err)
	}
}

func main() {
	http.HandleFunc("/ws", handler)
	log.Fatal(http.ListenAndServe(":8080", nil))
}
