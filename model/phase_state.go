package model

type Phase = string

const (
	WaitingPhase           = "WaitingPhase"
	GeneratePhase          = "GeneratePhase"
	GuessPhase             = "GuessPhase"
	ShowScorePhase         = "ShowScorePhase"
	ShowCorrectAnswerPhase = "ShowCorrectAnswerPhase"
	EndingPhase            = "EndingPhase"
)

type WaitingPhaseState struct {
	RoomInfo Room     `json:"roomInfo"`
	Players  []Player `json:"players"`
}

type GeneratePhaseState struct {
	PlayerStates []PlayerState `json:"playerStates"`
	Time         int           `json:"time"`
}

type GuessPhaseState struct {
	PlayerStates   []PlayerState             `json:"playerStates"`
	DealerPlayerId string                    `json:"dealerPlayerId"`
	Question       GeneratedQuestionForGuess `json:"question"`
	Time           int                       `json:"time"`
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

type PhaseStateResponse struct {
	Phase Phase `json:"phase"`
	State any   `json:"state"`
}
