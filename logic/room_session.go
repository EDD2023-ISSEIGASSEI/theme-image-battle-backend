package logic

import (
	"context"
	"edd2023-back/db"
	"edd2023-back/model"
	"encoding/json"
	"time"

	log "github.com/sirupsen/logrus"
)

type RoomSessionLogic struct {
	Session model.RoomSession
}

func (rl *RoomSessionLogic) CreateRoomSession(room model.Room, gameSessionId string) error {
	rl.Session = model.RoomSession{
		Room:          room,
		GameSessionId: gameSessionId,
	}

	var ctx = context.Background()
	jsonData, err := json.Marshal(rl.Session)
	if err != nil {
		log.Errorln("JsonMarshalError: ", err.Error())
		return err
	}
	db.RoomRedis.Set(ctx, rl.Session.Room.Id, jsonData, 24*time.Hour)
	return nil
}
