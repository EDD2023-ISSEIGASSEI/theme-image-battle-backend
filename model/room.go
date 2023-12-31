package model

type Room struct {
	Id            string `json:"id"` // 6桁
	Name          string `json:"name"`
	Password      string `json:"password"`
	PlayerNum     int    `json:"playerNum"`
	MaxPlayerNum  int    `json:"maxPlayerNum"`
	OwnerPlayerId string `json:"ownerPlayerId"`
}

type CreateRoomRequest struct {
	Name      string `json:"name"`
	Password  string `json:"password"`
	MaxMember int    `json:"maxMember"`
}

type RoomListResponse struct {
	Rooms []Room `json:"rooms"`
}

type JoinRoomRequest struct {
	Id       string `json:"id"`
	Password string `json:"password"`
}
