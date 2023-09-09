package model

import "golang.org/x/net/websocket"

type PlayerWsConn struct {
	Player Player
	Conn   websocket.Conn
}

type GameSession struct {
	Uuid               string              `json:"uuid"`
	Phase              string              `json:"phase"`
	RoundNum           int                 `json:"roundNum"`
	MaxRoundNum        int                 `json:"maxRoundNum"`
	Room               Room                `json:"room"`
	Players            []Player            `json:"players"`
	PlayerStates       []PlayerState       `json:"playerStates"`
	PlayerTopics       []PlayerTopic       `json:"playerTopics"`
	PlayerWsConns      []PlayerWsConn      `json:"playerWsConns"`
	Time               int                 `json:"time"`
	DealerPlayerId     string              `json:"dealerPlayerId"`
	ShowingPlayerId    string              `json:"showingPlayerId"`
	GeneratedQuestions []GeneratedQuestion `json:"generatedQuestions"`
	PlayerAnswers      []AnswerForQuestion `json:"playerAnswers"`
}
