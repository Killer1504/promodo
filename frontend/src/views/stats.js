/**
 * Stats View
 * Displays today's summary and a weekly bar chart.
 */
import { Chart, BarController, BarElement, CategoryScale, LinearScale, Tooltip } from '../../node_modules/chart.js/dist/chart.js';

// Register Chart.js components
Chart.register(BarController, BarElement, CategoryScale, LinearScale, Tooltip);

export class StatsView {
    constructor(container) {
        this.container = container;
        this.chart = null;
        this.render();
    }

    render() {
        this.container.innerHTML = `
      <div class="stats-view">
        <div class="stats-summary">
          <div class="stat-card">
            <div class="stat-value" id="stat-sessions">0</div>
            <div class="stat-label">Sessions</div>
          </div>
          <div class="stat-card">
            <div class="stat-value" id="stat-minutes">0</div>
            <div class="stat-label">Focus Min</div>
          </div>
        </div>
        <div class="stats-chart-section">
          <div class="stats-chart-header">
            <span class="stats-chart-title">This Week</span>
            <span class="stats-chart-subtitle" id="stats-week-range"></span>
          </div>
          <div class="stats-chart-container">
            <canvas id="weekly-chart"></canvas>
          </div>
        </div>
      </div>
    `;
    }

    async onActivate() {
        await this.loadStats();
    }

    async loadStats() {
        try {
            const app = window.go?.main?.App;
            if (!app) return;

            const [today, weekly] = await Promise.all([
                app.GetTodayStats(),
                app.GetWeeklyStats(),
            ]);

            this.updateSummary(today);
            this.updateChart(weekly || []);
        } catch (err) {
            console.error('Failed to load stats:', err);
        }
    }

    updateSummary(today) {
        const sessionsEl = document.getElementById('stat-sessions');
        const minutesEl = document.getElementById('stat-minutes');

        if (sessionsEl && today) {
            sessionsEl.textContent = today.totalSessions || 0;
        }
        if (minutesEl && today) {
            minutesEl.textContent = today.totalFocusMinutes || 0;
        }
    }

    updateChart(weekly) {
        const canvas = document.getElementById('weekly-chart');
        if (!canvas) return;

        const dayNames = ['Mon', 'Tue', 'Wed', 'Thu', 'Fri', 'Sat', 'Sun'];

        // Map data to labels
        const labels = weekly.map(d => {
            const date = new Date(d.date + 'T00:00:00');
            return dayNames[date.getDay() === 0 ? 6 : date.getDay() - 1];
        });

        const data = weekly.map(d => d.totalFocusMinutes || 0);
        const todayIndex = weekly.findIndex(d => d.isToday);

        const colors = data.map((_, i) =>
            i === todayIndex
                ? getComputedStyle(document.documentElement).getPropertyValue('--chart-bar-today').trim()
                : getComputedStyle(document.documentElement).getPropertyValue('--chart-bar').trim()
        );

        // Update week range subtitle
        const rangeEl = document.getElementById('stats-week-range');
        if (rangeEl && weekly.length >= 2) {
            const start = this.formatShortDate(weekly[0].date);
            const end = this.formatShortDate(weekly[weekly.length - 1].date);
            rangeEl.textContent = `${start} – ${end}`;
        }

        // Destroy old chart
        if (this.chart) {
            this.chart.destroy();
        }

        const textColor = getComputedStyle(document.documentElement).getPropertyValue('--text-tertiary').trim();
        const gridColor = getComputedStyle(document.documentElement).getPropertyValue('--chart-grid').trim();

        this.chart = new Chart(canvas, {
            type: 'bar',
            data: {
                labels,
                datasets: [{
                    data,
                    backgroundColor: colors,
                    borderRadius: 6,
                    borderSkipped: false,
                    barPercentage: 0.6,
                }],
            },
            options: {
                responsive: true,
                maintainAspectRatio: false,
                plugins: {
                    legend: { display: false },
                    tooltip: {
                        callbacks: {
                            label: (ctx) => `${ctx.raw} min`,
                        },
                    },
                },
                scales: {
                    y: {
                        beginAtZero: true,
                        ticks: {
                            color: textColor,
                            font: { size: 11 },
                            stepSize: 30,
                        },
                        grid: { color: gridColor },
                        border: { display: false },
                    },
                    x: {
                        ticks: {
                            color: textColor,
                            font: { size: 11 },
                        },
                        grid: { display: false },
                        border: { display: false },
                    },
                },
            },
        });
    }

    formatShortDate(dateStr) {
        const d = new Date(dateStr + 'T00:00:00');
        return d.toLocaleDateString('en-US', { month: 'short', day: 'numeric' });
    }
}
