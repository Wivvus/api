package events

import (
	"net/http"
	"strconv"
	"time"

	"github.com/Wivvus/api/internal/middleware"
	"github.com/Wivvus/api/internal/models"
	"github.com/gin-gonic/gin"
)

func ConfigureRouter(r *gin.Engine) {
	r.POST("/event", middleware.AuthRequired(), create)
	r.GET("/events", list)
	r.GET("/event/:id", get)
	r.PUT("/event/:id", middleware.AuthRequired(), update)
	r.DELETE("/event/:id", middleware.AuthRequired(), delete)
	r.GET("/event/:id/attendees", middleware.AuthRequired(), attendees)
	r.POST("/event/:id/attend", middleware.AuthRequired(), attend)
	r.DELETE("/event/:id/attend", middleware.AuthRequired(), drop)
}

func create(ctx *gin.Context) {
	user := ctx.MustGet("user").(*models.User)

	var newEvent models.Event
	if err := ctx.ShouldBindJSON(&newEvent); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := models.ValidateEvent(&newEvent); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	newEvent.CreatorUserID = user.ID

	er := models.EventRepo{}
	er.CreateOrUpdate(&newEvent)

	ar := models.AttendanceRepo{}
	ar.Attend(newEvent.ID, user.ID)

	ctx.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "event received succesfully",
		"event":   newEvent.ToAPI(),
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

func update(ctx *gin.Context) {
	user := ctx.MustGet("user").(*models.User)
	id := ctx.Param("id")

	er := models.EventRepo{}
	existing := er.FindByID(id)
	if existing == nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "event not found"})
		return
	}
	if existing.CreatorUserID != user.ID {
		ctx.JSON(http.StatusForbidden, gin.H{"error": "not the event creator"})
		return
	}

	var updated models.Event
	if err := ctx.ShouldBindJSON(&updated); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := models.ValidateEvent(&updated); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	updated.Model = existing.Model
	updated.CreatorUserID = existing.CreatorUserID
	er.CreateOrUpdate(&updated)

	ctx.JSON(http.StatusOK, updated.ToAPI())
}

func delete(ctx *gin.Context) {
	user := ctx.MustGet("user").(*models.User)
	id := ctx.Param("id")

	er := models.EventRepo{}
	existing := er.FindByID(id)
	if existing == nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "event not found"})
		return
	}
	if existing.CreatorUserID != user.ID {
		ctx.JSON(http.StatusForbidden, gin.H{"error": "not the event creator"})
		return
	}

	er.DeleteByID(id)
	ctx.JSON(http.StatusOK, gin.H{"success": true})
}

func attendees(ctx *gin.Context) {
	user := ctx.MustGet("user").(*models.User)
	id := ctx.Param("id")
	er := models.EventRepo{}
	event := er.FindByID(id)
	if event == nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "event not found"})
		return
	}
	ar := models.AttendanceRepo{}
	ctx.JSON(http.StatusOK, models.AttendeesResponse{
		Attendees:   ar.AttendeesForEvent(event.ID),
		IsAttending: ar.IsAttending(event.ID, user.ID),
	})
}

func attend(ctx *gin.Context) {
	user := ctx.MustGet("user").(*models.User)
	id := ctx.Param("id")
	er := models.EventRepo{}
	event := er.FindByID(id)
	if event == nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "event not found"})
		return
	}
	ar := models.AttendanceRepo{}
	if err := ar.Attend(event.ID, user.ID); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"success": true})
}

func drop(ctx *gin.Context) {
	user := ctx.MustGet("user").(*models.User)
	id := ctx.Param("id")
	er := models.EventRepo{}
	event := er.FindByID(id)
	if event == nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "event not found"})
		return
	}
	ar := models.AttendanceRepo{}
	if err := ar.Drop(event.ID, user.ID); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"success": true})
}

func parseFilters(ctx *gin.Context) models.Filters {
	f := models.Filters{}
	if v, err := strconv.ParseFloat(ctx.Query("min_pace"), 64); err == nil {
		f.MinPace = v
	}
	if v, err := strconv.ParseFloat(ctx.Query("max_pace"), 64); err == nil {
		f.MaxPace = v
	}
	if v, err := strconv.ParseFloat(ctx.Query("min_length"), 64); err == nil {
		f.MinLength = v
	}
	if v, err := strconv.ParseFloat(ctx.Query("max_length"), 64); err == nil {
		f.MaxLength = v
	}
	if v, err := strconv.ParseFloat(ctx.Query("max_radius"), 64); err == nil {
		f.MaxRadius = v
	}
	if v, err := strconv.ParseFloat(ctx.Query("user_lat"), 64); err == nil {
		f.UserLat = v
	}
	if v, err := strconv.ParseFloat(ctx.Query("user_lng"), 64); err == nil {
		f.UserLng = v
	}
	if v, err := time.Parse("2006-01-02", ctx.Query("date_from")); err == nil {
		f.DateFrom = &v
	}
	return f
}

func list(ctx *gin.Context) {
	er := models.EventRepo{}
	f := parseFilters(ctx)

	latMin, errA := strconv.ParseFloat(ctx.Query("lat_min"), 64)
	latMax, errB := strconv.ParseFloat(ctx.Query("lat_max"), 64)
	lngMin, errC := strconv.ParseFloat(ctx.Query("lng_min"), 64)
	lngMax, errD := strconv.ParseFloat(ctx.Query("lng_max"), 64)

	if errA == nil && errB == nil && errC == nil && errD == nil {
		bbox := models.BoundingBox{LatMin: latMin, LatMax: latMax, LngMin: lngMin, LngMax: lngMax}
		ctx.JSON(http.StatusOK, er.AllInBounds(bbox, f).ToAPI())
		return
	}

	ctx.JSON(http.StatusOK, er.All(f).ToAPI())
}
