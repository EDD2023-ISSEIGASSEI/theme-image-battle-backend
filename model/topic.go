package model

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
