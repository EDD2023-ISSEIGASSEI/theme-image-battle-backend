package handler

import (
	"edd2023-back/logic"
	"edd2023-back/model"
	"net/http"

	"github.com/gin-gonic/gin"
)

type PhaseHandler struct{}

func (*PhaseHandler) GetPhase(ctx *gin.Context) {
	sessionId, _ := ctx.Cookie("gameSessionId")
	gsl := logic.GameSessionLogic{
		Session: model.GameSession{Uuid: sessionId},
	}
	gsl.GetByUuid(sessionId)

	ctx.JSON(http.StatusOK, gin.H{"phase": gsl.Session.Phase})
}

func (*PhaseHandler) GetPhaseState(ctx *gin.Context) {
	sessionId, _ := ctx.Cookie("gameSessionId")
	gsl := logic.GameSessionLogic{
		Session: model.GameSession{Uuid: sessionId},
	}
	gsl.GetByUuid(sessionId)

	var state any
	switch gsl.Session.Phase {
	case model.WaitingPhase:
		wpsl := logic.WaitingPhaseStateLigic{}
		wpsl.FromGameSession(gsl.Session)
		state = wpsl.State
	case model.GeneratePhase:
		gpsl := logic.GeneratePhaseStateLigic{}
		gpsl.FromGameSession(gsl.Session)
		state = gpsl.State
	case model.GuessPhase:
		gpsl := logic.GuessPhaseStateLigic{}
		gpsl.FromGameSession(gsl.Session)
		state = gpsl.State
	case model.ShowScorePhase:
		spsl := logic.ShowScorePhaseStateLigic{}
		spsl.FromGameSession(gsl.Session)
		state = spsl.State
	case model.ShowCorrectAnswerPhase:
		spsl := logic.ShowCorrectAnswerPhaseStateLigic{}
		spsl.FromGameSession(gsl.Session)
		state = spsl.State
	case model.EndingPhase:
		epsl := logic.EndingPhaseStateLigic{}
		epsl.FromGameSession(gsl.Session)
		state = epsl.State
	}

	ctx.JSON(http.StatusOK, model.PhaseStateResponse{
		Phase: gsl.Session.Phase,
		State: state,
	})
}
