package models

import (
	"fmt"
	"net/http"
	"sync"

	"github.com/gorilla/websocket"
)

type Socket struct {
	upgrader   websocket.Upgrader
	clients    map[uint32]*websocket.Conn
	clientsMux sync.Mutex
}

func SocketInit() Socket {
	return Socket{
		upgrader: websocket.Upgrader{
			CheckOrigin: func(r *http.Request) bool {
				return true
			},
		},
		clients:    make(map[uint32]*websocket.Conn),
		clientsMux: sync.Mutex{},
	}
}

func (s *Socket) BroadcastMessage(id uint32, message string) {
	s.clientsMux.Lock()
	defer s.clientsMux.Unlock()

	for clientID, client := range s.clients {
		if clientID != id { // Don't send back to sender
			err := client.WriteMessage(websocket.TextMessage, []byte(fmt.Sprintf("%d: %s", id, message)))
			if err != nil {
				fmt.Printf("Error broadcasting to client %d: %v\n", clientID, err)
				client.Close()
				delete(s.clients, clientID)
			}
		}
	}
}

func (s *Socket) BroadcastSystemInfo(message string) {
	s.clientsMux.Lock()
	defer s.clientsMux.Unlock()
	for clientID, client := range s.clients {
		err := client.WriteMessage(websocket.TextMessage, []byte(message))
		if err != nil {
			fmt.Printf("Error broadcasting to client %d: %v\n", clientID, err)
			client.Close()
			delete(s.clients, clientID)
		}
	}
}
