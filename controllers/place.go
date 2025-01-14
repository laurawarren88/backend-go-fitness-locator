package controllers

import (
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"os"
	"strconv"
	"strings"
	"sync"

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

// type RateLimitData struct {
// 	LastRequestTime time.Time
// 	RequestCount    int
// }

// const (
// 	rateLimitWindow      = time.Minute * 1 // Time window for rate limiting
// 	rateLimitMaxRequests = 10              // Max requests per API key per window
// )

// type Place struct {
// 	Name     string  `json:"name"`
// 	Vicinity string  `json:"vicinity"`
// 	Lat      float64 `json:"lat"`
// 	Lng      float64 `json:"lng"`
// }

// func HandleError(ctx *gin.Context, status int, message string) {
// 	log.Printf("Error [%d]: %s\n", status, message)
// 	ctx.JSON(status, gin.H{"error": message})
// }

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
		Name         string `form:"name" json:"name"`
		Vicinity     string `form:"vicinity" json:"vicinity"`
		City         string `form:"city" json:"city"`
		Postcode     string `form:"postcode" json:"postcode"`
		Phone        string `form:"phone" json:"phone"`
		Email        string `form:"email" json:"email"`
		Website      string `form:"website" json:"website"`
		OpeningHours string `form:"opening_hours" json:"opening_hours"`
		Description  string `form:"description" json:"description"`
		TypeID       uint   `form:"typeID" json:"typeID"`
		Type         string `form:"type" json:"type"`
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
		typeID, _ := strconv.ParseUint(ctx.Request.FormValue("typeID"), 10, 32)
		placeFields.TypeID = uint(typeID)
		placeFields.Type = ctx.Request.FormValue("type")
	} else {
		// Unsupported Content-Type
		ctx.JSON(http.StatusUnsupportedMediaType, gin.H{"error": "Unsupported Content-Type"})
		return
	}

	// Debug: Log parsed fields
	log.Printf("Activity Text Fields: %+v\n", placeFields)

	// Handle file uploads if multipart/form-data
	var logoPath string
	if strings.HasPrefix(contentType, "multipart/form-data") {
		logoFile, logoHeader, err := ctx.Request.FormFile("logo")
		if err == nil && logoFile != nil {
			logoPath = "uploads/logos/" + logoHeader.Filename
			if err := saveFile(logoFile, logoPath); err != nil {
				ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save logo: " + err.Error()})
				return
			}
			log.Printf("Logo saved at: %s", logoPath)
		}
	}

	// Create activity object
	activities := models.Place{
		Name:         placeFields.Name,
		Vicinity:     placeFields.Vicinity,
		City:         placeFields.City,
		Postcode:     placeFields.Postcode,
		Phone:        placeFields.Phone,
		Email:        placeFields.Email,
		Website:      placeFields.Website,
		OpeningHours: placeFields.OpeningHours,
		Description:  placeFields.Description,
		TypeID:       placeFields.TypeID,
		Type:         placeFields.Type,
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
			"id":            activities.ID,
			"name":          activities.Name,
			"vicinity":      activities.Vicinity,
			"city":          activities.City,
			"postcode":      activities.Postcode,
			"phone":         activities.Phone,
			"email":         activities.Email,
			"website":       activities.Website,
			"opening_hours": activities.OpeningHours,
			"description":   activities.Description,
			"typeID":        activities.TypeID,
			"type":          activities.Type,
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

// func (pc *PlaceController) GetPlaceLocator(ctx *gin.Context) {
// 	ctx.JSON(http.StatusOK, gin.H{"message": "Locator Page"})
// }

// func (pc *PlaceController) GetActivityById(ctx *gin.Context) {
// 	id := ctx.Param("id")
// 	var activities models.Activities
// 	if err := pc.DB.First(&activities, id).Error; err != nil {
// 		ctx.JSON(http.StatusNotFound, gin.H{"error": "Activity not found"})
// 		return
// 	}
// 	ctx.JSON(http.StatusOK, activities)
// }

// func (pc *PlaceController) RenderEditActivityForm(ctx *gin.Context) {
// 	ctx.JSON(http.StatusOK, gin.H{
// 		"title": "Update Activity Form",
// 	})
// }

// func (pc *PlaceController) UpdateActivity(ctx *gin.Context) {
// 	id := ctx.Param("id")
// 	var activities models.Activities

// 	if err := pc.DB.First(&activities, id).Error; err != nil {
// 		ctx.JSON(http.StatusNotFound, gin.H{"error": "Activity not found"})
// 		return
// 	}

// 	if err := ctx.ShouldBind(&activities); err != nil {
// 		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
// 		return
// 	}

// 	if err := pc.DB.Save(&activities).Error; err != nil {
// 		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update Activity"})
// 		return
// 	}

// 	ctx.JSON(http.StatusOK, activities)
// }

// func (pc *PlaceController) RenderDeleteActivityForm(ctx *gin.Context) {
// 	ctx.JSON(http.StatusOK, gin.H{
// 		"title": "Delete Activity Form",
// 	})
// }

// func (pc *PlaceController) DeleteActivity(ctx *gin.Context) {
// 	id := ctx.Param("id")
// 	var activities models.Activities
// 	if err := pc.DB.First(&activities, id).Error; err != nil {
// 		ctx.JSON(http.StatusNotFound, gin.H{"error": "Activity not found"})
// 		return
// 	}
// 	pc.DB.Delete(&activities)
// 	ctx.JSON(http.StatusOK, gin.H{"message": "Activity deleted"})
// }
