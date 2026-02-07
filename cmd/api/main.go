package main

import (
	"log"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"

	"github.com/Wivvus/api/internal/middleware"

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

	r := gin.Default()

	// CORS
	r.Use(cors.New(cors.Config{
		AllowOrigins:     strings.Split(os.Getenv("ALLOWED_ORIGINS"), ","),
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
	}))

	// Public routes
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})

	// Protected routes
	r.GET("/api/user/data", middleware.AuthRequired(), getUserData)
	r.GET("/api/user/profile", middleware.AuthRequired(), getUserProfile)

	// // set up sessions and oauth
	// cookieStoreKey := os.Getenv("COOKIE_STORE_KEY")
	// cookieStoreSecret := os.Getenv("COOKIE_STORE_SECRET")
	// store := cookie.NewStore([]byte(cookieStoreSecret))
	// r.Use(sessions.Sessions(cookieStoreKey, store))
	// gothic.Store = store

	// // set up DB
	// dbHost := os.Getenv("PG_HOST")
	// dbPort := os.Getenv("PG_PORT")
	// dbUser := os.Getenv("PG_USER")
	// dbPass := os.Getenv("PG_PASSWORD")
	// dbDatabase := os.Getenv("PG_DB")

	// dsn := "host=" + dbHost + " user=" + dbUser + " password=" + dbPass + " dbname=" + dbDatabase + " port=" + dbPort + " sslmode=disable"
	// models.ConnectDB(dsn)

	// app.ConfigureRouter(r)

	if err := r.Run(); err != nil {
		log.Fatalf("failed to run server: %v", err)
	}
}

func getUserData(c *gin.Context) {
	logger := logrus.New()
	logger.Print("processing get user data")
	userEmail := c.GetString("user_email")
	userName := c.GetString("user_name")
	userID := c.GetString("user_id")

	c.JSON(200, gin.H{
		"message": "This is protected data",
		"user": gin.H{
			"email": userEmail,
			"name":  userName,
			"id":    userID,
		},
		"data": []string{
			"User-specific item 1",
			"User-specific item 2",
		},
	})
}

func getUserProfile(c *gin.Context) {
	userEmail := c.GetString("user_email")
	userName := c.GetString("user_name")
	userPicture := c.GetString("user_picture")

	c.JSON(200, gin.H{
		"email":   userEmail,
		"name":    userName,
		"picture": userPicture,
		"role":    "user",
	})
}
