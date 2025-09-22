package main

import (
	"log"
	"os"

	"github.com/88warren/lmw-fitness-backend/config"
	"github.com/88warren/lmw-fitness-backend/database"
	"go.uber.org/zap"
)

func init() {
	config.LoadEnv()
	config.SetGinMode()
}

func main() {
	database.InitLogger()
	defer database.SyncLogger()

	// Debug: Print environment variables
	log.Printf("GO_ENV: %s", os.Getenv("GO_ENV"))
	log.Printf("ALLOWED_ORIGIN: %s", os.Getenv("ALLOWED_ORIGIN"))
	log.Printf("PORT: %s", os.Getenv("PORT"))

	database.ConnectToDB()
	db := database.GetDB()

	router := config.SetupServer()

	config.SetupHandlers(router, db)

	port := config.GetEnv("PORT", "8082")

	if err := router.Run("0.0.0.0:" + port); err != nil {
		zap.L().Fatal("Failed to start the server", zap.Error(err))
	}
}
