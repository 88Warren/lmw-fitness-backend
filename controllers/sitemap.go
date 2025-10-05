package controllers

import (
	"fmt"
	"net/http"
	"time"

	"github.com/88warren/lmw-fitness-backend/models"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type SitemapController struct {
	DB *gorm.DB
}

func NewSitemapController(db *gorm.DB) *SitemapController {
	return &SitemapController{DB: db}
}

type SitemapURL struct {
	Loc        string `xml:"loc"`
	LastMod    string `xml:"lastmod"`
	ChangeFreq string `xml:"changefreq"`
	Priority   string `xml:"priority"`
}

type SitemapURLSet struct {
	XMLName string       `xml:"urlset"`
	Xmlns   string       `xml:"xmlns,attr"`
	URLs    []SitemapURL `xml:"url"`
}

func (sc *SitemapController) GenerateSitemap(c *gin.Context) {
	baseURL := "https://www.lmwfitness.co.uk"
	currentDate := time.Now().Format("2006-01-02")

	// Static pages
	urls := []SitemapURL{
		{
			Loc:        baseURL + "/",
			LastMod:    currentDate,
			ChangeFreq: "weekly",
			Priority:   "1.0",
		},
		{
			Loc:        baseURL + "/programs",
			LastMod:    currentDate,
			ChangeFreq: "weekly",
			Priority:   "0.9",
		},
		{
			Loc:        baseURL + "/blog",
			LastMod:    currentDate,
			ChangeFreq: "daily",
			Priority:   "0.8",
		},
		{
			Loc:        baseURL + "/calculator",
			LastMod:    currentDate,
			ChangeFreq: "monthly",
			Priority:   "0.7",
		},
		{
			Loc:        baseURL + "/login",
			LastMod:    currentDate,
			ChangeFreq: "monthly",
			Priority:   "0.3",
		},
		{
			Loc:        baseURL + "/register",
			LastMod:    currentDate,
			ChangeFreq: "monthly",
			Priority:   "0.3",
		},
	}

	// Get all published blog posts
	var blogPosts []models.Blog
	if err := sc.DB.Find(&blogPosts).Error; err == nil {
		for _, post := range blogPosts {
			blogURL := SitemapURL{
				Loc:        fmt.Sprintf("%s/blog/%d", baseURL, post.ID),
				LastMod:    post.UpdatedAt.Format("2006-01-02"),
				ChangeFreq: "monthly",
				Priority:   "0.7",
			}
			urls = append(urls, blogURL)
		}
	}

	sitemap := SitemapURLSet{
		Xmlns: "http://www.sitemaps.org/schemas/sitemap/0.9",
		URLs:  urls,
	}

	c.Header("Content-Type", "application/xml")
	c.XML(http.StatusOK, sitemap)
}
