package middleware

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func CORS() gin.HandlerFunc {
	config := cors.DefaultConfig()
	config.AllowAllOrigins = false
	config.AllowCredentials = true
	config.AllowOriginFunc = func(origin string) bool {
		return true
	}
	config.AllowMethods = []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"}
	config.AllowHeaders = []string{
		"X-Requested-With",
		"X-Timestamp",
		"Origin",
		"Accept",
		"Referer",
		"Cache-Control",
		"User-Agent",
		"Content-Length",
		"Content-Type",
		"Authorization",
		"Accept",
		"Connection",
		"Access-Control-Allow-Origin",
		"Access-Control-Allow-Headers",
		"Access-Control-Allow-Methods",
		"Access-Control-Allow-Credentials",
		"Access-Control-Request-Method",
		"Access-Control-Request-Headers",
	}
	return cors.New(config)
}
