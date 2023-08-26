package handler

import (
	"line-bot-otp-back/logic"
	"line-bot-otp-back/model"
	"line-bot-otp-back/util"
	"net/http"

	"github.com/line/line-bot-sdk-go/v7/linebot"
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

func (*UserHandler) SignIn(ctx *gin.Context, bot *linebot.Client) {
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
	f, err := ul.SelectById()
	if !f && err == nil {
		s := "InvalidId"
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

	sl := logic.SignInSessionLogic{
		Session: model.SignInSession{
			User: *ul.User,
		},
	}
	sl.CreateSession()
	message := linebot.NewTextMessage("↓ワンタイムパスワード↓\n" + sl.Session.Otp)
	_, err = bot.PushMessage(*sl.Session.User.LineUid, message).Do()
	if err != nil {
		log.Errorln(err.Error())
	}

	r := util.Ok(nil)
	ctx.JSON(r.StatusCode, gin.H{"sessionId": sl.Session.Uuid})
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
		return
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

	al := logic.AuthSessionLogic{
		Session: model.AuthSession{
			User: *ul.User,
		},
	}
	al.CreateSession()
	ctx.JSON(http.StatusOK, al.Session)
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
		al := logic.AuthSessionLogic{
			Session: model.AuthSession{
				User: sl.Session.User,
			},
		}
		al.CreateSession()
		ctx.JSON(http.StatusOK, al.Session)
		sl.DeleteSession()
	} else {
		s := "InvalidOTP"
		r := util.BadRequest(&s)
		ctx.JSON(r.StatusCode, r.Message)
		return
	}
}

type AuthSessionRequest struct {
	SessionId string `json:"sessionId"`
}

func (*UserHandler) ValidateSessionId(ctx *gin.Context) {
	var req AuthSessionRequest
	err := ctx.Bind(&req)
	if err != nil {
		s := err.Error()
		r := util.BadRequest(&s)
		log.Errorln("[Error]request parse error: ", s)
		ctx.JSON(r.StatusCode, r.Message)
		return
	}

	al := logic.AuthSessionLogic{
		Session: model.AuthSession{
			Uuid: req.SessionId,
		},
	}
	f, err := al.GetByUuid()
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

	ctx.JSON(http.StatusOK, al.Session)
}

func (*UserHandler) SignOut(ctx *gin.Context) {
	var req AuthSessionRequest
	err := ctx.Bind(&req)
	if err != nil {
		s := err.Error()
		r := util.BadRequest(&s)
		log.Errorln("[Error]request parse error: ", s)
		ctx.JSON(r.StatusCode, r.Message)
		return
	}

	al := logic.AuthSessionLogic{
		Session: model.AuthSession{
			Uuid: req.SessionId,
		},
	}
	al.DeleteSession()
	r := util.Ok(nil)
	ctx.JSON(r.StatusCode, r.Message)
}
