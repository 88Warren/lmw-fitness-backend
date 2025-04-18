package config

import (
	"log"
	"net/http"
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
	router.Use(middleware.CORSMiddleware())
	router.GET("/debug/images", func(c *gin.Context) {
		log.Println("Hit /debug/images route")
		files, err := os.ReadDir("./images")
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		var filenames []string
		for _, file := range files {
			filenames = append(filenames, file.Name())
		}

		c.JSON(http.StatusOK, gin.H{"files": filenames})
	})
	router.Use(middleware.DBMiddleware())
	return router
}

func SetupHandlers(router *gin.Engine, db *gorm.DB) {
	homeController := controllers.NewHomeController(db)
	healthController := controllers.NewHealthController(db)
	routes.RegisterHomeRoutes(router, homeController, healthController)
}
