package routes

import (
	"github.com/88warren/lmw-fitness-backend/controllers"
	"github.com/88warren/lmw-fitness-backend/middleware"
	"github.com/gin-gonic/gin"
)

func RegisterAMRAPRoutes(router *gin.Engine, ac *controllers.AMRAPController) {
	authenticated := router.Group("/api/amrap")
	authenticated.Use(middleware.AuthMiddleware())
	{
		authenticated.POST("/score", ac.SaveAMRAPScore)
		authenticated.GET("/score/:blockId", ac.GetAMRAPScore)
		authenticated.GET("/scores", ac.GetAllAMRAPScores)
	}
}
