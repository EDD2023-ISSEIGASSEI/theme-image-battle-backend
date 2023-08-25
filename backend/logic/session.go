package logic

import (
	"bytes"
	"context"
	"crypto/rand"
	"encoding/json"
	"line-bot-otp-back/db"
	"line-bot-otp-back/model"
	"line-bot-otp-back/util"
	"math/big"
	"time"

	log "github.com/sirupsen/logrus"
)

type SignUpSessionLigic struct {
	Session model.SignUpSession
}

type LineSessionLogic struct {
	Session model.LineSession
}

func generateOtp() (*string, error) {
	var buffer bytes.Buffer

	for i := 0; i < 6; i++ {
		n, err := rand.Int(rand.Reader, big.NewInt(10))
		if err != nil {
			return nil, err
		}
		buffer.WriteString(n.String())
	}

	randomString := buffer.String()
	return &randomString, nil
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

func (ll *LineSessionLogic) Create() error {
	var ctx = context.Background()
	otp, err := generateOtp()
	if err != nil {
		log.Errorln("GenerateOtpError:", err.Error())
		return err
	}
	ll.Session.Otp = *otp
	jsonData, err := json.Marshal(ll.Session)
	if err != nil {
		log.Errorln("JsonMarshalError: ", err.Error())
		return err
	}
	db.Redis.Set(ctx, ll.Session.Otp, jsonData, 0)
	return nil
}
