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

type SignUpRequest struct {
	Id       string `json:"id"`
	Name     string `json:"name"`
	Password string `json:"password"`
}

func (*UserHandler) SignUp(ctx *gin.Context) {
	var req SignUpRequest
	err := ctx.Bind(&req)
	if err != nil {
		s := err.Error()
		r := util.BadRequest(&s)
		log.Errorln("[Error]request parse error: ", s)
		ctx.JSON(r.StatusCode, r.Message)
		return
	}

	user := model.User{
		Id:       req.Id,
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

	cookie := http.Cookie{
		Name:     "sessionId",
		Value:    sl.Session.Uuid,
		MaxAge:   0,
		Path:     "/",
		Domain:   "",
		HttpOnly: true,
		SameSite: http.SameSiteNoneMode,
		Secure:   true,
	}
	http.SetCookie(ctx.Writer, &cookie)
	r := util.Ok(nil)
	ctx.JSON(r.StatusCode, r.Message)
}

type SignInRequest struct {
	Id       string `json:"id"`
	Password string `json:"password"`
}

func (*UserHandler) SignIn(ctx *gin.Context, bot *linebot.Client) {
	var req SignInRequest
	err := ctx.Bind(&req)
	if err != nil {
		s := err.Error()
		r := util.BadRequest(&s)
		log.Errorln("[Error]request parse error: ", s)
		ctx.JSON(r.StatusCode, r.Message)
		return
	}
	ul := logic.UserLigic{User: &model.User{Id: req.Id, Password: req.Password}}

	f, err := ul.IdIsExists()
	if err != nil {
		log.Errorln("[Error]exec error: ", err.Error())
		r := util.InternalServerError(nil)
		ctx.JSON(r.StatusCode, r.Message)
		return
	}
	if !f {
		s := "InvalidId"
		r := util.BadRequest(&s)
		log.Errorln("[Error]request parse error: ", s)
		ctx.JSON(r.StatusCode, r.Message)
		return
	}

	f, err = ul.SelectById()
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

	if !ul.VaridatePassword(req.Password) {
		s := "InvalidPassword"
		r := util.BadRequest(&s)
		ctx.JSON(r.StatusCode, r.Message)
		return
	}

	sl := logic.SignInSessionLogic{
		Session: model.SignInSession{
			User: *ul.User,
		},
	}
	sl.CreateSession()

	cookie := http.Cookie{
		Name:     "sessionId",
		Value:    sl.Session.Uuid,
		MaxAge:   0,
		Path:     "/",
		Domain:   "",
		HttpOnly: true,
		SameSite: http.SameSiteNoneMode,
		Secure:   true,
	}
	http.SetCookie(ctx.Writer, &cookie)

	message := linebot.NewTextMessage("↓ワンタイムパスワード↓\n" + sl.Session.Otp)
	_, err = bot.PushMessage(*sl.Session.User.LineUid, message).Do()
	if err != nil {
		log.Errorln(err.Error())
	}

	r := util.Ok(nil)
	ctx.JSON(r.StatusCode, r.Message)
}

type LineRegistrationRequest struct {
	Otp string `json:"otp"`
	// SessionId string `json:"sessionId"`
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

	sessionId, err := ctx.Cookie("sessionId")
	if err != nil {
		s := err.Error()
		r := util.BadRequest(&s)
		log.Errorln("[Error]request parse error: ", s)
		ctx.JSON(r.StatusCode, r.Message)
	}

	sl := logic.SignUpSessionLogic{
		Session: model.SignUpSession{Uuid: sessionId},
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

	cookie := http.Cookie{
		Name:     "sessionId",
		Value:    al.Session.Uuid,
		MaxAge:   0,
		Path:     "/",
		Domain:   "",
		HttpOnly: true,
		SameSite: http.SameSiteNoneMode,
		Secure:   true,
	}
	http.SetCookie(ctx.Writer, &cookie)

	ctx.JSON(http.StatusOK, al.Session.User)
}

type CheckOtpRequest struct {
	// SessionId string `json:"sessionId"`
	Otp string `json:"otp"`
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

	sessionId, err := ctx.Cookie("sessionId")
	if err != nil {
		s := err.Error()
		r := util.BadRequest(&s)
		log.Errorln("[Error]request parse error: ", s)
		ctx.JSON(r.StatusCode, r.Message)
	}

	sl := logic.SignInSessionLogic{
		Session: model.SignInSession{
			Uuid: sessionId,
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

		cookie := http.Cookie{
			Name:     "sessionId",
			Value:    al.Session.Uuid,
			MaxAge:   0,
			Path:     "/",
			Domain:   "",
			HttpOnly: true,
			SameSite: http.SameSiteNoneMode,
			Secure:   true,
		}
		http.SetCookie(ctx.Writer, &cookie)

		ctx.JSON(http.StatusOK, al.Session.User)
		sl.DeleteSession()
	} else {
		s := "InvalidOTP"
		r := util.BadRequest(&s)
		ctx.JSON(r.StatusCode, r.Message)
		return
	}
}

type AuthSessionRequest struct {
	// SessionId string `json:"sessionId"`
}

func (*UserHandler) ValidateSessionId(ctx *gin.Context) {
	// var req AuthSessionRequest
	// err := ctx.Bind(&req)
	// if err != nil {
	// 	s := err.Error()
	// 	r := util.BadRequest(&s)
	// 	log.Errorln("[Error]request parse error: ", s)
	// 	ctx.JSON(r.StatusCode, r.Message)
	// 	return
	// }

	sessionId, err := ctx.Cookie("sessionId")
	if err != nil {
		s := err.Error()
		r := util.BadRequest(&s)
		log.Errorln("[Error]request parse error: ", s)
		ctx.JSON(r.StatusCode, r.Message)
	}

	al := logic.AuthSessionLogic{
		Session: model.AuthSession{
			Uuid: sessionId,
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

	sessionId, err := ctx.Cookie("sessionId")
	if err != nil {
		s := err.Error()
		r := util.BadRequest(&s)
		log.Errorln("[Error]request parse error: ", s)
		ctx.JSON(r.StatusCode, r.Message)
	}

	al := logic.AuthSessionLogic{
		Session: model.AuthSession{
			Uuid: sessionId,
		},
	}
	al.DeleteSession()
	r := util.Ok(nil)
	ctx.JSON(r.StatusCode, r.Message)
}

type IdIsExistsRequest struct {
	Id string `json:"id"`
}

func (*UserHandler) IdIsExists(ctx *gin.Context) {
	var req IdIsExistsRequest
	err := ctx.Bind(&req)
	if err != nil {
		s := err.Error()
		r := util.BadRequest(&s)
		log.Errorln("[Error]request parse error: ", s)
		ctx.JSON(r.StatusCode, r.Message)
		return
	}

	ul := logic.UserLigic{
		User: &model.User{
			Id: req.Id,
		},
	}
	f, err := ul.IdIsExists()
	if err != nil {
		log.Errorln("[Error]exec error: ", err.Error())
		r := util.InternalServerError(nil)
		ctx.JSON(r.StatusCode, r.Message)
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"isExists": f})
}
