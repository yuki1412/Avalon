package main

import (
	"fmt"
	"net/http"
)

func main() {
	http.HandleFunc("/ws", handleConnection)
	fmt.Println("WebSocket server running on :8080")
	http.ListenAndServe(":8080", nil)
}
