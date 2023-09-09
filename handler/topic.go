package handler

import (
	"edd2023-back/logic"
	"edd2023-back/model"
	"edd2023-back/util"
	"net/http"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

type TopicHandler struct{}

func (*TopicHandler) GetTopic(ctx *gin.Context) {
	tl := logic.TopicLogic{}
	f, err := tl.GenerateTopic()
	if !f && err == nil {
		s := "InvalidId"
		r := util.BadRequest(&s)
		ctx.JSON(r.StatusCode, r.Message)
		return
	}
	if err != nil {
		log.Errorln("[Error]exec error: ", err.Error())
		r := util.InternalServerError(nil)
		ctx.JSON(r.StatusCode, r.Message)
		return
	}

	sessionId, err := ctx.Cookie("sessionId")
	asl := logic.AuthSessionLogic{
		Session: model.AuthSession{Uuid: sessionId},
	}
	asl.GetByUuid()
	player := model.UserToPlayer(asl.Session.User)

	gameSessionId, _ := ctx.Cookie("gameSessionId")
	gsl := logic.GameSessionLogic{}
	gsl.GetByUuid(gameSessionId)

	gsl.Session.PlayerTopics = append(gsl.Session.PlayerTopics, model.PlayerTopic{
		Player:   player,
		Question: tl.Topic,
	})
	err = gsl.UpdateByUuId()
	if err != nil {
		log.Errorln("[Error]exec error: ", err.Error())
		r := util.InternalServerError(nil)
		ctx.JSON(r.StatusCode, r.Message)
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"topic": tl.Topic})
}
