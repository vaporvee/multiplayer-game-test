package main

import (
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net"
	"os"

	"github.com/joho/godotenv"
)

type Payload struct {
	Type      string `json:"type"`
	Message   string `json:"msg,omitempty"`
	Direction *struct {
		X int `json:"x,omitempty"`
		Y int `json:"y,omitempty"`
	} `json:"direction,omitempty"`
	SessionID string `json:"session_id,omitempty"`
}

// Generate a unique session ID
func generateSessionID() string {
	b := make([]byte, 32)
	if _, err := io.ReadFull(rand.Reader, b); err != nil {
		return ""
	}
	return base64.URLEncoding.EncodeToString(b)
}

func main() {
	godotenv.Load()
	port := os.Getenv("PORT")
	if port == "" {
		port = "4477" // Default port if not specified
	}

	addr, err := net.ResolveUDPAddr("udp", ":"+port)
	if err != nil {
		fmt.Println(err)
		return
	}

	conn, err := net.ListenUDP("udp", addr)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer conn.Close()

	fmt.Println("Server running on " + port)

	buf := make([]byte, 1024)
	clients := make(map[string]string) // Map to keep track of clients and their session IDs

	for {
		n, addr, err := conn.ReadFromUDP(buf)
		if err != nil {
			fmt.Println(err)
			continue
		}

		var payload Payload
		err = json.Unmarshal(buf[:n], &payload)
		if err != nil {
			fmt.Println("Error parsing JSON:", err)
			continue
		}

		// Use the client's IP and port as the key to uniquely identify the connection
		key := addr.String()

		// Check if the client's connection already has a session ID
		if _, ok := clients[key]; !ok {
			// The client is new, generate a session ID
			sessionID := generateSessionID()
			clients[key] = sessionID
			// Send the session ID back to the client
			responsePayload := Payload{SessionID: sessionID}
			response, _ := json.Marshal(responsePayload)
			conn.WriteToUDP(response, addr)
		}

		switch payload.Type {
		case "init":
			fmt.Printf("Received initiation generated SessionID: %s", clients[key])
		case "move":
			fmt.Printf("Received move message from Session ID: %s : X=%d, Y=%d\n", clients[key], payload.Direction.X, payload.Direction.Y)
			// Handle movement logic here
		case "message":
			fmt.Printf("Received message from Session ID: %s : %s\n", clients[key], payload.Message)
			broadcastMessage(conn, clients, payload.Message)
		default:
			fmt.Printf("Received unknown message type Session ID: %s\n", clients[key])
		}
	}
}

func broadcastMessage(conn *net.UDPConn, clients map[string]string, message string) {
	for clientKey := range clients {
		// Parse the clientKey to get the *net.UDPAddr
		clientAddr, err := net.ResolveUDPAddr("udp", clientKey)
		if err != nil {
			fmt.Println("Error resolving UDP address:", err)
			continue
		}
		payload := Payload{
			Type:    "message",
			Message: message,
		}
		response, _ := json.Marshal(payload)
		conn.WriteToUDP(response, clientAddr)
	}
}
