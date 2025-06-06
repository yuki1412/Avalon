package models

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"sync/atomic"

	"github.com/gorilla/websocket"
)

// Request handler model needs to be updated to use string ID
type requestHandlerModel struct {
	id     string
	action string
	value  string
}

type ClientVar struct {
	// Round []*Round
	Round         *Round
	Socket        Socket
	Clients       Clients `json:"clients"`
	Lobby         *Lobby
	ConnIDCounter uint32 `json:"conn_id_counter"`
	// Store players by their persistent UserID for reconnection
	Players map[string]*Player `json:"players"`
	// Mapping between string UserIDs and uint32 IDs for Round system compatibility
	UserIDToRoundID map[string]uint32 `json:"user_id_to_round_id"`
	RoundIDToUserID map[string]string `json:"round_id_to_user_id"`
}

// Generate a random user ID
func generateUserID() string {
	bytes := make([]byte, 16)
	rand.Read(bytes)
	return hex.EncodeToString(bytes)
}

func ClientVarInit(clients []uint32) *ClientVar {
	return &ClientVar{
		Round:           RoundInit(ClientsInit(clients)),
		Socket:          SocketInit(),
		Clients:         ClientsInit(clients),
		Lobby:           LobbyInit(uint32(1)),
		Players:         make(map[string]*Player),
		UserIDToRoundID: make(map[string]uint32),
		RoundIDToUserID: make(map[string]string),
	}
}

func (c *ClientVar) HandleConnection(w http.ResponseWriter, req *http.Request) {
	conn, err := c.Socket.upgrader.Upgrade(w, req, nil)
	if err != nil {
		fmt.Println("Error upgrading connection:", err)
		return
	}

	// Get username and user ID from query parameters
	username := req.URL.Query().Get("username")
	existingUserID := req.URL.Query().Get("userID")

	var player *Player
	var isReconnection bool

	// Check if this is a reconnection
	if existingUserID != "" {
		if existingPlayer, exists := c.Players[existingUserID]; exists {
			// Reuse existing player with same ID
			player = existingPlayer
			isReconnection = true

			// Update socket clients map
			c.Socket.clientsMux.Lock()
			c.Socket.clients[player.ID] = conn
			c.Socket.clientsMux.Unlock()

			fmt.Printf("Player %s (%s) reconnected with UserID: %s\n", username, existingUserID[:8], existingUserID)
		} else {
			// User ID not found, treat as new player
			userID := generateUserID()
			roundID := atomic.AddUint32(&c.ConnIDCounter, 1)
			player = PlayerInitWithUserID(userID, username)
			c.Players[userID] = player

			// Create mapping for Round system
			c.UserIDToRoundID[userID] = roundID
			c.RoundIDToUserID[strconv.Itoa(int(roundID))] = userID

			// Add to socket clients map
			c.Socket.clientsMux.Lock()
			c.Socket.clients[userID] = conn
			c.Socket.clientsMux.Unlock()

			fmt.Printf("Player %s connected with new UserID: %s (old ID not found)\n", username, userID[:8])
		}
	} else {
		// New player
		userID := generateUserID()
		roundID := atomic.AddUint32(&c.ConnIDCounter, 1)
		player = PlayerInitWithUserID(userID, username)
		c.Players[userID] = player

		// Create mapping for Round system
		c.UserIDToRoundID[userID] = roundID
		c.RoundIDToUserID[strconv.Itoa(int(roundID))] = userID

		// Add to socket clients map
		c.Socket.clientsMux.Lock()
		c.Socket.clients[userID] = conn
		c.Socket.clientsMux.Unlock()

		fmt.Printf("New player %s connected with UserID: %s\n", username, userID[:8])
	}

	// Send user ID and connection info to client
	welcomeMsg := fmt.Sprintf("userID:%s;reconnection:%t", player.ID, isReconnection)
	conn.WriteMessage(websocket.TextMessage, []byte(welcomeMsg))

	// Only add to lobby if it's a new player (not reconnection)
	if !isReconnection {
		handleLobby(c, *player)
	}

	fmt.Printf("Connection %s established\n", player.ID[:8])
	// send comprehensive information to all subscribers when connection is updated
	c.Socket.BroadcastSystemInfo(fmt.Sprintf("clientList;%s", c.GetPlayersList()))
	c.Socket.BroadcastSystemInfo(c.GetPlayersSystemInfo())

	// Cleanup when connection closes
	defer func() {
		c.Socket.clientsMux.Lock()
		delete(c.Socket.clients, player.ID)
		c.Socket.clientsMux.Unlock()
		conn.Close()
		fmt.Printf("Connection %s closed (Player kept for reconnection)\n", player.ID[:8])
		// send comprehensive information to all subscribers when connection is updated
		c.Socket.BroadcastSystemInfo(fmt.Sprintf("clientList;%s", c.GetPlayersList()))
		c.Socket.BroadcastSystemInfo(c.GetPlayersSystemInfo())
	}()

	for {
		_, msg, err := conn.ReadMessage()
		if err != nil {
			fmt.Println("Read error:", err)
			break
		}
		fmt.Printf("Received from %s (%s): %s\n", player.Username, player.ID[:8], string(msg))

		c.handleMessage(player.ID, string(msg))
	}
}

func handleLobby(c *ClientVar, player Player) {
	c.Lobby.AddPlayer(player)
	fmt.Printf("Current Players: ")
	for _, p := range c.Lobby.Players {
		// Start of Selection
		idSuffix := p.ID
		if len(idSuffix) > 4 {
			idSuffix = idSuffix[len(idSuffix)-4:]
		}
		fmt.Printf("Player: ID(%s) Username(%s)\n", idSuffix, p.Username)
	}
}

func (c *ClientVar) handleMessage(id string, payload string) {
	parts := strings.Split(payload, ";")
	action := safeIndex(parts, 0)
	value := safeIndex(parts, 1)

	// Convert string UserID to uint32 for Round system compatibility
	roundID, exists := c.UserIDToRoundID[id]
	if !exists {
		fmt.Printf("Warning: No round ID mapping found for user %s\n", id[:8])
		return
	}

	requestHandler := requestHandlerModel{id: strconv.Itoa(int(roundID)), action: action, value: value}
	c.verifyAction(requestHandler)
}

func safeIndex(slice []string, index int) string {
	if index < len(slice) {
		return slice[index]
	}
	return ""
}

func (c *ClientVar) verifyAction(req requestHandlerModel) {
	switch req.action {
	case "start":
		// ignore if game is already started
		if c.Round.Status == "start" {
			return
		}
		clientList := make([]uint32, 0)
		for userID, _ := range c.Socket.clients {
			if roundID, exists := c.UserIDToRoundID[userID]; exists {
				clientList = append(clientList, roundID)
			}
		}
		fmt.Printf("clientList: %v\n", clientList)
		c.Round = RoundInit(ClientsInit(clientList))
		c.Round.InitRoundAssigner()
		c.Round.ChangeStatus(req.action)
		c.Round.UpdateStatus(req.action)
		c.Round.ClearRecord(req.action, "Game Started")
		c.Socket.BroadcastSystemInfo(c.GetPlayersSystemInfo())
	case "stop":
		if c.Round.Status == "stop" || c.Round.Status == "lobby" {
			return
		}
		c.Round.StopGame()
		c.Socket.BroadcastSystemInfo(c.GetPlayersSystemInfo())
	case "vote":
		// ignore if game is already stopped and not assigner
		if c.Round.Status == "stop" || !c.Round.SelectedClient[req.id] {
			return
		}
		value, err := strconv.ParseBool(req.value)
		if err != nil {
			return
		}
		if !c.Round.VoteList[req.id] {
			c.Round.VoteList[req.id] = value
		}
		c.Round.UpdateStatus(req.action)
		c.Socket.BroadcastSystemInfo(c.Round.handlerVote())
		if c.Round.IsVotingOver() {
			c.Socket.BroadcastSystemInfo(c.Round.getVoteResult())
			c.Round.nextRound()
			c.Round.nextRoundAssigner()
			if c.Round.AnyMajority() {
				c.Round.ChangeStatus("lobby")
				c.Round.UpdateStatus("Game Ended")
				c.Socket.BroadcastSystemInfo(c.GetPlayersSystemInfo())
				return
			}
		}
		c.Socket.BroadcastSystemInfo(c.GetPlayersSystemInfo())
	case "select":
		// ignore if not assigner
		if req.id != fmt.Sprint(c.Round.RoundAssigner) {
			return
		}
		c.Socket.BroadcastSystemInfo(c.HandlerSelectedClient(req.value))
	default:
		// Convert back to string UserID for broadcast
		if userID, exists := c.RoundIDToUserID[req.id]; exists {
			c.Socket.BroadcastMessage(userID, req.action)
		}
		// Broadcast current system info to all other clients
		c.Socket.BroadcastSystemInfo(c.GetPlayersSystemInfo())
	}
}

// GetPlayersList returns a formatted list of connected players with their information
func (c *ClientVar) GetPlayersList() string {
	playersList := ""

	// Get connected players (those in socket.clients)
	c.Socket.clientsMux.Lock()
	for userID := range c.Socket.clients {
		if _, exists := c.Players[userID]; exists {
			if roundID, hasRoundID := c.UserIDToRoundID[userID]; hasRoundID {
				selected := c.Round.SelectedClient[strconv.Itoa(int(roundID))]
				playersList += fmt.Sprintf("%d;%t;", roundID, selected)
			}
		}
	}
	c.Socket.clientsMux.Unlock()

	return playersList
}

// GetPlayersSystemInfo returns comprehensive system information including player details
func (c *ClientVar) GetPlayersSystemInfo() string {
	connectedCount := 0
	playersList := ""

	// Get connected players (those in socket.clients)
	c.Socket.clientsMux.Lock()
	for userID := range c.Socket.clients {
		if player, exists := c.Players[userID]; exists {
			// Show player info with username
			idSuffix := player.ID
			if len(idSuffix) > 8 {
				idSuffix = idSuffix[:8]
			}
			playersList += fmt.Sprintf("Player: %s (ID: %s)\n", player.Username, idSuffix)
			connectedCount++
		}
	}
	c.Socket.clientsMux.Unlock()

	// Build system info
	systemInfo := fmt.Sprintf("Status:%s\n", strings.ToTitle(c.Round.Status))
	systemInfo += fmt.Sprintf("GStatus:%s\n", strings.ToTitle(c.Round.Status))
	systemInfo += fmt.Sprintf("Round:%d\n", c.Round.MissionTracker+1)
	systemInfo += fmt.Sprintf("Round Assigner:%d\n", c.Round.RoundAssigner)
	systemInfo += playersList

	// Add mission status if game is active
	if c.Round.Status != "start" && c.Round.Status != "stop" {
		missionStatus := "Mission ["
		for _, v := range c.Round.MissionBoard {
			if v != nil {
				missionStatus += fmt.Sprintf("%t ", *v)
			} else {
				missionStatus += fmt.Sprintf("%t ", false)
			}
		}
		missionStatus += "]\n"
		systemInfo += missionStatus
	}

	return fmt.Sprintf("Number of Players: %d\n%s", connectedCount, systemInfo)
}

// HandlerSelectedClient handles client selection for missions
func (c *ClientVar) HandlerSelectedClient(client string) string {
	clientID, err := strconv.ParseUint(client, 10, 32)
	if err != nil {
		return ""
	}
	_, exist := c.Round.SelectedClient[string(clientID)]
	if !exist {
		c.Round.SelectedClient[string(clientID)] = true
	} else {
		delete(c.Round.SelectedClient, string(clientID))
	}
	return fmt.Sprintf("clientList;%s", c.GetPlayersList())
}
