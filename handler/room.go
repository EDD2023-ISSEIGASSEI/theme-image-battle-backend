package handler

import (
	"context"
	"edd2023-back/db"
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
	asl := logic.AuthSessionLogic{
		Session: model.AuthSession{Uuid: sessionId},
	}
	asl.GetByUuid()

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

type RoomListRes struct {
	Rooms []model.Room `json:"rooms"`
}

func (*RoomHandler) ReadAllRooms(ctx *gin.Context) {
	keys, err := db.RoomRedis.Keys(context.Background(), "*").Result()
	if err != nil {
		log.Fatal(err)
	}

	var rooms []model.Room
	rsl := logic.RoomSessionLogic{}
	for _, k := range keys {
		rsl.GetRoomSessionById(k)
		rooms = append(rooms, rsl.Session.Room)
	}

	ctx.JSON(http.StatusOK, RoomListRes{Rooms: rooms})
}
