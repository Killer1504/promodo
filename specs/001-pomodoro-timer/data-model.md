# Data Model: Pomodoro Focus Timer

**Branch**: `001-pomodoro-timer` | **Date**: 2026-03-10

## Entities

### FocusSession

Represents a single completed timer period (focus or break).

| Field | Type | Constraints | Notes |
|-------|------|-------------|-------|
| `id` | INTEGER | PRIMARY KEY, AUTOINCREMENT | Unique session ID |
| `session_type` | TEXT | NOT NULL, CHECK IN ('focus', 'short_break', 'long_break') | Type of timer period |
| `start_time` | TEXT (ISO 8601) | NOT NULL | When the session started |
| `end_time` | TEXT (ISO 8601) | NOT NULL | When the session ended |
| `duration_seconds` | INTEGER | NOT NULL, > 0 | Planned duration in seconds |
| `completed` | INTEGER | NOT NULL, DEFAULT 1 | 1 = completed, 0 = abandoned |
| `created_at` | TEXT (ISO 8601) | NOT NULL, DEFAULT CURRENT_TIMESTAMP | Record creation timestamp |

**Indexes**:
- `idx_session_start_time` on `start_time` — for date-range queries in Stats view
- `idx_session_type` on `session_type` — for filtering by focus/break

**Lifecycle**:
```
IDLE → RUNNING → PAUSED → RUNNING → COMPLETED
                       └──────────→ ABANDONED (if reset or app closed while running)
```

Note: State transitions happen in-memory in Go. Only completed/abandoned sessions are written to the database.

**Retention**: Records older than 30 days are auto-deleted on app launch (FR-013).

---

### UserSettings

Stores user preferences. Single-row table (settings are global).

| Field | Type | Constraints | Notes |
|-------|------|-------------|-------|
| `id` | INTEGER | PRIMARY KEY, DEFAULT 1 | Always 1 (single row) |
| `focus_duration` | INTEGER | NOT NULL, DEFAULT 1500 | Focus duration in seconds (25 min) |
| `short_break_duration` | INTEGER | NOT NULL, DEFAULT 300 | Short break in seconds (5 min) |
| `long_break_duration` | INTEGER | NOT NULL, DEFAULT 900 | Long break in seconds (15 min) |
| `sessions_before_long_break` | INTEGER | NOT NULL, DEFAULT 4 | Number of focus sessions before long break |
| `notification_sound` | TEXT | NOT NULL, DEFAULT 'bell' | Selected sound file name |
| `mute` | INTEGER | NOT NULL, DEFAULT 0 | 1 = muted, 0 = sound on |
| `theme` | TEXT | NOT NULL, DEFAULT 'system' | 'light', 'dark', or 'system' |

**Initialization**: If the settings row doesn't exist on app start, insert with defaults.

---

### DailyStats (View / Query)

Not a stored entity — computed via SQL aggregation from `FocusSession`.

```sql
SELECT
  DATE(start_time) AS date,
  COUNT(*) AS total_sessions,
  SUM(duration_seconds) / 60 AS total_focus_minutes
FROM focus_sessions
WHERE session_type = 'focus'
  AND completed = 1
  AND start_time >= DATE('now', '-30 days')
GROUP BY DATE(start_time)
ORDER BY date DESC;
```

---

### PausedSession (Transient Persistence)

Stored as a JSON file (`paused_session.json`) to handle FR-015 (resume after app restart).

| Field | Type | Notes |
|-------|------|-------|
| `session_type` | string | 'focus', 'short_break', or 'long_break' |
| `remaining_seconds` | int | Seconds left when paused |
| `cycle_position` | int | Which session in the 4-session cycle (0-3) |
| `paused_at` | string (ISO 8601) | When the session was paused |

**Lifecycle**: Created when app shuts down with a paused session. Deleted after user chooses to resume or discard on reopen.

## SQL Schema

```sql
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

-- Ensure settings row exists
INSERT OR IGNORE INTO user_settings (id) VALUES (1);
```
