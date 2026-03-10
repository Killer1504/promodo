package stats

import (
	"testing"
	"time"

	"pomodoro-timer/internal/storage"
)

func setupTestStats(t *testing.T) (*Service, *storage.Repository) {
	t.Helper()
	dir := t.TempDir()
	db, err := storage.NewDatabase(dir)
	if err != nil {
		t.Fatalf("create DB: %v", err)
	}
	t.Cleanup(func() { db.Close() })
	repo := storage.NewRepository(db)
	svc := NewService(repo)
	return svc, repo
}

func TestGetTodayStats_Empty(t *testing.T) {
	svc, _ := setupTestStats(t)
	stats := svc.GetTodayStats()

	if stats.TotalSessions != 0 {
		t.Errorf("expected 0 sessions, got %d", stats.TotalSessions)
	}
	if stats.TotalFocusMinutes != 0 {
		t.Errorf("expected 0 minutes, got %d", stats.TotalFocusMinutes)
	}
	if !stats.IsToday {
		t.Error("expected IsToday = true")
	}
}

func TestGetTodayStats_WithSessions(t *testing.T) {
	svc, repo := setupTestStats(t)

	now := time.Now()
	for i := 0; i < 3; i++ {
		_, err := repo.InsertSession(storage.FocusSession{
			SessionType:     "focus",
			StartTime:       now.Add(-time.Duration(i) * time.Hour),
			EndTime:         now.Add(-time.Duration(i)*time.Hour + 25*time.Minute),
			DurationSeconds: 1500,
			Completed:       true,
		})
		if err != nil {
			t.Fatalf("insert: %v", err)
		}
	}

	stats := svc.GetTodayStats()
	if stats.TotalSessions != 3 {
		t.Errorf("sessions = %d, want 3", stats.TotalSessions)
	}
	// 3 * 1500s = 4500s = 75min
	if stats.TotalFocusMinutes != 75 {
		t.Errorf("minutes = %d, want 75", stats.TotalFocusMinutes)
	}
}

func TestGetTodayStats_BreaksExcluded(t *testing.T) {
	svc, repo := setupTestStats(t)

	now := time.Now()
	// Insert one focus + one break
	_, _ = repo.InsertSession(storage.FocusSession{
		SessionType:     "focus",
		StartTime:       now.Add(-30 * time.Minute),
		EndTime:         now.Add(-5 * time.Minute),
		DurationSeconds: 1500,
		Completed:       true,
	})
	_, _ = repo.InsertSession(storage.FocusSession{
		SessionType:     "short_break",
		StartTime:       now.Add(-5 * time.Minute),
		EndTime:         now,
		DurationSeconds: 300,
		Completed:       true,
	})

	stats := svc.GetTodayStats()
	if stats.TotalSessions != 1 {
		t.Errorf("sessions = %d, want 1 (breaks should be excluded)", stats.TotalSessions)
	}
}

func TestGetWeeklyStats_Returns7Days(t *testing.T) {
	svc, _ := setupTestStats(t)
	weekly := svc.GetWeeklyStats()

	if len(weekly) != 7 {
		t.Fatalf("expected 7 days, got %d", len(weekly))
	}

	// Last entry should be today
	todayStr := time.Now().Format("2006-01-02")
	if weekly[6].Date != todayStr {
		t.Errorf("last day = %q, want %q", weekly[6].Date, todayStr)
	}
	if !weekly[6].IsToday {
		t.Error("last day should be IsToday")
	}
}

func TestGetWeeklyStats_GapFilling(t *testing.T) {
	svc, repo := setupTestStats(t)

	// Insert session 3 days ago
	threeDaysAgo := time.Now().AddDate(0, 0, -3)
	_, _ = repo.InsertSession(storage.FocusSession{
		SessionType:     "focus",
		StartTime:       threeDaysAgo,
		EndTime:         threeDaysAgo.Add(25 * time.Minute),
		DurationSeconds: 1500,
		Completed:       true,
	})

	weekly := svc.GetWeeklyStats()

	// Find the day 3 days ago
	targetDate := threeDaysAgo.Format("2006-01-02")
	found := false
	for _, d := range weekly {
		if d.Date == targetDate {
			found = true
			if d.TotalSessions != 1 {
				t.Errorf("3 days ago sessions = %d, want 1", d.TotalSessions)
			}
		}
	}
	if !found {
		t.Error("3-day-ago date not found in weekly stats")
	}

	// Today should be 0
	todayStats := weekly[6]
	if todayStats.TotalSessions != 0 {
		t.Errorf("today sessions = %d, want 0", todayStats.TotalSessions)
	}
}
