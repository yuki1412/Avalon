package models

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"sync/atomic"
)

type RoundClient struct {
	// Round []*Round
	Round         *Round
	Socket        Socket
	Clients       Clients `json:"clients"`
	ConnIDCounter uint32  `json:"conn_id_counter"`
}

func RoundClientInit(clients []uint32) *RoundClient {
	return &RoundClient{
		Round:   RoundInit(ClientsInit(clients)),
		Socket:  SocketInit(),
		Clients: ClientsInit(clients),
	}
}

func (r *RoundClient) HandleConnection(w http.ResponseWriter, req *http.Request) {
	conn, err := r.Socket.upgrader.Upgrade(w, req, nil)
	if err != nil {
		fmt.Println("Error upgrading connection:", err)
		return
	}

	// assign a unique ID to this connection
	id := atomic.AddUint32(&r.ConnIDCounter, 1)

	// Add to clients map
	r.Socket.clientsMux.Lock()
	r.Socket.clients[id] = conn
	r.Socket.clientsMux.Unlock()

	fmt.Printf("Connection %d established\n", id)
	// send clients information to subscribers
	r.Socket.BroadcastSystemInfo(fmt.Sprintf("clientList;%s", r.Round.GetClientList()))

	// Cleanup when connection closes
	defer func() {
		r.Socket.clientsMux.Lock()
		delete(r.Socket.clients, id)
		r.Socket.clientsMux.Unlock()
		conn.Close()
		fmt.Printf("Connection %d closed\n", id)
		r.Socket.BroadcastSystemInfo(fmt.Sprintf("clientList;%s", r.Round.GetClientList()))
	}()

	for {
		_, msg, err := conn.ReadMessage()
		if err != nil {
			fmt.Println("Read error:", err)
			break
		}
		fmt.Printf("Received from %d: %s\n", id, string(msg))

		r.handleMessage(id, string(msg))
	}
}

func (r *RoundClient) handleMessage(id uint32, payload string) {
	parts := strings.Split(payload, ";")
	action := safeIndex(parts, 0)
	value := safeIndex(parts, 1)
	requestHandler := requestHandlerModel{id: id, action: action, value: value}
	r.verifyAction(requestHandler)
}

func safeIndex(slice []string, index int) string {
	if index < len(slice) {
		return slice[index]
	}
	return ""
}

func (r *RoundClient) verifyAction(req requestHandlerModel) {
	switch req.action {
	case "start":
		// ignore if game is already started
		if r.Round.Status == "start" {
			return
		}
		clientList := make([]uint32, 0)
		for k, _ := range r.Socket.clients {
			clientList = append(clientList, uint32(k))
		}
		fmt.Printf("clientList: %v\n", clientList)
		r.Round = RoundInit(ClientsInit(clientList))
		r.Round.InitRoundAssigner()
		r.Round.ChangeStatus(req.action)
		r.Round.UpdateStatus(req.action)
		r.Round.ClearRecord(req.action, "Game Started")
		r.Socket.BroadcastSystemInfo(r.Round.GetSystemInfo())
	case "stop":
		if r.Round.Status == "stop" || r.Round.Status == "lobby" {
			return
		}
		r.Round.StopGame()
		r.Socket.BroadcastSystemInfo(r.Round.GetSystemInfo())
	case "vote":
		// ignore if game is already stopped and not assigner
		if r.Round.Status == "stop" || !r.Round.SelectedClient[req.id] {
			return
		}
		value, err := strconv.ParseBool(req.value)
		if err != nil {
			return
		}
		if !r.Round.VoteList[req.id] {
			r.Round.VoteList[req.id] = value
		}
		r.Round.UpdateStatus(req.action)
		r.Socket.BroadcastSystemInfo(r.Round.handlerVote())
		if r.Round.IsVotingOver() {
			r.Socket.BroadcastSystemInfo(r.Round.getVoteResult())
			r.Round.nextRound()
			r.Round.nextRoundAssigner()
			if r.Round.AnyMajority() {
				r.Round.ChangeStatus("lobby")
				r.Round.UpdateStatus("Game Ended")
				r.Socket.BroadcastSystemInfo(r.Round.GetSystemInfo())
				return
			}
		}
		r.Socket.BroadcastSystemInfo(r.Round.GetSystemInfo())
	case "select":
		// ignore if not assigner
		if req.id != r.Round.RoundAssigner {
			return
		}
		r.Socket.BroadcastSystemInfo(r.Round.HandlerSelectedClient(req.value))
	default:
		// Broadcast to all other clients
		r.Socket.BroadcastMessage(req.id, req.action)
		// Broadcast current system info to all other clients
		r.Socket.BroadcastSystemInfo(r.Round.GetSystemInfo())
	}
}
