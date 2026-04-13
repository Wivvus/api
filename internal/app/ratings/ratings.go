package ratings

import (
	"net/http"
	"time"

	"github.com/Wivvus/api/internal/middleware"
	"github.com/Wivvus/api/internal/models"
	"github.com/gin-gonic/gin"
)

func ConfigureRouter(r *gin.Engine) {
	r.POST("/event/:id/rate", middleware.AuthRequired(), rate)
	r.GET("/event/:id/ratings", getRatings)
}

func rate(c *gin.Context) {
	user := c.MustGet("user").(*models.User)
	id := c.Param("id")

	er := models.EventRepo{}
	event := er.FindByID(id)
	if event == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "event not found"})
		return
	}

	if event.CreatorUserID == user.ID {
		c.JSON(http.StatusForbidden, gin.H{"error": "you cannot rate your own event"})
		return
	}

	if event.StartTime.After(time.Now()) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "event has not happened yet"})
		return
	}

	ar := models.AttendanceRepo{}
	if !ar.IsAttending(event.ID, user.ID) {
		c.JSON(http.StatusForbidden, gin.H{"error": "you did not attend this event"})
		return
	}

	var body struct {
		Score   int    `json:"score" binding:"required"`
		Comment string `json:"comment"`
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "score is required"})
		return
	}
	if body.Score < 1 || body.Score > 5 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "score must be between 1 and 5"})
		return
	}

	rr := models.RatingRepo{}
	rating := &models.Rating{
		EventID: event.ID,
		UserID:  user.ID,
		Score:   body.Score,
		Comment: body.Comment,
	}
	if err := rr.Upsert(rating); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to save rating"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true})
}

func getRatings(c *gin.Context) {
	id := c.Param("id")

	er := models.EventRepo{}
	event := er.FindByID(id)
	if event == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "event not found"})
		return
	}

	rr := models.RatingRepo{}
	c.JSON(http.StatusOK, gin.H{
		"ratings": rr.ForEvent(event.ID),
		"average": rr.AverageForEvent(event.ID),
	})
}
