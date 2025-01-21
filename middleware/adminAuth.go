package middleware

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

func RequireAdmin() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		isAdmin, exists := ctx.Get("isAdmin")
		if !exists {
			fmt.Println("isAdmin context key missing")
			ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Authentication required"})
			ctx.Abort()
			return
		}

		isAdminBool, ok := isAdmin.(bool)
		if !ok || !isAdminBool {
			fmt.Println("Admin access denied")
			ctx.JSON(http.StatusForbidden, gin.H{"error": "Admin access required"})
			ctx.Abort()
			return
		}

		fmt.Println("Admin access granted")
		ctx.Next()
	}
}
