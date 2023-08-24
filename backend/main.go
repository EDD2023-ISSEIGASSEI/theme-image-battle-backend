package main

import (
	"fmt"
	"os"

	"github.com/gin-gonic/gin"
)

func main() {
	DbInit()
	engine := gin.Default()
	DefineRoutes(engine)
	fmt.Printf("hoge: %s\n", os.Getenv("SERVER_PORT"))
	engine.Run(fmt.Sprintf(":%s", os.Getenv("SERVER_PORT")))
}
