package router

import (
	"net/http"
	"wayne/service"
)

func AvalonRouter() *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("/lobby", service.AvalonHandler)
	mux.HandleFunc("/api/create-room", service.CreateRoom)
	mux.HandleFunc("/api/join-room", service.JoinRoom)
	mux.HandleFunc("/ws", service.HandleWebSocket)
	return mux
}
