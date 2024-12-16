package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/laurawarren88/go_spa_backend.git/controllers"
	"github.com/laurawarren88/go_spa_backend.git/middleware"
)

func RegisterUserRoutes(router *gin.Engine, uc *controllers.UserController) {
	userRoutes := router.Group("/users")
	{
		userRoutes.GET("/register", uc.GetSignupForm)
		userRoutes.POST("/register", uc.SignupUser)
		userRoutes.GET("/login", uc.GetLoginForm)
		userRoutes.POST("/login", uc.LoginUser)
		userRoutes.GET("/forgot_password", uc.ForgotPassword)
		userRoutes.POST("/forgot_password", uc.ResetPassword)
	}

	protected := router.Group("/users")
	protected.Use(middleware.AuthMiddleware())
	{
		protected.GET("/profile/:id", uc.GetProfile)
		protected.POST("/logout", uc.LogoutUser)
	}
}
