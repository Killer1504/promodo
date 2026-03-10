# Quickstart: Pomodoro Focus Timer

**Branch**: `001-pomodoro-timer`

## Prerequisites

- Go 1.21+
- Node.js 18+ (for frontend build)
- Wails CLI v2: `go install github.com/wailsapp/wails/v2/cmd/wails@latest`

## Setup

```bash
# 1. Initialize the Wails project (from repo root)
wails init -n pomodoro-timer -t vanilla

# 2. Install Go dependencies
go mod tidy

# 3. Install frontend dependencies
cd frontend && npm install && cd ..
```

## Development

```bash
# Run in dev mode (hot-reload)
wails dev
```

## Build

```bash
# Production build
wails build
```

The binary will be at `build/bin/pomodoro-timer.exe` (Windows).

## Project Structure

```
├── main.go                    # App entry point, Wails config
├── app.go                     # AppService (paused session, theme)
├── internal/
│   ├── timer/
│   │   └── service.go         # TimerService (countdown, state machine)
│   ├── stats/
│   │   └── service.go         # StatsService (daily/weekly aggregation)
│   ├── settings/
│   │   └── service.go         # SettingsService (CRUD user prefs)
│   ├── storage/
│   │   ├── database.go        # SQLite init, migrations, cleanup
│   │   └── repository.go      # Session & settings queries
│   └── notification/
│       └── notifier.go        # OS notifications via beeep
├── frontend/
│   ├── index.html             # Main HTML shell
│   ├── src/
│   │   ├── main.js            # App init, router, theme
│   │   ├── views/
│   │   │   ├── timer.js       # Timer view (ring, digits, dots)
│   │   │   ├── stats.js       # Stats view (chart, daily count)
│   │   │   └── settings.js    # Settings view (durations, sound, theme)
│   │   ├── components/
│   │   │   ├── progress-ring.js  # SVG circular progress
│   │   │   ├── tab-bar.js     # Bottom navigation tabs
│   │   │   └── session-dots.js   # 4-dot cycle indicator
│   │   └── styles/
│   │       ├── variables.css  # CSS custom properties (light/dark)
│   │       ├── base.css       # Reset, typography
│   │       ├── timer.css      # Timer view styles
│   │       ├── stats.css      # Stats view styles
│   │       └── settings.css   # Settings view styles
│   ├── assets/
│   │   └── sounds/            # Notification sound files (.mp3)
│   └── wailsjs/               # Auto-generated Wails bindings (DO NOT EDIT)
├── go.mod
├── go.sum
└── wails.json                 # Wails project config
```

## Testing

```bash
# Run Go tests
go test ./...

# Run specific service tests
go test ./internal/timer/...
go test ./internal/stats/...
go test ./internal/storage/...
```
