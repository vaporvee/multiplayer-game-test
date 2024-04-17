package main

import (
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	r "math/rand"
	"net"
	"os"
	"time"

	"github.com/joho/godotenv"
)

type Payload struct {
	SessionID string   `json:"session_id,omitempty"`
	Type      string   `json:"type"`
	Message   string   `json:"msg,omitempty"`
	Direction *Vector2 `json:"direction,omitempty"`
}

type InitPayload struct {
	Type         string    `json:"type"`
	PlayerClient *Client   `json:"player_client"`
	Clients      []*Client `json:"clients"`
}

type DisconnectPayload struct {
	Type      string `json:"type"`
	SessionID string `json:"session_id"`
}

type Client struct {
	SessionID string    `json:"session_id"`
	LastSeen  time.Time `json:"last_seen"`
	Positon   *Vector2  `json:"position"`
}

type Vector2 struct {
	X float64 `json:"x,omitempty"`
	Y float64 `json:"y,omitempty"`
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
	clients := make(map[string]*Client) // Map to keep track of clients and their session IDs

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
			clients[key] = &Client{SessionID: sessionID, LastSeen: time.Now(), Positon: &Vector2{X: randomFloatInRange(300, 600), Y: randomFloatInRange(300, 600)}} // Initialize a new Client struct and store its pointer
			// Send the session ID back to the client
			clientsSlice := make([]*Client, 0, len(clients))
			for _, client := range clients {
				clientsSlice = append(clientsSlice, client)
			}
			responsePayload := InitPayload{Type: "init_success", PlayerClient: clients[key], Clients: clientsSlice}
			broadcastMessage(conn, clients, responsePayload)
		}

		clients[key].LastSeen = time.Now()

		switch payload.Type {
		case "init":
			fmt.Printf("Received initiation generated SessionID: %s\n", clients[key].SessionID)
		case "move":
			fmt.Printf("Received move message from Session ID: %s : X=%f, Y=%f\n", clients[key].SessionID, payload.Direction.X, payload.Direction.Y)
			broadcastMessage(conn, clients, payload)

		case "disconnect":
			// Disconnect the client
			fmt.Printf("Client %s disconnected\n", clients[key].SessionID)
			broadcastMessage(conn, clients, DisconnectPayload{Type: "disconnect", SessionID: clients[key].SessionID})
			delete(clients, key)
		case "message":
			fmt.Printf("Received message from Session ID: %s : %s\n", clients[key].SessionID, payload.Message)
			broadcastMessage(conn, clients, payload)
		default:
			fmt.Printf("Received unknown message type Session ID: %s\n", clients[key].SessionID)
		}

		// Check for disconnected clients and reset their session ID
		for clientKey, client := range clients {
			if time.Since(client.LastSeen) > 5*time.Minute { // 5 minutes timeout
				delete(clients, clientKey)
				fmt.Printf("Client %s disconnected and session ID reset\n", clients[clientKey].SessionID)
			}
		}
	}
}

func broadcastMessage(conn *net.UDPConn, clients map[string]*Client, payload interface{}) {
	for clientKey := range clients {
		clientAddr, err := net.ResolveUDPAddr("udp", clientKey)
		if err != nil {
			fmt.Println("Error resolving UDP address:", err)
			continue
		}
		response, _ := json.Marshal(payload)
		conn.WriteToUDP(response, clientAddr)
	}
}

func randomFloatInRange(min, max float64) float64 {
	seed := time.Now().UnixNano()
	rf := r.New(r.NewSource(seed))
	return min + rf.Float64()*(max-min)
}
