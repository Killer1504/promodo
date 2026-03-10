package notification

// SessionCompleteMessage returns the notification title and body for a completed session.
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
