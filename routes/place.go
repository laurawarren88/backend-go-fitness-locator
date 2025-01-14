package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/laurawarren88/go_spa_backend.git/controllers"
)

func RegisterPlaceRoutes(router *gin.Engine, pc *controllers.PlaceController) {
	placeRoutes := router.Group("/activities")
	{
		placeRoutes.GET("/new", pc.RenderCreateActivityForm)
		placeRoutes.POST("/new", pc.CreateActivity)
		placeRoutes.GET("/locator", pc.GetPlaceLocator)
		placeRoutes.GET("/:id", pc.GetActivityById)
		placeRoutes.GET("/:id/edit", pc.RenderEditActivityForm)
		placeRoutes.PUT("/:id/edit", pc.UpdateActivity)
		placeRoutes.GET("/:id/delete", pc.RenderDeleteActivityForm)
		placeRoutes.DELETE("/:id/delete", pc.DeleteActivity)
	}
}
