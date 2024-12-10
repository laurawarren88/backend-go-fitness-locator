package main

import (
	"fmt"
	"log"

	"github.com/laurawarren88/go_spa_backend.git/config"
)

func init() {
	config.LoadEnv()
	config.SetGinMode()
}

func main() {
	router := config.SetupServer()

	config.SetupHandlers(router)

	fmt.Printf("Starting the server on port %s\n", config.GetEnv("PORT", "8000"))
	if err := router.Run(":" + config.GetEnv("PORT", "8000")); err != nil {
		log.Fatal("Failed to start the server:", err)
	}
}
