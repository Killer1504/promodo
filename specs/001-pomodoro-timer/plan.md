# Implementation Plan: Pomodoro Focus Timer

**Branch**: `001-pomodoro-timer` | **Date**: 2026-03-10 | **Spec**: [spec.md](file:///d:/122.Test-Claude/14.DemoWails/specs/001-pomodoro-timer/spec.md)
**Input**: Feature specification from `/specs/001-pomodoro-timer/spec.md`

## Summary

Build a Pomodoro Focus Timer desktop application using Go + Wails v2. The app features a circular progress timer with session cycling (4 focus sessions → long break), OS notifications with selectable sounds, daily/weekly statistics with a bar chart, and customizable durations. The Go backend owns all business logic (timer state machine, session persistence, settings), while the frontend is a thin vanilla HTML/CSS/JS renderer communicating via typed Wails bindings. Data persists in SQLite (pure Go driver). The app runs as a compact 400×600 px window with system tray support and dark/light theme.

## Technical Context

**Language/Version**: Go 1.21+ with Wails v2 framework
**Primary Dependencies**: `wailsapp/wails/v2`, `modernc.org/sqlite`, `gen2brain/beeep`, `getlantern/systray`, Chart.js (frontend)
**Storage**: SQLite via `modernc.org/sqlite` (pure Go, no CGo)
**Testing**: `go test` with table-driven tests (backend), manual + browser verification (frontend)
**Target Platform**: Windows (primary), macOS/Linux (secondary)
**Project Type**: Desktop application (Wails hybrid Go + WebView)
**Performance Goals**: Timer tick < 1ms, Stats view load < 1s, App launch < 2s
**Constraints**: < 100 MB memory, 400×600 px non-resizable window, 30-day data retention
**Scale/Scope**: Single user, ~30 sessions/day max, 3 views (Timer/Stats/Settings)

## Constitution Check

*GATE: Must pass before Phase 0 research. Re-check after Phase 1 design.*

| Principle | Status | Evidence |
|-----------|--------|----------|
| **I. Go-Backend Ownership** | ✅ PASS | Timer logic, session persistence, settings, notifications — all in Go. Frontend only renders. (See R-007) |
| **II. Typed Contract Bridge** | ✅ PASS | All 4 services use well-defined Go structs for input/output. Frontend uses auto-generated TS bindings. (See contracts/wails-bindings.md) |
| **III. Security-First Desktop Mindset** | ✅ PASS | No user input reaches file system. SQLite path is fixed. No external process execution. No credentials stored. |
| **IV. Frontend Quality & Accessibility** | ✅ PASS | Keyboard navigation planned for all views. Dark/light themes with WCAG AA contrast. Loading/error/empty states defined in spec. |
| **V. Test-First Verification** | ✅ PASS | Go unit tests for timer, storage, stats services. Table-driven test pattern. Integration tests for Wails bindings. |
| **VI. Simplicity & YAGNI** | ✅ PASS | Vanilla JS (no framework). 3 Go dependencies justified in research.md. No abstractions without concrete use cases. |
| **VII. Observability & Error Handling** | ✅ PASS | Structured logging via `slog`. User-friendly error messages in UI. Critical errors logged to file. |

**Gate Result**: ✅ ALL GATES PASS — proceed to implementation.

## Project Structure

### Documentation (this feature)

```text
specs/001-pomodoro-timer/
├── spec.md              # Feature specification (19 FRs, 10 clarifications)
├── plan.md              # This file
├── research.md          # Phase 0: 8 research decisions
├── data-model.md        # Phase 1: SQLite schema + entities
├── quickstart.md        # Phase 1: Setup & dev commands
├── contracts/
│   └── wails-bindings.md # Phase 1: 4 services, typed methods & events
├── checklists/
│   └── requirements.md  # Spec quality checklist
└── tasks.md             # Phase 2 output (not yet created)
```

### Source Code (repository root)

```text
├── main.go                      # Wails app entry, window config (400×600)
├── app.go                       # AppService: paused session, theme resolution
├── internal/
│   ├── timer/
│   │   ├── service.go           # TimerService: state machine, Go ticker, events
│   │   └── service_test.go      # Unit tests: start/pause/resume/reset/cycle
│   ├── stats/
│   │   ├── service.go           # StatsService: daily + weekly aggregation
│   │   └── service_test.go      # Unit tests: empty state, date ranges, aggregation
│   ├── settings/
│   │   ├── service.go           # SettingsService: CRUD, defaults, reset
│   │   └── service_test.go      # Unit tests: update, persist, reset
│   ├── storage/
│   │   ├── database.go          # SQLite init, migrations, 30-day cleanup
│   │   ├── repository.go        # Session & settings queries
│   │   └── database_test.go     # Integration tests: schema, CRUD, retention
│   └── notification/
│       ├── notifier.go          # OS notification via beeep
│       └── notifier_test.go     # Unit tests: message formatting
├── frontend/
│   ├── index.html               # HTML shell with tab bar structure
│   ├── src/
│   │   ├── main.js              # Init, router, theme application, event listeners
│   │   ├── views/
│   │   │   ├── timer.js         # Timer view: ring, digits, dots, controls
│   │   │   ├── stats.js         # Stats view: Chart.js bar chart, daily summary
│   │   │   └── settings.js      # Settings view: duration inputs, sound picker, theme toggle
│   │   ├── components/
│   │   │   ├── progress-ring.js # SVG circular progress ring component
│   │   │   ├── tab-bar.js       # Bottom navigation tab bar
│   │   │   └── session-dots.js  # 4-dot cycle indicator
│   │   └── styles/
│   │       ├── variables.css    # CSS custom properties (light + dark tokens)
│   │       ├── base.css         # Reset, fonts, global styles
│   │       ├── timer.css        # Timer view layout & animations
│   │       ├── stats.css        # Stats view chart & summary styles
│   │       └── settings.css     # Settings form & controls styles
│   ├── assets/
│   │   └── sounds/              # bell.mp3, chime.mp3, ding.mp3
│   └── wailsjs/                 # Auto-generated (DO NOT EDIT)
├── go.mod
├── go.sum
└── wails.json
```

**Structure Decision**: Single-project Wails v2 layout. Go services in `internal/` with clear domain separation (timer, stats, settings, storage, notification). Frontend in `frontend/` with vanilla JS views and components. This keeps the single-binary philosophy of Wails and respects Constitution Principle VI (Simplicity).

## Complexity Tracking

No constitution violations — table not needed.

## Verification Plan

### Automated Tests (Go)

```bash
# Run all Go tests
go test ./...

# Run with verbose output
go test -v ./internal/timer/...
go test -v ./internal/storage/...
go test -v ./internal/stats/...
go test -v ./internal/settings/...
```

**Test coverage targets**:
- `timer/service.go` — state transitions (idle→running→paused→running→completed), cycle counting (4 focus→long break), elapsed-time accuracy
- `storage/database.go` — schema creation, 30-day cleanup, session CRUD
- `stats/service.go` — daily/weekly aggregation, empty state, date boundary handling
- `settings/service.go` — get/update/reset, default values, persistence

### Manual Verification (Frontend + Integration)

1. **Timer flow**: Launch app → click Start Focus → verify ring animates and digits count down → pause → resume → let it finish → verify notification pops → verify 3-second break countdown → verify break starts
2. **System tray**: During a running session → minimize window → verify tray icon appears → click tray icon → verify window restores
3. **Session dots**: Start and complete 4 focus sessions → verify dots fill progressively → verify long break triggers after 4th
4. **Stats view**: After completing sessions → switch to Stats tab → verify daily count, minutes, and bar chart are correct
5. **Settings**: Switch to Settings → change focus duration to 50 min → go back → start session → verify 50:00 countdown. Change theme → verify UI switches. Select new sound → verify preview plays.
6. **Persistence**: Close and reopen app → verify settings persist. Pause a session → close app → reopen → verify resume prompt appears.
7. **Dark/light mode**: Test with OS in light mode → verify auto-detect. Toggle to dark in Settings → verify switch. Restart → verify preference persists.
