package model

type SignUpSession struct {
	Uuid string `json:"uuid"`
	User User   `json:"user"`
}

type LineSession struct {
	Otp     string `json:"otp"`
	LineUid string `json:"lineUid"`
}
