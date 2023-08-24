package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func DefineRoutes(r gin.IRouter) {
	r.GET("/", func(ctx *gin.Context) {
		ctx.JSONP(http.StatusOK, gin.H{
			"message": "ok",
		})
	})
}
