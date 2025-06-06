package models

type Lobby struct {
	ID      uint32   `json:"id"`
	Players []Player `json:"players"`
	Status  string   `json:"status"`
}

func LobbyInit(id uint32) *Lobby {
	return &Lobby{
		ID: id,
	}
}

func (l *Lobby) AddPlayer(player Player) {
	l.Players = append(l.Players, player)
}

func (l *Lobby) RemovePlayer(player Player) {
	for i, p := range l.Players {
		if p.ID == player.ID {
			l.Players = append(l.Players[:i], l.Players[i+1:]...)
			break
		}
	}
}

func (l *Lobby) GetPlayerList() []Player {
	return l.Players
}

func (l *Lobby) GetPlayerCount() int {
	return len(l.Players)
}
