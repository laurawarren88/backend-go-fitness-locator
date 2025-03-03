package main

import (
	"fmt"
	"log"

	"github.com/laurawarren88/go_spa_backend.git/config"
	"github.com/laurawarren88/go_spa_backend.git/database"
)

func init() {
	// config.LoadEnv()
	config.SetGinMode()
}

func main() {
	database.ConnectToDB()
	db := database.GetDB()

	if err := database.SetupAdminUser(db); err != nil {
		log.Printf("Error setting up admin user: %v", err)
	}

	router := config.SetupServer()

	config.SetupHandlers(router, db)

	fmt.Printf("Starting the server on port %s\n", config.GetEnv("PORT", "8081"))
	if err := router.Run(":" + config.GetEnv("PORT", "8081")); err != nil {
		log.Fatal("Failed to start the server:", err)
	}
}
