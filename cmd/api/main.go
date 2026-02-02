package main

import (
	"log"
	"os"

	"github.com/Wivvus/api/internal/app"
	"github.com/Wivvus/api/internal/models"
	"github.com/Wivvus/api/internal/oauth"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"

	"github.com/markbates/goth/gothic"
)

func main() {
	r := gin.Default()

	// set up sessions and oauth
	cookieStoreKey := os.Getenv("COOKIE_STORE_KEY")
	cookieStoreSecret := os.Getenv("COOKIE_STORE_SECRET")
	store := cookie.NewStore([]byte(cookieStoreSecret))
	r.Use(sessions.Sessions(cookieStoreKey, store))
	gothic.Store = store

	// set up DB
	dbHost := os.Getenv("PG_HOST")
	dbPort := os.Getenv("PG_PORT")
	dbUser := os.Getenv("PG_USER")
	dbPass := os.Getenv("PG_PASSWORD")
	dbDatabase := os.Getenv("PG_DB")

	dsn := "host=" + dbHost + " user=" + dbUser + " password=" + dbPass + " dbname=" + dbDatabase + " port=" + dbPort + " sslmode=disable"
	models.ConnectDB(dsn)

	// configure routes
	oauth.ConfigureRouter(r)
	app.ConfigureRouter(r)

	if err := r.Run(); err != nil {
		log.Fatalf("failed to run server: %v", err)
	}
}
