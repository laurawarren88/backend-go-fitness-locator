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
	type GymTextFields struct {
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

	var gymFields GymTextFields
	if err := ctx.ShouldBind(&gymFields); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input: " + err.Error()})
		log.Println("Payload binding error:", err)
		return
	}

	// Debug: Log the bound text fields
	log.Printf("Gym Text Fields After Binding: %+v\n", gymFields)

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

	// Create a complete Gym object
	gym := models.Gym{
		Business_name:     gymFields.Business_name,
		Address:           gymFields.Address,
		City:              gymFields.City,
		Postcode:          gymFields.Postcode,
		Phone:             gymFields.Phone,
		Email:             gymFields.Email,
		Website:           gymFields.Website,
		Opening_hours:     gymFields.Opening_hours,
		Activities:        gymFields.Activities,
		Facilities:        gymFields.Facilities,
		Logo:              logoPath,
		Facilities_images: facilitiesPath,
	}

	gym.CreatedAt = time.Now()
	gym.UpdatedAt = time.Now()

	// Save the gym object to the database
	if err := gc.DB.Create(&gym).Error; err != nil {
		log.Println("Error saving to database:", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create gym"})
		return
	}

	// Respond with the created gym details
	ctx.JSON(http.StatusCreated, gin.H{
		"message": "Gym created successfully",
		"gym": gin.H{
			"id":                gym.ID,
			"business_name":     gym.Business_name,
			"address":           gym.Address,
			"city":              gym.City,
			"postcode":          gym.Postcode,
			"phone":             gym.Phone,
			"email":             gym.Email,
			"website":           gym.Website,
			"opening_hours":     gym.Opening_hours,
			"activities":        gym.Activities,
			"facilities":        gym.Facilities,
			"logo":              gym.Logo,
			"facilities_images": gym.Facilities_images,
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

func (gc *GymController) GetGymById(ctx *gin.Context) {
	id := ctx.Param("id")
	var gym models.Gym
	if err := gc.DB.First(&gym, id).Error; err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "Gym not found"})
		return
	}
	ctx.JSON(http.StatusOK, gym)
}

func (gc *GymController) RenderEditGymForm(ctx *gin.Context) {
	id := ctx.Param("id")
	var gym models.Gym
	if err := gc.DB.First(&gym, id).Error; err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "Gym not found"})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"title": "Edit Gym",
		"gym":   gym,
	})
}

func (gc *GymController) UpdateGym(ctx *gin.Context) {
	id := ctx.Param("id")
	var gym models.Gym

	if err := gc.DB.First(&gym, id).Error; err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "Gym not found"})
		return
	}

	if err := ctx.ShouldBind(&gym); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	// // If files are uploaded, handle them separately
	// logo, err := ctx.FormFile("logo")
	// if err == nil {
	// 	// Process and save the logo file (for example, save to a folder)
	// 	if err := ctx.SaveUploadedFile(logo, "./uploads/"+logo.Filename); err != nil {
	// 		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to upload logo"})
	// 		return
	// 	}
	// 	gym.Logo = "./uploads/" + logo.Filename // Save the file path in the gym model
	// }

	// facilitiesImage, err := ctx.FormFile("facilities_image")
	// if err == nil {
	// 	// Process and save the facilities image
	// 	if err := ctx.SaveUploadedFile(facilitiesImage, "./uploads/"+facilitiesImage.Filename); err != nil {
	// 		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to upload facilities image"})
	// 		return
	// 	}
	// 	gym.Facilities_image = "./uploads/" + facilitiesImage.Filename // Save the file path in the gym model
	// }

	if err := gc.DB.Save(&gym).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update gym"})
		return
	}

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
