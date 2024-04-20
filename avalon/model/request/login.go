package request

type LoginRequest struct {
	Username string `json:"username" form:"username"`
	RoomCode string `json:"room_code" form:"room_code"`
}

type ChangeHostRequest struct {
	Username string `json:"username"`
	RoomCode string `json:"room_code"`
	NewHost  string `json:"new_host"`
}
