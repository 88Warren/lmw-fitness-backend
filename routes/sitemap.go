package routes

import (
	"github.com/88warren/lmw-fitness-backend/controllers"
	"github.com/gin-gonic/gin"
)

func RegisterSitemapRoutes(router *gin.Engine, sitemapController *controllers.SitemapController) {
	router.GET("/sitemap.xml", sitemapController.GenerateSitemap)
}
