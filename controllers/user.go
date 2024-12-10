package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func GetSignupForm(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, gin.H{"message": "Registration form"})
}

func SignupUser(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, gin.H{"message": "User created successfully"})
}

func GetLoginForm(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, gin.H{"message": "Login form"})
}

func LoginUser(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, gin.H{"message": "Login successful"})
}

func ForgotPassword(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, gin.H{"message": "Reset Password form"})
}

func ResetPassword(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, gin.H{"message": "Password reset instructions sent"})
}

func LogoutUser(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, gin.H{"message": "Logged out"})
}
