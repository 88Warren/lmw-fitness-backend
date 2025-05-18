package controllers

import (
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
)

type ImagesController struct {
}

func NewImagesController() *ImagesController {
	return &ImagesController{}
}

func (ic *ImagesController) GetImage(c *gin.Context) {
	filename := c.Param("filename")
	if filename == "" {
		c.Status(http.StatusBadRequest)
		return
	}

	// Ensure we only serve files from the images directory
	filepath := "./images/" + filename
	
	// Check if file exists
	if _, err := os.Stat(filepath); os.IsNotExist(err) {
		c.Status(http.StatusNotFound)
		return
	}

	// Serve the file
	c.File(filepath)
}
