package handler

import (
	"line-bot-otp-back/logic"
	"line-bot-otp-back/model"
	"line-bot-otp-back/util"

	log "github.com/sirupsen/logrus"

	"github.com/gin-gonic/gin"
)

type UserHandler struct{}

type UserRequest struct {
	Name     string `json:"name"`
	Password string `json:"password"`
}

func (*UserHandler) SignUp(ctx *gin.Context) {
	var req UserRequest
	err := ctx.Bind(&req)
	if err != nil {
		s := err.Error()
		r := util.BadRequest(&s)
		log.Errorln("[Error]request parse error: ", s)
		ctx.JSON(r.StatusCode, r.Message)
		return
	}

	user := model.User{
		Name:     req.Name,
		Password: req.Password,
	}
	sl := logic.SignUpSessionLigic{
		Session: model.SignUpSession{
			User: user,
		},
	}
	err = sl.Create()
	if err != nil {
		log.Errorln("[Error]exec error: ", err.Error())
		r := util.InternalServerError(nil)
		ctx.JSON(r.StatusCode, r.Message)
		return
	}
	r := util.Ok(nil)
	ctx.JSON(r.StatusCode, gin.H{"sessionId": sl.Session.Uuid})
}

func (*UserHandler) SignIn(ctx *gin.Context) {
	var req UserRequest
	err := ctx.Bind(&req)
	if err != nil {
		s := err.Error()
		r := util.BadRequest(&s)
		log.Errorln("[Error]request parse error: ", s)
		ctx.JSON(r.StatusCode, r.Message)
		return
	}
	ul := logic.UserLigic{User: &model.User{Name: req.Name, Password: req.Password}}
	err = ul.SelectByNameAndPass()
	if err != nil {
		log.Errorln("[Error]exec error: ", err.Error())
		r := util.InternalServerError(nil)
		ctx.JSON(r.StatusCode, r.Message)
		return
	}
	r := util.Ok(nil)
	ctx.JSON(r.StatusCode, gin.H{"user": ul.User})
}
