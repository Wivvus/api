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

	events := er.EventsInReminderWindow()
	for _, event := range events {
		attendees := ar.AttendeesNeedingReminderForEvent(event.ID)
		for _, user := range attendees {
			eventURL := fmt.Sprintf("%s/events/%d", appURL, event.ID)
			if err := email.SendEventReminder(user.Email, user.Name, event.Name, event.StartTime, eventURL); err != nil {
				log.Printf("failed to send reminder to %s for event %d: %v", user.Email, event.ID, err)
				continue
			}
			ar.MarkReminderSent(event.ID, user.ID)
		}
	}
}
