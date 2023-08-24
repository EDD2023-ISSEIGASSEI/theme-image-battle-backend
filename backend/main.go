package main

import (
	"fmt"
	"line-bot-otp-back/db"
	"os"

	log "github.com/sirupsen/logrus"

	"github.com/gin-gonic/gin"
)

func main() {
	log.SetLevel(log.DebugLevel)
	db.Init()
	engine := gin.Default()
	DefineRoutes(engine)
	fmt.Printf("hoge: %s\n", os.Getenv("SERVER_PORT"))
	engine.Run(fmt.Sprintf(":%s", os.Getenv("SERVER_PORT")))
}
