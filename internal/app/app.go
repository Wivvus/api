package app

import (
	"github.com/Wivvus/api/internal/app/events"
	"github.com/gin-gonic/gin"
)

func ConfigureRouter(r *gin.Engine) {
	events.ConfigureRouter(r)
}
