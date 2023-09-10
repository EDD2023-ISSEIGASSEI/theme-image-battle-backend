package handler

import (
	"edd2023-back/logic"
	"edd2023-back/util"
	"net/http"

	"github.com/gin-gonic/gin"
)

type ShwoScoreHandler struct{}

func (*ShwoScoreHandler) Next(ctx *gin.Context) {
	gameSessionId, _ := ctx.Cookie("gameSessionId")
	gsl := logic.GameSessionLogic{}
	gsl.GetByUuid(gameSessionId)

	showingNum := 0
	for _, p := range gsl.Session.Players {
		showingNum++
		if p.Id == gsl.Session.ShowingPlayerId {
			break
		}
	}
	// 表示されているユーザーの次のユーザーがいない
	// もしくは、表示されているユーザーの次のユーザーがオーナーである
	if showingNum+1 > len(gsl.Session.Players)-1 ||
		(showingNum+1 == len(gsl.Session.Players)-1 &&
			gsl.Session.Players[showingNum+1].Id == gsl.Session.DealerPlayerId) {
		s := "NoNextPlayer"
		r := util.BadRequest(&s)
		ctx.JSON(r.StatusCode, r.Message)
		return
	}

	nextNum := showingNum + 1
	if gsl.Session.Players[nextNum].Id == gsl.Session.DealerPlayerId {
		nextNum += 1
	}

	spsl := logic.ShowScorePhaseStateLigic{}
	spsl.FromGameSession(gsl.Session)
	gsl.Session.PlayerStates[showingNum].Score += spsl.State.PlayerAnswer.Answer.Score
	gsl.Session.ShowingPlayerId = gsl.Session.Players[nextNum].Id
	gsl.UpdateByUuId()
	// ToDo: broadcast
	ctx.JSON(http.StatusOK, gin.H{"message": "OK"})
}
