package main

import (
	"fmt"
	"line-bot-otp-back/db"
	"os"

	"github.com/line/line-bot-sdk-go/v7/linebot"
	log "github.com/sirupsen/logrus"

	"github.com/gin-gonic/gin"
)

func main() {
	log.SetLevel(log.DebugLevel)

	bot, err := linebot.New(
		os.Getenv("LINEBOT_CHANNEL_SECRET"),
		os.Getenv("LINEBOT_CHANNEL_TOKEN"),
	)
	if err != nil {
		log.Fatalln(err)
	}

	db.Init()
	db.InitRedis()
	engine := gin.Default()
	DefineRoutes(engine, bot)
	fmt.Printf("hoge: %s\n", os.Getenv("SERVER_PORT"))
	engine.Run(fmt.Sprintf(":%s", os.Getenv("SERVER_PORT")))
	defer db.Db.Close()
	defer db.Redis.Close()
}
