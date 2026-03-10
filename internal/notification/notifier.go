package notification

import (
	"fmt"
	"log/slog"

	"github.com/gen2brain/beeep"
)

// Notifier sends OS-level desktop notifications (visual only).
// Audio is handled by the frontend SoundPlayer (Web Audio API).
type Notifier struct{}

// NewNotifier creates a new Notifier.
func NewNotifier() *Notifier {
	return &Notifier{}
}

// Alert sends a visual OS notification.
// muted is accepted for API compatibility but audio is always
// handled by the frontend — so this call never produces a system beep.
func (n *Notifier) Alert(title, message string, muted bool) error {
	slog.Info("sending notification", "title", title, "muted", muted)
	if err := beeep.Notify(title, message, ""); err != nil {
		slog.Warn("notification failed", "error", err)
		return fmt.Errorf("send notification: %w", err)
	}
	return nil
}

// Notify is kept for backwards compatibility.
func (n *Notifier) Notify(title, message string) error {
	return n.Alert(title, message, false)
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
