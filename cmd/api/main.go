package main

import (
	"log"
	"os"

	"github.com/Wivvus/api/internal/app"
	"github.com/Wivvus/api/internal/oauth"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"

	"github.com/markbates/goth/gothic"
)

func main() {
	cookieStoreKey := os.Getenv("COOKIE_STORE_KEY")
	cookieStoreSecret := os.Getenv("COOKIE_STORE_SECRET")

	r := gin.Default()

	store := cookie.NewStore([]byte(cookieStoreSecret))
	r.Use(sessions.Sessions(cookieStoreKey, store))

	gothic.Store = store

	oauth.ConfigureRouter(r)

	app.ConfigureRouter(r)

	if err := r.Run(); err != nil {
		log.Fatalf("failed to run server: %v", err)
	}
}
