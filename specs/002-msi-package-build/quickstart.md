# Quickstart: MSI Package Build

**Branch**: `002-msi-package-build` | **Date**: 2026-03-10

## Prerequisites

1. **Go 1.24+** — [go.dev/dl](https://go.dev/dl/)
2. **Node.js 18+** — [nodejs.org](https://nodejs.org/)
3. **Wails CLI v2** — `go install github.com/wailsapp/wails/v2/cmd/wails@latest`
4. **.NET SDK 8+** — [dotnet.microsoft.com](https://dotnet.microsoft.com/download) (required for WiX)
5. **WiX Toolset v4** — `dotnet tool install --global wix` then `wix extension add WixToolset.UI.wixext`

## Build MSI

```powershell
# From repo root — single command
.\build\msi\build-msi.ps1

# With explicit version override
.\build\msi\build-msi.ps1 -Version "1.2.0"
```

Output: `build\bin\PomodoroFocusTimer-x64.msi`

## What the Script Does

1. Reads version from `wails.json`
2. Runs `wails build` to produce `PomodoroFocusTimer.exe`
3. Runs `wix build` to package the exe + bootstrapper into `.msi`
4. Outputs the MSI to `build\bin\`

## Manual Testing

1. **Install**: Double-click the `.msi` → follow wizard → verify app in Start Menu + Desktop
2. **Run**: Launch from Start Menu shortcut → timer should work normally
3. **Uninstall**: Apps & Features → Pomodoro Focus Timer → Uninstall → verify clean removal
4. **Upgrade**: Build a new version → run new `.msi` → verify old version replaced, data preserved
