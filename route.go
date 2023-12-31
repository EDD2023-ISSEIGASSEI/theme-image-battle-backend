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

	app.GET("/user", userHandler.ValidateSessionId)

	roomHandler := handler.RoomHandler{}
	app.POST("/room", roomHandler.CreateRoom)
	app.GET("/room/list", roomHandler.ReadAllRooms)
	app.POST("/room/join", roomHandler.JoinRoom)

	game := app.Group("game")
	game.Use(middleware.GameSessionCheck())

	phaseHandler := handler.PhaseHandler{}
	game.GET("gamePhase", phaseHandler.GetPhase)
	game.GET("phaseState", phaseHandler.GetPhaseState)

	wait := game.Group("")
	wait.Use(middleware.WaitingPhaseCheck())
	gameHandler := handler.GameHandler{}
	wait.POST("start", gameHandler.GameStart)
	wait.POST("closeRoom", roomHandler.CloseRoom)

	gene := game.Group("")
	gene.Use(middleware.GeneratePhaseCheck())
	topicHandler := handler.TopicHandler{}
	gene.GET("/topic", topicHandler.GetTopic)
	promptHandler := handler.PromptHandler{}
	gene.POST("/prompt", promptHandler.SubmitPrompt)

	guess := game.Group("")
	guess.Use(middleware.GuessPhaseCheck())
	answerHandler := handler.AnswerHandler{}
	guess.POST("/answer", answerHandler.SubmitAnswer)

	showScore := game.Group("showScore")
	showScore.Use(middleware.ShowScorePhaseCheck(), middleware.OwnerCheck())
	showScoreHandler := handler.ShwoScoreHandler{}
	showScore.POST("/next", showScoreHandler.Next)
	showScore.POST("/prev", showScoreHandler.Prev)
	showScore.POST("/end", showScoreHandler.End)

	showAns := game.Group("")
	showAns.Use(middleware.ShowCorrectAnswerPhaseCheck(), middleware.OwnerCheck())
	shwoCorrectAnswereHandler := handler.ShwoCorrectAnswereHandler{}
	showAns.POST("/showCorrectAnswer/next", shwoCorrectAnswereHandler.NextDealer)
	roundHandler := handler.RoundHandler{}
	showAns.POST("nextRound", roundHandler.NextRound)
	endingHandler := handler.EndingHandler{}
	showAns.POST("ending", endingHandler.Ending)

	end := game.Group("")
	end.Use(middleware.EndingPhaseCheck(), middleware.OwnerCheck())
	waitingHandler := handler.WaitingHandler{}
	end.POST("waiting", waitingHandler.Waiting)
}
