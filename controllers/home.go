package controllers

import "github.com/gin-gonic/gin"

func GetHome(ctx *gin.Context) {
	ctx.JSON(200, gin.H{
		"message": "Hello, World!",
	})
}
