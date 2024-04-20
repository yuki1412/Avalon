package service

import (
	"log"
	"net/http"
	"sync"

	"github.com/gorilla/websocket"
)

var (
	upgrader = websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}
	connections   = make(map[*websocket.Conn]string) // Map to store active connections and their room codes
	connectionsMu sync.RWMutex                       // Mutex for concurrent access to connections map
)

func HandleWebSocket(w http.ResponseWriter, r *http.Request) {
	// Upgrade the HTTP connection to a WebSocket connection
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("Error upgrading to WebSocket:", err)
		return
	}
	defer func() {
		// Remove connection from map on close
		connectionsMu.Lock()
		delete(connections, conn)
		connectionsMu.Unlock()
		conn.Close()
	}()

	// Get the room code from the query parameter
	roomCode := r.URL.Query().Get("roomCode")
	if roomCode == "" {
		log.Println("Room code is required")
		return
	}

	// Add connection to the map with its room code
	connectionsMu.Lock()
	connections[conn] = roomCode
	connectionsMu.Unlock()

	// Continuously read messages from the WebSocket connection
	for {
		// Read message from client
		_, message, err := conn.ReadMessage()
		if err != nil {
			log.Println("Error reading message:", err)
			break
		}

		// Log received message
		log.Printf("Received message in room %s: %s", roomCode, message)

		// Broadcast the message to all clients in the same room
		broadcastMessage(roomCode, message)
	}
}

func broadcastMessage(roomCode string, message []byte) {
	connectionsMu.RLock()
	defer connectionsMu.RUnlock()

	// Iterate over connections and send message to clients in the same room
	for conn, rc := range connections {
		if rc == roomCode {
			if err := conn.WriteMessage(websocket.TextMessage, message); err != nil {
				log.Println("Error broadcasting message:", err)
			}
		}
	}
}
