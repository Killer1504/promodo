# Data Model: MSI Package Build

**Branch**: `002-msi-package-build` | **Date**: 2026-03-10

## Entities

### MSI Package (WiX Product)

| Attribute | Source | Description |
|-----------|--------|-------------|
| `ProductName` | `wails.json → outputfilename` | Display name: "Pomodoro Focus Timer" |
| `Version` | `wails.json → version` (new field) | Semantic version: MAJOR.MINOR.PATCH |
| `Manufacturer` | `wails.json → author.name` | "hung le" |
| `UpgradeCode` | WiX `.wxs` (hardcoded GUID) | Fixed GUID — must never change across versions. Identifies the product family for major upgrades. |
| `ProductCode` | Auto-generated per build | Unique per version. WiX auto-generates with `*`. |
| `InstallDirectory` | User-selected (default: `C:\Program Files\Pomodoro Focus Timer\`) | Target path for application files |

### Installed Files

| File | Source | Install Location |
|------|--------|------------------|
| `PomodoroFocusTimer.exe` | `wails build` output | `[InstallDir]` |
| `MicrosoftEdgeWebview2Setup.exe` | Downloaded from Microsoft | Bundled in MSI, executed as custom action |

### Shortcuts Created

| Shortcut | Target | Location |
|----------|--------|----------|
| Start Menu | `[InstallDir]\PomodoroFocusTimer.exe` | `[ProgramMenuFolder]\Pomodoro Focus Timer` |
| Desktop (optional) | `[InstallDir]\PomodoroFocusTimer.exe` | `[DesktopFolder]` |

### Registry Entries (managed by Windows Installer)

| Purpose | Key |
|---------|-----|
| Add/Remove Programs | `HKLM\SOFTWARE\Microsoft\Windows\CurrentVersion\Uninstall\{ProductCode}` |

## State Transitions

```
Source Code → [wails build] → .exe Binary → [wix build] → .msi Package
                                                            ↓
                                          User runs .msi → Install → Running App
                                                            ↓
                                          User upgrades  → Major Upgrade (auto-uninstall old) → New Version Running
                                                            ↓
                                          User uninstalls → Clean removal (preserve user data)
```

## Data Not Managed by Installer

- `~/.pomodoro-timer/pomodoro.db` — SQLite database (user data)
- `~/.pomodoro-timer/paused_session.json` — saved state
- `~/.pomodoro-timer/pomodoro.log` — application log

These files are created/managed by the application at runtime and explicitly excluded from installer scope (FR-009).
