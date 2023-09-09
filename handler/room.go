package handler

import (
	"edd2023-back/logic"
	"edd2023-back/model"
	"edd2023-back/util"
	"net/http"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

type RoomHandler struct{}

func (*RoomHandler) CreateRoom(ctx *gin.Context) {
	var req model.CreateRoomRequest
	err := ctx.Bind(&req)
	if err != nil {
		s := err.Error()
		r := util.BadRequest(&s)
		log.Errorln("[Error]request parse error: ", s)
		ctx.JSON(r.StatusCode, r.Message)
		return
	}

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

	rl := logic.RoomLogic{}
	err = rl.CreateRoom(req, asl.Session.User.Id)
	if err != nil {
		log.Errorln("[Error]exec error: ", err.Error())
		r := util.InternalServerError(nil)
		ctx.JSON(r.StatusCode, r.Message)
		return
	}

	gsl := logic.GameSessionLogic{}
	err = gsl.CreateGameSession(model.UserToPlayer(asl.Session.User), rl.Room)
	if err != nil {
		log.Errorln("[Error]exec error: ", err.Error())
		r := util.InternalServerError(nil)
		ctx.JSON(r.StatusCode, r.Message)
		return
	}

	rsl := logic.RoomSessionLogic{}
	err = rsl.CreateRoomSession(rl.Room, gsl.Session.Uuid)
	if err != nil {
		log.Errorln("[Error]exec error: ", err.Error())
		r := util.InternalServerError(nil)
		ctx.JSON(r.StatusCode, r.Message)
		return
	}

	cookie := http.Cookie{
		Name:     "gameSessionId",
		Value:    gsl.Session.Uuid,
		MaxAge:   0,
		Path:     "/",
		Domain:   "",
		HttpOnly: true,
		SameSite: http.SameSiteNoneMode,
		Secure:   true,
	}
	http.SetCookie(ctx.Writer, &cookie)
	ctx.JSON(http.StatusOK, rl.Room)
}
