package events

import (
	"net/http"

	"github.com/Wivvus/api/internal/middleware"
	"github.com/Wivvus/api/internal/models"
	"github.com/gin-gonic/gin"
)

func ConfigureRouter(r *gin.Engine) {
	r.POST("/event", middleware.AuthRequired(), create)
	r.GET("/events", list)
	r.GET("/event/:id", get)
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

func get(ctx *gin.Context) {
	id := ctx.Param("id")
	er := models.EventRepo{}
	event := er.FindByID(id)
	if event == nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "event not found"})
		return
	}
	ctx.JSON(http.StatusOK, event.ToAPI())
}

func list(ctx *gin.Context) {
	eventsRepo := models.EventRepo{}
	events := eventsRepo.All()
	ctx.JSON(http.StatusOK, events.ToAPI())
}
