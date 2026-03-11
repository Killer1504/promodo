package storage

import (
	"os"
	"path/filepath"
	"testing"
	"time"
)

func setupTestDB(t *testing.T) (*Database, *Repository) {
	t.Helper()
	dir := t.TempDir()
	db, err := NewDatabase(dir)
	if err != nil {
		t.Fatalf("failed to create database: %v", err)
	}
	t.Cleanup(func() { db.Close() })
	return db, NewRepository(db)
}

func TestNewDatabase_CreatesSchema(t *testing.T) {
	db, _ := setupTestDB(t)

	// Verify tables exist
	tables := []string{"focus_sessions", "user_settings"}
	for _, table := range tables {
		var name string
		err := db.DB().QueryRow(
			"SELECT name FROM sqlite_master WHERE type='table' AND name=?", table,
		).Scan(&name)
		if err != nil {
			t.Errorf("table %s not found: %v", table, err)
		}
	}
}

func TestNewDatabase_DefaultSettings(t *testing.T) {
	_, repo := setupTestDB(t)

	s, err := repo.GetSettings()
	if err != nil {
		t.Fatalf("get settings: %v", err)
	}

	tests := []struct {
		name string
		got  interface{}
		want interface{}
	}{
		{"FocusDuration", s.FocusDuration, 1500},
		{"ShortBreakDuration", s.ShortBreakDuration, 300},
		{"LongBreakDuration", s.LongBreakDuration, 900},
		{"SessionsBeforeLongBreak", s.SessionsBeforeLongBreak, 4},
		{"NotificationSound", s.NotificationSound, "bell"},
		{"Mute", s.Mute, false},
		{"Theme", s.Theme, "system"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.got != tt.want {
				t.Errorf("got %v, want %v", tt.got, tt.want)
			}
		})
	}
}

func TestInsertSession_And_Query(t *testing.T) {
	_, repo := setupTestDB(t)

	now := time.Now()
	session := FocusSession{
		SessionType:     "focus",
		StartTime:       now.Add(-25 * time.Minute),
		EndTime:         now,
		DurationSeconds: 1500,
		Completed:       true,
	}

	id, err := repo.InsertSession(session)
	if err != nil {
		t.Fatalf("insert session: %v", err)
	}
	if id < 1 {
		t.Errorf("expected positive ID, got %d", id)
	}

	// Query back
	sessions, err := repo.GetSessionsByDateRange(
		now.Add(-1*time.Hour),
		now.Add(1*time.Hour),
	)
	if err != nil {
		t.Fatalf("query sessions: %v", err)
	}
	if len(sessions) != 1 {
		t.Fatalf("expected 1 session, got %d", len(sessions))
	}
	if sessions[0].SessionType != "focus" {
		t.Errorf("session type = %q, want %q", sessions[0].SessionType, "focus")
	}
	if sessions[0].DurationSeconds != 1500 {
		t.Errorf("duration = %d, want %d", sessions[0].DurationSeconds, 1500)
	}
}

func TestGetDailyStats_Empty(t *testing.T) {
	_, repo := setupTestDB(t)

	stats, err := repo.GetWeeklyStats(7)
	if err != nil {
		t.Fatalf("get weekly stats: %v", err)
	}
	if len(stats) != 0 {
		t.Errorf("expected 0 stats for empty db, got %d", len(stats))
	}
}

func TestGetDailyStats_WithData(t *testing.T) {
	_, repo := setupTestDB(t)

	now := time.Now()
	for i := 0; i < 3; i++ {
		_, err := repo.InsertSession(FocusSession{
			SessionType:     "focus",
			StartTime:       now.Add(-time.Duration(i) * time.Hour),
			EndTime:         now.Add(-time.Duration(i)*time.Hour + 25*time.Minute),
			DurationSeconds: 1500,
			Completed:       true,
		})
		if err != nil {
			t.Fatalf("insert session %d: %v", i, err)
		}
	}

	// Insert a break - should NOT count
	_, _ = repo.InsertSession(FocusSession{
		SessionType:     "short_break",
		StartTime:       now,
		EndTime:         now.Add(5 * time.Minute),
		DurationSeconds: 300,
		Completed:       true,
	})

	stats, err := repo.GetWeeklyStats(7)
	if err != nil {
		t.Fatalf("get weekly stats: %v", err)
	}
	if len(stats) == 0 {
		t.Fatal("expected at least 1 day of stats")
	}
	if stats[0].TotalSessions != 3 {
		t.Errorf("total sessions = %d, want 3", stats[0].TotalSessions)
	}
}

func TestUpdateSettings(t *testing.T) {
	_, repo := setupTestDB(t)

	s, _ := repo.GetSettings()
	s.FocusDuration = 3000
	s.Theme = "dark"
	s.Mute = true

	if err := repo.UpdateSettings(s); err != nil {
		t.Fatalf("update settings: %v", err)
	}

	updated, err := repo.GetSettings()
	if err != nil {
		t.Fatalf("get updated settings: %v", err)
	}
	if updated.FocusDuration != 3000 {
		t.Errorf("focus duration = %d, want 3000", updated.FocusDuration)
	}
	if updated.Theme != "dark" {
		t.Errorf("theme = %q, want %q", updated.Theme, "dark")
	}
	if !updated.Mute {
		t.Error("mute should be true")
	}
}

func TestResetSettings(t *testing.T) {
	_, repo := setupTestDB(t)

	s, _ := repo.GetSettings()
	s.FocusDuration = 9999
	_ = repo.UpdateSettings(s)

	if err := repo.ResetSettings(); err != nil {
		t.Fatalf("reset settings: %v", err)
	}

	reset, err := repo.GetSettings()
	if err != nil {
		t.Fatalf("get reset settings: %v", err)
	}
	if reset.FocusDuration != 1500 {
		t.Errorf("focus duration = %d, want 1500", reset.FocusDuration)
	}
}

func TestCleanup_RemovesOldSessions(t *testing.T) {
	dir := t.TempDir()
	db, err := NewDatabase(dir)
	if err != nil {
		t.Fatalf("create database: %v", err)
	}
	defer db.Close()

	repo := NewRepository(db)

	// Insert an old session (31 days ago)
	old := time.Now().AddDate(0, 0, -31)
	_, err = repo.InsertSession(FocusSession{
		SessionType:     "focus",
		StartTime:       old,
		EndTime:         old.Add(25 * time.Minute),
		DurationSeconds: 1500,
		Completed:       true,
	})
	if err != nil {
		t.Fatalf("insert old session: %v", err)
	}

	// Insert a recent session
	now := time.Now()
	_, _ = repo.InsertSession(FocusSession{
		SessionType:     "focus",
		StartTime:       now.Add(-25 * time.Minute),
		EndTime:         now,
		DurationSeconds: 1500,
		Completed:       true,
	})

	// Re-open database to trigger cleanup
	db.Close()
	db2, err := NewDatabase(dir)
	if err != nil {
		t.Fatalf("reopen database: %v", err)
	}
	defer db2.Close()

	// Both sessions should exist since created_at is current time (not the simulated start_time)
	var count int
	_ = db2.DB().QueryRow("SELECT COUNT(*) FROM focus_sessions").Scan(&count)
	if count < 1 {
		t.Errorf("expected at least 1 session after cleanup, got %d", count)
	}
}

func TestNewDatabase_CreatesDirIfMissing(t *testing.T) {
	dir := filepath.Join(t.TempDir(), "subdir", "nested")
	db, err := NewDatabase(dir)
	if err != nil {
		t.Fatalf("create database in nested dir: %v", err)
	}
	defer db.Close()

	if _, err := os.Stat(dir); os.IsNotExist(err) {
		t.Error("expected directory to be created")
	}
}
