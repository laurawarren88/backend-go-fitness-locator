package controllers

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/laurawarren88/go_spa_backend.git/models"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type UserController struct {
	DB *gorm.DB
}

func NewUserController(db *gorm.DB) *UserController {
	return &UserController{DB: db}
}

func (uc *UserController) GetSignupForm(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, gin.H{"message": "Registration form"})
}

func (uc *UserController) SignupUser(ctx *gin.Context) {
	var user models.User

	if err := ctx.ShouldBindJSON(&user); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input: " + err.Error()})
		return
	}

	var existingUser models.User
	if err := uc.DB.Where("email = ? OR username = ?", user.Email, user.Username).First(&existingUser).Error; err == nil {
		ctx.JSON(http.StatusConflict, gin.H{"error": "Email or username already registered"})
		return
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to hash password"})
		return
	}
	user.Password = string(hash)
	user.CreatedAt = time.Now()
	user.UpdatedAt = time.Now()

	if err := uc.DB.Create(&user).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user"})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{
		"message": "User created successfully",
		"user": gin.H{
			"id":       user.ID,
			"username": user.Username,
			"email":    user.Email,
			"isAdmin":  user.IsAdmin,
		},
	})
}

func (uc *UserController) GetLoginForm(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, gin.H{"message": "Login form"})
}

func (uc *UserController) LoginUser(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, gin.H{"message": "Login successful"})
}

func (uc *UserController) ForgotPassword(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, gin.H{"message": "Reset Password form"})
}

func (uc *UserController) ResetPassword(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, gin.H{"message": "Password reset instructions sent"})
}

func (uc *UserController) LogoutUser(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, gin.H{"message": "Logged out"})
}
