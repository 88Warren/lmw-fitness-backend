package config

import (
	"log"
	"net/http"
	"os"

	"github.com/88warren/lmw-fitness-backend/controllers"
	"github.com/88warren/lmw-fitness-backend/database"
	"github.com/88warren/lmw-fitness-backend/middleware"
	"github.com/88warren/lmw-fitness-backend/routes"
	"github.com/88warren/lmw-fitness-backend/workers"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
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
		if os.Getenv("KUBERNETES_SERVICE_HOST") != "" {
			// log.Printf("Running in Kubernetes, using environment variables from ConfigMap and Secrets")
		}
	} else {
		// log.Printf("Loaded environment variables from %s", envFile)
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

	router.Use(middleware.StructuredLoggingMiddleware())
	router.Use(middleware.MetricsMiddleware())
	router.Use(middleware.MetricsCollectionMiddleware())

	router.Use(middleware.CORSMiddleware())
	router.Static("/images", "./images")
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

	middleware.LogMetricsPeriodically()

	return router
}

func SetupHandlers(router *gin.Engine, db *gorm.DB) {
	database.SeedDB(db)
	// log.Println("Database models auto-migrated successfully.")

	homeController := controllers.NewHomeController(db)
	healthController := controllers.NewHealthController(db)
	paymentController := controllers.NewPaymentController(db)
	blogController := controllers.NewBlogController(db)
	userController := controllers.NewUserController(db)
	newsletterController := controllers.NewNewsletterController(db)
	workoutController := controllers.NewWorkoutController(db)
	monitoringController := controllers.NewMonitoringController(db)

	routes.RegisterHomeRoutes(router, homeController)
	routes.RegisterHealthRoutes(router, healthController)
	routes.RegisterPaymentRoutes(router, paymentController)
	routes.RegisterBlogRoutes(router, blogController)
	routes.RegisterUserRoutes(router, userController)
	routes.RegisterNewsletterRoutes(router, newsletterController)
	routes.RegisterWorkoutRoutes(router, workoutController)
	routes.RegisterMonitoringRoutes(router, monitoringController)

	go func() {
		workers.StartPaymentWorker(db, paymentController)
	}()
}
