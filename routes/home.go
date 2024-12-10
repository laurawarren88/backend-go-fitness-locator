package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/laurawarren88/go_spa_backend.git/controllers"
)

func RegisterHomeRoutes(router *gin.Engine) {

	homeRoutes := router.Group("/project")
	{
		homeRoutes.GET("/home", controllers.GetHome)
	}
}
