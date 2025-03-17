package middleware

import (
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func CORSMiddleware() gin.HandlerFunc {
	// env := os.Getenv("GO_ENV")

	// var allowedOrigins []string
	// if env == "development" {
	// 	allowedOrigins = []string{
	// 		os.Getenv("BASE_URL"),
	// 	}
	// } else {
	// 	allowedOrigins = []string{
	// 		os.Getenv("BASE_URL"),
	// 	}
	// }

	return cors.New(cors.Config{
		AllowOrigins: []string{
			"http://localhost:8081",
		},
		AllowMethods: []string{
			"GET",
			"POST",
			"PUT",
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
	})
}
