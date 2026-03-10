# Research: MSI Package Build

**Branch**: `002-msi-package-build` | **Date**: 2026-03-10

## Decision 1: MSI Authoring Tool

**Decision**: WiX Toolset v4 (via `dotnet tool install wix`)

**Rationale**:
- WiX is the industry-standard open-source MSI authoring tool.
- Version 4 uses a modern `dotnet tool` workflow — no separate installer needed.
- Native support for all spec requirements: license dialogs, major upgrades, shortcuts, custom install directory, Add/Remove Programs registration.
- XML-based `.wxs` source files are version-controllable and CI/CD friendly.
- Constitution Principle VI (Simplicity): WiX is the simplest tool that directly produces `.msi` — no intermediate formats or wrappers.

**Alternatives Considered**:
- **NSIS**: Already exists in the project (`build/windows/installer/project.nsi`) but produces `.exe` not `.msi`. NSIS cannot produce native MSI files. Keep existing NSIS as-is for users who prefer `.exe` installers.
- **Advanced Installer**: Commercial tool. Violates Simplicity principle — adds licensing cost and vendor dependency.
- **InstallShield**: Enterprise-grade, expensive, overkill for a personal desktop app.
- **Inno Setup**: Like NSIS, produces `.exe` only, not `.msi`.

## Decision 2: WebView2 Bootstrapper Integration

**Decision**: Bundle the Microsoft WebView2 online bootstrapper (`MicrosoftEdgeWebview2Setup.exe`, ~1.8 MB) as a custom action in the MSI.

**Rationale**:
- Clarification Q1 resolved: use online bootstrapper.
- The bootstrapper is a small download from Microsoft that auto-installs WebView2 if not present.
- Windows 11+ ships with WebView2 pre-installed — the bootstrapper is a no-op on those systems.
- The existing NSIS script already uses `!insertmacro wails.webview2runtime` — the WiX MSI should replicate this behavior.

**Alternatives Considered**:
- **Offline embedded runtime** (~150 MB): Rejected per clarification — too large for MSI size target.
- **No bundling**: Risky on Windows 10 where WebView2 may not be installed.

## Decision 3: Version Source

**Decision**: Read version from `wails.json` using a PowerShell build script.

**Rationale**:
- Clarification Q3 resolved: file-based version sourcing.
- `wails.json` already exists and is the project's canonical config file.
- The build script will parse `wails.json`, extract the version, and pass it to both `wails build` and `wix build`.
- Need to add a `version` field to `wails.json` (currently absent).

**Alternatives Considered**:
- **Dedicated `version.txt`**: Adds another file to maintain. `wails.json` is sufficient.
- **Git tags**: Requires git, fails in non-git builds, adds complexity.

## Decision 4: Coexistence with NSIS Installer

**Decision**: Keep the existing NSIS installer alongside the new MSI build.

**Rationale**:
- The NSIS `.exe` installer serves a different purpose — some users prefer `.exe` over `.msi`.
- Ripping out NSIS would break the existing `wails build --nsis` workflow.
- Constitution Principle VI: don't remove working functionality.
- Both can coexist — different build commands produce different outputs.

## Decision 5: License File

**Decision**: Create a `LICENSE` file at repo root with MIT License text, referenced by the WiX installer.

**Rationale**:
- Clarification Q2 resolved: show MIT license dialog.
- WiX's `WixUI_Mondo` or `WixUI_Advanced` dialog set includes a license agreement page.
- The license file needs to be in RTF format for WiX UI (WiX requirement).
- A `LICENSE` plaintext file at repo root is standard; the build script will convert to RTF for WiX.

## Decision 6: Build Script Technology

**Decision**: PowerShell script (`build-msi.ps1`) wrapping `wails build` + `wix build`.

**Rationale**:
- Project runs on Windows (primary build target per constitution).
- PowerShell is universally available on Windows 10+.
- Script orchestrates: parse version → `wails build` → `wix build` → output MSI.
- Constitution Principle VI: a single script satisfies FR-010 (single command build).
