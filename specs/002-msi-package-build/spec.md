# Feature Specification: MSI Package Build

**Feature Branch**: `002-msi-package-build`  
**Created**: 2026-03-10  
**Status**: Draft  
**Input**: User description: "package build to .msi file to user can easily setup"

## User Scenarios & Testing *(mandatory)*

### User Story 1 - One-Click MSI Installation (Priority: P1)

A new user downloads the Pomodoro Focus Timer `.msi` installer from the project releases. They double-click the file, follow a standard Windows installation wizard (welcome → MIT license acceptance → install location → install → finish), and the app is ready to use — with a Start Menu shortcut and optional Desktop shortcut.

**Why this priority**: This is the core deliverable. Without a working MSI installer, no other feature matters. Users expect a familiar Windows installation experience.

**Independent Test**: Run the MSI installer on a clean Windows machine. The app should install, launch from Start Menu, and function correctly.

**Acceptance Scenarios**:

1. **Given** a user has downloaded the `.msi` file, **When** they double-click it, **Then** a standard Windows Installer wizard appears with app name, version, and branding.
2. **Given** the wizard is open, **When** the user proceeds through all steps and clicks "Install", **Then** the application files are copied to the selected directory and shortcuts are created.
3. **Given** installation completes, **When** the user clicks "Finish", **Then** they can optionally launch the app immediately.

---

### User Story 2 - Clean Uninstallation (Priority: P1)

A user decides to remove Pomodoro Focus Timer. They go to "Apps & Features" (or "Add/Remove Programs"), find the app, and click Uninstall. The app is cleanly removed — all application files are deleted, shortcuts are removed, and no orphaned registry entries remain.

**Why this priority**: Equal to installation — a broken uninstall experience erodes user trust and violates Windows application standards.

**Independent Test**: After installing via MSI, uninstall from Windows Settings. Verify the install directory, shortcuts, and registry entries are cleaned up.

**Acceptance Scenarios**:

1. **Given** the app is installed, **When** the user opens "Apps & Features" and clicks Uninstall, **Then** the Windows Installer removes all application files and shortcuts.
2. **Given** uninstallation completes, **When** the user checks the install directory, **Then** the directory is either removed or contains only user data files.

---

### User Story 3 - Automated Build Pipeline (Priority: P2)

A developer runs a single build command (or script) that compiles the Wails app and produces a ready-to-distribute `.msi` file. This can be integrated into CI/CD pipelines for automated releases.

**Why this priority**: Enables reproducible builds and streamlines the release process. Without automation, every release requires manual packaging steps.

**Independent Test**: Run the build script from a fresh checkout. Verify the `.msi` file is generated in the expected output directory with correct metadata.

**Acceptance Scenarios**:

1. **Given** a developer has the required build tools installed, **When** they run the MSI build command, **Then** a valid `.msi` file is produced in the build output directory.
2. **Given** the build completes, **When** inspecting the MSI file properties, **Then** it contains correct product name, version, manufacturer, and icon.

---

### User Story 4 - Upgrade Over Existing Installation (Priority: P2)

A user already has Pomodoro Focus Timer installed. When a new version is released and they run the new MSI, the installer detects the existing installation and upgrades it in place — preserving user data stored in `~/.pomodoro-timer/`.

**Why this priority**: Users should not need to manually uninstall before upgrading. In-place upgrades are a standard expectation for Windows applications.

**Independent Test**: Install v1 of the MSI, create some timer sessions (to generate data), then install v2. Verify the app is updated and all session data is preserved.

**Acceptance Scenarios**:

1. **Given** an older version is installed, **When** the user runs a newer MSI, **Then** the installer upgrades the application without requiring manual uninstall.
2. **Given** an upgrade completes, **When** the user opens the app, **Then** all previous settings, sessions, and data in `~/.pomodoro-timer/` are intact.

---

### Edge Cases

- What happens when the user cancels the installation mid-way? The system should roll back any partial changes.
- How does the installer handle insufficient disk space? The user should see a clear error message before installation begins.
- What happens if the user tries to install without Administrator privileges? The installer should request elevation or display a clear permissions error.
- What happens when the app is currently running during an upgrade or uninstall? The installer should prompt the user to close the app first.
- What happens if the user data directory (`~/.pomodoro-timer/`) does not exist on first install? The app should create it on first launch (existing behavior), and the installer should NOT create it.

## Requirements *(mandatory)*

### Functional Requirements

- **FR-001**: The build system MUST produce a valid Windows Installer (`.msi`) package from the compiled Wails application binary.
- **FR-002**: The MSI MUST include the application executable, icon, and the WebView2 online bootstrapper (small ~1.8 MB component that downloads the runtime if not already present on the target machine).
- **FR-003**: The installer MUST create a Start Menu shortcut for the application.
- **FR-004**: The installer MUST provide an option to create a Desktop shortcut.
- **FR-005**: The installer MUST allow the user to choose a custom installation directory (with a sensible default like `C:\Program Files\Pomodoro Focus Timer\`).
- **FR-006**: The installer MUST register the application in Windows "Apps & Features" for clean uninstallation.
- **FR-007**: The MSI MUST support Major Upgrade — installing a newer version should automatically remove the previous version.
- **FR-008**: The uninstaller MUST remove all application files, shortcuts, and registry entries created during installation.
- **FR-009**: The uninstaller MUST NOT remove user data stored in `~/.pomodoro-timer/`.
- **FR-010**: The build process MUST be executable via a single command or script for CI/CD integration.
- **FR-011**: The MSI MUST embed correct product metadata: product name ("Pomodoro Focus Timer"), version, manufacturer ("hung le"), and application icon.
- **FR-012**: The installer MUST display the MIT License text in a license acceptance dialog during the installation wizard.
- **FR-013**: The build script MUST read the product version from a single source-of-truth project file (e.g., `wails.json` or a dedicated `version.txt` at repo root).

### Key Entities

- **MSI Package**: The distributable installer file containing the compiled app, metadata (name, version, manufacturer, icon), and installation logic (shortcuts, registry entries, upgrade behavior).
- **Build Script**: A reproducible automation script that compiles the Wails app and packages it into the MSI, reading the version number from a project file (e.g., `wails.json` or `version.txt`).

## Success Criteria *(mandatory)*

### Measurable Outcomes

- **SC-001**: Users can install the application in under 2 minutes via a standard Windows installer wizard.
- **SC-002**: The installed application launches correctly from both Start Menu and Desktop shortcuts.
- **SC-003**: Users can cleanly uninstall the application through Windows "Apps & Features" with zero orphaned files in the installation directory.
- **SC-004**: Upgrading from one version to the next preserves 100% of user data and settings.
- **SC-005**: The MSI build process completes in under 5 minutes from a single command invocation.
- **SC-006**: The MSI file size is under 20 MB (proportional to the ~6 MB executable).

## Assumptions

- The target platform is Windows 10+ (64-bit only).
- The Wails v2 `wails build` command is used to produce the Windows executable before MSI packaging.
- WiX Toolset v4+ or equivalent MSI authoring tool will be used (implementation detail — not a requirement for this spec).
- The application does not require a Windows Service — it runs as a standard user-space application.
- Code signing is out of scope for this feature (can be added as a separate enhancement). Users may see a SmartScreen warning on first launch.
- The existing user data directory (`~/.pomodoro-timer/`) is managed by the application at runtime, not by the installer.

## Clarifications

### Session 2026-03-10

- Q: How should WebView2 runtime dependency be handled — online bootstrapper, offline embedded, or no bundling? → A: Online bootstrapper (Option A). Keeps MSI small; downloads WebView2 only if missing. Windows 11+ has it pre-installed.
- Q: Should the installer show a license agreement dialog, and if so, which license? → A: Show MIT License (Option A). Display the MIT License text in a license acceptance dialog during installation.
- Q: Where should the MSI product version number come from? → A: File-based (Option A). Read version from a project file (`wails.json` or `version.txt`) as single source of truth.
