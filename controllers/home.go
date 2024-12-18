package controllers

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/laurawarren88/go_spa_backend.git/models"
	"gorm.io/gorm"
)

type HomeController struct {
	DB *gorm.DB
}

func NewHomeController(db *gorm.DB) *HomeController {
	return &HomeController{DB: db}
}

func (hc *HomeController) GetHome(ctx *gin.Context) {
	postcode := ctx.Query("postcode")
	radius := ctx.Query("radius")

	// Logging for debugging purposes
	fmt.Printf("Received query params - Postcode: %s, Radius: %s\n", postcode, radius)

	// Dummy data for now
	var activities []models.Activities
	ctx.JSON(http.StatusOK, activities)
}
