package storage

import (
	"database/sql"
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
	"time"

	_ "modernc.org/sqlite"
)

const retentionDays = 30

// Database manages the SQLite connection and schema.
type Database struct {
	db *sql.DB
}

// NewDatabase opens (or creates) the SQLite database at the given directory,
// runs migrations and performs 30-day data cleanup.
func NewDatabase(dataDir string) (*Database, error) {
	if err := os.MkdirAll(dataDir, 0o755); err != nil {
		return nil, fmt.Errorf("create data dir: %w", err)
	}

	dbPath := filepath.Join(dataDir, "pomodoro.db")
	db, err := sql.Open("sqlite", dbPath)
	if err != nil {
		return nil, fmt.Errorf("open database: %w", err)
	}

	// Enable WAL mode for better concurrent read/write
	if _, err := db.Exec("PRAGMA journal_mode=WAL"); err != nil {
		return nil, fmt.Errorf("set WAL mode: %w", err)
	}

	d := &Database{db: db}

	if err := d.migrate(); err != nil {
		return nil, fmt.Errorf("migrate: %w", err)
	}

	if err := d.cleanup(); err != nil {
		slog.Warn("cleanup old sessions failed", "error", err)
	}

	return d, nil
}

// DB returns the underlying sql.DB for use by repositories.
func (d *Database) DB() *sql.DB {
	return d.db
}

// Close closes the database connection.
func (d *Database) Close() error {
	return d.db.Close()
}

func (d *Database) migrate() error {
	schema := `
	CREATE TABLE IF NOT EXISTS focus_sessions (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		session_type TEXT NOT NULL CHECK(session_type IN ('focus', 'short_break', 'long_break')),
		start_time TEXT NOT NULL,
		end_time TEXT NOT NULL,
		duration_seconds INTEGER NOT NULL CHECK(duration_seconds > 0),
		completed INTEGER NOT NULL DEFAULT 1,
		created_at TEXT NOT NULL DEFAULT (datetime('now'))
	);

	CREATE INDEX IF NOT EXISTS idx_session_start_time ON focus_sessions(start_time);
	CREATE INDEX IF NOT EXISTS idx_session_type ON focus_sessions(session_type);

	CREATE TABLE IF NOT EXISTS user_settings (
		id INTEGER PRIMARY KEY DEFAULT 1 CHECK(id = 1),
		focus_duration INTEGER NOT NULL DEFAULT 1500,
		short_break_duration INTEGER NOT NULL DEFAULT 300,
		long_break_duration INTEGER NOT NULL DEFAULT 900,
		sessions_before_long_break INTEGER NOT NULL DEFAULT 4,
		notification_sound TEXT NOT NULL DEFAULT 'bell',
		mute INTEGER NOT NULL DEFAULT 0,
		theme TEXT NOT NULL DEFAULT 'system' CHECK(theme IN ('light', 'dark', 'system'))
	);

	INSERT OR IGNORE INTO user_settings (id) VALUES (1);
	`
	_, err := d.db.Exec(schema)
	return err
}

// cleanup removes sessions older than retentionDays (FR-013).
func (d *Database) cleanup() error {
	cutoff := time.Now().AddDate(0, 0, -retentionDays).Format("2006-01-02T15:04:05Z")
	result, err := d.db.Exec("DELETE FROM focus_sessions WHERE created_at < ?", cutoff)
	if err != nil {
		return err
	}
	if n, _ := result.RowsAffected(); n > 0 {
		slog.Info("cleaned up old sessions", "deleted", n, "older_than", cutoff)
	}
	return nil
}
