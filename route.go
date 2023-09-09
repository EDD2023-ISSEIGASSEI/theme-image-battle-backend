package main

import (
	"edd2023-back/handler"
	"edd2023-back/middleware"
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

	g.GET("/", func(ctx *gin.Context) {
		ctx.JSONP(http.StatusOK, gin.H{
			"message": "ok",
		})
	})

	auth := g.Group("auth")

	userHandler := handler.UserHandler{}
	auth.POST("/signUp", userHandler.SignUp)
	auth.POST("/lineRegistration", userHandler.LineRegistration)
	auth.POST("/signIn", func(ctx *gin.Context) { userHandler.SignIn(ctx, bot) })
	auth.POST("/idIsExists", userHandler.IdIsExists)
	auth.POST("/checkOtp", userHandler.CheckOtp)
	auth.POST("/varidateSessionId", userHandler.ValidateSessionId)
	auth.POST("/signOut", userHandler.SignOut)

	// lineDemoHandler := handler.LineDemoHandler{}
	// auth.POST("/lineDemo", lineDemoHandler.GenerateLineRegistrationOtp)

	linebotHandler := handler.LinebotHandler{}
	g.POST("/linebot", func(ctx *gin.Context) {
		linebotHandler.EventHandler(ctx, bot)
	})

	app := g.Group("")
	app.Use(middleware.AuthSessionCheck())

	roomHandler := handler.RoomHandler{}
	app.POST("/room", roomHandler.CreateRoom)
	app.GET("/room/list", roomHandler.ReadAllRooms)
	app.POST("/room/join", roomHandler.JoinRoom)
}
