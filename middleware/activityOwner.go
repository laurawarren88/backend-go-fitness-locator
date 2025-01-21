package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/laurawarren88/go_spa_backend.git/models"
	"gorm.io/gorm"
)

func ActivityOwner() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		userID, exists := ctx.Get("userID")
		if !exists {
			ctx.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
			ctx.Abort()
			return
		}

		activityID := ctx.Param("id")
		var place models.Place

		DB := ctx.MustGet("db").(*gorm.DB)
		if err := DB.Where("id = ?", activityID).First(&place).Error; err != nil {
			if err == gorm.ErrRecordNotFound {
				ctx.JSON(http.StatusNotFound, gin.H{"error": "Activity not found"})
			} else {
				ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Error retrieving activity"})
			}
			ctx.Abort()
			return
		}

		if place.UserID != userID.(uint) {
			var user models.User
			if err := DB.First(&user, userID).Error; err != nil {
				ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Error checking user permissions"})
				ctx.Abort()
				return
			}

			if !user.IsAdmin {
				ctx.JSON(http.StatusForbidden, gin.H{"error": "You are not authorized to perform this action"})
				ctx.Abort()
				return
			}
		}

		ctx.Next()
	}
}
