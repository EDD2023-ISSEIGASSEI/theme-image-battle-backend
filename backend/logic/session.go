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

type SignUpSessionLogic struct {
	Session model.SignUpSession
}

func generateOtp() (*string, error) {
	ctx := context.Background()
	var buffer bytes.Buffer
	for {
		for i := 0; i < 6; i++ {
			n, err := rand.Int(rand.Reader, big.NewInt(10))
			if err != nil {
				return nil, err
			}
			buffer.WriteString(n.String())
		}
		otp := buffer.String()
		res := db.Redis.Exists(ctx, otp).Val()
		if res == 0 {
			return &otp, nil
		}
	}
}

func (sl *SignUpSessionLogic) CreateSession() error {
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

// sessionが存在していたらtrue
func (sl *SignUpSessionLogic) GetByUuid() (bool, error) {
	ctx := context.Background()
	lineSessionByte, err := db.Redis.GetDel(ctx, sl.Session.Uuid).Bytes()
	if err != nil {
		log.Errorln("RedisReadError: ", err.Error())
		return false, err
	}

	var signUpSession model.SignUpSession
	err = json.Unmarshal(lineSessionByte, &signUpSession)
	if err != nil {
		log.Errorln("JsonUnmarshalError: ", err.Error())
		return true, err
	}

	sl.Session = signUpSession
	return true, nil
}

type LineSessionLogic struct {
	Session model.LineSession
}

// otpのsessionが存在していたらtrue
func (sl *SignUpSessionLogic) LineRegisterByOtp(otp string) (bool, error) {
	ctx := context.Background()
	lineSessionByte, err := db.Redis.GetDel(ctx, otp).Bytes()
	if err != nil {
		log.Errorln("RedisReadError: ", err.Error())
		return false, err
	}

	var lineSession model.LineSession
	err = json.Unmarshal(lineSessionByte, &lineSession)
	if err != nil {
		log.Errorln("JsonUnmarshalError: ", err.Error())
		return true, err
	}

	sl.Session.User.LineUid = &lineSession.LineUid
	return true, nil
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
	db.Redis.Set(ctx, ll.Session.Otp, jsonData, 5*60*time.Second)
	return nil
}

type SignInSessionLogic struct {
	Session model.SignInSession
}

func (sl *SignInSessionLogic) CreateSession() error {
	var ctx = context.Background()
	sl.Session.Uuid = util.GenerateUuid()
	otp, err := generateOtp()
	if err != nil {
		log.Errorln("GenerateOtpError:", err.Error())
		return err
	}
	sl.Session.Otp = *otp
	jsonData, err := json.Marshal(sl.Session)
	if err != nil {
		log.Errorln("JsonMarshalError: ", err.Error())
		return err
	}
	db.Redis.Set(ctx, sl.Session.Uuid, jsonData, 5*60*time.Second)
	return nil
}

// sessionが存在していたらtrue
func (sl *SignInSessionLogic) GetByUuid() (bool, error) {
	ctx := context.Background()
	log.Debugln(sl.Session.Uuid)
	signInSessionByte, err := db.Redis.Get(ctx, sl.Session.Uuid).Bytes()
	if signInSessionByte == nil {
		log.Debugln("Session not found")
		return false, nil
	}
	if err != nil {
		log.Errorln("RedisReadError: ", err.Error())
		return false, err
	}

	var signInSession model.SignInSession
	err = json.Unmarshal(signInSessionByte, &signInSession)
	if err != nil {
		log.Errorln("JsonUnmarshalError: ", err.Error())
		return true, err
	}

	sl.Session = signInSession
	return true, nil
}

func (sl *SignInSessionLogic) CheckOtp(otp string) bool {
	return sl.Session.Otp == otp
}

func (sl *SignInSessionLogic) DeleteSession() {
	ctx := context.Background()
	db.Redis.Del(ctx, sl.Session.Uuid)
}
