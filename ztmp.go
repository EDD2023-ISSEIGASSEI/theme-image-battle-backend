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
	Id           int    `json:"id"` // 6桁
	Name         string `json:"name"`
	Password     string `json:"password"`
	PlayerNum    int    `json:"playerNum"`
	MaxPlayerNum int    `json:"maxPlayerNum"`
}

type Player struct {
	Id           string `json:"id"`
	Name         string `json:"name"`
	IconImageUrl string `json:"iconImageUrl"`
}

type WaitingPhaseState struct {
	RoomInfo Room     `json:"roomInfo"`
	Players  []Player `json:"players"`
}

type Format = string // "<>をする<>"みたいにして、"<>"がblankであることを示す
type Genre = string

type Topic struct {
	Format Format   `json:"format"`
	Blanks []Genre  `json:"blanks"`
	Words  []string `json:"words"`
}

type TopicForGuess struct {
	Format Format  `json:"format"`
	Blanks []Genre `json:"blanks"`
}

type PlayerTopic struct {
	Player   Player `json:"player"`
	Question Topic  `json:"question"`
}

type PlayerState struct {
	Player      Player
	Score       int
	IsCompleted bool
}

type GeneratePhaseState struct {
	PlayerStates []PlayerState `json:"playerStates"`
	Time         int           `json:"time"`
}

type GeneratedQuestion struct {
	Player         Player `json:"player"`
	Topic          Topic  `json:"topic"`
	Prompt         string `json:"prompt"`
	ResultImageUrl string `json:"resultImageUrl"`
}

type GeneratedQuestionForGuess struct {
	TopicForGuess  TopicForGuess `json:"topicForGuess"`
	ResultImageUrl string        `json:"resultImageUrl"`
}

type GuessPhaseState struct {
	PlayerStates   []PlayerState             `json:"playerStates"`
	DealerPlayerId string                    `json:"dealerPlayerId"`
	Question       GeneratedQuestionForGuess `json:"question"`
	Time           int                       `json:"time"`
}

type Answer struct {
	Player       Player   `json:"player"`
	BlankAnswers []string `json:"blankAnswers"`
	Score        int      `json:"score"`
}

type AnswerForQuestion struct {
	DealerPlayerId   string   `json:"dealerPlayerId"`
	QuestionImageUrl string   `json:"questionImageUrl"`
	Answers          []Answer `json:"answers"`
}

type ShowScorePhaseState struct {
	PlayerStates    []PlayerState             `json:"playerStates"`
	DealerPlayerId  string                    `json:"dealerPlayerId"`
	ShowingPlayerId string                    `json:"showingPlayerId"`
	Question        GeneratedQuestionForGuess `json:"question"`
	PlayerAnswer    AnswerForQuestion         `json:"playerAnswer"`
}

type ShowCorrectAnswerPhaseState struct {
	Question    GeneratedQuestion `json:"question"`
	Answers     []Answer          `json:"answers"`
	DealerScore int               `json:"dealerScore"`
}

type EndingPhaseState struct {
	Ranking []PlayerState `json:"ranking"`
}

type Phase = string

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

// request response
type CreateRoomRequest struct {
	Name      string `json:"name"`
	Password  string `json:"password"`
	MaxMember int    `json:"maxMember"`
}

type CreateRoomResponse struct {
	Room Room `json:"room"`
}

type GetGamePhaseResponse struct {
	Phase Phase `json:"phase"`
}

type GetRoomListResponse struct {
	Rooms []Room `json:"rooms"`
}

type JoinRoomRequest struct {
	Id       int    `json:"id"`
	Password string `json:"password"`
}

type SubmitPromptRequest struct {
	Prompt string `json:"prompt"`
}

type SubmitPromptResponse struct {
	GeneratedImageUrl string `json:"generatedImageUrl"`
}

type SubmitAnswerRequest struct {
	Answers []string `json:"answers"`
}
