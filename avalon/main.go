package main

import (
	"fmt"
	"net/http"

	"wayne/initialize"
	"wayne/router"
)

func main() {
	initialize.InitialCache()
	// http.HandleFunc("/api/change-host", func(w http.ResponseWriter, r *http.Request) {
	// 	// Read the request body
	// 	body, err := io.ReadAll(r.Body)
	// 	if err != nil {
	// 		http.Error(w, "Error reading request body", http.StatusInternalServerError)
	// 		return
	// 	}

	// 	loginRequest := request.ChangeHostRequest{}
	// 	if err := json.Unmarshal(body, &loginRequest); err != nil {
	// 		fmt.Println("Err", err.Error())
	// 		return
	// 	}
	// 	p := model.Player{
	// 		Username: loginRequest.Username,
	// 	}
	// 	p.ChangeHost(loginRequest.RoomCode, loginRequest.NewHost)
	// })

	// http.HandleFunc("/api/start-game", func(w http.ResponseWriter, r *http.Request) {
	// 	// Read the request body
	// 	body, err := io.ReadAll(r.Body)
	// 	if err != nil {
	// 		http.Error(w, "Error reading request body", http.StatusInternalServerError)
	// 		return
	// 	}

	// 	loginRequest := request.LoginRequest{}
	// 	if err := json.Unmarshal(body, &loginRequest); err != nil {
	// 		fmt.Println("Err", err.Error())
	// 		return
	// 	}
	// 	room := model.GetRoom(loginRequest.RoomCode)
	// 	room.StartRoom(loginRequest.Username)
	// })

	// http.HandleFunc("/api/join-room", func(w http.ResponseWriter, r *http.Request) {
	// 	// Read the request body
	// 	body, err := io.ReadAll(r.Body)
	// 	if err != nil {
	// 		http.Error(w, "Error reading request body", http.StatusInternalServerError)
	// 		return
	// 	}

	// 	loginRequest := request.LoginRequest{}
	// 	if err := json.Unmarshal(body, &loginRequest); err != nil {
	// 		fmt.Println("Err", err.Error())
	// 		return
	// 	}
	// 	p := model.Player{
	// 		Username: loginRequest.Username,
	// 	}
	// 	p.JoinRoom(loginRequest.RoomCode)
	// })

	// http.HandleFunc("/api/get-role", func(w http.ResponseWriter, r *http.Request) {
	// 	// Simulate fetching data from a database or external service
	// 	res := map[string]interface{}{
	// 		"username": "Wayneng",
	// 	}
	// 	// Render a template with the data
	// 	tmpl := template.Must(template.New("pText").Parse(`<p id="role-details">Hello {{.username}}, Welcome to Avalon World</p>`))
	// 	tmpl.Execute(w, res)
	// })

	// lobbyTemp := func(w http.ResponseWriter, r *http.Request) {
	// 	tmpl := template.Must(template.ParseFiles("template/lobby.html"))
	// 	res := map[string]interface{}{
	// 		"username": "Wayne",
	// 	}
	// 	tmpl.Execute(w, res)
	// }
	// http.HandleFunc("/api/to-avalon", lobbyTemp)

	// // Start the HTTP server on port 8080
	// fmt.Println("Server started on http://localhost:8080")
	// http.ListenAndServe("localhost:8080", nil)

	server := &http.Server{
		Addr:    ":8080",
		Handler: router.HomeRouter(),
	}
	fmt.Println("Server running on port 8080")
	if err := server.ListenAndServe(); err != nil {
		fmt.Println("Error starting server:", err)
	}
}
