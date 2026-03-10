package timer

import (
	"testing"
	"time"

	"pomodoro-timer/internal/notification"
	"pomodoro-timer/internal/storage"
)

func setupTestService(t *testing.T) *Service {
	t.Helper()
	dir := t.TempDir()
	db, err := storage.NewDatabase(dir)
	if err != nil {
		t.Fatalf("create DB: %v", err)
	}
	t.Cleanup(func() { db.Close() })
	repo := storage.NewRepository(db)
	notifier := notification.NewNotifier()
	svc := NewService(repo, notifier)
	return svc
}

func TestNewService_DefaultsToIdle(t *testing.T) {
	svc := setupTestService(t)
	state := svc.GetState()

	if state.Status != StatusIdle {
		t.Errorf("status = %q, want %q", state.Status, StatusIdle)
	}
	if state.SessionType != SessionFocus {
		t.Errorf("sessionType = %q, want %q", state.SessionType, SessionFocus)
	}
}

func TestStartFocus(t *testing.T) {
	svc := setupTestService(t)
	state := svc.StartFocus()
	defer svc.Shutdown()

	if state.Status != StatusRunning {
		t.Errorf("status = %q, want %q", state.Status, StatusRunning)
	}
	if state.SessionType != SessionFocus {
		t.Errorf("sessionType = %q, want %q", state.SessionType, SessionFocus)
	}
	if state.TotalSeconds != 1500 {
		t.Errorf("totalSeconds = %d, want 1500", state.TotalSeconds)
	}
	if state.RemainingSeconds != 1500 {
		t.Errorf("remainingSeconds = %d, want 1500", state.RemainingSeconds)
	}
}

func TestPauseAndResume(t *testing.T) {
	svc := setupTestService(t)
	svc.StartFocus()
	defer svc.Shutdown()

	paused := svc.Pause()
	if paused.Status != StatusPaused {
		t.Errorf("after pause: status = %q, want %q", paused.Status, StatusPaused)
	}

	resumed := svc.Resume()
	if resumed.Status != StatusRunning {
		t.Errorf("after resume: status = %q, want %q", resumed.Status, StatusRunning)
	}
}

func TestReset(t *testing.T) {
	svc := setupTestService(t)
	svc.StartFocus()
	defer svc.Shutdown()

	state := svc.Reset()
	if state.Status != StatusIdle {
		t.Errorf("after reset: status = %q, want %q", state.Status, StatusIdle)
	}
}

func TestPauseOnlyWhenRunning(t *testing.T) {
	svc := setupTestService(t)
	defer svc.Shutdown()

	// Pausing when idle should do nothing
	state := svc.Pause()
	if state.Status != StatusIdle {
		t.Errorf("pause when idle: status = %q, want %q", state.Status, StatusIdle)
	}
}

func TestResumeOnlyWhenPaused(t *testing.T) {
	svc := setupTestService(t)
	defer svc.Shutdown()

	// Resuming when idle should do nothing
	state := svc.Resume()
	if state.Status != StatusIdle {
		t.Errorf("resume when idle: status = %q, want %q", state.Status, StatusIdle)
	}
}

func TestStartBreak(t *testing.T) {
	svc := setupTestService(t)
	defer svc.Shutdown()

	short := svc.StartBreak("short")
	if short.SessionType != SessionShortBreak {
		t.Errorf("short break: sessionType = %q, want %q", short.SessionType, SessionShortBreak)
	}
	if short.TotalSeconds != 300 {
		t.Errorf("short break: totalSeconds = %d, want 300", short.TotalSeconds)
	}

	long := svc.StartBreak("long")
	if long.SessionType != SessionLongBreak {
		t.Errorf("long break: sessionType = %q, want %q", long.SessionType, SessionLongBreak)
	}
	if long.TotalSeconds != 900 {
		t.Errorf("long break: totalSeconds = %d, want 900", long.TotalSeconds)
	}
}

func TestRestoreState(t *testing.T) {
	svc := setupTestService(t)
	defer svc.Shutdown()

	svc.RestoreState(SessionFocus, 600, 2)

	state := svc.GetState()
	if state.Status != StatusPaused {
		t.Errorf("restored status = %q, want %q", state.Status, StatusPaused)
	}
	if state.RemainingSeconds != 600 {
		t.Errorf("restored remaining = %d, want 600", state.RemainingSeconds)
	}
	if state.CyclePosition != 2 {
		t.Errorf("restored cycle = %d, want 2", state.CyclePosition)
	}
}

func TestFormatRemaining(t *testing.T) {
	tests := []struct {
		seconds int
		want    string
	}{
		{0, "00:00"},
		{59, "00:59"},
		{60, "01:00"},
		{1500, "25:00"},
		{3661, "61:01"},
	}
	for _, tt := range tests {
		got := FormatRemaining(tt.seconds)
		if got != tt.want {
			t.Errorf("FormatRemaining(%d) = %q, want %q", tt.seconds, got, tt.want)
		}
	}
}

func TestTickDecrementsRemaining(t *testing.T) {
	svc := setupTestService(t)
	state := svc.StartFocus()
	initialRemaining := state.RemainingSeconds
	defer svc.Shutdown()

	// Wait for 1-2 ticks
	time.Sleep(1500 * time.Millisecond)

	current := svc.GetState()
	if current.RemainingSeconds >= initialRemaining {
		t.Errorf("remaining should have decreased: was %d, now %d", initialRemaining, current.RemainingSeconds)
	}
}
