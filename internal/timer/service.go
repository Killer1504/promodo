package timer

import (
	"context"
	"fmt"
	"log/slog"
	"sync"
	"time"

	wailsRuntime "github.com/wailsapp/wails/v2/pkg/runtime"

	"pomodoro-timer/internal/notification"
	"pomodoro-timer/internal/storage"
)

// Status constants
const (
	StatusIdle           = "idle"
	StatusRunning        = "running"
	StatusPaused         = "paused"
	StatusBreakCountdown = "break_countdown"
)

// Session type constants
const (
	SessionFocus      = "focus"
	SessionShortBreak = "short_break"
	SessionLongBreak  = "long_break"
)

// TimerState is the typed response sent to the frontend.
type TimerState struct {
	Status           string `json:"status"`
	SessionType      string `json:"sessionType"`
	RemainingSeconds int    `json:"remainingSeconds"`
	TotalSeconds     int    `json:"totalSeconds"`
	CyclePosition    int    `json:"cyclePosition"`
}

// SessionComplete is emitted when a session finishes.
type SessionComplete struct {
	SessionType     string `json:"sessionType"`
	DurationSeconds int    `json:"durationSeconds"`
	NextBreakType   string `json:"nextBreakType"`
}

// BreakCountdown is emitted during the 3-second pre-break.
type BreakCountdown struct {
	SecondsRemaining int `json:"secondsRemaining"`
}

// Service manages the Pomodoro timer lifecycle.
type Service struct {
	mu sync.Mutex

	ctx      context.Context
	repo     *storage.Repository
	notifier *notification.Notifier

	status       string
	sessionType  string
	remaining    int
	total        int
	cycle        int // 0-3: which focus session we're on
	sessionStart time.Time

	ticker *time.Ticker
	done   chan struct{}

	// Settings cache (refreshed on each start)
	focusDuration   int
	shortBreakDur   int
	longBreakDur    int
	sessionsForLong int
}

// NewService creates a new timer Service.
func NewService(repo *storage.Repository, notifier *notification.Notifier) *Service {
	return &Service{
		repo:            repo,
		notifier:        notifier,
		status:          StatusIdle,
		sessionType:     SessionFocus,
		cycle:           0,
		focusDuration:   1500,
		shortBreakDur:   300,
		longBreakDur:    900,
		sessionsForLong: 4,
	}
}

// SetContext stores the Wails runtime context for event emission.
func (s *Service) SetContext(ctx context.Context) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.ctx = ctx
}

// refreshSettings loads the latest durations from the database.
func (s *Service) refreshSettings() {
	settings, err := s.repo.GetSettings()
	if err != nil {
		slog.Warn("failed to load settings, using cached values", "error", err)
		return
	}
	s.focusDuration = settings.FocusDuration
	s.shortBreakDur = settings.ShortBreakDuration
	s.longBreakDur = settings.LongBreakDuration
	s.sessionsForLong = settings.SessionsBeforeLongBreak
}

// durationForType returns the duration in seconds for the given session type.
func (s *Service) durationForType(sessionType string) int {
	switch sessionType {
	case SessionFocus:
		return s.focusDuration
	case SessionShortBreak:
		return s.shortBreakDur
	case SessionLongBreak:
		return s.longBreakDur
	default:
		return s.focusDuration
	}
}

// StartFocus begins a new focus session.
func (s *Service) StartFocus() TimerState {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.stopTickerLocked()
	s.refreshSettings()

	s.sessionType = SessionFocus
	s.total = s.focusDuration
	s.remaining = s.total
	s.status = StatusRunning
	s.sessionStart = time.Now()

	s.startTickerLocked()
	slog.Info("focus session started", "duration", s.total, "cycle", s.cycle)
	return s.stateLocked()
}

// StartBreak begins a break session.
func (s *Service) StartBreak(breakType string) TimerState {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.stopTickerLocked()

	switch breakType {
	case "short":
		s.sessionType = SessionShortBreak
		s.total = s.shortBreakDur
	case "long":
		s.sessionType = SessionLongBreak
		s.total = s.longBreakDur
	default:
		s.sessionType = SessionShortBreak
		s.total = s.shortBreakDur
	}

	s.remaining = s.total
	s.status = StatusRunning
	s.sessionStart = time.Now()

	s.startTickerLocked()
	slog.Info("break started", "type", s.sessionType, "duration", s.total)
	return s.stateLocked()
}

// Pause pauses the active timer.
func (s *Service) Pause() TimerState {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.status != StatusRunning {
		return s.stateLocked()
	}

	s.stopTickerLocked()
	s.status = StatusPaused
	slog.Info("timer paused", "remaining", s.remaining)
	return s.stateLocked()
}

// Resume resumes a paused timer.
func (s *Service) Resume() TimerState {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.status != StatusPaused {
		return s.stateLocked()
	}

	s.status = StatusRunning
	s.startTickerLocked()
	slog.Info("timer resumed", "remaining", s.remaining)
	return s.stateLocked()
}

// Reset stops the timer and returns to idle.
func (s *Service) Reset() TimerState {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.stopTickerLocked()
	s.status = StatusIdle
	s.remaining = 0
	s.total = 0
	slog.Info("timer reset")
	return s.stateLocked()
}

// SkipBreakCountdown skips the 3-second pre-break countdown.
func (s *Service) SkipBreakCountdown() TimerState {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.status != StatusBreakCountdown {
		return s.stateLocked()
	}

	s.stopTickerLocked()
	nextBreak := s.nextBreakTypeLocked()
	s.mu.Unlock()

	state := s.StartBreak(nextBreak)

	s.mu.Lock()
	return state
}

// GetState returns the current timer state.
func (s *Service) GetState() TimerState {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.stateLocked()
}

// GetCyclePosition returns the current cycle position (0-based).
func (s *Service) GetCyclePosition() int {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.cycle
}

// GetRemainingSeconds returns seconds left.
func (s *Service) GetRemainingSeconds() int {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.remaining
}

// GetSessionType returns the current session type.
func (s *Service) GetSessionType() string {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.sessionType
}

// GetStatus returns the current status.
func (s *Service) GetStatus() string {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.status
}

// RestoreState restores timer state from a paused session (for app reopen).
func (s *Service) RestoreState(sessionType string, remaining int, cycle int) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.sessionType = sessionType
	s.remaining = remaining
	s.total = s.durationForType(sessionType)
	s.cycle = cycle
	s.status = StatusPaused
}

// stateLocked builds a TimerState (must hold mu).
func (s *Service) stateLocked() TimerState {
	return TimerState{
		Status:           s.status,
		SessionType:      s.sessionType,
		RemainingSeconds: s.remaining,
		TotalSeconds:     s.total,
		CyclePosition:    s.cycle,
	}
}

// nextBreakTypeLocked determines the next break type based on cycle (must hold mu).
func (s *Service) nextBreakTypeLocked() string {
	if s.cycle+1 >= s.sessionsForLong {
		return "long"
	}
	return "short"
}

func (s *Service) startTickerLocked() {
	s.done = make(chan struct{})
	s.ticker = time.NewTicker(1 * time.Second)

	done := s.done     // capture locally for goroutine safety
	ticker := s.ticker // capture locally

	go func() {
		for {
			select {
			case <-done:
				return
			case <-ticker.C:
				s.tick()
			}
		}
	}()
}

func (s *Service) stopTickerLocked() {
	if s.ticker != nil {
		s.ticker.Stop()
		s.ticker = nil
	}
	if s.done != nil {
		select {
		case <-s.done:
			// already closed
		default:
			close(s.done)
		}
		s.done = nil
	}
}

func (s *Service) tick() {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.status == StatusBreakCountdown {
		s.remaining--
		if s.ctx != nil {
			wailsRuntime.EventsEmit(s.ctx, "timer:break-countdown", BreakCountdown{
				SecondsRemaining: s.remaining,
			})
		}
		if s.remaining <= 0 {
			s.stopTickerLocked()
			nextBreak := s.nextBreakTypeLocked()
			s.mu.Unlock()
			s.StartBreak(nextBreak)
			s.mu.Lock()
		}
		return
	}

	if s.status != StatusRunning {
		return
	}

	s.remaining--

	// Emit tick event
	if s.ctx != nil {
		wailsRuntime.EventsEmit(s.ctx, "timer:tick", s.stateLocked())
	}

	// Session complete
	if s.remaining <= 0 {
		s.stopTickerLocked()
		s.onSessionComplete()
	}
}

func (s *Service) onSessionComplete() {
	completedType := s.sessionType
	duration := s.total

	// Persist completed session
	session := storage.FocusSession{
		SessionType:     completedType,
		StartTime:       s.sessionStart,
		EndTime:         time.Now(),
		DurationSeconds: duration,
		Completed:       true,
	}
	if _, err := s.repo.InsertSession(session); err != nil {
		slog.Error("failed to persist session", "error", err)
	}

	// Send OS notification (respects mute setting)
	title, msg := notification.SessionCompleteMessage(completedType)
	go func() {
		muted := false
		if s.repo != nil {
			if settings, err := s.repo.GetSettings(); err == nil {
				muted = settings.Mute
			}
		}
		_ = s.notifier.Alert(title, msg, muted)
	}()

	// Determine next break type
	nextBreak := s.nextBreakTypeLocked()

	// Emit completion event
	if s.ctx != nil {
		wailsRuntime.EventsEmit(s.ctx, "timer:complete", SessionComplete{
			SessionType:     completedType,
			DurationSeconds: duration,
			NextBreakType:   nextBreak,
		})
	}

	if completedType == SessionFocus {
		// Advance cycle
		s.cycle++
		if s.cycle >= s.sessionsForLong {
			// After long break, reset cycle
		}

		// Start 3-second break countdown
		s.status = StatusBreakCountdown
		s.remaining = 3
		s.total = 3
		s.startTickerLocked()
	} else {
		// Break completed — reset to idle
		if completedType == SessionLongBreak {
			s.cycle = 0 // Reset cycle after long break
		}
		s.status = StatusIdle
		s.sessionType = SessionFocus
		s.remaining = 0
		s.total = 0
	}

	slog.Info("session complete",
		"type", completedType,
		"duration", duration,
		"cycle", s.cycle,
		"nextBreak", nextBreak,
	)
}

// Shutdown cleanly stops the timer.
func (s *Service) Shutdown() {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.stopTickerLocked()
}

// IsPaused returns true if the timer is currently paused.
func (s *Service) IsPaused() bool {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.status == StatusPaused
}

// FormatRemaining returns "MM:SS" string for the remaining time.
func FormatRemaining(seconds int) string {
	m := seconds / 60
	sec := seconds % 60
	return fmt.Sprintf("%02d:%02d", m, sec)
}
