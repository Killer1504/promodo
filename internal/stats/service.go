package stats

import (
	"time"

	"pomodoro-timer/internal/storage"
)

// DailyStatsResponse is the frontend-facing daily stats payload.
type DailyStatsResponse struct {
	Date              string `json:"date"`
	TotalSessions     int    `json:"totalSessions"`
	TotalFocusMinutes int    `json:"totalFocusMinutes"`
	IsToday           bool   `json:"isToday"`
}

// Service provides session statistics.
type Service struct {
	repo *storage.Repository
}

// NewService creates a new stats Service.
func NewService(repo *storage.Repository) *Service {
	return &Service{repo: repo}
}

// GetTodayStats returns today's focus session count and total minutes.
func (s *Service) GetTodayStats() DailyStatsResponse {
	stats, err := s.repo.GetDailyStats(1)
	if err != nil || len(stats) == 0 {
		return DailyStatsResponse{
			Date:    today(),
			IsToday: true,
		}
	}

	todayStr := today()
	for _, ds := range stats {
		if ds.Date == todayStr {
			return DailyStatsResponse{
				Date:              ds.Date,
				TotalSessions:     ds.TotalSessions,
				TotalFocusMinutes: ds.TotalFocusMinutes,
				IsToday:           true,
			}
		}
	}

	return DailyStatsResponse{Date: todayStr, IsToday: true}
}

// GetWeeklyStats returns the last 7 days of daily stats, filling gaps with zeros.
func (s *Service) GetWeeklyStats() []DailyStatsResponse {
	raw, err := s.repo.GetDailyStats(7)
	if err != nil {
		raw = nil
	}

	// Build a lookup map
	lookup := make(map[string]storage.DailyStats)
	for _, ds := range raw {
		lookup[ds.Date] = ds
	}

	todayStr := today()
	result := make([]DailyStatsResponse, 7)
	for i := 6; i >= 0; i-- {
		d := time.Now().AddDate(0, 0, -(6 - i))
		dateStr := d.Format("2006-01-02")

		resp := DailyStatsResponse{
			Date:    dateStr,
			IsToday: dateStr == todayStr,
		}

		if ds, ok := lookup[dateStr]; ok {
			resp.TotalSessions = ds.TotalSessions
			resp.TotalFocusMinutes = ds.TotalFocusMinutes
		}

		result[i] = resp
	}

	return result
}

func today() string {
	return time.Now().Format("2006-01-02")
}
