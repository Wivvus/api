package reminders

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/Wivvus/api/internal/email"
	"github.com/Wivvus/api/internal/models"
)

func Start() {
	go func() {
		ticker := time.NewTicker(5 * time.Minute)
		defer ticker.Stop()
		for range ticker.C {
			run()
		}
	}()
}

func run() {
	er := models.EventRepo{}
	ar := models.AttendanceRepo{}
	appURL := os.Getenv("APP_URL")

	// Pre-event reminders (12 hours before)
	for _, event := range er.EventsInReminderWindow() {
		eventURL := fmt.Sprintf("%s/run/%d", appURL, event.ID)
		for _, user := range ar.AttendeesNeedingReminderForEvent(event.ID) {
			if err := email.SendEventReminder(user.Email, user.Name, event.Name, event.StartTime, eventURL); err != nil {
				log.Printf("failed to send reminder to %s for event %d: %v", user.Email, event.ID, err)
				continue
			}
			ar.MarkReminderSent(event.ID, user.ID)
		}
	}

	// Post-event rating reminders (12 hours after)
	ur := models.UserRepo{}
	for _, event := range er.EventsForRatingReminder() {
		eventURL := fmt.Sprintf("%s/run/%d/review", appURL, event.ID)
		creatorName := ""
		if creator := ur.FindByID(fmt.Sprintf("%d", event.CreatorUserID)); creator != nil {
			creatorName = creator.Name
		}
		for _, user := range ar.AttendeesNeedingRatingReminderForEvent(event.ID, event.CreatorUserID) {
			if err := email.SendRatingReminder(user.Email, user.Name, event.Name, creatorName, eventURL); err != nil {
				log.Printf("failed to send rating reminder to %s for event %d: %v", user.Email, event.ID, err)
				continue
			}
			ar.MarkRatingReminderSent(event.ID, user.ID)
		}
	}
}
