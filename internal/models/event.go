package models

import "gorm.io/gorm"

type Event struct {
	gorm.Model
	Name string
}

type EventRepo struct {
}

func (e *EventRepo) FindByID(id string) *Event {
	event := &Event{}
	db.First(event, id)
	return event
}

func (e *EventRepo) CreateOrUpdate(event *Event) error {
	if db.Model(event).Updates(event).RowsAffected == 0 {
		return db.Create(event).Error
	}
	return nil
}

func (e *EventRepo) All() []*Event {
	events := []*Event{}
	db.Find(&events)
	return events
}
