package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/laurawarren88/LMW_Fitness/database"
)

func DBMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Set("db", database.GetDB())
		c.Next()
	}
}
