package models

type Rating struct {
	EventID uint   `gorm:"primaryKey"`
	UserID  uint   `gorm:"primaryKey"`
	Score   int    `json:"score"`
	Comment string `json:"comment"`
}

type RatingInfo struct {
	UserName  string  `json:"user_name"`
	AvatarURL string  `json:"avatar_url"`
	Score     int     `json:"score"`
	Comment   string  `json:"comment,omitempty"`
}

type RatingRepo struct{}

func (r *RatingRepo) Upsert(rating *Rating) error {
	return db.Save(rating).Error
}

func (r *RatingRepo) FindByEventAndUser(eventID, userID uint) *Rating {
	rating := &Rating{}
	if db.Where("event_id = ? AND user_id = ?", eventID, userID).First(rating).Error != nil {
		return nil
	}
	return rating
}

func (r *RatingRepo) ForEvent(eventID uint) []*RatingInfo {
	results := []*RatingInfo{}
	db.Model(&Rating{}).
		Select("users.name as user_name, users.avatar_url, ratings.score, ratings.comment").
		Joins("JOIN users ON users.id = ratings.user_id").
		Where("ratings.event_id = ?", eventID).
		Scan(&results)
	return results
}

func (r *RatingRepo) AverageForCreator(creatorUserID uint) *float64 {
	var avg *float64
	db.Model(&Rating{}).
		Select("AVG(ratings.score)").
		Joins("JOIN events ON events.id = ratings.event_id").
		Where("events.creator_user_id = ? AND events.deleted_at IS NULL", creatorUserID).
		Scan(&avg)
	return avg
}

func (r *RatingRepo) AverageForEvent(eventID uint) *float64 {
	var avg *float64
	db.Model(&Rating{}).Select("AVG(score)").Where("event_id = ?", eventID).Scan(&avg)
	return avg
}

type RatingInfoWithEvent struct {
	EventName string `json:"event_name"`
	EventID   uint   `json:"event_id"`
	UserName  string `json:"user_name"`
	AvatarURL string `json:"avatar_url"`
	Score     int    `json:"score"`
	Comment   string `json:"comment,omitempty"`
}

func (r *RatingRepo) ForCreator(creatorUserID uint) []*RatingInfoWithEvent {
	results := []*RatingInfoWithEvent{}
	db.Model(&Rating{}).
		Select("events.name as event_name, events.id as event_id, users.name as user_name, users.avatar_url, ratings.score, ratings.comment").
		Joins("JOIN users ON users.id = ratings.user_id").
		Joins("JOIN events ON events.id = ratings.event_id").
		Where("events.creator_user_id = ? AND events.deleted_at IS NULL", creatorUserID).
		Order("events.start_time DESC").
		Scan(&results)
	return results
}
