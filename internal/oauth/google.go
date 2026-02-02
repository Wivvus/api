package oauth

import (
	"net/http"
	"os"

	"github.com/Wivvus/api/internal/models"
	"github.com/gin-gonic/gin"
	"github.com/markbates/goth"
	"github.com/markbates/goth/gothic"
	"github.com/markbates/goth/providers/google"
)

const (
	GOOGLE_SIGNIN_ROUTE = "/auth/google/signup"
)

var (
	googleOauthRedirectPath string
)

func init() {
	googleOauthRedirectPath = os.Getenv("GOOGLE_OAUTH_REDIRECT_PATH")

	googleClientID := os.Getenv("GOOGLE_OAUTH_CLIENT_ID")
	googleSecret := os.Getenv("GOOGLE_OAUTH_SECRET")
	googleOauthRedirect := os.Getenv("GOOGLE_OAUTH_REDIRECT")

	// Register Google as an authentication provider
	goth.UseProviders(
		google.New(
			googleClientID,
			googleSecret,
			googleOauthRedirect,
		),
	)
}

func googleConfigureRouter(r *gin.Engine) {
	r.GET(GOOGLE_SIGNIN_ROUTE, SignupWithGoogle)
	r.GET(googleOauthRedirectPath, HandleGoogleAuth)
}

func LogoutWithGoogle(ctx *gin.Context) {
	gothic.Logout(ctx.Writer, ctx.Request)

	// store the user session
	session, err := gothic.Store.New(ctx.Request, "secret")
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Error storing user session",
			"error":   err.Error(),
		})
		return
	}

	// your logic for storing the user in database goes here

	session.Values[USER_SESSION_KEY] = nil

	// save the user session
	if err = session.Save(ctx.Request, ctx.Writer); err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Error saving user session",
			"error":   err.Error(),
		})
		return
	}

	ctx.Redirect(http.StatusTemporaryRedirect, "/auth/user")
}

func SignupWithGoogle(ctx *gin.Context) {
	query := ctx.Request.URL.Query()
	query.Add("provider", "google")
	ctx.Request.URL.RawQuery = query.Encode()

	gothic.BeginAuthHandler(ctx.Writer, ctx.Request)
}

func HandleGoogleAuth(ctx *gin.Context) {
	query := ctx.Request.URL.Query()
	query.Add("provider", "google")
	ctx.Request.URL.RawQuery = query.Encode()

	var err error
	user, err := gothic.CompleteUserAuth(ctx.Writer, ctx.Request)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Error authenticating user",
			"error":   err.Error(),
		})
		return
	}

	// store the user session
	session, err := gothic.Store.New(ctx.Request, "secret")
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Error storing user session",
			"error":   err.Error(),
		})
		return
	}

	// update or create in DB
	userRepo := &models.UserRepo{}
	userRepo.CreateOrUpdate(models.UserFromGoogleOAuth(user))

	session.Values[USER_SESSION_KEY] = user

	// save the user session
	if err = session.Save(ctx.Request, ctx.Writer); err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Error saving user session",
			"error":   err.Error(),
		})
		return
	}

	ctx.Redirect(http.StatusTemporaryRedirect, USER_ROUTE)
}
