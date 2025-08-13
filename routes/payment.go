package routes

import (
	"github.com/88warren/lmw-fitness-backend/controllers"
	"github.com/gin-gonic/gin"
)

func RegisterPaymentRoutes(router *gin.Engine, pc *controllers.PaymentController) {
	api := router.Group("/api")
	{
		api.POST("/create-checkout-session", pc.CreateCheckoutSession)
		api.POST("/stripe-webhook", pc.StripeWebhook)
		api.GET("/test-webhook", pc.TestWebhook) // Test endpoint
		api.POST("/get-workout-link", pc.GetWorkoutLink)
	}
}
