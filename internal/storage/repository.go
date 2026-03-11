package storage

import (
	"database/sql"
	"fmt"
	"time"
)

// FocusSession represents a completed timer period stored in the database.
type FocusSession struct {
	ID              int64
	SessionType     string
	StartTime       time.Time
	EndTime         time.Time
	DurationSeconds int
	Completed       bool
	CreatedAt       time.Time
}

// UserSettings represents the user's persisted preferences.
type UserSettings struct {
	FocusDuration           int    `json:"focusDuration"`
	ShortBreakDuration      int    `json:"shortBreakDuration"`
	LongBreakDuration       int    `json:"longBreakDuration"`
	SessionsBeforeLongBreak int    `json:"sessionsBeforeLongBreak"`
	NotificationSound       string `json:"notificationSound"`
	Mute                    bool   `json:"mute"`
	Theme                   string `json:"theme"`
}

// DailyStats holds aggregated focus data for a single day.
type DailyStats struct {
	Date              string `json:"date"`
	TotalSessions     int    `json:"totalSessions"`
	TotalFocusMinutes int    `json:"totalFocusMinutes"`
	IsToday           bool   `json:"isToday"`
}

// Repository provides data access methods for sessions and settings.
type Repository struct {
	db *sql.DB
}

// NewRepository creates a new Repository using the given Database.
func NewRepository(database *Database) *Repository {
	return &Repository{db: database.DB()}
}

// InsertSession stores a completed session.
func (r *Repository) InsertSession(s FocusSession) (int64, error) {
	result, err := r.db.Exec(
		`INSERT INTO focus_sessions (session_type, start_time, end_time, duration_seconds, completed)
		 VALUES (?, ?, ?, ?, ?)`,
		s.SessionType,
		s.StartTime.Format(time.RFC3339),
		s.EndTime.Format(time.RFC3339),
		s.DurationSeconds,
		boolToInt(s.Completed),
	)
	if err != nil {
		return 0, fmt.Errorf("insert session: %w", err)
	}
	return result.LastInsertId()
}

// GetSessionsByDateRange returns sessions within the given time range.
func (r *Repository) GetSessionsByDateRange(from, to time.Time) ([]FocusSession, error) {
	rows, err := r.db.Query(
		`SELECT id, session_type, start_time, end_time, duration_seconds, completed, created_at
		 FROM focus_sessions
		 WHERE start_time >= ? AND start_time < ?
		 ORDER BY start_time DESC`,
		from.Format(time.RFC3339),
		to.Format(time.RFC3339),
	)
	if err != nil {
		return nil, fmt.Errorf("query sessions: %w", err)
	}
	defer rows.Close()

	var sessions []FocusSession
	for rows.Next() {
		var s FocusSession
		var startStr, endStr, createdStr string
		var completedInt int
		if err := rows.Scan(&s.ID, &s.SessionType, &startStr, &endStr, &s.DurationSeconds, &completedInt, &createdStr); err != nil {
			return nil, fmt.Errorf("scan session: %w", err)
		}
		s.StartTime, _ = time.Parse(time.RFC3339, startStr)
		s.EndTime, _ = time.Parse(time.RFC3339, endStr)
		s.CreatedAt, _ = time.Parse(time.RFC3339, createdStr)
		s.Completed = completedInt == 1
		sessions = append(sessions, s)
	}
	return sessions, rows.Err()
}

// GetTodayStats returns focus stats for today (local timezone).
func (r *Repository) GetTodayStats() (DailyStats, error) {
	var ds DailyStats
	ds.IsToday = true
	err := r.db.QueryRow(
		`SELECT COALESCE(COUNT(*), 0),
		        COALESCE(SUM(duration_seconds) / 60, 0)
		 FROM focus_sessions
		 WHERE session_type = 'focus'
		   AND completed = 1
		   AND DATE(start_time, 'localtime') = DATE('now', 'localtime')`,
	).Scan(&ds.TotalSessions, &ds.TotalFocusMinutes)
	if err != nil {
		return ds, fmt.Errorf("get today stats: %w", err)
	}
	ds.Date = "today"
	return ds, nil
}

// GetWeeklyStats returns aggregated focus stats for the last N days (local timezone).
func (r *Repository) GetWeeklyStats(days int) ([]DailyStats, error) {
	today := "DATE('now', 'localtime')"
	rows, err := r.db.Query(
		`SELECT DATE(start_time, 'localtime') AS date,
		        COUNT(*) AS total_sessions,
		        COALESCE(SUM(duration_seconds) / 60, 0) AS total_focus_minutes
		 FROM focus_sessions
		 WHERE session_type = 'focus'
		   AND completed = 1
		   AND DATE(start_time, 'localtime') >= DATE('now', 'localtime', ?)
		 GROUP BY DATE(start_time, 'localtime')
		 ORDER BY date ASC`,
		fmt.Sprintf("-%d days", days-1),
	)
	_ = today // used conceptually above
	if err != nil {
		return nil, fmt.Errorf("query weekly stats: %w", err)
	}
	defer rows.Close()

	var stats []DailyStats
	for rows.Next() {
		var ds DailyStats
		if err := rows.Scan(&ds.Date, &ds.TotalSessions, &ds.TotalFocusMinutes); err != nil {
			return nil, fmt.Errorf("scan daily stats: %w", err)
		}
		stats = append(stats, ds)
	}
	return stats, rows.Err()
}

// GetSettings returns the current user settings.
func (r *Repository) GetSettings() (UserSettings, error) {
	var s UserSettings
	var muteInt int
	err := r.db.QueryRow(
		`SELECT focus_duration, short_break_duration, long_break_duration,
		        sessions_before_long_break, notification_sound, mute, theme
		 FROM user_settings WHERE id = 1`,
	).Scan(
		&s.FocusDuration, &s.ShortBreakDuration, &s.LongBreakDuration,
		&s.SessionsBeforeLongBreak, &s.NotificationSound, &muteInt, &s.Theme,
	)
	if err != nil {
		return s, fmt.Errorf("get settings: %w", err)
	}
	s.Mute = muteInt == 1
	return s, nil
}

// UpdateSettings updates user settings. Only non-nil fields are changed.
func (r *Repository) UpdateSettings(s UserSettings) error {
	_, err := r.db.Exec(
		`UPDATE user_settings SET
			focus_duration = ?,
			short_break_duration = ?,
			long_break_duration = ?,
			sessions_before_long_break = ?,
			notification_sound = ?,
			mute = ?,
			theme = ?
		 WHERE id = 1`,
		s.FocusDuration, s.ShortBreakDuration, s.LongBreakDuration,
		s.SessionsBeforeLongBreak, s.NotificationSound, boolToInt(s.Mute), s.Theme,
	)
	if err != nil {
		return fmt.Errorf("update settings: %w", err)
	}
	return nil
}

// ResetSettings restores all settings to their default values.
func (r *Repository) ResetSettings() error {
	_, err := r.db.Exec(
		`UPDATE user_settings SET
			focus_duration = 1500,
			short_break_duration = 300,
			long_break_duration = 900,
			sessions_before_long_break = 4,
			notification_sound = 'bell',
			mute = 0,
			theme = 'system'
		 WHERE id = 1`,
	)
	if err != nil {
		return fmt.Errorf("reset settings: %w", err)
	}
	return nil
}

func boolToInt(b bool) int {
	if b {
		return 1
	}
	return 0
}
