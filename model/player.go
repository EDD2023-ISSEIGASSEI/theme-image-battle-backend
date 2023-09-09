package model

type Player struct {
	Id           string `json:"id"`
	Name         string `json:"name"`
	IconImageUrl string `json:"iconImageUrl"`
}

type PlayerState struct {
	Player      Player
	Score       int
	IsCompleted bool
}

func UserToPlayer(user User) Player {
	return Player{
		Id:           user.Id,
		Name:         user.Name,
		IconImageUrl: user.IconImageUrl,
	}
}
