package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/laurawarren88/go_spa_backend.git/controllers"
	"github.com/laurawarren88/go_spa_backend.git/middleware"
)

func RegisterPlaceRoutes(router *gin.Engine, pc *controllers.PlaceController) {
	router.GET("/api/activities/:id/check-ownership", middleware.AuthMiddleware(), pc.CheckActivityOwnership)

	placeRoutes := router.Group("/api/activities")
	{
		placeRoutes.GET("/locator", pc.GetPlaceLocator)
		placeRoutes.GET("/:id", pc.GetActivityById)
	}

	protected := router.Group("/api/activities")
	protected.Use(middleware.AuthMiddleware())
	{
		protected.GET("/new", pc.RenderCreateActivityForm)
		protected.POST("/new", pc.CreateActivity)
	}
	userRoutes := router.Group("/api/activities")
	userRoutes.Use(middleware.AuthMiddleware(), middleware.ActivityOwner())
	{
		userRoutes.GET("/:id/edit", pc.RenderEditActivityForm)
		userRoutes.PUT("/:id/edit", pc.UpdateActivity)
		userRoutes.GET("/:id/delete", pc.RenderDeleteActivityForm)
		userRoutes.DELETE("/:id/delete", pc.DeleteActivity)
	}
}
