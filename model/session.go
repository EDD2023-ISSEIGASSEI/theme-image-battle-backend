package model

import "golang.org/x/net/websocket"

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

type AuthSession struct {
	Uuid string `json:"sessionId"`
	User User   `json:"user"`
}

// game

type PlayerWsConn struct {
	Player Player
	Conn   websocket.Conn
}

type GameSession struct {
	Phase              string
	RoundNum           int
	MaxRoundNum        int
	Room               Room
	Players            []Player
	PlayerStates       []PlayerState
	PlayerTopics       []PlayerTopic
	PlayerWsConns      []PlayerWsConn
	Time               int
	DealerPlayerId     string
	ShowingPlayerId    string
	GeneratedQuestions []GeneratedQuestion
	PlayerAnswers      []AnswerForQuestion
}

type RoomSession struct {
	Room          Room
	GameSessionId string
}
