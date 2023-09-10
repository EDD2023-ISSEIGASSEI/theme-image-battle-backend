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
		if p.Id == gsl.Session.ShowingPlayerId {
			break
		}
		showingNum++
	}

	nextNum := showingNum + 1
	for {
		if nextNum >= len(gsl.Session.Players) {
			s := "NoNextPlayer"
			r := util.BadRequest(&s)
			ctx.JSON(r.StatusCode, r.Message)
			return
		}
		if gsl.Session.Players[nextNum].Id != gsl.Session.DealerPlayerId {
			break
		}
		nextNum += 1
	}

	gsl.Session.ShowingPlayerId = gsl.Session.Players[nextNum].Id
	gsl.UpdateByUuId()
	spsl := logic.ShowScorePhaseStateLigic{}
	spsl.FromGameSession(gsl.Session)
	gsl.Session.PlayerStates[nextNum].Score += spsl.State.PlayerAnswer.Answer.Score
	gsl.UpdateByUuId()
	// ToDo: broadcast
	ctx.JSON(http.StatusOK, gin.H{"message": "OK"})
}

func (*ShwoScoreHandler) Prev(ctx *gin.Context) {
	gameSessionId, _ := ctx.Cookie("gameSessionId")
	gsl := logic.GameSessionLogic{}
	gsl.GetByUuid(gameSessionId)

	showingNum := 0
	for _, p := range gsl.Session.Players {
		if p.Id == gsl.Session.ShowingPlayerId {
			break
		}
		showingNum++
	}

	prevNum := showingNum - 1
	for {
		if prevNum < 0 {
			s := "NoPrevPlayer"
			r := util.BadRequest(&s)
			ctx.JSON(r.StatusCode, r.Message)
			return
		}
		if gsl.Session.Players[prevNum].Id != gsl.Session.DealerPlayerId {
			break
		}
		prevNum -= 1
	}

	gsl.Session.ShowingPlayerId = gsl.Session.Players[prevNum].Id
	gsl.UpdateByUuId()
	// ToDo: broadcast
	ctx.JSON(http.StatusOK, gin.H{"message": "OK"})
}
