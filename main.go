package main

import (
	"log"

	"github.com/88warren/lmw-fitness-backend/config"
	"github.com/88warren/lmw-fitness-backend/database"
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

	port := config.GetEnv("PORT", "8082")
	// fmt.Printf("Starting the server on port %s\n", port)

	if err := router.Run("0.0.0.0:" + port); err != nil {
		log.Fatal("Failed to start the server:", err)
	}
}
