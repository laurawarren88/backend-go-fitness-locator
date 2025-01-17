package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/laurawarren88/go_spa_backend.git/controllers"
)

func RegisterHomeRoutes(router *gin.Engine, hc *controllers.HomeController) {
	router.GET("/api/home", hc.GetHome)
}
