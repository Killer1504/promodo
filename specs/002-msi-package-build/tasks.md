# Tasks: MSI Package Build

**Input**: Design documents from `/specs/002-msi-package-build/`
**Prerequisites**: plan.md ✅, spec.md ✅, research.md ✅, data-model.md ✅, quickstart.md ✅

**Tests**: No automated test tasks — this feature produces build artifacts verified through manual install/uninstall testing.

**Organization**: Tasks grouped by user story for independent implementation.

## Format: `[ID] [P?] [Story] Description`

- **[P]**: Can run in parallel (different files, no dependencies)
- **[Story]**: Which user story this task belongs to (e.g., US1, US2, US3, US4)

---

## Phase 1: Setup (Prerequisites & Project Structure)

**Purpose**: Install tools and create the directory structure for MSI packaging.

- [x] T001 Install WiX Toolset v4 via `dotnet tool install --global wix`
- [x] T002 Install WiX UI extension via `wix extension add WixToolset.UI.wixext`
- [x] T003 Create MSI build directory structure at `build/msi/`
- [x] T004 [P] Add `info` block with `productVersion`, `companyName`, `productName`, `copyright`, and `comments` fields to `wails.json`
- [x] T005 [P] Create MIT License plaintext file at repo root `LICENSE`

**Checkpoint**: Tools installed, directory created, version source and license file ready.

---

## Phase 2: Foundational (WiX Definition & License)

**Purpose**: Create the WiX installer definition and license RTF — these are required by ALL user stories.

**⚠️ CRITICAL**: No user story work can begin until this phase is complete.

- [x] T006 Create MIT License RTF file at `build/msi/License.rtf` (WiX UI requires RTF format)
- [x] T007 Download WebView2 online bootstrapper (`MicrosoftEdgeWebview2Setup.exe`) to `build/msi/WebView2/`
- [x] T008 Create WiX installer definition at `build/msi/Product.wxs` with:
  - Product metadata (name, version, manufacturer, icon) using WiX variables
  - Fixed UpgradeCode GUID (generate once, never change)
  - MajorUpgrade element for automatic previous version removal
  - Directory structure: `ProgramFiles64Folder\Pomodoro Focus Timer\`
  - Component for main executable (`PomodoroFocusTimer.exe`)
  - Component for Start Menu shortcut
  - Component for Desktop shortcut (with user-selectable feature)
  - WixUI_Mondo dialog set (Welcome → License → Directory → Install → Finish)
  - Custom action to run WebView2 bootstrapper on install
  - Reference to `build/windows/icon.ico` for installer and ARP icon

**Checkpoint**: WiX definition complete — can now build and test the MSI.

---

## Phase 3: User Story 1 — One-Click MSI Installation (Priority: P1) 🎯 MVP

**Goal**: User double-clicks `.msi`, follows wizard, app installs with shortcuts.

**Independent Test**: Run MSI on Windows → wizard appears → install completes → app launches from Start Menu and Desktop shortcuts.

### Implementation for User Story 1

- [x] T009 [US1] Create build automation script at `build/msi/build-msi.ps1` with:
  - Parameter: `-Version` (optional override)
  - Step 1: Validate prerequisites (Go, Node, Wails CLI, .NET SDK, WiX)
  - Step 2: Parse `info.productVersion` from `wails.json` (or use `-Version`)
  - Step 3: Run `wails build` to produce exe in `build/bin/`
  - Step 4: Run `wix build` with `-d Version=X.Y.Z` and other variables → produce `.msi`
  - Step 5: Print output path and file size
  - Error handling: exit with clear message if any step fails
- [x] T010 [US1] Build the MSI by running `.\build\msi\build-msi.ps1` and verify `.msi` file is created in `build/bin/`
- [x] T011 [US1] Manual test: Double-click MSI → verify wizard shows Welcome → MIT License → Directory → Install → Finish
- [x] T012 [US1] Manual test: After install, verify `PomodoroFocusTimer.exe` exists in `C:\Program Files\Pomodoro Focus Timer\`
- [x] T013 [US1] Manual test: Verify Start Menu shortcut launches the app correctly
- [x] T014 [US1] Manual test: Verify Desktop shortcut launches the app correctly
- [x] T015 [US1] Manual test: Verify app functions normally (start a focus session, check themes)

**Checkpoint**: MSI installs correctly, app runs from shortcuts — MVP complete.

---

## Phase 4: User Story 2 — Clean Uninstallation (Priority: P1)

**Goal**: User uninstalls via Apps & Features, all app files and shortcuts removed, user data preserved.

**Independent Test**: Uninstall from Windows Settings → verify install dir removed, shortcuts gone, `~/.pomodoro-timer/` preserved.

### Implementation for User Story 2

- [x] T016 [US2] Verify WiX `Product.wxs` uninstall behavior: `RemoveFolder` elements for install directory and shortcut folders are correctly defined
- [x] T017 [US2] Verify `Product.wxs` does NOT reference or touch `~/.pomodoro-timer/` directory (FR-009)
- [x] T018 [US2] Manual test: Open Apps & Features → find "Pomodoro Focus Timer" → Uninstall
- [x] T019 [US2] Manual test: Verify `C:\Program Files\Pomodoro Focus Timer\` directory is removed
- [x] T020 [US2] Manual test: Verify Start Menu and Desktop shortcuts are removed
- [x] T021 [US2] Manual test: Verify `~/.pomodoro-timer/` user data directory is still intact (if it existed before uninstall)

**Checkpoint**: Clean uninstall verified — no orphaned files, user data preserved.

---

## Phase 5: User Story 3 — Automated Build Pipeline (Priority: P2)

**Goal**: Developer runs single command to produce `.msi` with correct metadata. CI/CD ready.

**Independent Test**: Run `build-msi.ps1` from clean checkout → MSI created with correct name, version, manufacturer.

### Implementation for User Story 3

- [x] T022 [US3] Verify `build-msi.ps1` works from a clean state (no prior build artifacts)
- [x] T023 [US3] Verify `build-msi.ps1 -Version "2.0.0"` correctly overrides version in MSI metadata
- [x] T024 [US3] Manual test: Right-click MSI → Properties → Details → verify Product Name, Product Version, Manufacturer fields
- [x] T025 [US3] Verify MSI file size is under 20 MB (SC-006)
- [x] T026 [US3] Verify build completes in under 5 minutes (SC-005)

**Checkpoint**: Build pipeline validated — reproducible, correct metadata, size constraints met.

---

## Phase 6: User Story 4 — Upgrade Over Existing Installation (Priority: P2)

**Goal**: Installing a newer MSI upgrades the existing installation without manual uninstall, preserving user data.

**Independent Test**: Install v1.0.0 → create sessions → install v1.1.0 → old version gone, app updated, data intact.

### Implementation for User Story 4

- [x] T027 [US4] Verify `Product.wxs` MajorUpgrade element is correctly configured (UpgradeCode matches, Schedule="afterInstallInitialize")
- [x] T028 [US4] Build MSI with version 1.0.0, install it, and create timer sessions to generate data in `~/.pomodoro-timer/`
- [x] T029 [US4] Change version to 1.1.0 in `wails.json`, rebuild MSI
- [x] T030 [US4] Manual test: Run v1.1.0 MSI → verify it auto-removes v1.0.0 without prompting for uninstall
- [x] T031 [US4] Manual test: After upgrade, verify app launches and all previous session data is intact in `~/.pomodoro-timer/`
- [x] T032 [US4] Manual test: Verify Apps & Features shows only v1.1.0 (no duplicate entries)

**Checkpoint**: Major upgrade works — seamless version transition with data preservation.

---

## Phase 7: Polish & Cross-Cutting Concerns

**Purpose**: Documentation, cleanup, and final validation.

- [x] T033 [P] Update `README.md` to add MSI build instructions section (prerequisites, build command, output)
- [x] T034 [P] Update `quickstart.md` with final verified paths and commands
- [x] T035 Run full quickstart.md validation (fresh build → install → run → uninstall → upgrade cycle)
- [x] T036 Update constitution `Technology Stack & Constraints` to mention MSI installer alongside NSIS

---

## Dependencies & Execution Order

### Phase Dependencies

- **Setup (Phase 1)**: No dependencies — can start immediately
- **Foundational (Phase 2)**: Depends on Setup completion — BLOCKS all user stories
- **User Story 1 (Phase 3)**: Depends on Phase 2 — core MVP
- **User Story 2 (Phase 4)**: Depends on Phase 3 (needs installed app to test uninstall)
- **User Story 3 (Phase 5)**: Depends on Phase 3 (build script already exists, validates pipeline)
- **User Story 4 (Phase 6)**: Depends on Phase 3 (needs working MSI to test upgrade)
- **Polish (Phase 7)**: Depends on all user stories being validated

### User Story Dependencies

- **US1 (P1)**: Start after Phase 2 — no other story dependencies
- **US2 (P1)**: Requires US1 complete (needs installed app to test uninstall)
- **US3 (P2)**: Can run after US1 (validates the build pipeline created in US1)
- **US4 (P2)**: Can run after US1 (needs working MSI to test upgrade flow)

### Within Each User Story

- Implementation tasks before manual testing tasks
- WiX definition verification before build execution
- Build execution before manual install testing

### Parallel Opportunities

- T004 + T005 (Phase 1): version field and LICENSE file are independent
- T033 + T034 (Phase 7): README and quickstart updates are independent
- US3 and US4 can run in parallel after US1 is complete

---

## Parallel Example: Phase 1

```
# These two tasks can be done simultaneously:
Task T004: Add version info to wails.json
Task T005: Create LICENSE file at repo root
```

## Parallel Example: After US1 Complete

```
# These two stories can run in parallel once US1 is validated:
Story US3: Automated Build Pipeline (validate build script)
Story US4: Upgrade Over Existing Installation (test major upgrade)
```

---

## Implementation Strategy

### MVP First (User Story 1 Only)

1. Complete Phase 1: Setup (T001–T005)
2. Complete Phase 2: Foundational (T006–T008)
3. Complete Phase 3: User Story 1 (T009–T015)
4. **STOP and VALIDATE**: Install via MSI, run app, verify shortcuts
5. If MSI installs and app works → MVP achieved ✅

### Incremental Delivery

1. Setup + Foundational → Tools and WiX definition ready
2. User Story 1 → MSI installs correctly → **MVP!**
3. User Story 2 → Clean uninstall verified
4. User Stories 3 + 4 (parallel) → Build pipeline + upgrade verified
5. Polish → Documentation updated, constitution amended
6. Feature complete ✅

---

## Notes

- This feature is entirely build tooling — no Go or frontend code changes
- All testing is manual (install/uninstall/upgrade on Windows)
- The UpgradeCode GUID in `Product.wxs` must be generated ONCE and NEVER changed across versions
- WebView2 bootstrapper should be downloaded from Microsoft's official URL and cached locally
- MSI file naming convention: `PomodoroFocusTimer-{version}-x64.msi`
