# 🍅 Pomodoro Focus Timer

A sleek desktop Pomodoro timer built with [Wails v2](https://wails.io/) — Go backend + vanilla JS/HTML/CSS frontend.

![Go](https://img.shields.io/badge/Go-1.24-00ADD8?logo=go&logoColor=white)
![Wails](https://img.shields.io/badge/Wails-v2.11-red)
![SQLite](https://img.shields.io/badge/SQLite-embedded-003B57?logo=sqlite)
![License](https://img.shields.io/badge/License-MIT-green)

## ✨ Features

- **Focus Timer** — 25/5/15 min Pomodoro cycles with visual ring countdown
- **Session Tracking** — automatic cycle progression (4 focus → long break)
- **Stats Dashboard** — daily session count, weekly Chart.js bar chart
- **Settings** — customizable durations, light/dark/system theme, mute toggle
- **Persistence** — SQLite storage, paused session restoration across restarts
- **OS Notifications** — alerts when focus sessions or breaks end
- **Keyboard Accessible** — visible focus indicators on all controls

## 📦 Prerequisites

| Tool | Version | Install |
|------|---------|---------|
| **Go** | 1.24+ | [go.dev/dl](https://go.dev/dl/) |
| **Node.js** | 18+ | [nodejs.org](https://nodejs.org/) |
| **Wails CLI** | v2 | `go install github.com/wailsapp/wails/v2/cmd/wails@latest` |

> **Windows**: also requires [WebView2 runtime](https://developer.microsoft.com/en-us/microsoft-edge/webview2/) (included in Windows 11+).

## 🚀 Quick Start

```bash
# 1. Clone
git clone <repo-url>
cd pomodoro-timer

# 2. Install frontend dependencies
cd frontend && npm install && cd ..

# 3. Development mode (hot-reload)
wails dev

# 4. Production build
wails build
```

The built binary will be in `build/bin/`.

### 📦 Build MSI Installer

To create a Windows Installer (`.msi`) package:

**Additional prerequisites**: [.NET SDK 8+](https://dotnet.microsoft.com/download), [WiX Toolset v4](https://wixtoolset.org/) (`dotnet tool install --global wix`)

```powershell
# Build MSI (includes wails build + WiX packaging)
.\build\msi\build-msi.ps1

# With version override
.\build\msi\build-msi.ps1 -Version "2.0.0"

# Skip wails build if exe already exists
.\build\msi\build-msi.ps1 -SkipWailsBuild
```

Output: `build/bin/PomodoroFocusTimer-{version}-x64.msi` (~7 MB)

## 🧪 Running Tests

```bash
# All Go tests (35 tests across 5 packages)
go test -v -count=1 ./internal/...

# Individual packages
go test -v ./internal/timer/...
go test -v ./internal/storage/...
go test -v ./internal/settings/...
go test -v ./internal/stats/...
go test -v ./internal/notification/...
```

## 📂 Project Structure

```
pomodoro-timer/
├── main.go                     # Wails app entry point
├── app.go                      # App struct — all frontend bindings
├── internal/
│   ├── timer/                  # Timer state machine
│   ├── storage/                # SQLite database + repository
│   ├── settings/               # Settings validation service
│   ├── stats/                  # Stats aggregation service
│   └── notification/           # OS notification wrapper
├── frontend/
│   ├── index.html              # App shell
│   └── src/
│       ├── main.js             # Router + theme detection
│       ├── components/         # Tab bar, progress ring, session dots
│       ├── views/              # Timer, stats, settings views
│       └── styles/             # CSS design system (light + dark)
├── specs/                      # Feature specifications
└── wails.json                  # Wails configuration
```

## 🎨 Themes

Supports **Light**, **Dark**, and **System** (auto-detect) themes. Configured in Settings and applied instantly without restart.

## 📊 Data Storage

All data is stored locally in `~/.pomodoro-timer/`:

| File | Purpose |
|------|---------|
| `pomodoro.db` | SQLite — sessions + settings |
| `paused_session.json` | Saved state on app close |
| `pomodoro.log` | Structured log output |

Data auto-cleans sessions older than 30 days on startup.

## 📝 License

MIT
