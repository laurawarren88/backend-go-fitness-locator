package controllers

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/laurawarren88/go_spa_backend.git/models"
	"gorm.io/gorm"
)

type PlaceController struct {
	RateLimiter sync.Map
	DB          *gorm.DB
}

func NewPlaceController(db *gorm.DB) *PlaceController {
	return &PlaceController{DB: db}
}

func (pc *PlaceController) RenderCreateActivityForm(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, gin.H{
		"title": "Create a New Activity",
	})
}

func (pc *PlaceController) CreateActivity(ctx *gin.Context) {
	// Handle preflight OPTIONS requests
	if ctx.Request.Method == "OPTIONS" {
		ctx.Status(http.StatusOK)
		return
	}

	// Define the structure for form and JSON binding
	type PlaceTextFields struct {
		Name            string `form:"name" json:"name"`
		Vicinity        string `form:"vicinity" json:"vicinity"`
		City            string `form:"city" json:"city"`
		Postcode        string `form:"postcode" json:"postcode"`
		Phone           string `form:"phone" json:"phone"`
		Email           string `form:"email" json:"email"`
		Website         string `form:"website" json:"website"`
		OpeningHours    string `form:"opening_hours" json:"opening_hours"`
		Description     string `form:"description" json:"description"`
		Type            string `form:"type" json:"type"`
		Latitude        string `form:"latitude" json:"latitude"`
		Longitude       string `form:"longitude" json:"longitude"`
		Logo            string `json:"logo" form:"logo" gorm:"size:255"`
		FacilitiesImage string `json:"facilities_image" form:"facilities_image" gorm:"size:255"`
	}

	var placeFields PlaceTextFields

	// Detect Content-Type and process accordingly
	contentType := ctx.GetHeader("Content-Type")
	if strings.HasPrefix(contentType, "application/json") {
		// Handle JSON payload
		if err := ctx.ShouldBindJSON(&placeFields); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON input: " + err.Error()})
			return
		}
	} else if strings.HasPrefix(contentType, "multipart/form-data") {
		// Handle form data with file uploads
		if err := ctx.Request.ParseMultipartForm(10 << 20); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "Failed to parse form: " + err.Error()})
			return
		}
		// Extract form fields
		placeFields.Name = ctx.Request.FormValue("name")
		placeFields.Vicinity = ctx.Request.FormValue("vicinity")
		placeFields.City = ctx.Request.FormValue("city")
		placeFields.Postcode = ctx.Request.FormValue("postcode")
		placeFields.Phone = ctx.Request.FormValue("phone")
		placeFields.Email = ctx.Request.FormValue("email")
		placeFields.Website = ctx.Request.FormValue("website")
		placeFields.OpeningHours = ctx.Request.FormValue("opening_hours")
		placeFields.Description = ctx.Request.FormValue("description")
		placeFields.Type = ctx.Request.FormValue("type")
		placeFields.Latitude = ctx.Request.FormValue("latitude")
		placeFields.Longitude = ctx.Request.FormValue("longitude")
		placeFields.Logo = ctx.Request.FormValue("logo")
		placeFields.FacilitiesImage = ctx.Request.FormValue("facilities_image")

		if logoFile, err := ctx.FormFile("logo"); err == nil {
			sanitizedFilename := fmt.Sprintf("%d_%s", time.Now().Unix(), filepath.Base(logoFile.Filename))
			logoFilePath := "./uploads/logos/" + sanitizedFilename

			if err := os.MkdirAll("./uploads/logos", os.ModePerm); err != nil {
				ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create directory"})
				return
			}

			if err := ctx.SaveUploadedFile(logoFile, logoFilePath); err != nil {
				ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save logo file"})
				return
			}

			placeFields.Logo = logoFilePath
		} else {
			log.Printf("No logo uploaded or error occurred: %v", err)
		}

		// Handle facilities image upload
		if facilitiesImageFile, err := ctx.FormFile("facilities_image"); err == nil {
			sanitizedFilename := fmt.Sprintf("%d_%s", time.Now().Unix(), filepath.Base(facilitiesImageFile.Filename))
			facilitiesImageFilePath := "./uploads/facilities/" + sanitizedFilename

			if err := os.MkdirAll("./uploads/facilities", os.ModePerm); err != nil {
				ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create directory"})
				return
			}

			if err := ctx.SaveUploadedFile(facilitiesImageFile, facilitiesImageFilePath); err != nil {
				ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save logo file"})
				return
			}

			placeFields.FacilitiesImage = facilitiesImageFilePath
		} else {
			log.Printf("No logo uploaded or error occurred: %v", err)
		}
	}

	// Debug: Log parsed fields
	log.Printf("Activity Text Fields: %+v\n", placeFields)

	// Create activity object
	activities := models.Place{
		Name:            placeFields.Name,
		Vicinity:        placeFields.Vicinity,
		City:            placeFields.City,
		Postcode:        placeFields.Postcode,
		Phone:           placeFields.Phone,
		Email:           placeFields.Email,
		Website:         placeFields.Website,
		OpeningHours:    placeFields.OpeningHours,
		Description:     placeFields.Description,
		Type:            placeFields.Type,
		Latitude:        placeFields.Latitude,
		Longitude:       placeFields.Longitude,
		Logo:            placeFields.Logo,
		FacilitiesImage: placeFields.FacilitiesImage,
	}

	// Save activity to the database
	if err := pc.DB.Create(&activities).Error; err != nil {
		log.Println("Error saving to database:", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create activity"})
		return
	}

	// Respond with success
	ctx.JSON(http.StatusCreated, gin.H{
		"message": "Activity created successfully",
		"activities": gin.H{
			"id":               activities.ID,
			"name":             activities.Name,
			"vicinity":         activities.Vicinity,
			"city":             activities.City,
			"postcode":         activities.Postcode,
			"phone":            activities.Phone,
			"email":            activities.Email,
			"website":          activities.Website,
			"opening_hours":    activities.OpeningHours,
			"description":      activities.Description,
			"type":             activities.Type,
			"latitude":         activities.Latitude,
			"longitude":        activities.Longitude,
			"logo":             activities.Logo,
			"facilities_image": activities.FacilitiesImage,
		},
	})
}

func (pc *PlaceController) GetPlaceLocator(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, gin.H{"message": "Locator Page"})
}

func (pc *PlaceController) GetActivityById(ctx *gin.Context) {
	id := ctx.Param("id")

	// Ensure the ID is a valid integer
	var places models.Place
	if err := pc.DB.First(&places, "id = ?", id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			ctx.JSON(http.StatusNotFound, gin.H{"error": "Activity not found"})
		} else {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve activity"})
		}
		return
	}

	// Return the activity as JSON
	ctx.JSON(http.StatusOK, gin.H{
		"id":               places.ID,
		"name":             places.Name,
		"vicinity":         places.Vicinity,
		"city":             places.City,
		"postcode":         places.Postcode,
		"phone":            places.Phone,
		"email":            places.Email,
		"website":          places.Website,
		"opening_hours":    places.OpeningHours,
		"description":      places.Description,
		"type":             places.Type,
		"latitude":         places.Latitude,
		"longitude":        places.Longitude,
		"logo":             places.Logo,
		"facilities_image": places.FacilitiesImage,
	})
}

func (pc *PlaceController) RenderEditActivityForm(ctx *gin.Context) {
	id := ctx.Param("id")
	var existingPlace models.Place

	// Find the activity by ID
	if err := pc.DB.First(&existingPlace, "id = ?", id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			ctx.JSON(http.StatusNotFound, gin.H{"error": "Activity not found"})
		} else {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve activity"})
		}
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"title":    "Update Activity Form",
		"activity": existingPlace,
	})
}

func (pc *PlaceController) UpdateActivity(ctx *gin.Context) {
	if ctx.Request.Method == "OPTIONS" {
		ctx.Status(http.StatusOK)
		return
	}

	id := ctx.Param("id")
	var existingPlace models.Place

	// Find the activity by ID
	if err := pc.DB.First(&existingPlace, "id = ?", id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			ctx.JSON(http.StatusNotFound, gin.H{"error": "Activity not found"})
		} else {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve activity"})
		}
		return
	}

	// Bind JSON input to a new Place object to avoid overwriting the entire existingPlace
	var input models.Place
	if err := ctx.ShouldBindJSON(&input); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input: " + err.Error()})
		return
	}

	// Update only the fields provided in the request
	if input.Name != "" {
		existingPlace.Name = input.Name
	}
	if input.Vicinity != "" {
		existingPlace.Vicinity = input.Vicinity
	}
	if input.City != "" {
		existingPlace.City = input.City
	}
	if input.Postcode != "" {
		existingPlace.Postcode = input.Postcode
	}
	if input.Phone != "" {
		existingPlace.Phone = input.Phone
	}
	if input.Email != "" {
		existingPlace.Email = input.Email
	}
	if input.Website != "" {
		existingPlace.Website = input.Website
	}
	if input.OpeningHours != "" {
		existingPlace.OpeningHours = input.OpeningHours
	}
	if input.Description != "" {
		existingPlace.Description = input.Description
	}
	if input.Type != "" {
		existingPlace.Type = input.Type
	}
	if input.Latitude != "" {
		existingPlace.Latitude = input.Latitude
	}
	if input.Longitude != "" {
		existingPlace.Longitude = input.Longitude
	}
	if input.Logo != "" {
		existingPlace.Logo = input.Logo
	}
	if input.FacilitiesImage != "" {
		existingPlace.FacilitiesImage = input.FacilitiesImage
	}

	// Save the updated activity
	if err := pc.DB.Save(&existingPlace).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update activity"})
		return
	}

	// Respond with the updated activity
	ctx.JSON(http.StatusOK, gin.H{
		"message":  "Activity updated successfully",
		"activity": existingPlace,
	})
}

func (pc *PlaceController) RenderDeleteActivityForm(ctx *gin.Context) {
	id := ctx.Param("id")
	var existingPlace models.Place

	if err := pc.DB.First(&existingPlace, "id = ?", id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			ctx.JSON(http.StatusNotFound, gin.H{"error": "Activity not found"})
		} else {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve activity"})
		}
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"title":    "Delete Activity Form",
		"activity": existingPlace,
	})
}

func (pc *PlaceController) DeleteActivity(ctx *gin.Context) {
	if ctx.Request.Method == "OPTIONS" {
		ctx.Status(http.StatusOK)
		return
	}

	id := ctx.Param("id")
	var place models.Place

	// Check if the activity exists
	if err := pc.DB.First(&place, "id = ?", id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			ctx.JSON(http.StatusNotFound, gin.H{"error": "Activity not found"})
		} else {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve activity"})
		}
		return
	}

	// Delete the activity
	if err := pc.DB.Delete(&place).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete activity"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Activity deleted successfully"})
}
