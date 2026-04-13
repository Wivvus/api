package app

import (
	"github.com/Wivvus/api/internal/app/auth"
	"github.com/Wivvus/api/internal/app/events"
	"github.com/Wivvus/api/internal/app/ratings"
	"github.com/gin-gonic/gin"
)

func ConfigureRouter(r *gin.Engine) {
	auth.ConfigureRouter(r)
	events.ConfigureRouter(r)
	ratings.ConfigureRouter(r)
}
