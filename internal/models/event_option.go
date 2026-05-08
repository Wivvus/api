package models

import "gorm.io/gorm"

type EventOption struct {
	ID        uint   `json:"id" gorm:"primaryKey;autoIncrement"`
	EventID   uint   `json:"event_id" gorm:"index"`
	Text      string `json:"text"`
	SortOrder int    `json:"sort_order"`
}

type EventOptionAPI struct {
	ID   uint   `json:"id"`
	Text string `json:"text"`
}

type EventOptionRepo struct{}

func (r *EventOptionRepo) ForEvent(eventID uint) []*EventOptionAPI {
	var opts []*EventOption
	db.Where("event_id = ?", eventID).Order("sort_order").Find(&opts)
	result := make([]*EventOptionAPI, len(opts))
	for i, o := range opts {
		result[i] = &EventOptionAPI{ID: o.ID, Text: o.Text}
	}
	return result
}

func (r *EventOptionRepo) ReplaceForEvent(eventID uint, texts []string) error {
	return db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Delete(&EventOption{}, "event_id = ?", eventID).Error; err != nil {
			return err
		}
		var firstID *uint
		for i, text := range texts {
			if text == "" {
				continue
			}
			opt := &EventOption{EventID: eventID, Text: text, SortOrder: i}
			if err := tx.Create(opt).Error; err != nil {
				return err
			}
			if firstID == nil {
				id := opt.ID
				firstID = &id
			}
		}
		// Remap attendances whose option was deleted to the first remaining option (or NULL if none).
		if err := tx.Model(&Attendance{}).
			Where("event_id = ?", eventID).
			Update("event_option_id", firstID).Error; err != nil {
			return err
		}
		return nil
	})
}
