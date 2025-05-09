package main

import (
	"fmt"
	"strconv"
	"strings"
)

type requestHandlerModel struct {
	id     uint32
	action string
	value  string
}

func handleMessage(id uint32, payload string) {
	parts := strings.Split(payload, ";")
	action := safeIndex(parts, 0)
	value := safeIndex(parts, 1)
	requestHandler := requestHandlerModel{id: id, action: action, value: value}
	verifyAction(requestHandler)
}

func verifyAction(r requestHandlerModel) {
	switch r.action {
	case "start":
		// ignore if game is already started
		if status == "start" {
			return
		}
		changeStatus(r.action)
		changeGameStatus(r.action)
		clearRecord(r.action, "Game Started")
		handlerRoundAssigner()
		broadcastSystemInfo(getSystemInfo())
	case "stop":
		if status == "stop" || status == "lobby" {
			return
		}
		stopGame()
		broadcastSystemInfo(getSystemInfo())
	case "vote":
		// ignore if game is already stopped and not assigner
		if status == "stop" || !selectedClient[r.id] {
			return
		}
		value, err := strconv.ParseBool(r.value)
		if err != nil {
			return
		}
		if !voteList[r.id] {
			voteList[r.id] = value
		}
		changeGameStatus(r.action)
		broadcastSystemInfo(handlerVote())
		if isVotingOver() {
			broadcastSystemInfo(getVoteResult())
			nextRound()
			handlerRoundAssigner()
			if checkMajorityWin() {
				changeStatus("lobby")
				changeGameStatus("Game Ended")
				broadcastSystemInfo(getSystemInfo())
				return
			}
		}
		broadcastSystemInfo(getSystemInfo())
	case "select":
		// ignore if not assigner
		if r.id != roundAssigner {
			return
		}
		handlerSelectedClient(r.value)
	default:
		// Broadcast to all other clients
		broadcastMessage(r.id, r.action)
		// Broadcast current system info to all other clients
		broadcastSystemInfo(getSystemInfo())
	}
}

func handlerRoundAssigner() {
	// Find the next client ID in the clients map
	found := false
	for k := range clients {
		if found {
			roundAssigner = k
			break
		}
		if k == roundAssigner {
			found = true
		}
	}
	// If we didn't find a next client, wrap around to the first client
	if !found {
		for k := range clients {
			roundAssigner = k
			break
		}
	}
}

func handlerSelectedClient(client string) {
	clientID, err := strconv.ParseUint(client, 10, 32)
	if err != nil {
		return
	}
	_, exist := selectedClient[uint32(clientID)]
	if !exist {
		selectedClient[uint32(clientID)] = true
	} else {
		delete(selectedClient, uint32(clientID))
	}
	broadcastSystemInfo(fmt.Sprintf("clientList;%s", getClientList()))
}
func changeStatus(msg string) {
	status = msg
}

func changeGameStatus(mgs string) {
	gameStatus = mgs
}

func stopGame() {
	clearRecord("stop", "Game Stopped")
}

func checkMajorityWin() bool {
	trueCount := 0
	falseCount := 0
	for _, v := range missionBoard {
		if v == nil {
			continue
		}
		if *v {
			trueCount++
		} else {
			falseCount++
		}
		if trueCount > len(missionBoard)/2 || falseCount > len(missionBoard)/2 {
			return true
		}
	}
	return false
}

func getSystemInfo() string {
	stringList := fmt.Sprintf("Status:%s\n", strings.ToTitle(status))
	stringList += fmt.Sprintf("GStatus:%s\n", strings.ToTitle(gameStatus))
	for id, _ := range clients {
		stringList += fmt.Sprintf("- %d\n", id)
	}

	if status != "start" && status != "stop" {
		missionStatus := "Mission ["
		for _, v := range missionBoard {
			if v == &missionTrue || v == &missionFalse {
				missionStatus += fmt.Sprintf("%t ", *v)
			} else {
				missionStatus += fmt.Sprintf("%t ", false)
			}
		}
		missionStatus += "]\n"
		stringList += missionStatus
	}

	return fmt.Sprintf("Number of Connected Clients: %d\n%s", len(clients), stringList)
}

func isVotingOver() bool {
	return len(voteList) >= len(clients)
}

func handlerVote() string {
	// if all clients have voted, return the number of votes
	if isVotingOver() {
		status = "Vote Result"
		changeGameStatus(status)
	}

	return fmt.Sprintf("Votes: [%d/%d]\n", len(voteList), len(clients))
}

func getVoteResult() string {
	var result string = "Success"
	missionBoard[missionRound] = &missionTrue
	for _, v := range voteList {
		if !v {
			result = "Failed"
			missionBoard[missionRound] = &missionFalse
		}
	}
	return fmt.Sprintf("Vote Result: %s\n", result)
}

func isMissionOver() bool {
	return missionRound >= len(missionBoard)-1
}

func nextRound() {
	if isMissionOver() {
		status = "lobby"
		changeStatus(status)
		return
	}

	missionRound++
	voteList = make(map[uint32]bool)
	status = "vote"
}

func clearRecord(action string, gStatus string) {
	missionBoard = make([]*bool, 5)
	missionRound = 0
	voteList = make(map[uint32]bool)
	changeStatus(action)
	changeGameStatus(gStatus)
}

func safeIndex(slice []string, index int) string {
	if index < len(slice) {
		return slice[index]
	}
	return ""
}
