//go:build !windows

package notification

import (
	"fmt"
	"log/slog"

	"github.com/gen2brain/beeep"
)

// Notifier sends OS-level desktop notifications.
type Notifier struct{}

// NewNotifier creates a new Notifier.
func NewNotifier() *Notifier {
	return &Notifier{}
}

// Alert sends a visual notification. muted suppresses any system sound.
func (n *Notifier) Alert(title, message string, muted bool) error {
	slog.Info("sending notification", "title", title, "muted", muted)
	if muted {
		if err := beeep.Notify(title, message, ""); err != nil {
			slog.Warn("notification failed", "error", err)
			return fmt.Errorf("send notification: %w", err)
		}
		return nil
	}
	if err := beeep.Alert(title, message, ""); err != nil {
		slog.Warn("notification failed", "error", err)
		return fmt.Errorf("send notification: %w", err)
	}
	return nil
}

// Notify sends a notification with sound (mute=false).
func (n *Notifier) Notify(title, message string) error {
	return n.Alert(title, message, false)
}
