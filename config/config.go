package config

import (
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/laurawarren88/LMW_Fitness/controllers"
	"github.com/laurawarren88/LMW_Fitness/middleware"
	"github.com/laurawarren88/LMW_Fitness/routes"
	"gorm.io/gorm"
)

func LoadEnv() {
	env := os.Getenv("GO_ENV")
	var envFile string

	if env == "production" {
		envFile = ".env.production"
	} else {
		envFile = ".env.development"
	}

	err := godotenv.Load(envFile)
	if err != nil {
		log.Printf("Warning: No %s file found, relying on system environment variables", envFile)
	} else {
		log.Printf("Loaded environment variables from %s", envFile)
	}
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
	router.Static("/images", "./images")
	router.Use(middleware.DBMiddleware())
	return router
}

func SetupHandlers(router *gin.Engine, db *gorm.DB) {
	homeController := controllers.NewHomeController(db)

	routes.RegisterHomeRoutes(router, homeController)
}
