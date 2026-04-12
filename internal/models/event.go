package models

import (
	"fmt"
	"time"

	"gorm.io/gorm"
)

type Events []*Event

func (e Events) ToAPI() []*EventAPIDecorator {
	ret := []*EventAPIDecorator{}

	for _, ev := range e {
		ret = append(ret, ev.ToAPI())
	}
	return ret
}

type Event struct {
	gorm.Model
	Name          string    `json:"name"`
	Description   string    `json:"description"`
	StartTime     time.Time `json:"start_time"`
	DistanceKm    float64   `json:"distance_km"`
	PaceMinKm     float64   `json:"pace_min_km"`
	AllPaces      bool      `json:"all_paces"`
	Location      Location  `json:"location" gorm:"embedded"`
	CreatorUserID uint      `json:"creator_id"`
}

type EventAPIDecorator struct {
	ID            uint      `json:"id"`
	Name          string    `json:"name"`
	Description   string    `json:"description"`
	StartTime     time.Time `json:"start_time"`
	DistanceKm    float64   `json:"distance_km"`
	PaceMinKm     float64   `json:"pace_min_km"`
	AllPaces      bool      `json:"all_paces"`
	Location      Location  `json:"location"`
	CreatorID     uint      `json:"creator_id"`
	AttendeeCount int64     `json:"attendee_count"`
}

func (e *Event) ToAPI() *EventAPIDecorator {
	ar := &AttendanceRepo{}
	return &EventAPIDecorator{
		ID:            e.ID,
		Name:          e.Name,
		Description:   e.Description,
		StartTime:     e.StartTime,
		DistanceKm:    e.DistanceKm,
		PaceMinKm:     e.PaceMinKm,
		AllPaces:      e.AllPaces,
		Location:      e.Location,
		CreatorID:     e.CreatorUserID,
		AttendeeCount: ar.CountForEvent(e.ID),
	}
}

type EventRepo struct {
}

func (e *EventRepo) FindByID(id string) *Event {
	event := &Event{}
	if result := db.First(event, id); result.Error != nil {
		return nil
	}
	return event
}

func (e *EventRepo) CreateOrUpdate(event *Event) error {
	if db.Model(event).Updates(event).RowsAffected == 0 {
		return db.Create(event).Error
	}
	return nil
}

type Filters struct {
	MinLength  float64
	MaxLength  float64
	MinPace    float64
	MaxPace    float64
	MaxRadius  float64
	UserLat    float64
	UserLng    float64
	DateFrom   *time.Time
}

func applyFilters(q *gorm.DB, f Filters) *gorm.DB {
	now := time.Now()
	start := now
	if f.DateFrom != nil && f.DateFrom.After(now) {
		start = *f.DateFrom
	}
	q = q.Where("start_time >= ?", start)
	if f.MinLength > 0 {
		q = q.Where("distance_km >= ?", f.MinLength)
	}
	if f.MaxLength > 0 {
		q = q.Where("distance_km <= ?", f.MaxLength)
	}
	if f.MinPace > 0 {
		q = q.Where("all_paces = true OR pace_min_km >= ?", f.MinPace)
	}
	if f.MaxPace > 0 {
		q = q.Where("all_paces = true OR pace_min_km <= ?", f.MaxPace)
	}
	if f.MaxRadius > 0 && (f.UserLat != 0 || f.UserLng != 0) {
		q = q.Where(`(6371 * acos(
			cos(radians(?)) * cos(radians(lat)) * cos(radians(long) - radians(?)) +
			sin(radians(?)) * sin(radians(lat))
		)) <= ?`, f.UserLat, f.UserLng, f.UserLat, f.MaxRadius)
	}
	return q
}

func (e *EventRepo) All(f Filters) Events {
	events := Events{}
	applyFilters(db, f).Find(&events)
	return events
}

type BoundingBox struct {
	LatMin, LatMax, LngMin, LngMax float64
}

func (e *EventRepo) AllInBounds(b BoundingBox, f Filters) Events {
	events := Events{}
	q := applyFilters(db, f)
	q.Where("lat BETWEEN ? AND ? AND long BETWEEN ? AND ?",
		b.LatMin, b.LatMax, b.LngMin, b.LngMax,
	).Find(&events)
	return events
}

func ValidateEvent(e *Event) error {
	if e.Name == "" {
		return fmt.Errorf("event name is required")
	}
	if e.Location.Lat == 0 && e.Location.Long == 0 {
		return fmt.Errorf("location is required")
	}
	if e.StartTime.IsZero() {
		return fmt.Errorf("start time is required")
	}
	if e.StartTime.Before(time.Now()) {
		return fmt.Errorf("start time must be in the future")
	}
	return nil
}

func (e *EventRepo) AllByCreator(userID uint) Events {
	events := Events{}
	db.Where("creator_user_id = ?", userID).Order("start_time DESC").Find(&events)
	return events
}

func (e *EventRepo) DeleteByID(id string) error {
	return db.Delete(&Event{}, id).Error
}
