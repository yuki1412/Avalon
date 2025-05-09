package main

import (
	"fmt"
	"net/http"
	"sort"
	"sync"
	"sync/atomic"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

// Store all active connections
var (
	clients    = make(map[uint32]*websocket.Conn)
	clientsMux sync.Mutex
)

func broadcastMessage(id uint32, message string) {
	clientsMux.Lock()
	defer clientsMux.Unlock()

	for clientID, client := range clients {
		if clientID != id { // Don't send back to sender
			err := client.WriteMessage(websocket.TextMessage, []byte(fmt.Sprintf("%d: %s", id, message)))
			if err != nil {
				fmt.Printf("Error broadcasting to client %d: %v\n", clientID, err)
				client.Close()
				delete(clients, clientID)
			}
		}
	}
}

func broadcastSystemInfo(message string) {
	clientsMux.Lock()
	defer clientsMux.Unlock()
	for clientID, client := range clients {
		err := client.WriteMessage(websocket.TextMessage, []byte(message))
		if err != nil {
			fmt.Printf("Error broadcasting to client %d: %v\n", clientID, err)
			client.Close()
			delete(clients, clientID)
		}
	}
}

func getClientList() string {
	clientsMux.Lock()
	defer clientsMux.Unlock()

	keys := make([]uint32, 0, len(clients))
	for k := range clients {
		keys = append(keys, k)
	}
	sort.Slice(keys, func(i, j int) bool {
		return keys[i] < keys[j]
	})

	clientList := ""
	for _, clientID := range keys {
		val, exist := selectedClient[clientID]
		if exist {
			clientList += fmt.Sprintf("%d;%t;", clientID, val)
		} else {
			clientList += fmt.Sprintf("%d;%t;", clientID, false)
		}
	}

	return clientList
}

func handleConnection(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		fmt.Println("Error upgrading connection:", err)
		return
	}

	// assign a unique ID to this connection
	id := atomic.AddUint32(&connIDCounter, 1)

	// Add to clients map
	clientsMux.Lock()
	clients[id] = conn
	clientsMux.Unlock()

	fmt.Printf("Connection %d established\n", id)
	// send clients information to subscribers
	broadcastSystemInfo(fmt.Sprintf("clientList;%s", getClientList()))

	// Cleanup when connection closes
	defer func() {
		clientsMux.Lock()
		delete(clients, id)
		clientsMux.Unlock()
		conn.Close()
		fmt.Printf("Connection %d closed\n", id)
		broadcastSystemInfo(fmt.Sprintf("clientList;%s", getClientList()))
	}()

	for {
		_, msg, err := conn.ReadMessage()
		if err != nil {
			fmt.Println("Read error:", err)
			break
		}
		fmt.Printf("Received from %d: %s\n", id, string(msg))

		handleMessage(id, string(msg))
	}
}
