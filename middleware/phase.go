package middleware

import (
	"edd2023-back/logic"
	"edd2023-back/model"
	"edd2023-back/util"

	"github.com/gin-gonic/gin"
)

func WaitingPhaseCheck() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		sessionId, _ := ctx.Cookie("gameSessionId")

		gsl := logic.GameSessionLogic{
			Session: model.GameSession{Uuid: sessionId},
		}
		gsl.GetByUuid(sessionId)

		if gsl.Session.Phase != model.WaitingPhase {
			s := "NotWaitingPhases"
			r := util.BadRequest(&s)
			ctx.JSON(r.StatusCode, r.Message)
			ctx.Abort()
		}
	}
}

func GeneratePhaseCheck() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		sessionId, _ := ctx.Cookie("gameSessionId")

		gsl := logic.GameSessionLogic{
			Session: model.GameSession{Uuid: sessionId},
		}
		gsl.GetByUuid(sessionId)

		if gsl.Session.Phase != model.GeneratePhase {
			s := "NotGeneratePhases"
			r := util.BadRequest(&s)
			ctx.JSON(r.StatusCode, r.Message)
			ctx.Abort()
		}
	}
}

func GuessPhaseCheck() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		sessionId, _ := ctx.Cookie("gameSessionId")

		gsl := logic.GameSessionLogic{
			Session: model.GameSession{Uuid: sessionId},
		}
		gsl.GetByUuid(sessionId)

		if gsl.Session.Phase != model.GuessPhase {
			s := "NotGuessPhases"
			r := util.BadRequest(&s)
			ctx.JSON(r.StatusCode, r.Message)
			ctx.Abort()
		}
	}
}

func ShowScorePhaseCheck() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		sessionId, _ := ctx.Cookie("gameSessionId")

		gsl := logic.GameSessionLogic{
			Session: model.GameSession{Uuid: sessionId},
		}
		gsl.GetByUuid(sessionId)

		if gsl.Session.Phase != model.ShowScorePhase {
			s := "NotShowScorePhases"
			r := util.BadRequest(&s)
			ctx.JSON(r.StatusCode, r.Message)
			ctx.Abort()
		}
		ctx.Next()
	}
}

func ShowCorrectAnswerPhaseCheck() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		sessionId, _ := ctx.Cookie("gameSessionId")

		gsl := logic.GameSessionLogic{
			Session: model.GameSession{Uuid: sessionId},
		}
		gsl.GetByUuid(sessionId)

		if gsl.Session.Phase != model.ShowCorrectAnswerPhase {
			s := "NotShowCorrectAnswerPhases"
			r := util.BadRequest(&s)
			ctx.JSON(r.StatusCode, r.Message)
			ctx.Abort()
		}
		ctx.Next()
	}
}

func EndingPhaseCheck() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		sessionId, _ := ctx.Cookie("gameSessionId")

		gsl := logic.GameSessionLogic{
			Session: model.GameSession{Uuid: sessionId},
		}
		gsl.GetByUuid(sessionId)

		if gsl.Session.Phase != model.EndingPhase {
			s := "NotEndingPhases"
			r := util.BadRequest(&s)
			ctx.JSON(r.StatusCode, r.Message)
			ctx.Abort()
		}
		ctx.Next()
	}
}
