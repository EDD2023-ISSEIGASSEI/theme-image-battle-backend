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
