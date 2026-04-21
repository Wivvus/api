package auth

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/Wivvus/api/internal/email"
	"github.com/Wivvus/api/internal/metrics"
	"github.com/Wivvus/api/internal/middleware"
	"github.com/Wivvus/api/internal/models"
	"github.com/Wivvus/api/internal/storage"
	"github.com/Wivvus/api/internal/tokens"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

func ConfigureRouter(r *gin.Engine) {
	r.POST("/auth/register", register)
	r.POST("/auth/set-password", setPassword)
	r.POST("/auth/forgot-password", forgotPassword)
	r.POST("/auth/change-password", middleware.AuthRequired(), changePassword)
	r.POST("/auth/login", login)
	r.POST("/auth/google", googleAuth)
	r.POST("/auth/google/code", googleCodeAuth)
	r.POST("/user/avatar", middleware.AuthRequired(), uploadAvatar)
	r.GET("/user/events", middleware.AuthRequired(), getUserEvents)
	r.DELETE("/user", middleware.AuthRequired(), deleteAccount)
}

func register(c *gin.Context) {
	var body struct {
		Email string `json:"email" binding:"required"`
		Name  string `json:"name" binding:"required"`
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "email and name are required"})
		return
	}

	vr := &models.VerificationRepo{}
	token, err := vr.Create(body.Email, body.Name)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create verification token"})
		return
	}

	appURL := os.Getenv("APP_URL")
	link := fmt.Sprintf("%s/set-password?token=%s", appURL, token.Token)

	if err := email.SendVerification(body.Email, body.Name, link); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to send verification email"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Check your email to continue"})
}

func setPassword(c *gin.Context) {
	var body struct {
		Token    string `json:"token" binding:"required"`
		Password string `json:"password" binding:"required"`
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "token and password are required"})
		return
	}

	if len(body.Password) < 8 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "password must be at least 8 characters"})
		return
	}

	vr := &models.VerificationRepo{}
	vt, err := vr.FindValid(body.Token)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid or expired token"})
		return
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(body.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to hash password"})
		return
	}

	ur := &models.UserRepo{}
	user, err := ur.UpsertLocalPassword(vt.Email, string(hash), vt.Name)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to save account"})
		return
	}

	vr.MarkUsed(vt.ID)
	metrics.UserRegistered(user.ID, user.Email)

	jwt, err := tokens.Sign(user.ID, user.Email, user.Name)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to sign token"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"token": jwt,
		"user":  user.ToAPI(),
	})
}

func uploadAvatar(c *gin.Context) {
	user := c.MustGet("user").(*models.User)

	fileHeader, err := c.FormFile("avatar")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "avatar file is required"})
		return
	}

	if fileHeader.Size > 5*1024*1024 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "file must be under 5MB"})
		return
	}

	ct := fileHeader.Header.Get("Content-Type")
	if !strings.HasPrefix(ct, "image/") {
		c.JSON(http.StatusBadRequest, gin.H{"error": "file must be an image"})
		return
	}

	file, err := fileHeader.Open()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to read file"})
		return
	}
	defer file.Close()

	ext := filepath.Ext(fileHeader.Filename)
	filename := fmt.Sprintf("%d-%d%s", user.ID, time.Now().UnixMilli(), ext)

	url, err := storage.UploadAvatar(file, filename, ct)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to upload avatar"})
		return
	}

	ur := &models.UserRepo{}
	if err := ur.UpdateAvatar(user.ID, url); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to save avatar"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"avatar_url": url})
}

func googleAuth(c *gin.Context) {
	var body struct {
		IDToken string `json:"id_token" binding:"required"`
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "id_token is required"})
		return
	}

	user, err := middleware.VerifyGoogleToken(c.Request.Context(), body.IDToken)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid Google token"})
		return
	}

	jwt, err := tokens.Sign(user.ID, user.Email, user.Name)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to sign token"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"token": jwt,
		"user":  user.ToAPI(),
	})
}

func getUserEvents(c *gin.Context) {
	user := c.MustGet("user").(*models.User)
	er := models.EventRepo{}
	c.JSON(http.StatusOK, er.AllByCreator(user.ID).ToAPI())
}

func deleteAccount(c *gin.Context) {
	user := c.MustGet("user").(*models.User)
	ur := &models.UserRepo{}
	if err := ur.Delete(user.ID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to delete account"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true})
}

func changePassword(c *gin.Context) {
	user := c.MustGet("user").(*models.User)

	var body struct {
		CurrentPassword string `json:"current_password" binding:"required"`
		NewPassword     string `json:"new_password" binding:"required"`
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "current_password and new_password are required"})
		return
	}

	if len(body.NewPassword) < 8 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "password must be at least 8 characters"})
		return
	}

	ur := &models.UserRepo{}
	existing, err := ur.FindByEmailWithPassword(user.Email)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "no password set on this account"})
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(existing.PasswordHash), []byte(body.CurrentPassword)); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "current password is incorrect"})
		return
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(body.NewPassword), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to hash password"})
		return
	}

	existing.PasswordHash = string(hash)
	if err := ur.UpdatePassword(existing.ID, string(hash)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update password"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "password updated"})
}

func forgotPassword(c *gin.Context) {
	var body struct {
		Email string `json:"email" binding:"required"`
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "email is required"})
		return
	}

	// Always return 200 to avoid revealing whether an account exists
	ur := &models.UserRepo{}
	user, err := ur.FindByEmailWithPassword(body.Email)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"message": "If an account exists, a reset email has been sent"})
		return
	}

	vr := &models.VerificationRepo{}
	token, err := vr.Create(user.Email, user.Name)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create reset token"})
		return
	}

	appURL := os.Getenv("APP_URL")
	link := fmt.Sprintf("%s/set-password?token=%s&reset=true", appURL, token.Token)

	if err := email.SendPasswordReset(user.Email, link); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to send reset email"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "If an account exists, a reset email has been sent"})
}

func login(c *gin.Context) {
	var body struct {
		Email    string `json:"email" binding:"required"`
		Password string `json:"password" binding:"required"`
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "email and password are required"})
		return
	}

	ur := &models.UserRepo{}
	user, err := ur.FindByEmailWithPassword(body.Email)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid email or password"})
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(body.Password)); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid email or password"})
		return
	}

	jwt, err := tokens.Sign(user.ID, user.Email, user.Name)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to sign token"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"token": jwt,
		"user":  user.ToAPI(),
	})
}

// googleCodeAuth exchanges an OAuth2 authorization code for a Wivvus JWT.
// Used by the redirect-based Google sign-in flow (works in all browsers including Telegram).
func googleCodeAuth(c *gin.Context) {
	var body struct {
		Code        string `json:"code" binding:"required"`
		RedirectURI string `json:"redirect_uri" binding:"required"`
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "code and redirect_uri are required"})
		return
	}

	idToken, err := exchangeGoogleCode(body.Code, body.RedirectURI)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "failed to exchange Google code"})
		return
	}

	user, err := middleware.VerifyGoogleToken(c.Request.Context(), idToken)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid Google token"})
		return
	}

	jwt, err := tokens.Sign(user.ID, user.Email, user.Name)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to sign token"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"token": jwt,
		"user":  user.ToAPI(),
	})
}

func exchangeGoogleCode(code, redirectURI string) (string, error) {
	data := url.Values{}
	data.Set("code", code)
	data.Set("client_id", os.Getenv("GOOGLE_OAUTH_CLIENT_ID"))
	data.Set("client_secret", os.Getenv("GOOGLE_OAUTH_CLIENT_SECRET"))
	data.Set("redirect_uri", redirectURI)
	data.Set("grant_type", "authorization_code")

	resp, err := http.PostForm("https://oauth2.googleapis.com/token", data)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("google token exchange failed: %s", string(body))
	}

	var result struct {
		IDToken string `json:"id_token"`
	}
	if err := json.Unmarshal(body, &result); err != nil || result.IDToken == "" {
		return "", fmt.Errorf("no id_token in response")
	}

	return result.IDToken, nil
}
