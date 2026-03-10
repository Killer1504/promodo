package settings

import (
	"testing"

	"pomodoro-timer/internal/storage"
)

func setupTestSettings(t *testing.T) *Service {
	t.Helper()
	dir := t.TempDir()
	db, err := storage.NewDatabase(dir)
	if err != nil {
		t.Fatalf("create DB: %v", err)
	}
	t.Cleanup(func() { db.Close() })
	repo := storage.NewRepository(db)
	return NewService(repo)
}

func TestGetSettings_Defaults(t *testing.T) {
	svc := setupTestSettings(t)
	s, err := svc.GetSettings()
	if err != nil {
		t.Fatalf("get settings: %v", err)
	}

	if s.FocusDuration != 1500 {
		t.Errorf("focus = %d, want 1500", s.FocusDuration)
	}
	if s.ShortBreakDuration != 300 {
		t.Errorf("short break = %d, want 300", s.ShortBreakDuration)
	}
	if s.LongBreakDuration != 900 {
		t.Errorf("long break = %d, want 900", s.LongBreakDuration)
	}
	if s.Theme != "system" {
		t.Errorf("theme = %q, want system", s.Theme)
	}
}

func TestUpdateSettings_Valid(t *testing.T) {
	svc := setupTestSettings(t)

	err := svc.UpdateSettings(storage.UserSettings{
		FocusDuration:           3000,
		ShortBreakDuration:      600,
		LongBreakDuration:       1800,
		SessionsBeforeLongBreak: 4,
		NotificationSound:       "bell",
		Mute:                    false,
		Theme:                   "dark",
	})
	if err != nil {
		t.Fatalf("update: %v", err)
	}

	s, _ := svc.GetSettings()
	if s.FocusDuration != 3000 {
		t.Errorf("focus = %d, want 3000", s.FocusDuration)
	}
	if s.Theme != "dark" {
		t.Errorf("theme = %q, want dark", s.Theme)
	}
}

func TestUpdateSettings_InvalidFocus(t *testing.T) {
	svc := setupTestSettings(t)

	tests := []struct {
		name  string
		focus int
	}{
		{"zero", 0},
		{"negative", -5},
		{"too_high", 7201},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := svc.UpdateSettings(storage.UserSettings{
				FocusDuration:           tt.focus,
				ShortBreakDuration:      300,
				LongBreakDuration:       900,
				SessionsBeforeLongBreak: 4,
				Theme:                   "system",
			})
			if err == nil {
				t.Errorf("expected error for focus=%d", tt.focus)
			}
		})
	}
}

func TestUpdateSettings_InvalidTheme(t *testing.T) {
	svc := setupTestSettings(t)

	err := svc.UpdateSettings(storage.UserSettings{
		FocusDuration:           1500,
		ShortBreakDuration:      300,
		LongBreakDuration:       900,
		SessionsBeforeLongBreak: 4,
		Theme:                   "purple",
	})
	if err == nil {
		t.Error("expected error for invalid theme")
	}
}

func TestResetToDefaults(t *testing.T) {
	svc := setupTestSettings(t)

	// Change settings
	_ = svc.UpdateSettings(storage.UserSettings{
		FocusDuration:           3000,
		ShortBreakDuration:      600,
		LongBreakDuration:       1800,
		SessionsBeforeLongBreak: 6,
		Theme:                   "dark",
	})

	// Reset
	s, err := svc.ResetToDefaults()
	if err != nil {
		t.Fatalf("reset: %v", err)
	}
	if s.FocusDuration != 1500 {
		t.Errorf("focus = %d, want 1500", s.FocusDuration)
	}
	if s.Theme != "system" {
		t.Errorf("theme = %q, want system", s.Theme)
	}
}

func TestPersistence(t *testing.T) {
	dir := t.TempDir()

	// First service instance: update
	db1, _ := storage.NewDatabase(dir)
	repo1 := storage.NewRepository(db1)
	svc1 := NewService(repo1)
	_ = svc1.UpdateSettings(storage.UserSettings{
		FocusDuration:           2700,
		ShortBreakDuration:      600,
		LongBreakDuration:       1200,
		SessionsBeforeLongBreak: 3,
		Theme:                   "light",
	})
	db1.Close()

	// Second service instance: read back
	db2, _ := storage.NewDatabase(dir)
	t.Cleanup(func() { db2.Close() })
	repo2 := storage.NewRepository(db2)
	svc2 := NewService(repo2)
	s, _ := svc2.GetSettings()

	if s.FocusDuration != 2700 {
		t.Errorf("persisted focus = %d, want 2700", s.FocusDuration)
	}
	if s.Theme != "light" {
		t.Errorf("persisted theme = %q, want light", s.Theme)
	}
}
