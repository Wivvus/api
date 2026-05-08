package models

import "time"

type Attendance struct {
	EventID            uint  `gorm:"primaryKey"`
	UserID             uint  `gorm:"primaryKey"`
	EventOptionID      *uint `gorm:"default:null"`
	ReminderSent       bool
	RatingReminderSent bool
}

type AttendeeInfo struct {
	Name       string  `json:"name"`
	AvatarURL  string  `json:"avatar_url"`
	OptionText *string `json:"option_text,omitempty"`
}

type AttendeesResponse struct {
	Attendees    []*AttendeeInfo `json:"attendees"`
	IsAttending  bool            `json:"is_attending"`
	MyOptionID   *uint           `json:"my_option_id"`
}

type AttendanceRepo struct{}

func (a *AttendanceRepo) Attend(eventID, userID uint, optionID *uint) error {
	var event Event
	db.First(&event, eventID)
	reminderSent := !event.StartTime.IsZero() && event.StartTime.Before(time.Now().Add(12*time.Hour))

	var existing Attendance
	result := db.Where("event_id = ? AND user_id = ?", eventID, userID).First(&existing)
	if result.Error != nil {
		// Not found — create
		return db.Create(&Attendance{
			EventID:       eventID,
			UserID:        userID,
			EventOptionID: optionID,
			ReminderSent:  reminderSent,
		}).Error
	}
	// Found — update option
	return db.Model(&existing).Update("event_option_id", optionID).Error
}

func (a *AttendanceRepo) AttendeesNeedingReminderForEvent(eventID uint) []*User {
	var users []*User
	db.Joins("JOIN attendances ON attendances.user_id = users.id AND attendances.event_id = ? AND attendances.reminder_sent = false", eventID).Find(&users)
	return users
}

func (a *AttendanceRepo) MarkReminderSent(eventID, userID uint) {
	db.Model(&Attendance{}).Where("event_id = ? AND user_id = ?", eventID, userID).Update("reminder_sent", true)
}

func (a *AttendanceRepo) AttendeesNeedingRatingReminderForEvent(eventID uint, creatorUserID uint) []*User {
	var users []*User
	db.Joins("JOIN attendances ON attendances.user_id = users.id AND attendances.event_id = ? AND attendances.rating_reminder_sent = false", eventID).
		Where("users.id != ?", creatorUserID).
		Find(&users)
	return users
}

func (a *AttendanceRepo) MarkRatingReminderSent(eventID, userID uint) {
	db.Model(&Attendance{}).Where("event_id = ? AND user_id = ?", eventID, userID).Update("rating_reminder_sent", true)
}

func (a *AttendanceRepo) Drop(eventID, userID uint) error {
	return db.Delete(&Attendance{}, "event_id = ? AND user_id = ?", eventID, userID).Error
}

func (a *AttendanceRepo) IsAttending(eventID, userID uint) bool {
	var count int64
	db.Model(&Attendance{}).Where("event_id = ? AND user_id = ?", eventID, userID).Count(&count)
	return count > 0
}

func (a *AttendanceRepo) CountForEvent(eventID uint) int64 {
	var count int64
	db.Model(&Attendance{}).Where("event_id = ?", eventID).Count(&count)
	return count
}

func (a *AttendanceRepo) MyOptionID(eventID, userID uint) *uint {
	var att Attendance
	if db.Where("event_id = ? AND user_id = ?", eventID, userID).First(&att).Error != nil {
		return nil
	}
	return att.EventOptionID
}

type attendeeRow struct {
	Name      string
	AvatarURL string
	OptionText *string
}

func (a *AttendanceRepo) AttendeesForEvent(eventID uint) []*AttendeeInfo {
	var rows []attendeeRow
	db.Table("users").
		Select("users.name, users.avatar_url, event_options.text AS option_text").
		Joins("JOIN attendances ON attendances.user_id = users.id AND attendances.event_id = ?", eventID).
		Joins("LEFT JOIN event_options ON event_options.id = attendances.event_option_id").
		Scan(&rows)
	result := make([]*AttendeeInfo, len(rows))
	for i, r := range rows {
		result[i] = &AttendeeInfo{Name: r.Name, AvatarURL: r.AvatarURL, OptionText: r.OptionText}
	}
	return result
}
