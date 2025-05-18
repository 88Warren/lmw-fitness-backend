package controllers

import (
	"net/http"

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

	// Serve the file from the images directory
	c.FileFromFS("./images/"+filename, http.Dir("."))
}
