package service

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"text/template"
	"wayne/model"
)

func AvalonHandler(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.ParseFiles("template/lobby.html"))
	res := map[string]interface{}{}
	tmpl.Execute(w, res)
	// fmt.Fprintf(w, "Welcome to the home page!")
}

func CreateRoom(w http.ResponseWriter, r *http.Request) {
	// Check Content-Type header
	// contentType := r.Header.Get("Content-Type")
	// if contentType != "application/json" {
	// 	http.Error(w, "Unsupported Content-Type", http.StatusUnsupportedMediaType)
	// 	return
	// }

	// Read the request body
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Error reading request body", http.StatusInternalServerError)
		return
	}
	defer r.Body.Close() // Close the request body
	fmt.Printf("Resp: %+v\n", string(body))
	// loginRequest := request.LoginRequest{}
	// if err := json.Unmarshal(body, &loginRequest); err != nil {
	// 	fmt.Println("Err", err.Error())
	// 	return
	// }
	// fmt.Printf("Resp2: %+v\n", loginRequest)
	values, err := url.ParseQuery(string(body))
	if err != nil {
		fmt.Println("Error parsing query:", err)
		return
	}
	username := values.Get("username")
	roomCode := values.Get("room_code")
	// Create Room
	// model.CreateRoom(loginRequest.RoomCode, loginRequest.Username)
	roomObj := model.CreateRoom(roomCode, username)
	var usernameList []string
	for _, item := range roomObj.GetPlayerList() {
		usernameList = append(usernameList, item.Username)
	}

	tmpl := template.Must(template.ParseFiles("template/lobby_player.html"))
	res := map[string]interface{}{
		"host_name":     roomObj.Host,
		"username_list": usernameList,
	}
	tmpl.Execute(w, res)
}

func JoinRoom(w http.ResponseWriter, r *http.Request) {
	// Read the request body
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Error reading request body", http.StatusInternalServerError)
		return
	}
	defer r.Body.Close() // Close the request body
	fmt.Printf("Resp: %+v\n", string(body))
	values, err := url.ParseQuery(string(body))
	if err != nil {
		fmt.Println("Error parsing query:", err)
		return
	}
	username := values.Get("username")
	roomCode := values.Get("room_code")
	// Join Room
	player := model.Player{
		Username: username,
	}
	roomObj, err := player.JoinRoom(roomCode)
	if err != nil {
		fmt.Println("SDF: ", err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	fmt.Println("Room: ", roomObj)
	var usernameList []string
	for _, item := range roomObj.GetPlayerList() {
		usernameList = append(usernameList, item.Username)
	}

	tmpl := template.Must(template.ParseFiles("template/lobby_player.html"))
	res := map[string]interface{}{
		"host_name":     roomObj.Host,
		"username_list": usernameList,
	}
	tmpl.Execute(w, res)
}
