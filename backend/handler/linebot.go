package handler

import (
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
			uid := event.Source.UserID
			switch message := event.Message.(type) {
			case *linebot.TextMessage:
				_, err = bot.ReplyMessage(event.ReplyToken, linebot.NewTextMessage(message.Text+uid)).Do()
				if err != nil {
					log.Errorln(err.Error())
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
