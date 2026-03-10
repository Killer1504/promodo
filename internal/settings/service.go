package settings

import (
	"fmt"
	"log/slog"

	"pomodoro-timer/internal/storage"
)

// Service manages user settings with validation.
type Service struct {
	repo *storage.Repository
}

// NewService creates a new settings Service.
func NewService(repo *storage.Repository) *Service {
	return &Service{repo: repo}
}

// GetSettings returns current user settings.
func (s *Service) GetSettings() (storage.UserSettings, error) {
	return s.repo.GetSettings()
}

// UpdateSettings validates and persists new settings.
func (s *Service) UpdateSettings(req storage.UserSettings) error {
	if err := validate(req); err != nil {
		return err
	}

	if err := s.repo.UpdateSettings(req); err != nil {
		return fmt.Errorf("save settings: %w", err)
	}

	slog.Info("settings updated",
		"focus", req.FocusDuration,
		"shortBreak", req.ShortBreakDuration,
		"longBreak", req.LongBreakDuration,
		"theme", req.Theme,
	)
	return nil
}

// ResetToDefaults restores all settings to defaults.
func (s *Service) ResetToDefaults() (storage.UserSettings, error) {
	if err := s.repo.ResetSettings(); err != nil {
		return storage.UserSettings{}, fmt.Errorf("reset settings: %w", err)
	}

	slog.Info("settings reset to defaults")
	return s.repo.GetSettings()
}

// validate checks settings values are within acceptable ranges.
// Durations are stored in seconds.
func validate(s storage.UserSettings) error {
	if s.FocusDuration < 60 || s.FocusDuration > 7200 {
		return fmt.Errorf("focus duration must be 60–7200 seconds, got %d", s.FocusDuration)
	}
	if s.ShortBreakDuration < 60 || s.ShortBreakDuration > 3600 {
		return fmt.Errorf("short break must be 60–3600 seconds, got %d", s.ShortBreakDuration)
	}
	if s.LongBreakDuration < 60 || s.LongBreakDuration > 3600 {
		return fmt.Errorf("long break must be 60–3600 seconds, got %d", s.LongBreakDuration)
	}
	if s.SessionsBeforeLongBreak < 1 || s.SessionsBeforeLongBreak > 10 {
		return fmt.Errorf("sessions before long break must be 1–10, got %d", s.SessionsBeforeLongBreak)
	}

	validThemes := map[string]bool{"light": true, "dark": true, "system": true}
	if !validThemes[s.Theme] {
		return fmt.Errorf("invalid theme %q, must be light/dark/system", s.Theme)
	}

	return nil
}
