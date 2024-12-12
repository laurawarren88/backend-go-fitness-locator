package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/laurawarren88/go_spa_backend.git/controllers"
)

func RegisterUserRoutes(router *gin.Engine, uc *controllers.UserController) {
	userRoutes := router.Group("/project/users")
	{
		userRoutes.GET("/register", uc.GetSignupForm)
		userRoutes.POST("/register", uc.SignupUser)
		userRoutes.GET("/login", uc.GetLoginForm)
		userRoutes.POST("/login", uc.LoginUser)
		userRoutes.GET("/forgot_password", uc.ForgotPassword)
		userRoutes.POST("/forgot_password", uc.ResetPassword)
		userRoutes.POST("/logout", uc.LogoutUser)
	}
}
