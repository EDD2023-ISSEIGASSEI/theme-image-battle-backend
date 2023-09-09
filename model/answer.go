package model

type Answer struct {
	Player       Player   `json:"player"`
	BlankAnswers []string `json:"blankAnswers"`
	Score        int      `json:"score"`
}

type AnswerForQuestion struct {
	DealerPlayerId   string `json:"dealerPlayerId"`
	QuestionImageUrl string `json:"questionImageUrl"`
	Answer           Answer `json:"answer"`
}
