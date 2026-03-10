package models

import (
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
	Name          string     `json:"name"`
	Description   string     `json:"description"`
	StartTime     *time.Time `json:"start_time,omitempty"`
	DistanceKm    float64    `json:"distance_km"`
	PaceMinKm     float64    `json:"pace_min_km"`
	Location      Location   `json:"location" gorm:"embedded"`
	CreatorUserID  uint      `json:"creator_id"`
}

type EventAPIDecorator struct {
	ID          uint       `json:"id"`
	Name        string     `json:"name"`
	Description string     `json:"description"`
	StartTime   *time.Time `json:"start_time,omitempty"`
	DistanceKm  float64    `json:"distance_km"`
	PaceMinKm   float64    `json:"pace_min_km"`
	Location    Location   `json:"location"`
	CreatorID   uint       `json:"creator_id"`
}

func (e *Event) ToAPI() *EventAPIDecorator {
	return &EventAPIDecorator{
		ID:          e.ID,
		Name:        e.Name,
		Description: e.Description,
		StartTime:   e.StartTime,
		DistanceKm:  e.DistanceKm,
		PaceMinKm:   e.PaceMinKm,
		Location:    e.Location,
		CreatorID:   e.CreatorUserID,
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
	db.Find(&events)
	return events
}

func (e *EventRepo) DeleteByID(id string) error {
	return db.Delete(&Event{}, id).Error
}
