package main

import (
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"

	"auth-server/config"
	"auth-server/handlers"
	"auth-server/middleware"
	"auth-server/models"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file")
	}

	config.ConnectDatabase()

	r := gin.Default()

	// Public routes
	r.POST("/auth/register", handlers.Register)
	r.POST("/auth/login", handlers.Login)
	r.POST("/auth/refresh", handlers.RefreshToken)

	// Protected routes group
	protected := r.Group("/api")
	protected.Use(middleware.AuthMiddleware())
	{
		// Add protected routes here
		protected.GET("/profile", func(c *gin.Context) {
			userId, _ := c.Get("user_id")
			var user models.User
			if err := config.DB.First(&user, userId).Error; err != nil {
				c.JSON(404, gin.H{"error": "User not found"})
				return
			}
			c.JSON(200, gin.H{"user": user})
		})
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	if err := r.Run(":" + port); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}
