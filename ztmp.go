package main

import "golang.org/x/net/websocket"

type User struct {
	Id           string
	Name         string
	Password     string
	LineUid      string
	IconImageUrl string
}

type Room struct {
	Id           int // 6桁
	Name         string
	Password     string
	PlayerNum    int
	MaxPlayerNum int
}

type Player struct {
	Id           string
	Name         string
	IconImageUrl string
}

type WaitingPhaseState struct {
	RoomInfo Room
	Players  []Player
}

type Format = string // "<>をする<>"みたいにして、"<>"がblankであることを示す
type Genre = string

type Topic struct {
	Format Format
	Blanks []Genre
}

type PlayerTopic struct {
	Player   Player
	Question Topic
}

type PlayerState struct {
	Player      Player
	Score       int
	IsCompleted bool
}

type GeneratePhaseState struct {
	PlayerStates []PlayerState
	Time         int
}

type GeneratedQuestion struct {
	Player         Player
	Topic          Topic
	Prompt         string
	ResultImageUrl string
}

type GeneratedQuestionForGuess struct {
	Format         Format
	ResultImageUrl string
}

type GuessPhaseState struct {
	PlayerStates   []PlayerState
	DealerPlayerId string
	Question       GeneratedQuestionForGuess
	Time           int
}

type Answer struct {
	Player       Player
	BlankAnswers []string
	Score        int
}

type AnswerForQuestion struct {
	DealerPlayerId   string
	QuestionImageUrl string
	Answers          []Answer
}

type ShowScorePhaseState struct {
	PlayerStates    []PlayerState
	DealerPlayerId  string
	ShowingPlayerId string
	Question        GeneratedQuestionForGuess
	PlayerAnswer    AnswerForQuestion
}

type ShowCorrectAnswerPhaseState struct {
	Question    GeneratedQuestion
	Answers     []Answer
	DealerScore int
}

type EndingPhaseState struct {
	Ranking []PlayerState
}

const (
	WaitingPhase           = "WaitingPhase"
	GeneratePhase          = "GeneratePhase"
	GuessPhase             = "GuessPhase"
	ShowScorePhase         = "ShowScorePhase"
	ShowCorrectAnswerPhase = "ShowCorrectAnswerPhase"
	EndingPhase            = "EndingPhase"
)

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
