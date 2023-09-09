package model

import "golang.org/x/net/websocket"

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
