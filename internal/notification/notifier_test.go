package notification

import (
	"testing"
)

func TestSessionCompleteMessage_Focus(t *testing.T) {
	title, msg := SessionCompleteMessage("focus")
	if title != "Focus Session Complete!" {
		t.Errorf("title = %q, want %q", title, "Focus Session Complete!")
	}
	if msg != "Great work! Time for a break." {
		t.Errorf("message = %q", msg)
	}
}

func TestSessionCompleteMessage_ShortBreak(t *testing.T) {
	title, msg := SessionCompleteMessage("short_break")
	if title != "Break Over!" {
		t.Errorf("title = %q, want %q", title, "Break Over!")
	}
	if msg != "Ready to focus again?" {
		t.Errorf("message = %q", msg)
	}
}

func TestSessionCompleteMessage_LongBreak(t *testing.T) {
	title, msg := SessionCompleteMessage("long_break")
	if title != "Long Break Over!" {
		t.Errorf("title = %q, want %q", title, "Long Break Over!")
	}
	if msg != "Refreshed and ready to go!" {
		t.Errorf("message = %q", msg)
	}
}

func TestSessionCompleteMessage_Unknown(t *testing.T) {
	title, msg := SessionCompleteMessage("unknown")
	if title != "Timer Complete" {
		t.Errorf("title = %q, want %q", title, "Timer Complete")
	}
	if msg != "Session finished." {
		t.Errorf("message = %q", msg)
	}
}

func TestNewNotifier(t *testing.T) {
	n := NewNotifier()
	if n == nil {
		t.Fatal("NewNotifier() returned nil")
	}
}
