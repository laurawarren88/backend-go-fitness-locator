package controllers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type PlaceController struct {
	RateLimiter sync.Map
	DB          *gorm.DB
}

func NewPlaceController(db *gorm.DB) *PlaceController {
	return &PlaceController{DB: db}
}

type RateLimitData struct {
	LastRequestTime time.Time
	RequestCount    int
}

const (
	rateLimitWindow      = time.Minute * 1 // Time window for rate limiting
	rateLimitMaxRequests = 10              // Max requests per API key per window
)

type Place struct {
	Name     string  `json:"name"`
	Vicinity string  `json:"vicinity"`
	Lat      float64 `json:"lat"`
	Lng      float64 `json:"lng"`
}

func HandleError(ctx *gin.Context, status int, message string) {
	log.Printf("Error [%d]: %s\n", status, message)
	ctx.JSON(status, gin.H{"error": message})
}

func (pc *PlaceController) RenderCreateActivityForm(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, gin.H{
		"title": "Create a New Activity",
	})
}

func (pc *PlaceController) CreateActivity(ctx *gin.Context) {
	if ctx.Request.Method == "OPTIONS" {
		return
	}

	// Parse multipart form with a 10 MB limit
	if err := ctx.Request.ParseMultipartForm(10 << 20); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Failed to parse form: " + err.Error()})
		return
	}

	type ActivityTextFields struct {
		Name        string `form:"name"`
		Address     string `form:"address"`
		City        string `form:"city"`
		Postcode    string `form:"postcode"`
		Category    string `form:"category"`
		Description string `form:"description"`
		TypeID      uint   `form:"typeID"`
		Type        string `form:"type"`
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
		Name:        activityFields.Name,
		Address:     activityFields.Address,
		City:        activityFields.City,
		Postcode:    activityFields.Postcode,
		Description: activityFields.Description,
		TypeID:      activityFields.TypeID,
		Type:        activityFields.Type,
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
			"id":          activities.ID,
			"name":        activities.Name,
			"address":     activities.Address,
			"city":        activities.City,
			"postcode":    activities.Postcode,
			"description": activities.Description,
			"typeID":      activities.TypeID,
			"type":        activities.Type,
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

func (pc *PlaceController) GetPlaceLocator(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, gin.H{"message": "Locator Page"})
}

func (pc *PlaceController) RenderMapPage(ctx *gin.Context) {
	lat := ctx.Query("lat")
	lng := ctx.Query("lng")

	Only allow "lat" and "lng" query parameters
	if lat == "" || lng == "" {
		HandleError(ctx, http.StatusBadRequest, "Missing required query parameters: lat, lng")
		return
	}

	Enforce "gym" as the activity type
	activityType := ctx.DefaultQuery("type", "gym")
	apiKey := os.Getenv("GOOGLE_PLACES_API_KEY")

	if apiKey == "" {
		HandleError(ctx, http.StatusInternalServerError, "Google Places API key is not set")
		return
	}

	Build the Google Places API URL with hardcoded "type=gym"
	placesAPIURL := fmt.Sprintf(
		"https://maps.googleapis.com/maps/api/place/nearbysearch/json?location=%s,%s&radius=1500&type=%s&key=%s",
		lat, lng, activityType, apiKey,
	)

	log.Printf("Places API URL: %s", placesAPIURL) // Debugging log

	Fetch data from Google Places API
	resp, err := http.Get(placesAPIURL)
	if err != nil {
		HandleError(ctx, http.StatusInternalServerError, "Failed to fetch data from Google Places API")
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		HandleError(ctx, http.StatusInternalServerError, "Google Places API returned an error")
		return
	}

	Parse the API response
	var placesResponse struct {
		Results []struct {
			Name     string `json:"name"`
			Vicinity string `json:"vicinity"`
			Geometry struct {
				Location struct {
					Lat float64 `json:"lat"`
					Lng float64 `json:"lng"`
				} `json:"location"`
			} `json:"geometry"`
		} `json:"results"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&placesResponse); err != nil {
		HandleError(ctx, http.StatusInternalServerError, "Failed to decode Google Places API response")
		return
	}

	Collect the results
	var places []Place
	for _, result := range placesResponse.Results {
		places = append(places, Place{
			Name:     result.Name,
			Vicinity: result.Vicinity,
			Lat:      result.Geometry.Location.Lat,
			Lng:      result.Geometry.Location.Lng,
		})
	}
	log.Printf("Activities to be sent: %+v\n", places) // Debugging log

	Respond with the activities
	ctx.JSON(http.StatusOK, gin.H{
		"message": "Gyms found successfully",
		"places":  places,
	})
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
	ctx.JSON(http.StatusOK, gin.H{
		"title": "Update Activity Form",
	})
}

func (pc *PlaceController) UpdateActivity(ctx *gin.Context) {
	id := ctx.Param("id")
	var activities models.Activities

	if err := pc.DB.First(&activities, id).Error; err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "Activity not found"})
		return
	}

	if err := ctx.ShouldBind(&activities); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	if err := pc.DB.Save(&activities).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update Activity"})
		return
	}

	ctx.JSON(http.StatusOK, activities)
}

func (pc *PlaceController) RenderDeleteActivityForm(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, gin.H{
		"title": "Delete Activity Form",
	})
}

func (pc *PlaceController) DeleteActivity(ctx *gin.Context) {
	id := ctx.Param("id")
	var activities models.Activities
	if err := pc.DB.First(&activities, id).Error; err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "Activity not found"})
		return
	}
	pc.DB.Delete(&activities)
	ctx.JSON(http.StatusOK, gin.H{"message": "Activity deleted"})
}
