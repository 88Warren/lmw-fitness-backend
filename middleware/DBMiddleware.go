package middleware

import (
	"github.com/88warren/lmw-fitness-backend/database"
	"github.com/gin-gonic/gin"
)

func DBMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Set("db", database.GetDB())
		c.Next()
	}
}
