/**
 * SVG Progress Ring Component
 * Circular countdown visualization for the timer.
 */

export class ProgressRing {
    constructor(container) {
        this.container = container;
        this.size = 220;
        this.strokeWidth = 8;
        this.radius = (this.size - this.strokeWidth) / 2;
        this.circumference = 2 * Math.PI * this.radius;
        this.render();
        this.progressCircle = this.container.querySelector('.progress-ring-circle');
        this.timeDisplay = this.container.querySelector('.progress-ring-time');
        this.labelDisplay = this.container.querySelector('.progress-ring-label');
    }

    render() {
        const cx = this.size / 2;
        const cy = this.size / 2;

        this.container.innerHTML = `
      <div class="progress-ring-wrapper">
        <svg class="progress-ring" width="${this.size}" height="${this.size}" viewBox="0 0 ${this.size} ${this.size}">
          <!-- Background ring -->
          <circle
            class="progress-ring-bg"
            cx="${cx}"
            cy="${cy}"
            r="${this.radius}"
            stroke="var(--ring-bg)"
            stroke-width="${this.strokeWidth}"
            fill="none"
          />
          <!-- Progress ring -->
          <circle
            class="progress-ring-circle"
            cx="${cx}"
            cy="${cy}"
            r="${this.radius}"
            stroke="var(--ring-focus)"
            stroke-width="${this.strokeWidth}"
            fill="none"
            stroke-linecap="round"
            stroke-dasharray="${this.circumference}"
            stroke-dashoffset="0"
            transform="rotate(-90 ${cx} ${cy})"
          />
        </svg>
        <div class="progress-ring-content">
          <span class="progress-ring-time">25:00</span>
          <span class="progress-ring-label">Focus</span>
        </div>
      </div>
    `;
    }

    /**
     * Update the ring's progress.
     * @param {number} remaining - Seconds remaining
     * @param {number} total - Total seconds for the session
     * @param {string} sessionType - 'focus', 'short_break', or 'long_break'
     */
    update(remaining, total, sessionType) {
        if (!this.progressCircle) return;

        // Update progress arc
        const progress = total > 0 ? remaining / total : 0;
        const offset = this.circumference * (1 - progress);
        this.progressCircle.style.strokeDashoffset = offset;

        // Update ring color based on session type
        const colorVar = this.getColorVar(sessionType);
        this.progressCircle.setAttribute('stroke', `var(${colorVar})`);

        // Update time display
        if (this.timeDisplay) {
            this.timeDisplay.textContent = this.formatTime(remaining);
        }

        // Update label
        if (this.labelDisplay) {
            this.labelDisplay.textContent = this.getLabel(sessionType);
        }
    }

    /**
     * Set the ring to idle state.
     * @param {number} defaultDuration - Default duration in seconds
     */
    setIdle(defaultDuration) {
        if (this.progressCircle) {
            this.progressCircle.style.strokeDashoffset = '0';
            this.progressCircle.setAttribute('stroke', 'var(--ring-focus)');
        }
        if (this.timeDisplay) {
            this.timeDisplay.textContent = this.formatTime(defaultDuration || 1500);
        }
        if (this.labelDisplay) {
            this.labelDisplay.textContent = 'Ready';
        }
    }

    getColorVar(sessionType) {
        switch (sessionType) {
            case 'short_break': return '--ring-short-break';
            case 'long_break': return '--ring-long-break';
            default: return '--ring-focus';
        }
    }

    getLabel(sessionType) {
        switch (sessionType) {
            case 'short_break': return 'Short Break';
            case 'long_break': return 'Long Break';
            default: return 'Focus';
        }
    }

    formatTime(seconds) {
        const m = Math.floor(seconds / 60);
        const s = seconds % 60;
        return `${String(m).padStart(2, '0')}:${String(s).padStart(2, '0')}`;
    }
}
