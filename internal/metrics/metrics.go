package metrics

import (
	"os"

	"github.com/posthog/posthog-go"
)

var client posthog.Client

func Init() {
	key := os.Getenv("POSTHOG_KEY")
	host := os.Getenv("POSTHOG_HOST")
	if key == "" {
		return
	}
	if host == "" {
		host = "https://eu.posthog.com"
	}
	client, _ = posthog.NewWithConfig(key, posthog.Config{Endpoint: host})
}

func Close() {
	if client != nil {
		client.Close()
	}
}

func UserRegistered(userID uint, email string) {
	capture(userID, "user_registered", posthog.Properties{"email": email})
}

func EventCreated(userID uint, eventID uint) {
	capture(userID, "event_created", posthog.Properties{"event_id": eventID})
}

func capture(userID uint, event string, props posthog.Properties) {
	if client == nil {
		return
	}
	client.Enqueue(posthog.Capture{
		DistinctId: userIDString(userID),
		Event:      event,
		Properties: props,
	})
}

func userIDString(id uint) string {
	return "user_" + itoa(id)
}

func itoa(n uint) string {
	if n == 0 {
		return "0"
	}
	b := make([]byte, 0, 10)
	for n > 0 {
		b = append([]byte{byte('0' + n%10)}, b...)
		n /= 10
	}
	return string(b)
}
