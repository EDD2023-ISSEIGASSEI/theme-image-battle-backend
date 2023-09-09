package model

type Room struct {
	Id           int    `json:"id"` // 6桁
	Name         string `json:"name"`
	Password     string `json:"password"`
	PlayerNum    int    `json:"playerNum"`
	MaxPlayerNum int    `json:"maxPlayerNum"`
}
