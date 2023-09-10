package handler

import (
	"edd2023-back/logic"
	"edd2023-back/model"
	"net/http"

	"github.com/gin-gonic/gin"
)

type EndingHandler struct{}

func (*EndingHandler) Ending(ctx *gin.Context) {
	gameSessionId, _ := ctx.Cookie("gameSessionId")
	gsl := logic.GameSessionLogic{}
	gsl.GetByUuid(gameSessionId)

	gsl.Session.Phase = model.EndingPhase
	gsl.UpdateByUuId()
	// ToDo: broadcast
	ctx.JSON(http.StatusOK, gin.H{"message": "OK"})
}
