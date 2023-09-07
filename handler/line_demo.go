package handler

import (
	"edd2023-back/logic"
	"edd2023-back/model"
	"edd2023-back/util"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

type LineDemoHandler struct{}

type LineDemoRequest struct {
	LineUid string `json:"lineUid"`
}

func (*LineDemoHandler) GenerateLineRegistrationOtp(ctx *gin.Context) {
	var req LineDemoRequest
	err := ctx.Bind(&req)
	if err != nil {
		s := err.Error()
		r := util.BadRequest(&s)
		log.Errorln("[Error]request parse error: ", s)
		ctx.JSON(r.StatusCode, r.Message)
		return
	}

	ll := logic.LineSessionLogic{
		Session: model.LineSession{
			LineUid: req.LineUid,
		},
	}
	err = ll.Create()
	if err != nil {
		log.Errorln("[Error]exec error: ", err.Error())
		r := util.InternalServerError(nil)
		ctx.JSON(r.StatusCode, r.Message)
		return
	}
	r := util.Ok(nil)
	ctx.JSON(r.StatusCode, gin.H{"otp": ll.Session.Otp})
}
