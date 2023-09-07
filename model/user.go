package model

type User struct {
	Id           string  `json:"id"`
	Name         string  `json:"name"`
	Password     string  `json:"password"`
	LineUid      *string `json:"lineUid"`
	IconImageUrl string  `json:"iconImageUrl"`
}
