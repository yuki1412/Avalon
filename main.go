package main

import (
	"fmt"
	"log"
	"net/http"
	"playground/models"
)

var roundClient models.ClientVar = *models.ClientVarInit([]uint32{})

func main() {
	http.HandleFunc("/ws", roundClient.HandleConnection)
	fmt.Println("WebSocket server running on :8080")
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal("ListenAndServe error:", err)
	}
}
