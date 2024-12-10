package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/laurawarren88/go_spa_backend.git/controllers"
)

func RegisterUserRoutes(router *gin.Engine) {
	userRoutes := router.Group("/project/users")
	{
		userRoutes.GET("/register", controllers.GetSignupForm)
		userRoutes.POST("/register", controllers.SignupUser)
		userRoutes.GET("/login", controllers.GetLoginForm)
		userRoutes.POST("/login", controllers.LoginUser)
		userRoutes.GET("/forgot_password", controllers.ForgotPassword)
		userRoutes.POST("/forgot_password", controllers.ResetPassword)
		userRoutes.POST("/logout", controllers.LogoutUser)
	}
}
