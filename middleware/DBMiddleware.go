package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/laurawarren88/go_spa_backend.git/database"
)

func DBMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Set("db", database.GetDB())
		c.Next()
	}
}
