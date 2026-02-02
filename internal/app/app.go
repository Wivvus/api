package app

import (
	"net/http"

	"github.com/Wivvus/api/internal/oauth"
	"github.com/gin-gonic/gin"
)

func ConfigureRouter(r *gin.Engine) {
	r.GET("/secure", oauth.RequireAuth(), secure)
	r.GET("/insecure", insecure)
}

func secure(ctx *gin.Context) {
	// Return user info
	ctx.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "User fetched successfully",
		"data":    ctx.Keys[oauth.USER_CONTEXT_KEY],
	})
}

func insecure(ctx *gin.Context) {
	// Return user info
	ctx.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Page loaded successfully",
		"data":    "insecure",
	})
}
