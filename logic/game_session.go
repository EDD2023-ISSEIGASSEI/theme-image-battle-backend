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
	// TODO wsConn追加
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

func (gl *GameSessionLogic) GetByUuid(uuid string) (bool, error) {
	ctx := context.Background()
	gameSessionByte, err := db.Redis.Get(ctx, uuid).Bytes()
	if err != nil {
		log.Errorln("RedisReadError: ", err.Error())
		return false, err
	}

	var gameSession model.GameSession
	err = json.Unmarshal(gameSessionByte, &gameSession)
	if err != nil {
		log.Errorln("JsonUnmarshalError: ", err.Error())
		return true, err
	}

	gl.Session = gameSession
	return true, nil
}

func (gl *GameSessionLogic) UpdateByUuId() error {
	var ctx = context.Background()
	jsonData, err := json.Marshal(gl.Session)
	if err != nil {
		log.Errorln("JsonMarshalError: ", err.Error())
		return err
	}
	db.Redis.GetSet(ctx, gl.Session.Uuid, jsonData)
	return nil
}

func (gl *GameSessionLogic) IsExistsPlayer(player model.Player) bool {
	for _, p := range gl.Session.Players {
		if p.Id == player.Id {
			return true
		}
	}
	return false
}

func (gl *GameSessionLogic) JoinPlayer(player model.Player) {
	// TODO wsConn追加
	gl.Session.Players = append(gl.Session.Players, player)
	gl.Session.Room.PlayerNum += 1
}

func (gl *GameSessionLogic) DeleteSession() {
	ctx := context.Background()
	db.Redis.Del(ctx, gl.Session.Uuid)
}
