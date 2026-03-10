package middleware

import (
	"context"
	"errors"
	"net/http"
	"strings"

	"github.com/Wivvus/api/internal/models"
	"github.com/coreos/go-oidc/v3/oidc"
	"github.com/gin-gonic/gin"
)

var (
	verifier *oidc.IDTokenVerifier
)

func InitAuth(clientID string) error {
	provider, err := oidc.NewProvider(context.Background(), "https://accounts.google.com")
	if err != nil {
		return err
	}

	verifier = provider.Verifier(&oidc.Config{
		ClientID: clientID,
	})

	return nil
}

func AuthRequired() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header required"})
			c.Abort()
			return
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		if tokenString == authHeader {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid authorization format"})
			c.Abort()
			return
		}

		// Verify the ID token
		idToken, err := verifier.Verify(c.Request.Context(), tokenString)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			c.Abort()
			return
		}

		// Extract claims
		var claims struct {
			Email         string `json:"email"`
			EmailVerified bool   `json:"email_verified"`
			Name          string `json:"name"`
			Picture       string `json:"picture"`
		}

		if err := idToken.Claims(&claims); err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Failed to parse claims"})
			c.Abort()
			return
		}

		user := &models.User{
			OauthID:   idToken.Subject,
			Name:      claims.Name,
			Email:     claims.Email,
			AvatarURL: claims.Picture,
		}

		ur := &models.UserRepo{}
		err = ur.Create(user)
		if err != nil && !errors.Is(models.UserExists, err) {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "error storing user data"})
			c.Abort()
			return
		}

		// Store user info in context
		c.Set("user", user)
		c.Next()
	}
}
