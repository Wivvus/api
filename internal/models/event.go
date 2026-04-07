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
	DistanceKm    float64    `json:"distance_km"`
	PaceMinKm     float64    `json:"pace_min_km"`
	Location      Location   `json:"location" gorm:"embedded"`
	CreatorUserID  uint      `json:"creator_id"`
}

type EventAPIDecorator struct {
	ID            uint      `json:"id"`
	Name          string    `json:"name"`
	Description   string    `json:"description"`
	StartTime     time.Time `json:"start_time"`
	DistanceKm    float64    `json:"distance_km"`
	PaceMinKm     float64    `json:"pace_min_km"`
	Location      Location   `json:"location"`
	CreatorID     uint       `json:"creator_id"`
	AttendeeCount int64      `json:"attendee_count"`
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

func (e *EventRepo) All() Events {
	events := Events{}
	db.Where("start_time >= ?", time.Now()).Find(&events)
	return events
}

type BoundingBox struct {
	LatMin, LatMax, LngMin, LngMax float64
}

func (e *EventRepo) AllInBounds(b BoundingBox) Events {
	events := Events{}
	db.Where("start_time >= ? AND lat BETWEEN ? AND ? AND long BETWEEN ? AND ?",
		time.Now(), b.LatMin, b.LatMax, b.LngMin, b.LngMax,
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

func (e *EventRepo) DeleteByID(id string) error {
	return db.Delete(&Event{}, id).Error
}
