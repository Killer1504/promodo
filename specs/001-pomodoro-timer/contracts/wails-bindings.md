# Wails Binding Contracts: Pomodoro Focus Timer

**Branch**: `001-pomodoro-timer` | **Date**: 2026-03-10

This document defines the Go methods exposed to the frontend via Wails bindings. These form the typed contract bridge (Constitution Principle II).

## TimerService

Controls the Pomodoro timer lifecycle. All timer logic runs in Go.

### Methods

| Method | Input | Output | Description |
|--------|-------|--------|-------------|
| `StartFocus()` | — | `TimerState` | Start a new focus session |
| `StartBreak(breakType string)` | `"short"` or `"long"` | `TimerState` | Start a break session |
| `Pause()` | — | `TimerState` | Pause the active timer |
| `Resume()` | — | `TimerState` | Resume a paused timer |
| `Reset()` | — | `TimerState` | Reset the timer to idle |
| `GetState()` | — | `TimerState` | Get current timer state |
| `SkipBreakCountdown()` | — | `TimerState` | Skip the 3-second pre-break countdown |

### Events (Go → Frontend via `runtime.EventsEmit`)

| Event | Payload | Frequency | Description |
|-------|---------|-----------|-------------|
| `timer:tick` | `TimerState` | Every 1s | Timer countdown update |
| `timer:complete` | `SessionComplete` | On completion | Session ended, triggers notification |
| `timer:break-countdown` | `BreakCountdown` | Every 1s (3s total) | Pre-break countdown (3...2...1) |

### Types

```go
type TimerState struct {
    Status          string `json:"status"`           // "idle", "running", "paused", "break_countdown"
    SessionType     string `json:"sessionType"`      // "focus", "short_break", "long_break"
    RemainingSeconds int   `json:"remainingSeconds"`
    TotalSeconds    int    `json:"totalSeconds"`
    CyclePosition   int   `json:"cyclePosition"`    // 0-3 (which dot is active)
}

type SessionComplete struct {
    SessionType     string `json:"sessionType"`
    DurationSeconds int    `json:"durationSeconds"`
    NextBreakType   string `json:"nextBreakType"`    // "short" or "long"
}

type BreakCountdown struct {
    SecondsRemaining int `json:"secondsRemaining"` // 3, 2, 1
}
```

---

## StatsService

Provides session statistics for the Stats view.

### Methods

| Method | Input | Output | Description |
|--------|-------|--------|-------------|
| `GetTodayStats()` | — | `DailyStatsResponse` | Today's session count and focus minutes |
| `GetWeeklyStats()` | — | `[]DailyStatsResponse` | Last 7 days of daily stats |

### Types

```go
type DailyStatsResponse struct {
    Date             string `json:"date"`             // "2026-03-10"
    TotalSessions    int    `json:"totalSessions"`
    TotalFocusMinutes int   `json:"totalFocusMinutes"`
    IsToday          bool   `json:"isToday"`
}
```

---

## SettingsService

Manages user preferences.

### Methods

| Method | Input | Output | Description |
|--------|-------|--------|-------------|
| `GetSettings()` | — | `UserSettingsResponse` | Get current settings |
| `UpdateSettings(s UserSettingsRequest)` | `UserSettingsRequest` | `UserSettingsResponse` | Save updated settings |
| `ResetToDefaults()` | — | `UserSettingsResponse` | Reset all settings to defaults |
| `GetAvailableSounds()` | — | `[]SoundOption` | List available notification sounds |

### Types

```go
type UserSettingsResponse struct {
    FocusDuration          int    `json:"focusDuration"`           // seconds
    ShortBreakDuration     int    `json:"shortBreakDuration"`      // seconds
    LongBreakDuration      int    `json:"longBreakDuration"`       // seconds
    SessionsBeforeLongBreak int   `json:"sessionsBeforeLongBreak"`
    NotificationSound      string `json:"notificationSound"`
    Mute                   bool   `json:"mute"`
    Theme                  string `json:"theme"`                   // "light", "dark", "system"
}

type UserSettingsRequest struct {
    FocusDuration          *int    `json:"focusDuration,omitempty"`
    ShortBreakDuration     *int    `json:"shortBreakDuration,omitempty"`
    LongBreakDuration      *int    `json:"longBreakDuration,omitempty"`
    SessionsBeforeLongBreak *int   `json:"sessionsBeforeLongBreak,omitempty"`
    NotificationSound      *string `json:"notificationSound,omitempty"`
    Mute                   *bool   `json:"mute,omitempty"`
    Theme                  *string `json:"theme,omitempty"`
}

type SoundOption struct {
    Name     string `json:"name"`      // "bell", "chime", "ding"
    Filename string `json:"filename"`  // "bell.mp3"
}
```

---

## AppService

Application-level operations.

### Methods

| Method | Input | Output | Description |
|--------|-------|--------|-------------|
| `GetPausedSession()` | — | `*PausedSessionResponse` | Check for a paused session from last run (nil if none) |
| `ResumePausedSession()` | — | `TimerState` | Resume the saved paused session |
| `DiscardPausedSession()` | — | `bool` | Discard the saved paused session |
| `GetTheme()` | — | `string` | Get resolved theme ("light" or "dark") |

### Types

```go
type PausedSessionResponse struct {
    SessionType      string `json:"sessionType"`
    RemainingSeconds int    `json:"remainingSeconds"`
    CyclePosition    int    `json:"cyclePosition"`
    PausedAt         string `json:"pausedAt"`
}
```
