package main

import (
	"line-bot-otp-back/handler"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/line/line-bot-sdk-go/v7/linebot"
)

func DefineRoutes(r gin.IRouter, bot *linebot.Client) {
	r.GET("/", func(ctx *gin.Context) {
		ctx.JSONP(http.StatusOK, gin.H{
			"message": "ok",
		})
	})

	userHandler := handler.UserHandler{}
	r.POST("/signUp", userHandler.SignUp)
	r.POST("/lineRegistration", userHandler.LineRegistration)
	r.POST("/signIn", func(ctx *gin.Context) { userHandler.SignIn(ctx, bot) })
	r.POST("/checkOtp", userHandler.CheckOtp)

	lineDemoHandler := handler.LineDemoHandler{}
	r.POST("/lineDemo", lineDemoHandler.GenerateLineRegistrationOtp)

	linebotHandler := handler.LinebotHandler{}
	r.POST("/linebot", func(ctx *gin.Context) {
		linebotHandler.EventHandler(ctx, bot)
	})
}
