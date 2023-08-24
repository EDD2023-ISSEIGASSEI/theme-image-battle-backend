package main

import (
	"line-bot-otp-back/handler"
	"net/http"

	"github.com/gin-gonic/gin"
)

func DefineRoutes(r gin.IRouter) {
	r.GET("/", func(ctx *gin.Context) {
		ctx.JSONP(http.StatusOK, gin.H{
			"message": "ok",
		})
	})

	userHandler := handler.UserHandler{}
	r.POST("/signUp", userHandler.SignUp)
	r.POST("/signIn", userHandler.SignIn)
}
