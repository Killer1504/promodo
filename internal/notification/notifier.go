package notification

import (
	"fmt"
	"log/slog"

	"github.com/gen2brain/beeep"
)

// Notifier sends OS-level notifications with sound.
type Notifier struct{}

// NewNotifier creates a new Notifier.
func NewNotifier() *Notifier {
	return &Notifier{}
}

// Alert sends an OS notification WITH an audible beep.
// Uses beeep.Alert (visual toast + Beep syscall).
// If muted is true, sends a silent visual-only notification.
func (n *Notifier) Alert(title, message string, muted bool) error {
	slog.Info("sending notification", "title", title, "muted", muted)

	if muted {
		// Visual only — no sound
		if err := beeep.Notify(title, message, ""); err != nil {
			slog.Warn("notification failed", "error", err)
			return fmt.Errorf("send notification: %w", err)
		}
		return nil
	}

	// Visual + audible beep
	if err := beeep.Alert(title, message, ""); err != nil {
		slog.Warn("notification failed", "error", err)
		return fmt.Errorf("send notification: %w", err)
	}
	return nil
}

// Notify is kept for backwards compatibility (no mute check, uses Alert for sound).
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
