package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/laurawarren88/go_spa_backend.git/controllers"
)

func RegisterGymRoutes(router *gin.Engine, gc *controllers.GymController) {
	gymRoutes := router.Group("/gyms")
	{
		gymRoutes.GET("/", gc.GetAllGyms)
		gymRoutes.GET("/new", gc.RenderCreateGymForm)
		gymRoutes.POST("/new", gc.CreateGym)
		gymRoutes.GET("/:id", gc.GetGymById)
		gymRoutes.GET("/:id/edit", gc.RenderEditGymForm)
		gymRoutes.PUT("/:id/edit", gc.UpdateGym)
		gymRoutes.DELETE("/:id", gc.DeleteGym)
	}
}
