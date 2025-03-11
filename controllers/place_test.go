package controllers_test

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/laurawarren88/go_spa_backend.git/controllers"
	"github.com/laurawarren88/go_spa_backend.git/models"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// Setup test DB
func setupTestDB() (*gorm.DB, error) {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	// Auto migrate the test database
	err = db.AutoMigrate(&models.Place{}, &models.User{})
	if err != nil {
		return nil, err
	}

	return db, nil
}

// Test helper to create test data
func createTestData(db *gorm.DB) error {
	// Create test user
	user := models.User{
		Username: "testuser",
		Email:    "test@example.com",
		Password: "testpassword",
		IsAdmin:  false,
	}
	if err := db.Create(&user).Error; err != nil {
		return err
	}

	// Create test place
	place := models.Place{
		Name:        "Test Place",
		Description: "Test Description",
		Phone:       "1234567890",
		UserID:      user.ID,
	}
	return db.Create(&place).Error
}
func TestGetActivityById(t *testing.T) {
	db, err := setupTestDB()
	assert.NoError(t, err)

	err = createTestData(db)
	assert.NoError(t, err)

	var place models.Place
	err = db.First(&place).Error
	assert.NoError(t, err)

	gin.SetMode(gin.TestMode)
	r := gin.Default()
	controller := controllers.NewPlaceController(db)

	r.GET("/activity/:id", controller.GetActivityById)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/activity/1", nil)

	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]any
	err = json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, "Test Place", response["name"])
}

func TestCheckActivityOwnership(t *testing.T) {
	db, err := setupTestDB()
	assert.NoError(t, err)

	err = createTestData(db)
	assert.NoError(t, err)

	var place models.Place
	err = db.First(&place).Error
	assert.NoError(t, err)

	gin.SetMode(gin.TestMode)
	r := gin.Default()
	controller := controllers.NewPlaceController(db)

	r.Use(func(c *gin.Context) {
		c.Set("userID", uint(1))
		c.Next()
	})

	r.GET("/activity/:id/check-ownership", controller.CheckActivityOwnership)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/activity/"+fmt.Sprint(place.ID)+"/check-ownership", nil)

	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]any
	err = json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.True(t, response["isOwner"].(bool))
}
