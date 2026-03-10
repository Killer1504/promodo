//go:build windows

package notification

import (
	"fmt"
	"log/slog"

	toast "git.sr.ht/~jackmordaunt/go-toast"
)

const appName = "Pomodoro Focus Timer"

// Notifier sends OS-level desktop notifications.
type Notifier struct{}

// NewNotifier creates a new Notifier.
func NewNotifier() *Notifier {
	return &Notifier{}
}

// Alert sends a visual Windows toast notification with the correct app name.
// muted is honoured (silent notification when true).
func (n *Notifier) Alert(title, message string, muted bool) error {
	slog.Info("sending notification", "title", title, "muted", muted)

	audio := toast.Default
	if muted {
		audio = toast.Silent
	}

	notification := toast.Notification{
		AppID:    appName,
		Title:    title,
		Body:     message,
		Audio:    audio,
		Duration: toast.Short,
	}

	if err := notification.Push(); err != nil {
		slog.Warn("notification failed", "error", err)
		return fmt.Errorf("send notification: %w", err)
	}
	return nil
}

// Notify sends a notification with sound (mute=false).
func (n *Notifier) Notify(title, message string) error {
	return n.Alert(title, message, false)
}
