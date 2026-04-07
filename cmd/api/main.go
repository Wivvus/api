package main

import (
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/gin-gonic/gin"

	"github.com/Wivvus/api/internal/app"
	"github.com/Wivvus/api/internal/middleware"
	"github.com/Wivvus/api/internal/models"
	"github.com/Wivvus/api/internal/tokens"

	"github.com/gin-contrib/cors"
)

func main() {
	clientID := os.Getenv("GOOGLE_OAUTH_CLIENT_ID")
	if clientID == "" {
		log.Fatal("GOOGLE_OAUTH_CLIENT_ID environment variable is required")
	}

	// Initialize authentication with go-oidc
	if err := middleware.InitAuth(clientID); err != nil {
		log.Fatalf("Failed to initialize auth: %v", err)
	}

	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		log.Fatal("JWT_SECRET environment variable is required")
	}
	tokens.Init(jwtSecret)

	r := gin.Default()

	// CORS
	r.Use(cors.New(cors.Config{
		AllowOrigins:     strings.Split(os.Getenv("ALLOWED_ORIGINS"), ","),
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization", "X-Auth-Provider"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
	}))

	// Public routes
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})

	// Protected routes
	r.GET("/user/data", middleware.AuthRequired(), getUserData)
	r.GET("/user/profile", middleware.AuthRequired(), getUserProfile)

	// set up DB
	dbHost := os.Getenv("PG_HOST")
	dbPort := os.Getenv("PG_PORT")
	dbUser := os.Getenv("PG_USER")
	dbPass := os.Getenv("PG_PASSWORD")
	dbDatabase := os.Getenv("PG_DB")

	sslMode := os.Getenv("PG_SSLMODE")
	if sslMode == "" {
		sslMode = "disable"
	}
	dsn := "host=" + dbHost + " user=" + dbUser + " password=" + dbPass + " dbname=" + dbDatabase + " port=" + dbPort + " sslmode=" + sslMode
	models.ConnectDB(dsn)

	app.ConfigureRouter(r)

	if err := r.Run(); err != nil {
		log.Fatalf("failed to run server: %v", err)
	}
}

func getUserData(c *gin.Context) {
	user, exists := c.Get("user")
	if !exists {
		c.JSON(http.StatusGone, gin.H{
			"message": "the user is no longer available",
		})
		c.Abort()
		return
	}
	userModel, ok := user.(*models.User)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "unexpected error processing user data",
		})
		c.Abort()
		return
	}
	c.JSON(200, gin.H{
		"message": "This is protected data",
		"user":    userModel.ToAPI(),
		"data": []string{
			"User-specific item 1",
			"User-specific item 2",
		},
	})
}

func getUserProfile(c *gin.Context) {
	user, exists := c.Get("user")
	if !exists {
		c.JSON(http.StatusGone, gin.H{
			"message": "the user is no longer available",
		})
		c.Abort()
		return
	}
	userModel, ok := user.(*models.User)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "unexpected error processing user data",
		})
		c.Abort()
		return
	}

	c.JSON(200, gin.H{
		"email":   userModel.Email,
		"name":    userModel.Name,
		"picture": userModel.AvatarURL,
		"role":    "user",
	})
}
