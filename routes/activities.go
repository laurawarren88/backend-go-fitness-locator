package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/laurawarren88/go_spa_backend.git/controllers"
)

func RegisterActivitiesRoutes(router *gin.Engine, ac *controllers.ActivitiesController) {
	activitiesRoutes := router.Group("/activities")
	{
		activitiesRoutes.GET("/", ac.GetAllActivities)
		activitiesRoutes.GET("/new", ac.RenderCreateActivityForm)
		activitiesRoutes.POST("/new", ac.CreateActivity)
		activitiesRoutes.GET("/locator", ac.GetActivitiesLocator)
		activitiesRoutes.GET("/:id", ac.GetActivityById)
		activitiesRoutes.GET("/:id/edit", ac.RenderEditActivityForm)
		activitiesRoutes.PUT("/:id/edit", ac.UpdateActivity)
		activitiesRoutes.DELETE("/:id", ac.DeleteActivity)
	}
}
