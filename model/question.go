package model

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
