package main

import (
	"context"
	"encoding/json"
	"io"
	"log/slog"
	"os"
	"path/filepath"

	"pomodoro-timer/internal/notification"
	"pomodoro-timer/internal/stats"
	"pomodoro-timer/internal/storage"
	"pomodoro-timer/internal/timer"

	wailsRuntime "github.com/wailsapp/wails/v2/pkg/runtime"
)

// PausedSessionResponse is the frontend-facing paused session data.
type PausedSessionResponse struct {
	SessionType      string `json:"sessionType"`
	RemainingSeconds int    `json:"remainingSeconds"`
	CyclePosition    int    `json:"cyclePosition"`
	PausedAt         string `json:"pausedAt"`
}

// App is the main application struct bound to the frontend.
type App struct {
	ctx      context.Context
	db       *storage.Database
	repo     *storage.Repository
	timer    *timer.Service
	stats    *stats.Service
	notifier *notification.Notifier
	dataDir  string
	logFile  *os.File
}

// NewApp creates a new App.
func NewApp() *App {
	homeDir, _ := os.UserHomeDir()
	dataDir := filepath.Join(homeDir, ".pomodoro-timer")
	return &App{dataDir: dataDir}
}

// startup is called by Wails on app launch.
func (a *App) startup(ctx context.Context) {
	a.ctx = ctx

	// Ensure data directory exists
	os.MkdirAll(a.dataDir, 0o755)

	// T036: Set up structured file logging
	logPath := filepath.Join(a.dataDir, "pomodoro.log")
	f, err := os.OpenFile(logPath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0o644)
	if err == nil {
		a.logFile = f
		w := io.MultiWriter(os.Stderr, f)
		slog.SetDefault(slog.New(slog.NewTextHandler(w, &slog.HandlerOptions{Level: slog.LevelInfo})))
	}

	// Initialize database
	db, err := storage.NewDatabase(a.dataDir)
	if err != nil {
		slog.Error("failed to initialize database", "error", err)
		return
	}
	a.db = db
	a.repo = storage.NewRepository(db)

	// Initialize notification service
	a.notifier = notification.NewNotifier()

	// Initialize timer service
	a.timer = timer.NewService(a.repo, a.notifier)
	a.timer.SetContext(ctx)

	// Initialize stats service
	a.stats = stats.NewService(a.repo)

	// Start system tray (FR-011)
	go startTray(ctx)

	slog.Info("app started", "dataDir", a.dataDir)
}

// shutdown is called by Wails on app close.
func (a *App) shutdown(ctx context.Context) {
	// Save paused session if timer is running/paused
	if a.timer != nil {
		if a.timer.IsPaused() {
			a.savePausedSession()
		}
		a.timer.Shutdown()
	}

	if a.db != nil {
		a.db.Close()
	}

	if a.logFile != nil {
		a.logFile.Close()
	}

	slog.Info("app shutdown")
}

// QuitApp performs a clean exit — saves state then quits the process.
func (a *App) QuitApp() {
	wailsRuntime.Quit(a.ctx)
}

// ShowWindow restores the window if it was hidden via X button.
func (a *App) ShowWindow() {
	wailsRuntime.WindowShow(a.ctx)
}

// --- Timer bindings ---

// StartFocus starts a new focus session.
func (a *App) StartFocus() timer.TimerState {
	return a.timer.StartFocus()
}

// PauseTimer pauses the running timer.
func (a *App) PauseTimer() timer.TimerState {
	return a.timer.Pause()
}

// ResumeTimer resumes a paused timer.
func (a *App) ResumeTimer() timer.TimerState {
	return a.timer.Resume()
}

// ResetTimer resets the timer to idle.
func (a *App) ResetTimer() timer.TimerState {
	a.deletePausedSession()
	return a.timer.Reset()
}

// SkipBreakCountdown skips the pre-break countdown.
func (a *App) SkipBreakCountdown() timer.TimerState {
	return a.timer.SkipBreakCountdown()
}

// GetTimerState returns the current timer state.
func (a *App) GetTimerState() timer.TimerState {
	return a.timer.GetState()
}

// --- Settings bindings ---

// GetSettings returns user settings.
func (a *App) GetSettings() (*storage.UserSettings, error) {
	settings, err := a.repo.GetSettings()
	if err != nil {
		return nil, err
	}
	return &settings, nil
}

// UpdateSettings saves updated settings.
func (a *App) UpdateSettings(s storage.UserSettings) (*storage.UserSettings, error) {
	if err := a.repo.UpdateSettings(s); err != nil {
		return nil, err
	}
	updated, err := a.repo.GetSettings()
	if err != nil {
		return nil, err
	}
	return &updated, nil
}

// ResetToDefaults resets all settings.
func (a *App) ResetToDefaults() (*storage.UserSettings, error) {
	if err := a.repo.ResetSettings(); err != nil {
		return nil, err
	}
	settings, err := a.repo.GetSettings()
	if err != nil {
		return nil, err
	}
	return &settings, nil
}

// --- Stats bindings ---

// GetTodayStats returns today's focus stats.
func (a *App) GetTodayStats() stats.DailyStatsResponse {
	return a.stats.GetTodayStats()
}

// GetWeeklyStats returns 7 days of stats with gap-filling.
func (a *App) GetWeeklyStats() []stats.DailyStatsResponse {
	return a.stats.GetWeeklyStats()
}

// --- Theme ---

// GetTheme returns the resolved theme ("light" or "dark").
func (a *App) GetTheme() string {
	if a.repo == nil {
		return "system"
	}
	settings, err := a.repo.GetSettings()
	if err != nil {
		return "system"
	}
	return settings.Theme
}

// --- Paused Session Persistence (FR-015) ---

func (a *App) pausedSessionPath() string {
	return filepath.Join(a.dataDir, "paused_session.json")
}

func (a *App) savePausedSession() {
	ps := PausedSessionResponse{
		SessionType:      a.timer.GetSessionType(),
		RemainingSeconds: a.timer.GetRemainingSeconds(),
		CyclePosition:    a.timer.GetCyclePosition(),
		PausedAt:         "now",
	}

	data, err := json.Marshal(ps)
	if err != nil {
		slog.Error("marshal paused session", "error", err)
		return
	}

	if err := os.WriteFile(a.pausedSessionPath(), data, 0o644); err != nil {
		slog.Error("save paused session", "error", err)
		return
	}

	slog.Info("paused session saved", "remaining", ps.RemainingSeconds)
}

// GetPausedSession checks for a saved paused session.
func (a *App) GetPausedSession() *PausedSessionResponse {
	data, err := os.ReadFile(a.pausedSessionPath())
	if err != nil {
		return nil
	}

	var ps PausedSessionResponse
	if err := json.Unmarshal(data, &ps); err != nil {
		return nil
	}

	return &ps
}

// ResumePausedSession restores and resumes a saved paused session.
func (a *App) ResumePausedSession() timer.TimerState {
	ps := a.GetPausedSession()
	if ps == nil {
		return a.timer.GetState()
	}

	a.timer.RestoreState(ps.SessionType, ps.RemainingSeconds, ps.CyclePosition)
	a.deletePausedSession()
	return a.timer.Resume()
}

// DiscardPausedSession deletes the saved paused session.
func (a *App) DiscardPausedSession() bool {
	a.deletePausedSession()
	return true
}

func (a *App) deletePausedSession() {
	os.Remove(a.pausedSessionPath())
}
