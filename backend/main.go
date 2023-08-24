package main

import (
	"fmt"
	"os"

	"line-bot-otp-back/db"

	"github.com/gin-gonic/gin"
)

func main() {
	db.Init()
	engine := gin.Default()
	DefineRoutes(engine)
	fmt.Printf("hoge: %s\n", os.Getenv("SERVER_PORT"))
	engine.Run(fmt.Sprintf(":%s", os.Getenv("SERVER_PORT")))
}
