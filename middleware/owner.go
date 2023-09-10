package middleware

import (
	"edd2023-back/logic"
	"edd2023-back/model"
	"edd2023-back/util"

	"github.com/gin-gonic/gin"
)

func OwnerCheck() gin.HandlerFunc {

	return func(ctx *gin.Context) {
		sessionId, _ := ctx.Cookie("sessionId")
		asl := logic.AuthSessionLogic{
			Session: model.AuthSession{Uuid: sessionId},
		}
		asl.GetByUuid()

		gameSessionId, _ := ctx.Cookie("gameSessionId")
		gsl := logic.GameSessionLogic{
			Session: model.GameSession{Uuid: gameSessionId},
		}
		gsl.GetByUuid(gameSessionId)

		if asl.Session.User.Id != gsl.Session.Room.OwnerPlayerId {
			s := "OnlyRoomOwner"
			r := util.BadRequest(&s)
			ctx.JSON(r.StatusCode, r.Message)
			ctx.Abort()
		}

		ctx.Next()
	}
}
