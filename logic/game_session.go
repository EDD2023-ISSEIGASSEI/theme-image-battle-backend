package logic

import (
	"context"
	"edd2023-back/db"
	"edd2023-back/model"
	"edd2023-back/util"
	"encoding/json"
	"time"

	log "github.com/sirupsen/logrus"
)

type GameSessionLogic struct {
	Session model.GameSession
}

func (rl *GameSessionLogic) CreateGameSession(player model.Player, room model.Room) error {
	uuid := util.GenerateUuid()
	rl.Session = model.GameSession{
		Uuid:               uuid,
		Phase:              model.WaitingPhase,
		RoundNum:           0,
		MaxRoundNum:        0,
		Room:               room,
		Players:            []model.Player{player},
		PlayerStates:       []model.PlayerState{},
		PlayerTopics:       []model.PlayerTopic{},
		PlayerWsConns:      []model.PlayerWsConn{},
		Time:               0,
		DealerPlayerId:     "",
		ShowingPlayerId:    "",
		GeneratedQuestions: []model.GeneratedQuestion{},
		PlayerAnswers:      []model.AnswerForQuestion{},
	}

	var ctx = context.Background()
	jsonData, err := json.Marshal(rl.Session)
	if err != nil {
		log.Errorln("JsonMarshalError: ", err.Error())
		return err
	}
	db.Redis.Set(ctx, rl.Session.Uuid, jsonData, 24*time.Hour)
	return nil
}
