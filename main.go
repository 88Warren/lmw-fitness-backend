package main

import (
	"fmt"
	"log"

	"github.com/laurawarren88/LMW_Fitness/config"
	"github.com/laurawarren88/LMW_Fitness/database"
)

func init() {
	config.LoadEnv()
	config.SetGinMode()
}

func main() {
	database.ConnectToDB()
	db := database.GetDB()

	router := config.SetupServer()

	config.SetupHandlers(router, db)

	port := config.GetEnv("PORT", "8080")
	fmt.Printf("Starting the server on port %s\n", port)

	// Start the server on the specified port
	if err := router.Run("0.0.0.0:" + port); err != nil {
		log.Fatal("Failed to start the server:", err)
	}
}
