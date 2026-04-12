package middleware

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/Wivvus/api/internal/models"
	"github.com/Wivvus/api/internal/storage"
	"github.com/Wivvus/api/internal/tokens"
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

		provider := c.GetHeader("X-Auth-Provider")

		var user *models.User
		var err error

		if provider == "local" {
			user, err = verifyLocalToken(tokenString)
		} else {
			user, err = verifyGoogleToken(c.Request.Context(), tokenString)
		}

		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
			c.Abort()
			return
		}

		spacesHost := fmt.Sprintf("%s.%s.digitaloceanspaces.com", os.Getenv("DO_SPACES_BUCKET"), os.Getenv("DO_SPACES_REGION"))
		if user.AvatarURL != "" && !strings.Contains(user.AvatarURL, spacesHost) {
			ur := &models.UserRepo{}
			go copyAvatarToSpaces(user, ur)
		}

		c.Set("user", user)
		c.Next()
	}
}

func VerifyGoogleToken(ctx context.Context, tokenString string) (*models.User, error) {
	return verifyGoogleToken(ctx, tokenString)
}

func verifyGoogleToken(ctx context.Context, tokenString string) (*models.User, error) {
	idToken, err := verifier.Verify(ctx, tokenString)
	if err != nil {
		return nil, errors.New("Invalid token")
	}

	var claims struct {
		Email   string `json:"email"`
		Name    string `json:"name"`
		Picture string `json:"picture"`
	}
	if err := idToken.Claims(&claims); err != nil {
		return nil, errors.New("Failed to parse claims")
	}

	user := &models.User{
		OauthID:   idToken.Subject,
		Name:      claims.Name,
		Email:     claims.Email,
		AvatarURL: claims.Picture,
		Provider:  "google",
	}

	ur := &models.UserRepo{}
	err = ur.Create(user)
	if err != nil && !errors.Is(models.UserExists, err) {
		return nil, errors.New("error storing user data")
	}

	return user, nil
}

func copyAvatarToSpaces(user *models.User, ur *models.UserRepo) {
	filename := fmt.Sprintf("%d-avatar.jpg", user.ID)
	url, err := storage.CopyFromURL(user.AvatarURL, filename)
	if err != nil {
		return
	}
	ur.UpdateAvatar(user.ID, url)
	user.AvatarURL = url
}

func verifyLocalToken(tokenString string) (*models.User, error) {
	claims, err := tokens.Verify(tokenString)
	if err != nil {
		return nil, errors.New("Invalid token")
	}

	ur := &models.UserRepo{}
	user := ur.FindByID(claims.Subject)
	if user.ID == 0 {
		return nil, errors.New("User not found")
	}
	return user, nil
}
