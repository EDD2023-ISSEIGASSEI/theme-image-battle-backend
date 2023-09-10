package handler

import (
	"edd2023-back/logic"
	"edd2023-back/model"
	"edd2023-back/util"
	"net/http"

	"github.com/gin-gonic/gin"
)

type RoundHandler struct{}

func (*RoundHandler) NextRound(ctx *gin.Context) {
	gameSessionId, _ := ctx.Cookie("gameSessionId")
	gsl := logic.GameSessionLogic{}
	gsl.GetByUuid(gameSessionId)

	nextRound := gsl.Session.RoundNum + 1
	if nextRound > gsl.Session.MaxRoundNum {
		s := "NoNextRound"
		r := util.BadRequest(&s)
		ctx.JSON(r.StatusCode, r.Message)
		return
	}
	gsl.Session.RoundNum = nextRound
	gsl.Session.Phase = model.GeneratePhase
	gsl.Session.GeneratedQuestions = []model.GeneratedQuestion{}
	gsl.Session.PlayerAnswers = []model.AnswerForQuestion{}
	gsl.Session.PlayerTopics = []model.PlayerTopic{}
	for idx := range gsl.Session.PlayerStates {
		gsl.Session.PlayerStates[idx].IsCompleted = false
	}
	gsl.UpdateByUuId()
	// ToDo: broadcast
	ctx.JSON(http.StatusOK, gin.H{"message": "OK"})
}
