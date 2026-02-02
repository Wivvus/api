package oauth

import (
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/markbates/goth"
	"github.com/markbates/goth/gothic"
)

const (
	LOGOUT_ROUTE = "/auth/logout"
	USER_ROUTE   = "/auth/user"

	USER_SESSION_KEY = "USER"
	USER_CONTEXT_KEY = "USER"
)

var (
	cookieStoreKey    string
	cookieStoreSecret string
)

func init() {
	cookieStoreKey = os.Getenv("COOKIE_STORE_KEY")
	cookieStoreSecret = os.Getenv("COOKIE_STORE_SECRET")
}

func ConfigureRouter(r *gin.Engine) {
	r.GET(LOGOUT_ROUTE, LogoutWithGoogle)
	r.GET(USER_ROUTE, GetAuthUser)

	googleConfigureRouter(r)
}

func GetAuthUser(ctx *gin.Context) {
	user, err := GetSessionUser(ctx)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Error retrieving user session",
			"error":   err.Error(),
		})
		return
	}

	ctx.Keys[USER_CONTEXT_KEY] = user

	// Return user info
	ctx.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "User fetched successfully",
		"data":    user,
	})
}

func GetSessionUser(ctx *gin.Context) (goth.User, error) {
	// Retrieve the session
	session, err := gothic.Store.Get(ctx.Request, cookieStoreSecret)
	if err != nil {
		return goth.User{}, err
	}

	var ok bool
	// Get user data from session
	user, ok := session.Values[USER_SESSION_KEY].(goth.User)
	if !ok {
		return user, ErrNoUserInSession
	}
	return user, nil

}

func RequireAuth() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// gets the user session from the request+
		user, err := GetSessionUser(ctx)
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"success": false,
				"message": "Unauthorized User",
				"error":   err.Error(),
			})
			return
		}
		ctx.Keys[USER_CONTEXT_KEY] = user
		// calls the next middleware or handler
		ctx.Next()
	}
}
