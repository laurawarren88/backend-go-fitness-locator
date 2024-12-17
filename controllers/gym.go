package controllers

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/laurawarren88/go_spa_backend.git/models"
	"gorm.io/gorm"
)

type GymController struct {
	DB *gorm.DB
}

func NewGymController(db *gorm.DB) *GymController {
	return &GymController{DB: db}
}

func (gc *GymController) GetAllGyms(ctx *gin.Context) {
	var gyms []models.Gym
	gc.DB.Find(&gyms)
	ctx.JSON(http.StatusOK, gyms)
}

func (gc *GymController) RenderCreateGymForm(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, gin.H{
		"title": "Create a New Gym",
	})
}

func (gc *GymController) CreateGym(ctx *gin.Context) {
	var gym models.Gym
	if err := ctx.ShouldBindJSON(&gym); err != nil {
		log.Println("Error binding JSON:", err)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	// Save to the database
	if err := gc.DB.Create(&gym).Error; err != nil {
		log.Println("Error saving to database:", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create gym"})
		return
	}

	ctx.JSON(http.StatusCreated, gym) // Respond with the created gym
}

func (gc *GymController) GetGymById(ctx *gin.Context) {
	id := ctx.Param("id")
	var gym models.Gym
	if err := gc.DB.First(&gym, id).Error; err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "Gym not found"})
		return
	}
	ctx.JSON(http.StatusOK, gym)
}

func (gc *GymController) UpdateGym(ctx *gin.Context) {
	id := ctx.Param("id")
	var gym models.Gym
	if err := gc.DB.First(&gym, id).Error; err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "Gym not found"})
		return
	}
	if err := ctx.ShouldBindJSON(&gym); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}
	gc.DB.Save(&gym)
	ctx.JSON(http.StatusOK, gym)
}

func (gc *GymController) DeleteGym(ctx *gin.Context) {
	id := ctx.Param("id")
	var gym models.Gym
	if err := gc.DB.First(&gym, id).Error; err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "Gym not found"})
		return
	}
	gc.DB.Delete(&gym)
	ctx.JSON(http.StatusOK, gin.H{"message": "Gym deleted"})
}
