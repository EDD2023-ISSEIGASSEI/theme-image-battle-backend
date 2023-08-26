package main

import (
	"line-bot-otp-back/handler"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/line/line-bot-sdk-go/v7/linebot"
)

func DefineRoutes(r gin.IRouter, bot *linebot.Client) {
	// r.Use(cors.New(cors.Config{
	// 	AllowOrigins: []string{
	// 		"https://localhost:5173",
	// 		// "*",
	// 	},
	// 	AllowMethods: []string{
	// 		"POST",
	// 		"PUT",
	// 		"GET",
	// 		"OPTIONS",
	// 	},
	// 	AllowHeaders: []string{
	// 		"Access-Control-Allow-Credentials",
	// 		"Access-Control-Allow-Headers",
	// 		"Content-Type",
	// 		"Content-Length",
	// 		"Accept-Encoding",
	// 		"Authorization",
	// 		"Origin",
	// 		"Cookie",
	// 		"Set-Cookie",
	// 	},
	// 	AllowCredentials: true,
	// }))

	g := r.Group("api")

	r.GET("/", func(ctx *gin.Context) {
		ctx.JSONP(http.StatusOK, gin.H{
			"message": "ok",
		})
	})

	userHandler := handler.UserHandler{}
	g.POST("/signUp", userHandler.SignUp)
	g.POST("/lineRegistration", userHandler.LineRegistration)
	g.POST("/signIn", func(ctx *gin.Context) { userHandler.SignIn(ctx, bot) })
	g.POST("/idIsExists", userHandler.IdIsExists)
	g.POST("/checkOtp", userHandler.CheckOtp)
	g.POST("/varidateSessionId", userHandler.ValidateSessionId)
	g.POST("/signOut", userHandler.SignOut)

	// lineDemoHandler := handler.LineDemoHandler{}
	// g.POST("/lineDemo", lineDemoHandler.GenerateLineRegistrationOtp)

	linebotHandler := handler.LinebotHandler{}
	g.POST("/linebot", func(ctx *gin.Context) {
		linebotHandler.EventHandler(ctx, bot)
	})
}
