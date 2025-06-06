package models

type Player struct {
	ID       string `json:"id"` // Persistent user ID (primary identifier)
	Username string `json:"username"`
}

func PlayerInit(username string) *Player {
	return &Player{
		Username: username,
	}
}

func PlayerInitWithUserID(userID string, username string) *Player {
	return &Player{
		ID:       userID,
		Username: username,
	}
}
