package config

import (
	"log"
	"net/http"
	"os"

	"github.com/88warren/lmw-fitness-backend/controllers"
	"github.com/88warren/lmw-fitness-backend/database"
	"github.com/88warren/lmw-fitness-backend/middleware"
	"github.com/88warren/lmw-fitness-backend/models"
	"github.com/88warren/lmw-fitness-backend/routes"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"golang.org/x/crypto/bcrypt"
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
	return router
}

func SetupAdminUser(db *gorm.DB) error {
	adminEmail := os.Getenv("ADMIN_EMAIL")
	adminPassword := os.Getenv("ADMIN_PASSWORD")

	if adminEmail == "" || adminPassword == "" {
		// log.Println("ADMIN_EMAIL or ADMIN_PASSWORD environment variables not set. Skipping admin user setup.")
		return nil
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(adminPassword), bcrypt.DefaultCost)
	if err != nil {
		log.Printf("Failed to hash admin password: %v", err)
		return err
	}

	adminUser := models.User{
		Email:        adminEmail,
		PasswordHash: string(hashedPassword),
		Role:         "admin",
	}

	var existingUser models.User
	result := db.Where("email = ?", adminUser.Email).First(&existingUser)

	switch {
	case result.Error == gorm.ErrRecordNotFound:
		if err := db.Create(&adminUser).Error; err != nil {
			log.Printf("Failed to create admin user: %v", err)
			return err
		}
		// log.Printf("Admin user '%s' created successfully.", adminUser.Email)
	case result.Error == nil:
		existingUser.PasswordHash = adminUser.PasswordHash
		existingUser.Role = "admin"
		if err := db.Save(&existingUser).Error; err != nil {
			log.Printf("Failed to update admin user '%s': %v", existingUser.Email, err)
			return err
		}
		// log.Printf("Admin user '%s' updated successfully.", existingUser.Email)
	default:
		// log.Printf("Database error checking for admin user: %v", result.Error)
		return result.Error
	}

	return nil
}

func SetupHandlers(router *gin.Engine, db *gorm.DB) {
	// Run database migration
	database.MigrateDB()
	database.SeedDB(db)
	// log.Println("Database models auto-migrated successfully.")

	if err := SetupAdminUser(db); err != nil {
		log.Fatalf("Failed to setup admin user: %v", err)
	}

	homeController := controllers.NewHomeController(db)
	healthController := controllers.NewHealthController(db)
	paymentController := controllers.NewPaymentController()
	blogController := controllers.NewBlogController(db)
	userController := controllers.NewUserController(db)
	newsletterController := controllers.NewNewsletterController(db)
	workoutController := controllers.NewWorkoutController(db)

	routes.RegisterHomeRoutes(router, homeController)
	routes.RegisterHealthRoutes(router, healthController)
	routes.RegisterPaymentRoutes(router, paymentController)
	routes.RegisterBlogRoutes(router, blogController)
	routes.RegisterUserRoutes(router, userController)
	routes.RegisterNewsletterRoutes(router, newsletterController)
	routes.RegisterWorkoutRoutes(router, workoutController)
}
