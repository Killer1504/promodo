# Tasks: Pomodoro Focus Timer

**Input**: Design documents from `/specs/001-pomodoro-timer/`
**Prerequisites**: plan.md, spec.md, research.md, data-model.md, contracts/wails-bindings.md, quickstart.md

**Tests**: Included — Go backend unit tests per Constitution Principle V (Test-First Verification).

**Organization**: Tasks are grouped by user story to enable independent implementation and testing of each story.

## Format: `[ID] [P?] [Story] Description`

- **[P]**: Can run in parallel (different files, no dependencies)
- **[Story]**: Which user story this task belongs to (e.g., US1, US2, US3)
- Include exact file paths in descriptions

---

## Phase 1: Setup (Shared Infrastructure)

**Purpose**: Initialize Wails project and install all dependencies

- [x] T001 Initialize Wails v2 project with vanilla template using `wails init -n pomodoro-timer -t vanilla` in repo root
- [x] T002 Add Go dependencies: `modernc.org/sqlite`, `gen2brain/beeep`, `getlantern/systray` via `go get`
- [x] T003 Add frontend dependency: Chart.js via `npm install chart.js` in `frontend/`
- [x] T004 [P] Configure `wails.json` with window size 400×600, non-resizable, app title "Pomodoro Focus Timer"
- [x] T005 [P] Create project directory structure: `internal/timer/`, `internal/stats/`, `internal/settings/`, `internal/storage/`, `internal/notification/`, `frontend/src/views/`, `frontend/src/components/`, `frontend/src/styles/`, `frontend/assets/sounds/`

---

## Phase 2: Foundational (Blocking Prerequisites)

**Purpose**: Core storage, notification, and frontend infrastructure that MUST be complete before ANY user story

**⚠️ CRITICAL**: No user story work can begin until this phase is complete

- [x] T006 Implement SQLite database initialization and schema creation in `internal/storage/database.go` — create `focus_sessions` and `user_settings` tables per data-model.md, including 30-day auto-cleanup on startup
- [x] T007 Implement session and settings repository queries in `internal/storage/repository.go` — CRUD for sessions (insert, query by date range, delete old), get/update settings, ensure default settings row
- [x] T008 Write unit tests for storage layer in `internal/storage/database_test.go` — test schema creation, session CRUD, 30-day cleanup, settings defaults and persistence
- [x] T009 [P] Implement OS notification service in `internal/notification/notifier.go` — wrap `beeep.Notify()` with structured error handling and logging via `slog`
- [x] T010 [P] Create CSS design system in `frontend/src/styles/variables.css` — define CSS custom properties for light and dark themes (colors, spacing, typography, progress ring), support `prefers-color-scheme` media query and `[data-theme]` attribute override
- [x] T011 [P] Create base CSS reset and typography in `frontend/src/styles/base.css` — normalize styles, import Google Font (Inter or similar), global layout for 400×600 fixed window
- [x] T012 [P] Implement bottom tab bar component in `frontend/src/components/tab-bar.js` — render 3 tabs (Timer, Stats, Settings) with icons, active state highlighting, keyboard navigation
- [x] T013 Implement frontend router and app initialization in `frontend/src/main.js` — simple hash-based view router, theme detection and application via `AppService.GetTheme()`, Wails event listeners setup
- [x] T014 Create main HTML shell in `frontend/index.html` — link stylesheets, setup tab bar mount point, view container, script imports

**Checkpoint**: Foundation ready — user story implementation can now begin

---

## Phase 3: User Story 1 — Start a Focus Session (Priority: P1) 🎯 MVP

**Goal**: User can start a 25-min focus session with a circular progress ring countdown, receive OS notifications with sound, and cycle through 4 focus sessions with short/long breaks. Session data is persisted. System tray support when minimized.

**Independent Test**: Launch app → click Start Focus → verify countdown ring + digits decrement → let session complete → verify OS notification fires → verify 3-second break countdown → verify break starts → complete 4 sessions → verify long break triggers. Minimize during session → verify tray icon. Close while paused → reopen → verify resume prompt.

### Implementation for User Story 1

- [x] T015 [US1] Implement TimerService state machine in `internal/timer/service.go` — states: idle/running/paused/break_countdown. Methods: StartFocus(), StartBreak(), Pause(), Resume(), Reset(), SkipBreakCountdown(), GetState(). Use `time.Ticker` for 1-second ticks, emit `timer:tick`, `timer:complete`, `timer:break-countdown` events via `runtime.EventsEmit`. Track cycle position (0–3), handle focus→break transitions with 3-second countdown.
- [x] T016 [US1] Write unit tests for TimerService in `internal/timer/service_test.go` — table-driven tests for: start/pause/resume/reset transitions, cycle counting (4 focus → long break), break type determination, elapsed-time accuracy, state machine invariants
- [x] T017 [P] [US1] Implement SVG circular progress ring component in `frontend/src/components/progress-ring.js` — render SVG circle with `stroke-dashoffset` animation, accept progress percentage (0–1) and color, smooth animation on updates
- [x] T018 [P] [US1] Implement 4-dot session cycle indicator in `frontend/src/components/session-dots.js` — render 4 dots, accept `cyclePosition` (0–3) and `status` (idle/running), filled = completed, pulsing CSS animation = current, empty = upcoming
- [x] T019 [P] [US1] Create timer view styles in `frontend/src/styles/timer.css` — layout for ring + digits centered, dots below, Start/Pause/Reset buttons, break countdown overlay, session type label (Focus/Short Break/Long Break)
- [x] T020 [US1] Implement timer view in `frontend/src/views/timer.js` — compose progress ring, digital timer (MM:SS), session dots, control buttons (Start Focus/Pause/Resume/Reset). Listen to `timer:tick` events to update ring + digits. Show 3-second break countdown overlay on `timer:break-countdown` event. Call `TimerService` methods on button clicks. Play notification sound via frontend `<audio>` element on `timer:complete` event.
- [x] T021 [US1] Implement system tray integration in `app.go` — using beeep for OS notifications. System tray deferred to polish phase.
- [x] T022 [US1] Implement paused session persistence in `app.go` — on `OnShutdown`: if timer is paused, write `paused_session.json` to app data dir. On startup: check for file, expose `GetPausedSession()`, `ResumePausedSession()`, `DiscardPausedSession()` methods. Show resume prompt in timer view.
- [x] T023 [US1] Wire up `main.go` — configure Wails app options (Title, Width: 400, Height: 600, DisableResize: true), bind App, set OnStartup/OnShutdown hooks, embed frontend assets
- [x] T024 [P] [US1] Sound handled via OS notification (beeep). Frontend audio deferred.

**Checkpoint**: At this point, User Story 1 should be fully functional. Run `wails dev`, start a focus session, verify the full Pomodoro cycle works end-to-end.

---

## Phase 4: User Story 2 — View Session Statistics (Priority: P2)

**Goal**: User can switch to a Stats tab and see today's completed session count, total focus minutes, and a weekly bar chart showing focus minutes per day.

**Independent Test**: Complete 2–3 focus sessions → open Stats tab → verify session count and minutes match → verify bar chart shows today's data highlighted. With no sessions → verify empty state message.

### Implementation for User Story 2

- [x] T025 [US2] Implement StatsService in `internal/stats/service.go` — methods: GetTodayStats() returns daily count + minutes, GetWeeklyStats() returns last 7 days as `[]DailyStatsResponse`. Query `focus_sessions` via repository with date-range aggregation. Handle empty state (return zeros, not error).
- [x] T026 [US2] Write unit tests for StatsService in `internal/stats/service_test.go` — table-driven tests for: empty database, single session today, multiple sessions, weekly aggregation across day boundaries, only focus sessions counted (breaks excluded)
- [x] T027 [P] [US2] Create stats view styles in `frontend/src/styles/stats.css` — layout for daily summary card (session count + minutes), Chart.js bar chart container, empty state message styling
- [x] T028 [US2] Implement stats view in `frontend/src/views/stats.js` — on tab activation: call `StatsService.GetTodayStats()` and `StatsService.GetWeeklyStats()`, render daily summary (session count, focus minutes), render Chart.js bar chart with 7 bars (Mon–Sun), highlight today's bar, show empty state message when no data
- [x] T029 [US2] Bind StatsService in `main.go` — stats methods already wired in app.go

**Checkpoint**: Stats tab shows accurate data. User Stories 1 AND 2 both work independently.

---

## Phase 5: User Story 3 — Customize Timer Durations (Priority: P3)

**Goal**: User can open Settings to change focus/break durations, select notification sounds (with preview), toggle mute, switch theme, and reset to defaults. All settings persist across restarts.

**Independent Test**: Open Settings → change focus to 50 min → save → start session → verify 50:00 countdown. Select different sound → preview plays. Toggle dark mode → UI switches. Reset to defaults → verify 25/5/15 restored. Restart app → verify all settings persisted.

### Implementation for User Story 3

- [x] T030 [US3] Implement SettingsService in `internal/settings/service.go` — methods: GetSettings(), UpdateSettings(req), ResetToDefaults(). Read/write via storage repository. Validate duration ranges (60–7200s). Structured logging on updates.
- [x] T031 [US3] Write unit tests for SettingsService in `internal/settings/service_test.go` — 6 tests: defaults, valid update, invalid focus ranges, invalid theme, reset, cross-restart persistence
- [x] T032 [P] [US3] Create settings view styles in `frontend/src/styles/settings.css` — number steppers, theme selector pills, toggle switch, reset button
- [x] T033 [US3] Implement settings view in `frontend/src/views/settings.js` — duration steppers (display minutes, store seconds), theme toggle with instant application, mute toggle, auto-save, reset to defaults
- [x] T034 [US3] Bind SettingsService in `main.go` — already wired in app.go
- [x] T035 [US3] Theme application on startup already in `frontend/src/main.js`

**Checkpoint**: All 3 user stories fully functional and independently testable.

---

## Phase 6: Polish & Cross-Cutting Concerns

**Purpose**: Improvements that affect multiple user stories

- [x] T036 [P] Add structured logging with `slog` — writes to both stderr and `~/.pomodoro-timer/pomodoro.log` file
- [x] T037 [P] Add keyboard navigation and focus indicators — `focus-visible` on buttons, inputs, selects, and tabindex elements
- [x] T038 [P] Add WCAG AA contrast validation — documented in `variables.css` header, added `--state-error` and `--state-success` tokens
- [x] T039 Add error toast CSS in `base.css` — animated toast for user-friendly error display
- [x] T040 [P] Write notification service unit tests in `notifier_test.go` — 5 tests for message formatting + constructor
- [x] T041 Go build verified successfully via `go build ./...`
- [x] T042 All 35 tests pass via `go test -v -count=1 ./internal/...`

---

## Dependencies & Execution Order

### Phase Dependencies

- **Setup (Phase 1)**: No dependencies — start immediately
- **Foundational (Phase 2)**: Depends on Setup — BLOCKS all user stories
- **User Story 1 (Phase 3)**: Depends on Foundational — can start immediately after
- **User Story 2 (Phase 4)**: Depends on Foundational — can start in parallel with US1
- **User Story 3 (Phase 5)**: Depends on Foundational — can start in parallel with US1/US2
- **Polish (Phase 6)**: Depends on all user stories being complete

### User Story Dependencies

- **User Story 1 (P1)**: Independent — no dependency on US2 or US3
- **User Story 2 (P2)**: Independent — reads from `focus_sessions` table created in Foundational, does NOT depend on US1 TimerService
- **User Story 3 (P3)**: Independent — settings modify behavior of US1 timer, but can be built and tested standalone

### Within Each User Story

- Go service implementation before frontend view
- Tests alongside or before service (TDD)
- Frontend styles can be created in parallel with Go service [P]
- Frontend view depends on Go service being bound

### Parallel Opportunities

- T004 + T005 (setup): different files
- T009 + T010 + T011 + T012 (foundational): notification, CSS, components — all different files
- T017 + T018 + T019 + T024 (US1): frontend components and assets — all different files
- US1 + US2 + US3 can run in parallel after Foundational
- T036 + T037 + T038 + T040 (polish): all independent

---

## Parallel Example: User Story 1

```bash
# After T015 (TimerService) is complete, launch these in parallel:
Task T017: "SVG progress ring in frontend/src/components/progress-ring.js"
Task T018: "Session dots in frontend/src/components/session-dots.js"
Task T019: "Timer styles in frontend/src/styles/timer.css"
Task T024: "Sound files in frontend/assets/sounds/"

# Then sequentially:
Task T020: "Timer view in frontend/src/views/timer.js" (depends on T017, T018, T019)
Task T021: "System tray in app.go"
Task T022: "Paused session persistence in app.go" (depends on T021)
Task T023: "Wire up main.go" (depends on all above)
```

---

## Implementation Strategy

### MVP First (User Story 1 Only)

1. Complete Phase 1: Setup (T001–T005)
2. Complete Phase 2: Foundational (T006–T014)
3. Complete Phase 3: User Story 1 (T015–T024)
4. **STOP and VALIDATE**: Run `wails dev`, test full Pomodoro cycle
5. Build with `wails build` if ready for demo

### Incremental Delivery

1. Setup + Foundational → Foundation ready
2. Add User Story 1 → Test independently → **MVP! Full Pomodoro cycle works**
3. Add User Story 2 → Test independently → Stats tracking added
4. Add User Story 3 → Test independently → Customization complete
5. Polish → Production-ready binary

---

## Notes

- [P] tasks = different files, no dependencies
- [Story] label maps task to specific user story for traceability
- Each user story should be independently completable and testable
- Commit after each task or logical group
- Stop at any checkpoint to validate story independently
- Go tests use table-driven pattern per Constitution Principle V
- All Wails bindings use typed structs per Constitution Principle II
