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
		eventURL := fmt.Sprintf("%s/events/%d", appURL, event.ID)
		for _, user := range ar.AttendeesNeedingReminderForEvent(event.ID) {
			if err := email.SendEventReminder(user.Email, user.Name, event.Name, event.StartTime, eventURL); err != nil {
				log.Printf("failed to send reminder to %s for event %d: %v", user.Email, event.ID, err)
				continue
			}
			ar.MarkReminderSent(event.ID, user.ID)
		}
	}

	// Post-event rating reminders (12 hours after)
	for _, event := range er.EventsForRatingReminder() {
		eventURL := fmt.Sprintf("%s/events/%d", appURL, event.ID)
		for _, user := range ar.AttendeesNeedingRatingReminderForEvent(event.ID, event.CreatorUserID) {
			if err := email.SendRatingReminder(user.Email, user.Name, event.Name, eventURL); err != nil {
				log.Printf("failed to send rating reminder to %s for event %d: %v", user.Email, event.ID, err)
				continue
			}
			ar.MarkRatingReminderSent(event.ID, user.ID)
		}
	}
}
