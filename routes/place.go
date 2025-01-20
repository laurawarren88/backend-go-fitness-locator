package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/laurawarren88/go_spa_backend.git/controllers"
	"github.com/laurawarren88/go_spa_backend.git/middleware"
)

func RegisterPlaceRoutes(router *gin.Engine, pc *controllers.PlaceController) {
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
	adminRoutes := router.Group("/api/activities")
	adminRoutes.Use(middleware.AuthMiddleware(), middleware.RequireAdmin())
	{
		adminRoutes.GET("/:id/edit", pc.RenderEditActivityForm)
		adminRoutes.PUT("/:id/edit", pc.UpdateActivity)
		adminRoutes.GET("/:id/delete", pc.RenderDeleteActivityForm)
		adminRoutes.DELETE("/:id/delete", pc.DeleteActivity)
	}
}
