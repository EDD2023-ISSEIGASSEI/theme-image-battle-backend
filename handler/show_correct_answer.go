package handler

import (
	"edd2023-back/logic"
	"edd2023-back/util"
	"net/http"

	"github.com/gin-gonic/gin"
)

type ShwoCorrectAnswereHandler struct{}

func (*ShwoCorrectAnswereHandler) NextDealer(ctx *gin.Context) {
	gameSessionId, _ := ctx.Cookie("gameSessionId")
	gsl := logic.GameSessionLogic{}
	gsl.GetByUuid(gameSessionId)

	dealerNum := 0
	for _, p := range gsl.Session.Players {
		if p.Id == gsl.Session.ShowingPlayerId {
			break
		}
		dealerNum++
	}

	nextNum := dealerNum + 1
	if nextNum >= len(gsl.Session.Players) {
		s := "NoPrevDealer"
		r := util.BadRequest(&s)
		ctx.JSON(r.StatusCode, r.Message)
		return
	}

	nextShowingNum := 0

	gsl.Session.DealerPlayerId = gsl.Session.Players[nextNum].Id
	gsl.Session.ShowingPlayerId = gsl.Session.Players[nextShowingNum].Id
	gsl.UpdateByUuId()
	spsl := logic.ShowScorePhaseStateLigic{}
	spsl.FromGameSession(gsl.Session)
	gsl.Session.PlayerStates[nextShowingNum].Score += spsl.State.PlayerAnswer.Answer.Score
	gsl.UpdateByUuId()
	// ToDo: broadcast
	ctx.JSON(http.StatusOK, gin.H{"message": "OK"})
}
