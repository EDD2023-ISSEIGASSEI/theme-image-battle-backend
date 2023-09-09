package middleware

import (
	"edd2023-back/logic"
	"edd2023-back/model"
	"edd2023-back/util"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

func AuthSessionCheck() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		sessionId, err := ctx.Cookie("sessionId")
		if err != nil {
			s := err.Error()
			r := util.BadRequest(&s)
			log.Errorln("[Error]request parse error: ", s)
			ctx.JSON(r.StatusCode, r.Message)
		}

		asl := logic.AuthSessionLogic{
			Session: model.AuthSession{Uuid: sessionId},
		}
		f, err := asl.GetByUuid()
		if !f && err != nil {
			s := "InvalidSessionId"
			r := util.BadRequest(&s)
			log.Errorln("[Error]request parse error: ", err.Error())
			ctx.JSON(r.StatusCode, r.Message)
			return
		}
		if f && err != nil {
			log.Errorln("[Error]exec error: ", err.Error())
			r := util.InternalServerError(nil)
			ctx.JSON(r.StatusCode, r.Message)
		}

		ctx.Next()
	}
}
