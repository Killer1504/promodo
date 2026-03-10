# Technical Documentation — Pomodoro Focus Timer

## Table of Contents

- [Architecture Overview](#architecture-overview)
- [Backend (Go)](#backend-go)
- [Frontend (JavaScript)](#frontend-javascript)
- [Data Model](#data-model)
- [API Reference](#api-reference)
- [Design System](#design-system)
- [Testing Strategy](#testing-strategy)
- [Troubleshooting](#troubleshooting)

---

## Architecture Overview

```
┌────────────────────────────────────────────────────┐
│                   Wails v2 Runtime                  │
│                                                     │
│  ┌──────────────┐          ┌─────────────────────┐ │
│  │   Frontend    │ ◄──────► │     Go Backend      │ │
│  │  (WebView2)   │  Wails   │                     │ │
│  │              │  Bindings │  ┌───────────────┐  │ │
│  │  HTML/JS/CSS │          │  │ TimerService  │  │ │
│  │  + Chart.js  │          │  │ StatsService  │  │ │
│  │              │          │  │ SettingsService│ │ │
│  │              │          │  │ Notifier      │  │ │
│  └──────────────┘          │  └───────┬───────┘  │ │
│                             │          │          │ │
│                             │  ┌───────▼───────┐  │ │
│                             │  │   Repository  │  │ │
│                             │  │   (SQLite)    │  │ │
│                             │  └───────────────┘  │ │
│                             └─────────────────────┘ │
└────────────────────────────────────────────────────┘
```

**Communication Model**: The frontend calls Go methods directly via Wails bindings (`window.go.main.App.*`). The backend emits events to the frontend using `runtime.EventsEmit()` for real-time updates (timer ticks, session completions).

---

## Backend (Go)

### Package Dependency Graph

```
main.go / app.go
  ├── internal/timer      → timer state machine
  │     └── internal/storage  (insert sessions)
  │     └── internal/notification (send OS alerts)
  ├── internal/settings   → validate & persist settings
  │     └── internal/storage
  ├── internal/stats      → aggregate session data
  │     └── internal/storage
  ├── internal/storage    → SQLite database + repository
  └── internal/notification → beeep wrapper
```

### `internal/timer` — Timer State Machine

The core timer operates as a finite state machine with 4 states:

| State | Description | Transitions To |
|-------|-------------|----------------|
| `idle` | No active session | `running` (Start) |
| `running` | Countdown active | `paused` (Pause), `break_countdown` (Complete) |
| `paused` | Timer frozen | `running` (Resume), `idle` (Reset) |
| `break_countdown` | 3s pre-break | `running` (break starts) |

**Key Implementation Details**:

- Uses `time.NewTicker(1s)` goroutine for countdown
- Thread-safe via `sync.Mutex`
- Goroutine race condition fix: captures `done` and `ticker` channels in local variables before spawning goroutine to prevent nil pointer dereferences after `Shutdown()`
- Emits Wails events: `timer:tick`, `timer:complete`, `timer:break-countdown`

**Session Cycle Logic**:
```
Focus → Short Break → Focus → Short Break → Focus → Short Break → Focus → Long Break
  1         ↕           2         ↕           3         ↕           4         ↕
                                                              (cycle resets)
```

### `internal/storage` — Database Layer

**Technology**: Pure Go SQLite via `modernc.org/sqlite` (no CGO required).

**Schema** (`user_settings` table):

| Column | Type | Default |
|--------|------|---------|
| `focus_duration` | INTEGER | 1500 (25 min) |
| `short_break_duration` | INTEGER | 300 (5 min) |
| `long_break_duration` | INTEGER | 900 (15 min) |
| `sessions_before_long_break` | INTEGER | 4 |
| `notification_sound` | TEXT | "bell" |
| `mute` | BOOLEAN | false |
| `theme` | TEXT | "system" |

**Schema** (`focus_sessions` table):

| Column | Type | Description |
|--------|------|-------------|
| `id` | INTEGER PK | Auto-increment |
| `session_type` | TEXT | "focus" / "short_break" / "long_break" |
| `duration_seconds` | INTEGER | Actual duration completed |
| `completed_at` | DATETIME | UTC timestamp |

**Startup Cleanup**: Automatically removes sessions older than 30 days.

### `internal/settings` — Validation Service

Validates all settings before persistence:

| Field | Valid Range | Unit |
|-------|-------------|------|
| Focus Duration | 60 – 7200 | seconds (1–120 min) |
| Short Break | 60 – 3600 | seconds (1–60 min) |
| Long Break | 60 – 3600 | seconds (1–60 min) |
| Sessions Before Long Break | 1 – 10 | count |
| Theme | light / dark / system | enum |

### `internal/stats` — Aggregation Service

- `GetTodayStats()`: Returns session count and total focus minutes for today
- `GetWeeklyStats()`: Returns 7 days of stats with gap-filling for empty days, marks today with `isToday` flag
- Breaks are excluded from stats (only `focus` session types counted)

### `internal/notification` — OS Notifications

Wraps `github.com/gen2brain/beeep` for cross-platform OS notifications. Provides pre-formatted messages per session type:

| Session Type | Title | Message |
|--------------|-------|---------|
| focus | "Focus Session Complete!" | "Great work! Time for a break." |
| short_break | "Break Over!" | "Ready to focus again?" |
| long_break | "Long Break Over!" | "Refreshed and ready to go!" |

---

## Frontend (JavaScript)

### Router Architecture

Hash-based single-page routing in `main.js`:

```
#timer    →  views/timer.js     (default)
#stats    →  views/stats.js
#settings →  views/settings.js
```

Views are lazy-loaded modules. Each exports `mount(container)` and `unmount()`.

### Component Hierarchy

```
index.html
  ├── TabBar (tab-bar.js)       # 3 tabs with ARIA roles
  ├── TimerView
  │     ├── ProgressRing        # SVG circle countdown
  │     └── SessionDots         # Cycle position indicators
  ├── StatsView
  │     ├── SummaryCards         # Sessions + minutes
  │     └── WeeklyChart         # Chart.js bar chart
  └── SettingsView
        ├── Number Steppers     # Duration controls
        ├── Theme Selector      # Light / Dark / System pills
        ├── Mute Toggle         # Switch control
        └── Reset Button        # Restore defaults
```

### Wails Event Handling

| Event | Emitted By | Handled In | Payload |
|-------|------------|------------|---------|
| `timer:tick` | TimerService goroutine | timer.js | `{remaining, formatted, sessionType}` |
| `timer:complete` | TimerService | timer.js | `{sessionType, next}` |
| `timer:break-countdown` | TimerService | timer.js | `{seconds}` |

### Data Flow: Settings

```
Frontend (minutes) → × 60 → Backend (seconds) → SQLite
SQLite (seconds) → Backend → ÷ 60 → Frontend (minutes)
```

The conversion boundary is in `settings.js` (`loadSettings` and `saveSettings`).

---

## Data Model

### `app.go` — Wails Bindings (API Surface)

All methods on the `App` struct are callable from the frontend:

| Method | Return | Description |
|--------|--------|-------------|
| `StartFocus()` | `TimerState` | Start a focus session |
| `PauseTimer()` | `TimerState` | Pause running timer |
| `ResumeTimer()` | `TimerState` | Resume paused timer |
| `ResetTimer()` | `TimerState` | Reset to idle |
| `SkipBreakCountdown()` | `TimerState` | Skip the 3s break countdown |
| `GetTimerState()` | `TimerState` | Get current state snapshot |
| `GetSettings()` | `UserSettings, error` | Read all settings |
| `UpdateSettings(s)` | `UserSettings, error` | Update settings with validation |
| `ResetToDefaults()` | `UserSettings, error` | Reset settings to defaults |
| `GetTodayStats()` | `DailyStats, error` | Today's sessions + minutes |
| `GetWeeklyStats()` | `[]DailyStats, error` | 7-day history |
| `GetTheme()` | `string` | Resolved theme name |
| `GetPausedSession()` | `*PausedSessionResponse` | Check for saved state |
| `ResumePausedSession()` | `TimerState` | Restore & resume |
| `DiscardPausedSession()` | `bool` | Delete saved state |

### `TimerState` Struct

```go
type TimerState struct {
    Status        string // "idle" | "running" | "paused" | "break_countdown"
    SessionType   string // "focus" | "short_break" | "long_break"
    Remaining     int    // seconds remaining
    Formatted     string // "25:00" display format
    CyclePosition int    // 0-based position in cycle
    MaxCycles     int    // sessions before long break
}
```

---

## Design System

### CSS Architecture

```
styles/
├── variables.css     # Design tokens (colors, spacing, fonts)
├── base.css          # Reset, typography, buttons, cards
├── tab-bar.css       # Tab bar component
├── timer.css         # Progress ring, session dots, controls
├── stats.css         # Chart container, summary cards
└── settings.css      # Steppers, theme pills, toggle switch
```

### Theme Tokens

Theme switching uses CSS custom properties on `[data-theme]`:

- **Light**: `#FAFBFC` background, `#1A1D23` text, `#3B82F6` accent
- **Dark**: `#0F1117` background, `#F0F2F5` text, `#60A5FA` accent

All color pairs validated for WCAG AA contrast (≥ 4.5:1 for text).

### Typography

Font: **Inter** (Google Fonts), fallback to system sans-serif.

| Token | Size |
|-------|------|
| `--font-size-timer` | 3.5rem |
| `--font-size-2xl` | 2rem |
| `--font-size-xl` | 1.5rem |
| `--font-size-base` | 1rem |
| `--font-size-sm` | 0.875rem |
| `--font-size-xs` | 0.75rem |

---

## Testing Strategy

### Test Pyramid

```
35 Total Tests
├── Unit Tests (Go)
│   ├── timer/service_test.go      (10 tests) — state machine transitions
│   ├── storage/database_test.go   ( 9 tests) — schema, CRUD, cleanup
│   ├── settings/service_test.go   ( 6 tests) — validation, persistence
│   ├── stats/service_test.go      ( 5 tests) — aggregation, gap-filling
│   └── notification/notifier_test.go (5 tests) — message formatting
└── Integration: go build ./... verifies all packages compile together
```

### Key Test Patterns

- **In-memory SQLite**: Tests use `t.TempDir()` for isolated DB instances
- **Table-driven tests**: Invalid input validation uses subtests
- **Cross-restart persistence**: Settings tests create two separate DB instances on the same directory
- **Race condition coverage**: Timer tick test uses real `time.Sleep` to verify goroutine behavior

### Running Tests

```bash
# Full suite with verbose output
go test -v -count=1 ./internal/...

# Single package with race detection
go test -race -v ./internal/timer/...
```

---

## Troubleshooting

### `wails` command not found

```bash
# Ensure GOPATH/bin is in PATH
go install github.com/wailsapp/wails/v2/cmd/wails@latest

# Verify
wails version
```

### Build fails on Windows

Ensure you have the WebView2 runtime installed. On Windows 11 it's built-in. On Windows 10:
[Download WebView2](https://developer.microsoft.com/en-us/microsoft-edge/webview2/).

### SQLite errors

The app uses pure Go SQLite (`modernc.org/sqlite`), no CGO required. If you see import errors:

```bash
go mod tidy
```

### Logging

Check `~/.pomodoro-timer/pomodoro.log` for structured log output. All timer events, storage operations, and settings changes are logged with `slog`.
