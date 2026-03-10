package notification

import (
	"fmt"
	"log/slog"

	"github.com/gen2brain/beeep"
)

// Notifier sends OS-level notifications.
type Notifier struct{}

// NewNotifier creates a new Notifier.
func NewNotifier() *Notifier {
	return &Notifier{}
}

// Notify sends an OS notification with the given title and message.
func (n *Notifier) Notify(title, message string) error {
	slog.Info("sending notification", "title", title, "message", message)
	if err := beeep.Notify(title, message, ""); err != nil {
		slog.Warn("notification failed", "error", err)
		return fmt.Errorf("send notification: %w", err)
	}
	return nil
}

// SessionCompleteMessage returns the notification message for a completed session.
func SessionCompleteMessage(sessionType string) (title, message string) {
	switch sessionType {
	case "focus":
		return "Focus Session Complete!", "Great work! Time for a break."
	case "short_break":
		return "Break Over!", "Ready to focus again?"
	case "long_break":
		return "Long Break Over!", "Refreshed and ready to go!"
	default:
		return "Timer Complete", "Session finished."
	}
}
