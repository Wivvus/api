package models

type Attendance struct {
	EventID uint `gorm:"primaryKey"`
	UserID  uint `gorm:"primaryKey"`
}

type AttendeeInfo struct {
	Name      string `json:"name"`
	AvatarURL string `json:"avatar_url"`
}

type AttendeesResponse struct {
	Attendees   []*AttendeeInfo `json:"attendees"`
	IsAttending bool            `json:"is_attending"`
}

type AttendanceRepo struct{}

func (a *AttendanceRepo) Attend(eventID, userID uint) error {
	return db.FirstOrCreate(&Attendance{}, Attendance{EventID: eventID, UserID: userID}).Error
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

func (a *AttendanceRepo) AttendeesForEvent(eventID uint) []*AttendeeInfo {
	var users []User
	db.Joins("JOIN attendances ON attendances.user_id = users.id AND attendances.event_id = ?", eventID).Find(&users)
	result := make([]*AttendeeInfo, len(users))
	for i, u := range users {
		result[i] = &AttendeeInfo{Name: u.Name, AvatarURL: u.AvatarURL}
	}
	return result
}
