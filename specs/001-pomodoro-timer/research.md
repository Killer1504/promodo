# Research: Pomodoro Focus Timer

**Branch**: `001-pomodoro-timer` | **Date**: 2026-03-10

## R-001: Wails v2 System Tray Support

- **Decision**: Use `github.com/wailsapp/wails/v2/pkg/runtime` for window management and the `systray` package (or Wails v2 built-in tray support if available) for tray icon functionality.
- **Rationale**: Wails v2 exposes `runtime.WindowMinimise`, `runtime.WindowShow`, `runtime.WindowHide` methods. For system tray, Wails v2 doesn't have built-in tray support — we'll use `github.com/getlantern/systray` which is the standard Go system tray library. It integrates cleanly with Wails by running in a separate goroutine.
- **Alternatives considered**:
  - `fyne.io/systray` — fork of getlantern, lighter but less maintained.
  - Custom Windows-only tray via `syscall` — too platform-specific, violates cross-platform goal.

## R-002: OS Notifications

- **Decision**: Use `github.com/gen2brain/beeep` for cross-platform OS notifications.
- **Rationale**: Lightweight, no CGo dependency, supports Windows (toast), macOS (NSUserNotification), and Linux (libnotify). Ideal for a simple notification with title + message.
- **Alternatives considered**:
  - `github.com/go-toast/toast` — Windows-only.
  - Wails runtime events + frontend `Notification` API — requires browser notification permissions, less native feel.

## R-003: Audio Playback for Notification Sounds

- **Decision**: Play notification sounds from the frontend using the Web Audio API (`<audio>` element). Sound files embedded as frontend assets.
- **Rationale**: The Web Audio API is the simplest approach for playing short audio clips. Embedding 3–5 small `.mp3` files (< 100 KB each) in the frontend assets keeps the implementation thin and avoids Go audio library complexity.
- **Alternatives considered**:
  - `github.com/hajimehoshi/oto` (Go audio) — adds CGo complexity and binary size for a simple use case.
  - `github.com/faiface/beep` — similar CGo concerns.

## R-004: Local Storage / Persistence

- **Decision**: Use SQLite via `github.com/mattn/go-sqlite3` (CGo) or `modernc.org/sqlite` (pure Go) for session data and settings.
- **Rationale**: SQLite is the standard for local desktop persistence. Pure Go `modernc.org/sqlite` avoids CGo build complexity while providing full SQL capabilities for session queries and aggregations needed by the Stats view.
- **Alternatives considered**:
  - JSON files — fragile for querying 30 days of session data, no aggregation support.
  - BoltDB/bbolt — key-value only, poor fit for time-range queries and aggregations.
  - SQLite via CGo (`mattn/go-sqlite3`) — faster but requires C compiler in build chain.

## R-005: Frontend Stack

- **Decision**: Vanilla HTML + CSS + JavaScript (no framework).
- **Rationale**: Constitution Principle VI (Simplicity & YAGNI) — with only 3 views (Timer, Stats, Settings) and minimal interactivity, a framework adds unnecessary complexity. Vanilla JS with a simple router pattern keeps the app lean. Chart.js will be used for the weekly bar chart.
- **Alternatives considered**:
  - React — overkill for 3 static views with minimal state.
  - Svelte — lighter but still adds build tooling complexity for little benefit.
  - Vue — same reasoning as React.

## R-006: Chart Library for Stats View

- **Decision**: Use Chart.js (v4) via CDN/bundled for the weekly bar chart.
- **Rationale**: Lightweight, well-documented, and renders canvas-based charts that work perfectly in WebView. Only one chart type (bar) is needed.
- **Alternatives considered**:
  - D3.js — too low-level for a single bar chart.
  - Custom canvas drawing — unnecessary effort when Chart.js handles it with 10 lines.

## R-007: Timer Implementation Strategy

- **Decision**: Go-side timer using `time.Ticker` with Wails runtime events to push updates to the frontend.
- **Rationale**: Constitution Principle I (Go-Backend Ownership) requires all business logic in Go. The timer is business logic — it drives session state transitions, break scheduling, and completion recording. Using `runtime.EventsEmit` to push tick updates to the frontend every second keeps the frontend as a pure renderer.
- **Alternatives considered**:
  - Frontend `setInterval` — violates Constitution Principle I; timer accuracy affected by tab throttling in WebView.
  - Hybrid (Go tracks time, frontend renders independently) — complexity without benefit.

## R-008: Dark/Light Theme Implementation

- **Decision**: CSS custom properties with `prefers-color-scheme` media query for auto-detection, plus a `data-theme` attribute on `<html>` for manual override.
- **Rationale**: Standard CSS approach. Theme preference stored in UserSettings (Go-side) and applied on startup via a Wails binding.
- **Alternatives considered**:
  - CSS-in-JS — no JS framework, so not applicable.
  - Separate CSS files per theme — harder to maintain than custom properties.
