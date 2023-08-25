package logic

import (
	"context"
	"encoding/json"
	"line-bot-otp-back/db"
	"line-bot-otp-back/model"
	"line-bot-otp-back/util"
	"time"

	log "github.com/sirupsen/logrus"
)

type SignUpSessionLigic struct {
	Session model.SignUpSession
}

func (sl *SignUpSessionLigic) Create() error {
	sl.Session.Uuid = util.GenerateUuid()
	ctx := context.Background()
	jsonData, err := json.Marshal(sl.Session)
	if err != nil {
		log.Errorln("JsonMarshalError: ", err.Error())
		return err
	}
	db.Redis.Set(ctx, sl.Session.Uuid, jsonData, 5*60*time.Second)
	return nil
}
