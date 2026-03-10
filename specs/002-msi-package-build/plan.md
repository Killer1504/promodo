# Implementation Plan: MSI Package Build

**Branch**: `002-msi-package-build` | **Date**: 2026-03-10 | **Spec**: [spec.md](file:///d:/122.Test-Claude/14.DemoWails/specs/002-msi-package-build/spec.md)
**Input**: Feature specification from `/specs/002-msi-package-build/spec.md`

## Summary

Package the Pomodoro Focus Timer Wails application into a Windows Installer (`.msi`) file using WiX Toolset v4. The MSI will provide a standard Windows installation wizard with MIT license acceptance, installable directory selection, Start Menu + Desktop shortcuts, Add/Remove Programs registration, major upgrade support, and WebView2 online bootstrapper. A PowerShell build script will automate the entire pipeline: parse version → compile app → package MSI.

## Technical Context

**Language/Version**: Go 1.24+ (app), PowerShell 5.1+ (build script), XML (WiX `.wxs`)
**Primary Dependencies**: Wails v2 CLI, WiX Toolset v4 (`dotnet tool`), .NET SDK 8+
**Storage**: N/A (build tooling only)
**Testing**: Manual install/uninstall/upgrade testing on Windows
**Target Platform**: Windows 10+ (64-bit only, x64)
**Project Type**: Desktop app installer / build tooling
**Performance Goals**: MSI build < 5 minutes, MSI file < 20 MB
**Constraints**: No code signing (out of scope), WebView2 via online bootstrapper
**Scale/Scope**: Single binary packaging, single platform

## Constitution Check

*GATE: Must pass before Phase 0 research. Re-check after Phase 1 design.*

| Principle | Status | Notes |
|-----------|--------|-------|
| I. Go-Backend Ownership | ✅ N/A | No Go logic changes |
| II. Typed Contract Bridge | ✅ N/A | No Wails binding changes |
| III. Security-First | ✅ Pass | Build script uses explicit paths, no secrets, no shell interpolation |
| IV. Frontend Quality | ✅ N/A | No frontend changes |
| V. Test-First | ✅ Pass | Manual test plan defined; build script verifiable via dry run |
| VI. Simplicity & YAGNI | ✅ Pass | WiX is simplest tool for MSI; single script; no abstractions |
| VII. Observability | ✅ Pass | Build script outputs progress to console; WiX logs build warnings |

**Constitution note**: The constitution mentions "NSIS installer for Windows distribution" in Technology Stack. This plan **adds** MSI alongside NSIS — it does not replace NSIS. The constitution could be amended later to include "MSI installer" — flagged but not blocking.

## Project Structure

### Documentation (this feature)

```text
specs/002-msi-package-build/
├── spec.md              # Feature specification (done)
├── plan.md              # This file
├── research.md          # Phase 0 output (done)
├── data-model.md        # Phase 1 output (done)
├── quickstart.md        # Phase 1 output (done)
└── tasks.md             # Phase 2 output (/speckit.tasks)
```

### Source Code (new files)

```text
build/
├── msi/                     # [NEW] MSI packaging directory
│   ├── build-msi.ps1        # [NEW] Build automation script
│   ├── Product.wxs          # [NEW] WiX installer definition
│   ├── License.rtf          # [NEW] MIT License in RTF format (for WiX UI)
│   └── WebView2/            # [NEW] WebView2 bootstrapper download cache
└── windows/
    └── installer/           # [KEEP] Existing NSIS installer (unchanged)

LICENSE                      # [NEW] MIT License plaintext at repo root
wails.json                   # [MODIFY] Add "version" field
```

**Structure Decision**: All MSI-related files go under `build/msi/` — separated from the existing NSIS installer at `build/windows/installer/`. This ensures clean coexistence and follows the existing `build/` directory convention.

---

## Proposed Changes

### Component 1: Version Management

#### [MODIFY] [wails.json](file:///d:/122.Test-Claude/14.DemoWails/wails.json)

Add a `version` field to serve as the single source of truth for product versioning:

```diff
 {
   "$schema": "https://wails.io/schemas/config.v2.json",
   "name": "pomodoro-timer",
   "outputfilename": "PomodoroFocusTimer",
+  "info": {
+    "productVersion": "1.0.0",
+    "companyName": "hung le",
+    "productName": "Pomodoro Focus Timer",
+    "copyright": "Copyright 2026 hung le",
+    "comments": "A sleek desktop Pomodoro timer"
+  },
   "frontend:install": "npm install",
```

> The `info` block already exists conceptually (referenced by `build/windows/info.json` template) but needs explicit values in `wails.json`. The build script will read `info.productVersion` from here.

---

### Component 2: License File

#### [NEW] [LICENSE](file:///d:/122.Test-Claude/14.DemoWails/LICENSE)

Standard MIT License plaintext file at repo root. Content: MIT License with `hung le` as copyright holder, year 2026.

#### [NEW] [License.rtf](file:///d:/122.Test-Claude/14.DemoWails/build/msi/License.rtf)

RTF-formatted version of the MIT License for the WiX installer UI. WiX requires RTF format for the license agreement dialog. The build script can generate this from `LICENSE`, or we maintain it as a static file.

**Decision**: Maintain as a static RTF file for simplicity — the license text rarely changes.

---

### Component 3: WiX Installer Definition

#### [NEW] [Product.wxs](file:///d:/122.Test-Claude/14.DemoWails/build/msi/Product.wxs)

The core WiX v4 installer definition file. Key features:

1. **Product metadata**: Name, version, manufacturer, icon — sourced via variables from build script
2. **UpgradeCode**: Fixed GUID for major upgrade detection (must never change)
3. **MajorUpgrade element**: Auto-uninstalls previous versions (FR-007)
4. **Directory structure**: `ProgramFiles64Folder\Pomodoro Focus Timer\`
5. **Components**:
   - Main executable (`PomodoroFocusTimer.exe`)
   - Start Menu shortcut (FR-003)
   - Desktop shortcut with condition (FR-004)
6. **WixUI_Mondo dialog set**: Welcome → License (MIT) → Directory → Install → Finish (FR-012)
7. **Custom action**: Run WebView2 bootstrapper on install (FR-002)
8. **Icon**: References `build/windows/icon.ico`

---

### Component 4: Build Automation

#### [NEW] [build-msi.ps1](file:///d:/122.Test-Claude/14.DemoWails/build/msi/build-msi.ps1)

PowerShell build script that orchestrates the entire MSI build pipeline:

```
Input: -Version (optional override) or reads from wails.json
Steps:
  1. Validate prerequisites (Go, Node, Wails CLI, .NET SDK, WiX)
  2. Parse version from wails.json (or use -Version override)
  3. Run `wails build` to produce .exe
  4. Download WebView2 bootstrapper if not cached
  5. Run `wix build` with version variables → produce .msi
  6. Print output path and file size
Output: build/bin/PomodoroFocusTimer-{version}-x64.msi
```

Error handling:
- Exit with clear error if any prerequisite is missing
- Exit if `wails build` fails
- Exit if `wix build` fails with WiX error output

---

## Verification Plan

### Manual Verification (User Testing)

Since this feature produces installer artifacts, automated unit tests are not applicable. Verification is done through manual installation testing:

#### Test 1: Fresh Install

1. Run `.\build\msi\build-msi.ps1` from repo root
2. Verify MSI is created at `build\bin\PomodoroFocusTimer-1.0.0-x64.msi`
3. Verify MSI file size is under 20 MB
4. Double-click the MSI → verify installation wizard appears
5. Verify wizard shows: Welcome → MIT License → Directory selection → Install → Finish
6. After install, check `C:\Program Files\Pomodoro Focus Timer\` for `PomodoroFocusTimer.exe`
7. Verify Start Menu shortcut exists and launches the app
8. Verify Desktop shortcut exists and launches the app
9. Verify the timer app functions normally (start a focus session)

#### Test 2: Clean Uninstall

1. Open Windows Settings → Apps & Features
2. Find "Pomodoro Focus Timer" in the list
3. Click Uninstall → verify the app is removed
4. Check `C:\Program Files\Pomodoro Focus Timer\` — directory should be gone
5. Check Start Menu and Desktop — shortcuts should be gone
6. Check `~/.pomodoro-timer/` — user data should still be there (if it exists)

#### Test 3: Major Upgrade

1. Build MSI with version 1.0.0 and install
2. Create a timer session to generate user data
3. Modify version to 1.1.0 in `wails.json`
4. Build new MSI and run it
5. Verify old version is auto-removed, new version installed
6. Verify user data in `~/.pomodoro-timer/` is preserved
7. Verify app launches and previous session data is intact

#### Test 4: Build Script Validation

1. Run `.\build\msi\build-msi.ps1` from a clean checkout
2. Verify it completes without errors
3. Verify the output MSI file path is printed
4. Right-click MSI → Properties → Details → verify product name, version, manufacturer

## Complexity Tracking

No constitution violations — no complexity justification needed.
