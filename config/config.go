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
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file")
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
	blogController := controllers.NewBlogController(db)
	userController := controllers.NewUserController(db)

	routes.RegisterHomeRoutes(router, homeController)
	routes.RegisterBlogRoutes(router, blogController)
	routes.RegisterUserRoutes(router, userController)
}
