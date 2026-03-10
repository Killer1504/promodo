# Feature Specification: Pomodoro Focus Timer

**Feature Branch**: `001-pomodoro-timer`
**Created**: 2026-03-10
**Status**: Draft
**Input**: User description: "Pomodoro Focus Timer — a desktop timer with session tracking, daily/weekly stats, and OS notifications"

## Clarifications

### Session 2026-03-10

- Q: Should the app minimize to the system tray or the taskbar? → A: Minimize to system tray with a tray icon showing timer status.
- Q: Should the app play an audible alert when a session ends? → A: Yes, user-selectable sound from multiple built-in options, with a mute toggle.
- Q: How long should session history be retained? → A: Keep last 30 days; auto-delete older data.
- Q: Should the break timer auto-start after a focus session? → A: Auto-start after a 3-second skippable countdown.
- Q: What happens to a paused session when the app is closed and reopened? → A: Ask user on reopen: "Resume previous session or start fresh?"

### Session 2026-03-10 (UI/UX)

- Q: How should the user navigate between Timer, Stats, and Settings? → A: Bottom tab bar with 3 tabs (Timer, Stats, Settings).
- Q: How should the countdown timer be visually presented? → A: Circular progress ring with large digital digits inside.
- Q: Should the app support dark mode, light mode, or both? → A: Both, with system preference auto-detection and manual toggle in Settings.
- Q: How should the app show which session the user is on in the 4-session cycle? → A: Four dots below the timer (filled = completed, pulsing = current).
- Q: What should the default window size be? → A: Compact 400×600 px, non-resizable.

## User Scenarios & Testing *(mandatory)*

### User Story 1 - Start a Focus Session (Priority: P1)

A user opens the app and starts a 25-minute focus session. A circular progress ring with large digital digits appears, counting down. When the session ends, the system sends an OS-level notification and automatically suggests a 5-minute break. After 4 focus sessions, a longer 15-minute break is suggested.

**Why this priority**: The core Pomodoro loop is the entire purpose of the app. Without it, nothing else has value.

**Independent Test**: Launch the app → click Start → verify the countdown decrements correctly → verify notification fires at zero → verify break begins automatically.

**Acceptance Scenarios**:

1. **Given** the app is open and idle, **When** the user clicks "Start Focus", **Then** a 25-minute countdown begins with a visible timer display updating every second.
2. **Given** a focus session is running, **When** the timer reaches 00:00, **Then** an OS notification is shown, a notification sound plays, and a 3-second countdown appears ("Break starting in 3... 2... 1..."). The user can skip to start the break immediately or let it auto-start.
3. **Given** a focus session is running, **When** the user clicks "Pause", **Then** the timer stops and can be resumed or reset.
4. **Given** 4 consecutive focus sessions have completed, **When** the 4th session ends, **Then** the app suggests a 15-minute long break instead of a 5-minute short break. All four progress dots are filled.
5. **Given** a focus session is running, **When** the user looks at the timer view, **Then** four dots below the progress ring indicate which session they are on (filled = completed, pulsing = current, empty = upcoming).
6. **Given** a focus session is running, **When** the user minimizes the window, **Then** the app minimizes to the system tray with an icon displaying the remaining time, and can be restored by clicking the tray icon.
7. **Given** the app was closed while a session was paused, **When** the user reopens the app, **Then** a prompt appears asking "Resume previous session (X:XX remaining) or start fresh?" with two clear action buttons.

---

### User Story 2 - View Session Statistics (Priority: P2)

A user wants to see how many focus sessions they completed today and this week. They navigate to a stats view showing daily session counts, total focus time, and a simple weekly bar chart.

**Why this priority**: Tracking progress motivates continued use. Stats give the timer long-term value beyond a single session.

**Independent Test**: Complete 2–3 focus sessions → open the Stats view → verify session count, total minutes, and chart data match the completed sessions.

**Acceptance Scenarios**:

1. **Given** the user has completed 3 focus sessions today, **When** they open the Stats view, **Then** they see "3 sessions" and the correct total focus minutes for today.
2. **Given** the user has session data spanning the current week, **When** they view the weekly chart, **Then** a bar chart shows focus minutes per day with today highlighted.
3. **Given** the user has no session data, **When** they open the Stats view, **Then** they see an empty state message encouraging them to start their first session.

---

### User Story 3 - Customize Timer Durations (Priority: P3)

A user wants to change the default focus duration from 25 minutes to 50 minutes, the short break from 5 to 10, and the long break from 15 to 30. Changes persist across app restarts.

**Why this priority**: Personalization is important for adoption, but the app is fully functional with defaults. This is an enhancement, not a core requirement.

**Independent Test**: Open Settings → change focus duration to 50 minutes → start a session → verify the timer starts at 50:00 → restart the app → verify settings persist.

**Acceptance Scenarios**:

1. **Given** the user opens Settings, **When** they change focus duration to 50 minutes and save, **Then** the next focus session starts with a 50-minute countdown.
2. **Given** the user has customized durations, **When** they restart the app, **Then** the custom durations are preserved and used for new sessions.
3. **Given** the user is in Settings, **When** they click "Reset to Defaults", **Then** all durations revert to 25/5/15 minutes.
4. **Given** the user is in Settings, **When** they browse notification sounds, **Then** they can preview and select from multiple built-in sounds, or mute notifications entirely.

---

### Edge Cases

- What happens when the user closes the app mid-session? If the session was paused, its state (remaining time, session type) is persisted. On reopen, the user is prompted to resume or discard. If the session was actively running (not paused), it is discarded and NOT counted as complete.
- What happens when the user's system clock changes (e.g., daylight saving)? The timer uses elapsed-time tracking, not wall-clock comparison, so it is unaffected.
- What happens when OS notifications are disabled? The app still transitions to break mode; notification failure is logged but does not block the timer flow.
- What happens when the stats database is corrupted or deleted? The app starts with a fresh empty database and shows the empty state. A warning is logged.
- What happens when session data older than 30 days exists? The system auto-deletes records older than 30 days on app launch. Stats view only shows data within the retention window.

## Requirements *(mandatory)*

### Functional Requirements

- **FR-001**: System MUST provide a countdown timer displayed as a circular progress ring with large digital digits (MM:SS) centered inside, updating every second.
- **FR-002**: System MUST send an OS-level notification when a focus session or break ends.
- **FR-003**: System MUST automatically transition between focus → short break → focus cycles, with a long break after every 4th focus session.
- **FR-004**: Users MUST be able to pause, resume, and reset an active timer.
- **FR-005**: System MUST persist completed session data (timestamp, duration, type) to local storage.
- **FR-006**: System MUST display daily session count and total focus minutes on a Stats view.
- **FR-007**: System MUST display a weekly bar chart showing focus minutes per day.
- **FR-008**: Users MUST be able to customize focus, short break, and long break durations.
- **FR-009**: Custom duration settings MUST persist across app restarts.
- **FR-010**: System MUST use elapsed-time tracking for the countdown (not wall-clock), ensuring accuracy regardless of system clock changes.
- **FR-011**: System MUST minimize to the system tray when the window is minimized, showing a tray icon with timer status. Clicking the tray icon MUST restore the window.
- **FR-012**: System MUST play a user-selectable notification sound when a session ends. Users MUST be able to choose from multiple built-in sounds and mute audio entirely.
- **FR-013**: System MUST auto-delete session data older than 30 days. Deletion MUST occur on app launch.
- **FR-014**: After a focus session ends, the system MUST display a 3-second countdown before auto-starting the break. Users MUST be able to skip the countdown to start the break immediately.
- **FR-015**: If a paused session exists when the app is closed, the system MUST persist timer state (remaining time, session type). On reopen, the system MUST prompt the user to resume or start fresh.
- **FR-016**: The app MUST use a bottom tab bar with three tabs: Timer (default), Stats, and Settings. The active tab MUST be visually highlighted.
- **FR-017**: The app MUST support light and dark themes. On first launch, the theme MUST auto-detect the OS preference. Users MUST be able to override the theme via a toggle in Settings.
- **FR-018**: The Timer view MUST display four progress dots below the circular timer ring, indicating position in the 4-session cycle (filled = completed, pulsing = current, empty = upcoming). Dots reset after a long break.
- **FR-019**: The app window MUST be a compact 400×600 px non-resizable window.

### Key Entities

- **FocusSession**: Represents a single completed focus or break period. Key attributes: start time, end time, duration, session type (focus/short-break/long-break), completion status.
- **UserSettings**: Stores user preferences. Key attributes: focus duration, short break duration, long break duration, sessions before long break count, selected notification sound, mute toggle, theme preference (light/dark/system).
- **DailyStats**: Aggregated view of sessions per day. Key attributes: date, total sessions, total focus minutes.

## Success Criteria *(mandatory)*

### Measurable Outcomes

- **SC-001**: Users can start a focus session within 2 seconds of opening the app (no setup required).
- **SC-002**: Timer display updates smoothly with no visible lag or skipped seconds during a session.
- **SC-003**: OS notification appears within 1 second of a session ending on supported platforms.
- **SC-004**: Stats view loads within 1 second with up to 30 days of session history.
- **SC-005**: Customized duration settings persist correctly across 100% of app restarts.
- **SC-006**: App memory usage stays below 100 MB during normal operation.
