package handler

import (
	"edd2023-back/logic"
	"edd2023-back/model"
	"net/http"

	"github.com/gin-gonic/gin"
)

type WaitingHandler struct{}

func (*WaitingHandler) Waiting(ctx *gin.Context) {
	gameSessionId, _ := ctx.Cookie("gameSessionId")
	gsl := logic.GameSessionLogic{}
	gsl.GetByUuid(gameSessionId)

	gsl.Session.Phase = model.WaitingPhase
	gsl.UpdateByUuId()
	// TODO broadcast
	ctx.JSON(http.StatusOK, gin.H{"message": "OK"})
}
