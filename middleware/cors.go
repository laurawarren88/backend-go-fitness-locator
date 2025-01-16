package middleware

import (
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func CORSMiddleware() gin.HandlerFunc {
	return cors.New(cors.Config{
		// AllowOrigins: []string{"*"},
		AllowOrigins: []string{"http://localhost:5050", "http://127.0.0.1:5050"},
		AllowMethods: []string{
			"GET",
			"POST",
			"PUT",
			"PATCH",
			"DELETE",
			"OPTIONS",
		},
		AllowHeaders: []string{
			"Content-Type",
			"Content-Length",
			"Accept-Encoding",
			"Authorization",
			"Accept",
			"Origin",
			"Cache-Control",
			"X-Requested-With",
		},
		ExposeHeaders: []string{
			"Content-Length",
			"Content-Type",
			"Content-Disposition",
		},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
		AllowWildcard:    true,
	})
}
