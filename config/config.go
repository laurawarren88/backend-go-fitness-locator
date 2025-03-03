package config

import (
	"os"

	"github.com/gin-gonic/gin"
	"github.com/laurawarren88/go_spa_backend.git/controllers"
	"github.com/laurawarren88/go_spa_backend.git/middleware"
	"github.com/laurawarren88/go_spa_backend.git/routes"
	"gorm.io/gorm"
)

func LoadEnv() {
	// if err := godotenv.Load(); err != nil {
	// 	log.Fatal("Error loading .env file")
	// }
}

func SetGinMode() {
	gin.SetMode(gin.ReleaseMode)
}

func GetEnv(key string, fallback string) string {
	value, exists := os.LookupEnv(key)
	if !exists {
		return fallback
	}
	return value
}

func SetupServer() *gin.Engine {
	router := gin.Default()
	router.Static("/images", "./uploads")
	router.Use(middleware.DBMiddleware())
	return router
}

func SetupHandlers(router *gin.Engine, db *gorm.DB) {
	homeController := controllers.NewHomeController(db)
	placeController := controllers.NewPlaceController(db)
	userController := controllers.NewUserController(db)

	routes.RegisterHomeRoutes(router, homeController)
	routes.RegisterPlaceRoutes(router, placeController)
	routes.RegisterUserRoutes(router, userController)
}
