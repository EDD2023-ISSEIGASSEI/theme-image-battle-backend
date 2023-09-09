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

type PlayerStateByScore []PlayerState

func (p PlayerStateByScore) Len() int           { return len(p) }
func (p PlayerStateByScore) Swap(i, j int)      { p[i], p[j] = p[j], p[i] }
func (p PlayerStateByScore) Less(i, j int) bool { return p[i].Score < p[j].Score }
