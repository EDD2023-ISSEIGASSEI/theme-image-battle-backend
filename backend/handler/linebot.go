package handler

import (
	"line-bot-otp-back/logic"
	"line-bot-otp-back/model"

	"github.com/gin-gonic/gin"
	"github.com/line/line-bot-sdk-go/v7/linebot"
	log "github.com/sirupsen/logrus"
)

type LinebotHandler struct{}

func (*LinebotHandler) EventHandler(ctx *gin.Context, bot *linebot.Client) {
	events, err := bot.ParseRequest(ctx.Request)
	if err != nil {
		if err == linebot.ErrInvalidSignature {
			ctx.Writer.WriteHeader(400)
		} else {
			ctx.Writer.WriteHeader(500)
		}
		return
	}
	for _, event := range events {
		if event.Type == linebot.EventTypeMessage {
			switch message := event.Message.(type) {
			case *linebot.TextMessage:
				if message.Text == "登録" || message.Text == "とうろく" {
					text := ""
					otp, err := generateLineRegistrationOtp(event.Source.UserID)
					if err != nil {
						log.Errorln(err.Error())
						text = "サーバで問題が発生しました\nしばらくしてからもう一度試しやがれ"
					} else {
						text = "↓ワンタイムパスワード↓\n" + *otp
					}
					_, err = bot.ReplyMessage(event.ReplyToken, linebot.NewTextMessage(text)).Do()
					if err != nil {
						log.Errorln(err.Error())
					}
				} else {
					text := "おいおいタイポかよーー"
					_, err = bot.ReplyMessage(event.ReplyToken, linebot.NewTextMessage(text)).Do()
					if err != nil {
						log.Errorln(err.Error())
					}
				}
			default:
				replyMessage := "can not"
				_, err = bot.ReplyMessage(event.ReplyToken, linebot.NewTextMessage(replyMessage)).Do()
				if err != nil {
					log.Errorln(err.Error())
				}
			}
		}
	}
}

func generateLineRegistrationOtp(uid string) (*string, error) {
	ll := logic.LineSessionLogic{
		Session: model.LineSession{
			LineUid: uid,
		},
	}
	err := ll.Create()
	if err != nil {
		log.Errorln("GenerateLineOptSession: ", err.Error())
		return nil, err
	}
	return &ll.Session.Otp, nil
}
