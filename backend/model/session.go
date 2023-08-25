package model

type SignUpSession struct {
	Uuid string `json:"uuid"`
	User User   `json:"user"`
}

type LineSession struct {
	Otp     string `json:"otp"`
	LineUid string `json:"lineUid"`
}

type SignInSession struct {
	Uuid string `json:"uuid"`
	Otp  string `json:"otp"`
	User User   `json:"user"`
}
