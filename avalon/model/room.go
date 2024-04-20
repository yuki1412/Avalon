package model

import (
	"fmt"
	"math/rand"
	"wayne/consts"
)

var ROOM_MAP = map[string]*Room{}

type Room struct {
	Players         []Player
	RoomCode        string
	Host            string
	Leader          Player
	MissionTeam     []Player
	RoundAssignment []int
}

func (r *Room) GetPlayerList() []Player {
	return r.Players
}

func (r *Room) getLeader() Player {
	return r.Leader
}

func (r *Room) getAssignee() []Player {
	return r.MissionTeam
}

func (r *Room) getNumberOfPlayers() int {
	return len(r.Players)
}

func getRoom(roomCode string) *Room {
	room, ok := ROOM_MAP[roomCode]
	if !ok {
		fmt.Printf("Room %s does not exist\n", roomCode)
	}
	return room
}

func (r *Room) checkDuplicateName(newUsername string) bool {
	for _, p := range r.GetPlayerList() {
		if p.Username == newUsername {
			return true
		}
	}
	return false
}

func CreateRoom(roomCode string, hostUsername string) *Room {
	fmt.Println("Creating")
	ROOM_MAP[roomCode] = &Room{
		Players: []Player{{
			Username: hostUsername,
		}},
		Host:     hostUsername,
		RoomCode: roomCode,
	}
	fmt.Printf("Room %s is created\n", roomCode)
	return ROOM_MAP[roomCode]
}

func (r *Room) StartRoom(userTrigger string) {
	if r.Host != userTrigger {
		fmt.Println("Only Host can start game")
		return
	}
	roundAssignment, ok := consts.ROUND_ASSIGNMENT_MAP[r.getNumberOfPlayers()]
	if !ok {
		fmt.Printf("Current Player: %d, insufficient/overloaded players\n", r.getNumberOfPlayers())
		return
	}
	r.RoundAssignment = roundAssignment
	r.shufflePlayersArrangement()
	r.Leader = r.Players[0]
	r.randomAssignRole()
	fmt.Printf("Room %s is starting now\n", r.RoomCode)
	fmt.Printf("Room Details: %+v\n", r)
}

func (r *Room) getTeamCount() map[string]int {
	teamMap, ok := consts.ROLE_SHUFFLE_MAP[r.getNumberOfPlayers()]
	if !ok {
		fmt.Println("Unable to get team map")
		return teamMap
	}
	return teamMap
}

func (r *Room) randomAssignRole() {
	teamMap := r.getTeamCount()
	playerIdx := 0
	for _, role := range consts.EVIL_ROLE[:teamMap["EVIL"]] {
		r.Players[playerIdx].SetRole(role)
		playerIdx++
	}
	for _, role := range consts.GOOD_ROLE[:teamMap["GOOD"]] {
		r.Players[playerIdx].SetRole(role)
		playerIdx++
	}
}

func (r *Room) shufflePlayersArrangement() {
	// Shuffle the slice
	rand.Shuffle(r.getNumberOfPlayers(), func(i, j int) {
		r.Players[i], r.Players[j] = r.Players[j], r.Players[i]
	})
}
