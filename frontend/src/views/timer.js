/**
 * Timer View
 * Main timer interface — progress ring, session dots, and controls.
 */
import { ProgressRing } from '../components/progress-ring.js';
import { SessionDots } from '../components/session-dots.js';
import { SoundPlayer } from '../components/sound-player.js';

export class TimerView {
    constructor(container) {
        this.container = container;
        this.sound = new SoundPlayer();
        this.state = {
            status: 'idle',
            sessionType: 'focus',
            remainingSeconds: 1500,
            totalSeconds: 1500,
            cyclePosition: 0,
        };
        this.render();
        this.bindEvents();
        this.loadInitialState();
    }

    render() {
        this.container.innerHTML = `
      <div class="timer-view">
        <div id="progress-ring-container"></div>
        <div id="session-dots-container"></div>
        <div id="timer-controls" class="timer-controls"></div>
      </div>
    `;

        this.ring = new ProgressRing(document.getElementById('progress-ring-container'));
        this.dots = new SessionDots(document.getElementById('session-dots-container'));
        this.controlsEl = document.getElementById('timer-controls');

        this.updateControls();
    }

    bindEvents() {
        // Timer tick
        if (window.runtime) {
            window.runtime.EventsOn('timer:tick', (state) => {
                this.state = state;
                this.updateDisplay();
            });

            window.runtime.EventsOn('timer:complete', (data) => {
                this.handleSessionComplete(data);
            });

            window.runtime.EventsOn('timer:break-countdown', (data) => {
                this.state.status = 'break_countdown';
                this.state.remainingSeconds = data.secondsRemaining;
                this.updateDisplay();
            });
        }

        // Controls delegation
        this.container.addEventListener('click', (e) => {
            const btn = e.target.closest('[data-action]');
            if (btn) this.handleAction(btn.dataset.action);
        });
    }

    async loadInitialState() {
        const app = await this._waitForApp();
        if (!app) return;
        try {
            const state = await app.GetTimerState();
            if (state) {
                this.state = state;
                this.updateDisplay();
            }
        } catch (err) {
            console.error('Failed to load timer state:', err);
        }
    }

    /** Waits up to 2s for window.go.main.App to be bound by Wails. */
    _waitForApp(retries = 20, delayMs = 100) {
        return new Promise((resolve) => {
            const check = (n) => {
                const app = window.go?.main?.App;
                if (app) return resolve(app);
                if (n <= 0) return resolve(null);
                setTimeout(() => check(n - 1), delayMs);
            };
            check(retries);
        });
    }

    async handleAction(action) {
        try {
            let state;
            const app = window.go?.main?.App;
            if (!app) return;

            switch (action) {
                case 'start':
                    state = await app.StartFocus();
                    break;
                case 'pause':
                    state = await app.PauseTimer();
                    break;
                case 'resume':
                    state = await app.ResumeTimer();
                    break;
                case 'reset':
                    state = await app.ResetTimer();
                    break;
                case 'skip':
                    state = await app.SkipBreakCountdown();
                    break;
            }

            if (state) {
                this.state = state;
                this.updateDisplay();
            }
        } catch (err) {
            console.error(`Timer action failed: ${action}`, err);
        }
    }

    async handleSessionComplete(data) {
        // Play the appropriate chime
        const type = data?.sessionType || 'focus';
        await this._syncMute();
        if (type === 'focus') {
            this.sound.play('focus');
        } else if (data?.nextBreakType === 'long') {
            this.sound.play('longBreak');
        } else {
            this.sound.play('break');
        }
        // Refresh timer state
        this.loadInitialState();
    }

    async _syncMute() {
        try {
            const settings = await window.go?.main?.App?.GetSettings();
            if (settings) this.sound.muted = !!settings.mute;
        } catch { /* keep current mute state */ }
    }

    updateDisplay() {
        const { status, sessionType, remainingSeconds, totalSeconds, cyclePosition } = this.state;

        // Update progress ring
        if (status === 'idle') {
            this.ring.setIdle(totalSeconds || 1500);
        } else if (status === 'break_countdown') {
            // Show countdown number instead of ring
        } else {
            this.ring.update(remainingSeconds, totalSeconds, sessionType);
        }

        // Update session dots
        this.dots.update(cyclePosition);

        // Update controls
        this.updateControls();

        // Apply paused blinking
        const timeEl = this.container.querySelector('.progress-ring-time');
        if (timeEl) {
            timeEl.classList.toggle('paused', status === 'paused');
        }
    }

    updateControls() {
        const { status } = this.state;

        switch (status) {
            case 'idle':
                this.controlsEl.innerHTML = `
          <button class="btn-timer-primary" data-action="start">▶ Start Focus</button>
        `;
                break;

            case 'running':
                this.controlsEl.innerHTML = `
          <button class="btn-timer-secondary" data-action="reset" title="Reset">↺</button>
          <button class="btn-timer-primary" data-action="pause">⏸ Pause</button>
        `;
                break;

            case 'paused':
                this.controlsEl.innerHTML = `
          <button class="btn-timer-secondary" data-action="reset" title="Reset">↺</button>
          <button class="btn-timer-primary" data-action="resume">▶ Resume</button>
        `;
                break;

            case 'break_countdown':
                this.controlsEl.innerHTML = `
          <div class="break-countdown">
            <span class="break-countdown-number">${this.state.remainingSeconds}</span>
            <span class="break-countdown-text">Break starting...</span>
            <button class="btn-skip" data-action="skip">Skip →</button>
          </div>
        `;
                break;
        }
    }

    onActivate() {
        this.loadInitialState();
    }
}
