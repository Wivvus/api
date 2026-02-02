package events

import (
	"net/http"

	"github.com/Wivvus/api/internal/models"
	"github.com/Wivvus/api/internal/oauth"
	"github.com/gin-gonic/gin"
)

func ConfigureRouter(r *gin.Engine) {
	r.POST("/event", oauth.RequireAuth(), create)
	r.GET("/events", list)
}

func create(ctx *gin.Context) {
	var newEvent models.Event
	if err := ctx.ShouldBindJSON(&newEvent); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	er := models.EventRepo{}
	er.CreateOrUpdate(&newEvent)

	ctx.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "event received succesfully",
		"event":   newEvent,
	})
}

func list(ctx *gin.Context) {
	eventsRepo := models.EventRepo{}

	events := eventsRepo.All()

	ctx.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Page loaded successfully",
		"data":    events,
	})
}
