package controllers

import (
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/laurawarren88/go_spa_backend.git/models"
	"gorm.io/gorm"
)

type ActivitiesController struct {
	DB *gorm.DB
}

func NewActivitiesController(db *gorm.DB) *ActivitiesController {
	return &ActivitiesController{DB: db}
}

func (ac *ActivitiesController) GetAllActivities(ctx *gin.Context) {
	var activities []models.Activities
	ac.DB.Find(&activities)
	ctx.JSON(http.StatusOK, activities)
}

func (ac *ActivitiesController) RenderCreateActivityForm(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, gin.H{
		"title": "Create a New Activity",
	})
}

func (ac *ActivitiesController) CreateActivity(ctx *gin.Context) {
	if ctx.Request.Method == "OPTIONS" {
		return
	}

	// Parse multipart form with a 10 MB limit
	if err := ctx.Request.ParseMultipartForm(10 << 20); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Failed to parse form: " + err.Error()})
		return
	}

	// Debug: Log all form data received
	log.Printf("Form Data: %+v\n", ctx.Request.Form)

	// Define a temporary struct to bind only text fields
	type ActivityTextFields struct {
		Business_name string `form:"business_name"`
		Address       string `form:"address"`
		City          string `form:"city"`
		Postcode      string `form:"postcode"`
		Phone         string `form:"phone"`
		Email         string `form:"email"`
		Website       string `form:"website"`
		Opening_hours string `form:"opening_hours"`
		Activities    string `form:"activities"`
		Facilities    string `form:"facilities"`
	}

	var activityFields ActivityTextFields
	if err := ctx.ShouldBind(&activityFields); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input: " + err.Error()})
		log.Println("Payload binding error:", err)
		return
	}

	// Debug: Log the bound text fields
	log.Printf("Activity Text Fields After Binding: %+v\n", activityFields)

	// Handle file uploads for logo
	var logoPath string
	logoFile, logoHeader, err := ctx.Request.FormFile("logo")
	if err == nil && logoFile != nil {
		logoPath = "uploads/logos/" + logoHeader.Filename
		if err := saveFile(logoFile, logoPath); err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save logo: " + err.Error()})
			return
		}
		log.Printf("Logo successfully saved at: %s", logoPath)
	} else if err != nil && err != http.ErrMissingFile {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Error processing logo file: " + err.Error()})
		return
	}

	// Handle file uploads for facilities images
	var facilitiesPath string
	facilitiesFile, facilitiesHeader, err := ctx.Request.FormFile("facilities_images")
	if err == nil && facilitiesFile != nil {
		facilitiesPath = "uploads/facilities/" + facilitiesHeader.Filename
		if err := saveFile(facilitiesFile, facilitiesPath); err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save facilities image: " + err.Error()})
			return
		}
		log.Printf("Facilities images successfully saved at: %s", facilitiesPath)
	} else if err != nil && err != http.ErrMissingFile {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Error processing facilities images file: " + err.Error()})
		return
	}

	// Create a complete Activities object
	activities := models.Activities{
		Business_name:     activityFields.Business_name,
		Address:           activityFields.Address,
		City:              activityFields.City,
		Postcode:          activityFields.Postcode,
		Phone:             activityFields.Phone,
		Email:             activityFields.Email,
		Website:           activityFields.Website,
		Opening_hours:     activityFields.Opening_hours,
		Activities:        activityFields.Activities,
		Facilities:        activityFields.Facilities,
		Logo:              logoPath,
		Facilities_images: facilitiesPath,
	}

	activities.CreatedAt = time.Now()
	activities.UpdatedAt = time.Now()

	// Save the activities object to the database
	if err := ac.DB.Create(&activities).Error; err != nil {
		log.Println("Error saving to database:", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create activity"})
		return
	}

	// Respond with the created activity details
	ctx.JSON(http.StatusCreated, gin.H{
		"message": "Activity created successfully",
		"activities": gin.H{
			"id":                activities.ID,
			"business_name":     activities.Business_name,
			"address":           activities.Address,
			"city":              activities.City,
			"postcode":          activities.Postcode,
			"phone":             activities.Phone,
			"email":             activities.Email,
			"website":           activities.Website,
			"opening_hours":     activities.Opening_hours,
			"activities":        activities.Activities,
			"facilities":        activities.Facilities,
			"logo":              activities.Logo,
			"facilities_images": activities.Facilities_images,
		},
	})
}

func saveFile(file multipart.File, path string) error {
	out, err := os.Create(path)
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, file)
	return err
}

func (ac *ActivitiesController) GetActivitiesLocator(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, gin.H{"message": "Locator Page"})
}

func (ac *ActivitiesController) GetActivityById(ctx *gin.Context) {
	id := ctx.Param("id")
	var activities models.Activities
	if err := ac.DB.First(&activities, id).Error; err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "Activity not found"})
		return
	}
	ctx.JSON(http.StatusOK, activities)
}

func (ac *ActivitiesController) RenderEditActivityForm(ctx *gin.Context) {
	id := ctx.Param("id")
	var activities models.Activities
	if err := ac.DB.First(&activities, id).Error; err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "Activity not found"})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"title":      "Edit Activity",
		"activities": activities,
	})
}

func (ac *ActivitiesController) UpdateActivity(ctx *gin.Context) {
	id := ctx.Param("id")
	var activities models.Activities

	if err := ac.DB.First(&activities, id).Error; err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "Activity not found"})
		return
	}

	if err := ctx.ShouldBind(&activities); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	if err := ac.DB.Save(&activities).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update Activity"})
		return
	}

	ctx.JSON(http.StatusOK, activities)
}

func (ac *ActivitiesController) DeleteActivity(ctx *gin.Context) {
	id := ctx.Param("id")
	var activities models.Activities
	if err := ac.DB.First(&activities, id).Error; err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "Activity not found"})
		return
	}
	ac.DB.Delete(&activities)
	ctx.JSON(http.StatusOK, gin.H{"message": "Activity deleted"})
}
