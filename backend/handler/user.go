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
	sl := logic.SignUpSessionLogic{
		Session: model.SignUpSession{
			User: user,
		},
	}
	err = sl.CreateSession()
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

	sl := logic.SignInSessionLogic{
		Session: model.SignInSession{
			User: *ul.User,
		},
	}
	sl.CreateSession()
	log.Debugln("OTP:: ", sl.Session.Otp)

	r := util.Ok(nil)
	ctx.JSON(r.StatusCode, gin.H{"user": ul.User})
}

type LineRegistrationRequest struct {
	Otp       string `json:"otp"`
	SessionId string `json:"sessionId"`
}

func (*UserHandler) LineRegistration(ctx *gin.Context) {
	var req LineRegistrationRequest
	err := ctx.Bind(&req)
	if err != nil {
		s := err.Error()
		r := util.BadRequest(&s)
		log.Errorln("[Error]request parse error: ", s)
		ctx.JSON(r.StatusCode, r.Message)
		return
	}

	sl := logic.SignUpSessionLogic{
		Session: model.SignUpSession{Uuid: req.SessionId},
	}
	f, err := sl.GetByUuid()
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

	f, err = sl.LineRegisterByOtp(req.Otp)
	if !f && err == nil {
		s := "InvalidOTP"
		r := util.BadRequest(&s)
		ctx.JSON(r.StatusCode, r.Message)
		return
	}
	if err != nil {
		log.Errorln("[Error]exec error: ", err.Error())
		r := util.InternalServerError(nil)
		ctx.JSON(r.StatusCode, r.Message)
	}

	sl.DeleteSession()

	ul := logic.UserLigic{User: &sl.Session.User}
	err = ul.Create()
	if err != nil {
		log.Errorln("[Error]exec error: ", err.Error())
		r := util.InternalServerError(nil)
		ctx.JSON(r.StatusCode, r.Message)
		return
	}
	r := util.Ok(nil)
	ctx.JSON(r.StatusCode, gin.H{"user": ul.User})
}

type CheckOtpRequest struct {
	SessionId string `json:"sessionId"`
	Otp       string `json:"otp"`
}

func (*UserHandler) CheckOtp(ctx *gin.Context) {
	var req CheckOtpRequest
	err := ctx.Bind(&req)
	if err != nil {
		s := err.Error()
		r := util.BadRequest(&s)
		log.Errorln("[Error]request parse error: ", s)
		ctx.JSON(r.StatusCode, r.Message)
		return
	}

	sl := logic.SignInSessionLogic{
		Session: model.SignInSession{
			Uuid: req.SessionId,
		},
	}
	f, err := sl.GetByUuid()
	if !f && err == nil {
		s := "InvalidSessionID"
		r := util.BadRequest(&s)
		ctx.JSON(r.StatusCode, r.Message)
		return
	}
	if err != nil {
		log.Errorln("[Error]exec error: ", err.Error())
		r := util.InternalServerError(nil)
		ctx.JSON(r.StatusCode, r.Message)
		return
	}

	f = sl.CheckOtp(req.Otp)
	if f {
		r := util.Ok(nil)
		ctx.JSON(r.StatusCode, gin.H{"user": sl.Session.User})
		sl.DeleteSession()
	} else {
		s := "InvalidOTP"
		r := util.BadRequest(&s)
		ctx.JSON(r.StatusCode, r.Message)
		return
	}
}
