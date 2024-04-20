package model

import (
	"errors"
	"fmt"
)

type Player struct {
	Username string
	Role     string
}

func (p *Player) getRole() string {
	return p.Role
}

func (p *Player) SetRole(role string) {
	p.Role = role
}

func (p *Player) JoinRoom(roomCode string) (room *Room, err error) {
	room, ok := ROOM_MAP[roomCode]
	if !ok {
		message := fmt.Sprintf("room %s does not exist", roomCode)
		return nil, errors.New(message)
	}
	if room.checkDuplicateName(p.Username) {
		fmt.Printf("Room %s does not exist\n", roomCode)
		return nil, fmt.Errorf("username %s is used", p.Username)
	}
	room.Players = append(room.Players, *p)
	ROOM_MAP[roomCode] = room
	fmt.Printf("Player %s have joined the room\n", p.Username)
	fmt.Printf("Current Number of Player: %d \n", room.getNumberOfPlayers())
	return
}

func (p *Player) ChangeHost(roomCode string, newHost string) {
	room := getRoom(roomCode)
	room.Host = newHost
	ROOM_MAP[roomCode] = room
	fmt.Printf("%s is now host\n", newHost)
}
