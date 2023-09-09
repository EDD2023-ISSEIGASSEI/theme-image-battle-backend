package handler

import (
	"edd2023-back/logic"
	"edd2023-back/model"
	"edd2023-back/util"
	"net/http"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

type GameHandler struct{}

func ownerCheck(
	authSessionId string,
	gameSessionId string,
) (
	bool,
	logic.AuthSessionLogic,
	logic.GameSessionLogic,
) {
	asl := logic.AuthSessionLogic{
		Session: model.AuthSession{Uuid: authSessionId},
	}
	asl.GetByUuid()

	gsl := logic.GameSessionLogic{
		Session: model.GameSession{Uuid: gameSessionId},
	}
	gsl.GetByUuid(gameSessionId)

	return asl.Session.User.Id == gsl.Session.Room.OwnerPlayerId, asl, gsl
}

type GameStartRequest struct {
	Round int `json:"round"`
}

func (*GameHandler) GameStart(ctx *gin.Context) {
	sessionId, _ := ctx.Cookie("sessionId")
	gameSessionId, _ := ctx.Cookie("gameSessionId")
	f, _, gsl := ownerCheck(sessionId, gameSessionId)
	if !f {
		s := "NotOwner"
		log.Infoln("[Forbiden] ", s)
		ctx.JSON(http.StatusForbidden, gin.H{"message": s})
		return
	}

	var req GameStartRequest
	err := ctx.Bind(&req)
	if err != nil {
		s := err.Error()
		r := util.BadRequest(&s)
		log.Errorln("[Error]request parse error: ", s)
		ctx.JSON(r.StatusCode, r.Message)
		return
	}

	ps := []model.PlayerState{}
	if len(gsl.Session.PlayerStates) == 0 {
		for _, p := range gsl.Session.Players {
			ps = append(ps, model.PlayerState{
				Player:      p,
				Score:       0,
				IsCompleted: false,
			})
		}
	} else {
		for idx := range gsl.Session.PlayerStates {
			gsl.Session.PlayerStates[idx].IsCompleted = false
		}
	}
	gsl.Session.Phase = model.GeneratePhase
	gsl.Session.PlayerStates = ps
	gsl.Session.RoundNum = req.Round
	gsl.Session.GeneratedQuestions = []model.GeneratedQuestion{}
	gsl.Session.PlayerTopics = []model.PlayerTopic{}

	err = gsl.UpdateByUuId()
	if err != nil {
		log.Errorln("[Error]exec error: ", err.Error())
		r := util.InternalServerError(nil)
		ctx.JSON(r.StatusCode, r.Message)
		return
	}

	// ToDo: broadcast

	ctx.JSON(http.StatusOK, gin.H{"message": "OK"})
}
