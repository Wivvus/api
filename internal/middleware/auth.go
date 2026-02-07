package middleware

import (
	"context"
	"net/http"
	"strings"

	"github.com/coreos/go-oidc/v3/oidc"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
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
		logger := logrus.New()
		logger.Print("processing auth required")
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			logger.Printf("no auth header")
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

		// Store user info in context
		c.Set("user_email", claims.Email)
		c.Set("user_name", claims.Name)
		c.Set("user_id", idToken.Subject)
		c.Set("user_picture", claims.Picture)

		c.Next()
	}
}
