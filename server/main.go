package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"github.com/gorilla/websocket"
	"github.com/joho/godotenv"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

var clients = make(map[*websocket.Conn]bool) // connected clients
var broadcast = make(chan []byte)            // broadcast channel

func main() {
	godotenv.Load()
	http.HandleFunc("/ws", handleConnections)
	go handleMessages()
	fmt.Println("Server running on " + os.Getenv("PORT"))
	http.ListenAndServeTLS(os.Getenv("PORT"), "cert.pem", "key.pem", nil)
}

func handleConnections(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer conn.Close() // Closes the connection at the end

	clients[conn] = true // Add a new WebSocket connection to the clients map.

	for {
		_, msg, err := conn.ReadMessage()
		if err != nil {
			delete(clients, conn)
			break
		}

		// Parse the received message as JSON
		var payload map[string]interface{}
		err = json.Unmarshal(msg, &payload)
		if err != nil {
			fmt.Println("Error parsing JSON:", err)
			continue
		}

		// Check if the JSON object contains the key "type" with the value "broadcast"
		if payload["type"] == "broadcast" {
			broadcast <- msg // Send received WebSocket Message to Broadcast
		}
	}
}

func handleMessages() {
	for {
		msg := <-broadcast
		for client := range clients {
			err := client.WriteMessage(websocket.TextMessage, msg)
			if err != nil {
				client.Close()
				delete(clients, client)
			}
		}
	}
}
