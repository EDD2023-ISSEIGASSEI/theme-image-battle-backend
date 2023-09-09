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

func (*RoomHandler) ReadAllRooms(ctx *gin.Context) {
	keys, err := db.RoomRedis.Keys(context.Background(), "*").Result()
	if err != nil {
		log.Errorln("[Error]exec error: ", err.Error())
		r := util.InternalServerError(nil)
		ctx.JSON(r.StatusCode, r.Message)
		return
	}

	var rooms []model.Room
	rsl := logic.RoomSessionLogic{}
	for _, k := range keys {
		f, err := rsl.GetRoomSessionById(k)
		if f && err != nil {
			log.Errorln("[Error]exec error: ", err.Error())
			r := util.InternalServerError(nil)
			ctx.JSON(r.StatusCode, r.Message)
		}
		rooms = append(rooms, rsl.Session.Room)
	}

	ctx.JSON(http.StatusOK, model.RoomListResponse{Rooms: rooms})
}

func (*RoomHandler) JoinRoom(ctx *gin.Context) {
	var req model.JoinRoomRequest
	err := ctx.Bind(&req)
	if err != nil {
		s := err.Error()
		r := util.BadRequest(&s)
		log.Errorln("[Error]request parse error: ", s)
		ctx.JSON(r.StatusCode, r.Message)
		return
	}

	rsl := logic.RoomSessionLogic{}
	f, err := rsl.GetRoomSessionById(req.Id)
	if !f && err != nil {
		s := "InvalidRoomId"
		r := util.BadRequest(&s)
		log.Errorln("[Error]request parse error: ", err.Error())
		ctx.JSON(r.StatusCode, r.Message)
		return
	}
	if f && err != nil {
		log.Errorln("[Error]exec error: ", err.Error())
		r := util.InternalServerError(nil)
		ctx.JSON(r.StatusCode, r.Message)
		return
	}

	// sessionIdからPlayer取得
	sessionId, err := ctx.Cookie("sessionId")
	asl := logic.AuthSessionLogic{
		Session: model.AuthSession{Uuid: sessionId},
	}
	asl.GetByUuid()
	player := model.UserToPlayer(asl.Session.User)

	// RoomSessionに紐づいたGameSessionを取得
	gsl := logic.GameSessionLogic{}
	f, err = gsl.GetByUuid(rsl.Session.GameSessionId)
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
		return
	}

	// パスワードとキャパの確認
	rl := logic.RoomLogic{
		Room: rsl.Session.Room,
	}
	if !rl.VaridatePassword(req.Password) {
		s := "InvalidPassword"
		r := util.BadRequest(&s)
		ctx.JSON(r.StatusCode, r.Message)
		return
	}
	log.Debugln(rl.CanJoin())
	log.Debugln(rl.Room.PlayerNum)
	log.Debugln(rl.Room.MaxPlayerNum)
	log.Debugln(rl.Room.PlayerNum + 1)
	log.Debugln(rl.Room.PlayerNum+1 <= rl.Room.MaxPlayerNum)
	if !rl.CanJoin() {
		s := "CapacityOver"
		r := util.BadRequest(&s)
		ctx.JSON(r.StatusCode, r.Message)
		return
	}

	// GameにUserが重複しないか確認
	if gsl.IsExistsPlayer(player) {
		s := "ExistsUser"
		r := util.BadRequest(&s)
		ctx.JSON(r.StatusCode, r.Message)
		return
	}

	// RoomにJoinしてセッションアップデート
	rl.Join()
	err = rsl.UpdateRoomInfo(rl.Room)
	if err != nil {
		log.Errorln("[Error]exec error: ", err.Error())
		r := util.InternalServerError(nil)
		ctx.JSON(r.StatusCode, r.Message)
		return
	}

	// GameにJoinしてセッションアップデート
	gsl.JoinPlayer(player)
	err = gsl.UpdateByUuId()
	if err != nil {
		log.Errorln("[Error]exec error: ", err.Error())
		r := util.InternalServerError(nil)
		ctx.JSON(r.StatusCode, r.Message)
		return
	}

	// broadcast WaitiongPhase
	ctx.JSON(http.StatusOK, gin.H{"message": "OK"})
}
