/**
 * Settings View
 * Duration steppers, theme toggle, mute switch, and reset.
 */
export class SettingsView {
    constructor(container) {
        this.container = container;
        this.settings = null;
        this.render();
    }

    render() {
        this.container.innerHTML = `
      <div class="settings-view">
        <div class="settings-section">
          <div class="settings-section-title">Durations</div>
          <div class="setting-row">
            <div>
              <div class="setting-label">Focus</div>
              <div class="setting-sublabel">Work session length</div>
            </div>
            <div class="number-stepper" data-field="focusDuration" data-min="1" data-max="120">
              <button data-step="-5">−</button>
              <span class="stepper-value" id="val-focus">25</span>
              <span class="stepper-unit">min</span>
              <button data-step="5">+</button>
            </div>
          </div>
          <div class="setting-row">
            <div>
              <div class="setting-label">Short Break</div>
              <div class="setting-sublabel">Between sessions</div>
            </div>
            <div class="number-stepper" data-field="shortBreakDuration" data-min="1" data-max="60">
              <button data-step="-1">−</button>
              <span class="stepper-value" id="val-short">5</span>
              <span class="stepper-unit">min</span>
              <button data-step="1">+</button>
            </div>
          </div>
          <div class="setting-row">
            <div>
              <div class="setting-label">Long Break</div>
              <div class="setting-sublabel">After full cycle</div>
            </div>
            <div class="number-stepper" data-field="longBreakDuration" data-min="1" data-max="60">
              <button data-step="-5">−</button>
              <span class="stepper-value" id="val-long">15</span>
              <span class="stepper-unit">min</span>
              <button data-step="5">+</button>
            </div>
          </div>
          <div class="setting-row">
            <div>
              <div class="setting-label">Sessions / Cycle</div>
              <div class="setting-sublabel">Before long break</div>
            </div>
            <div class="number-stepper" data-field="sessionsBeforeLongBreak" data-min="1" data-max="10">
              <button data-step="-1">−</button>
              <span class="stepper-value" id="val-sessions">4</span>
              <span class="stepper-unit"></span>
              <button data-step="1">+</button>
            </div>
          </div>
        </div>

        <div class="settings-section">
          <div class="settings-section-title">Appearance</div>
          <div class="setting-row">
            <div class="setting-label">Theme</div>
            <div class="theme-selector">
              <button class="theme-option" data-theme="light">☀ Light</button>
              <button class="theme-option" data-theme="dark">🌙 Dark</button>
              <button class="theme-option active" data-theme="system">⚙ System</button>
            </div>
          </div>
          <div class="setting-row">
            <div class="setting-label">Mute Notifications</div>
            <label class="toggle-switch">
              <input type="checkbox" id="mute-toggle">
              <span class="toggle-slider"></span>
            </label>
          </div>
        </div>

        <div class="settings-reset">
          <button class="btn-reset" id="btn-reset-defaults">Reset to Defaults</button>
        </div>
      </div>
    `;

        this.bindEvents();
    }

    bindEvents() {
        // Number steppers
        this.container.querySelectorAll('.number-stepper button').forEach(btn => {
            btn.addEventListener('click', () => {
                const stepper = btn.closest('.number-stepper');
                const field = stepper.dataset.field;
                const min = parseInt(stepper.dataset.min);
                const max = parseInt(stepper.dataset.max);
                const step = parseInt(btn.dataset.step);
                const valueEl = stepper.querySelector('.stepper-value');
                const current = parseInt(valueEl.textContent);
                const next = Math.min(max, Math.max(min, current + step));
                valueEl.textContent = next;
                this.saveSettings();
            });
        });

        // Theme selector
        this.container.querySelectorAll('.theme-option').forEach(btn => {
            btn.addEventListener('click', () => {
                this.container.querySelectorAll('.theme-option').forEach(b => b.classList.remove('active'));
                btn.classList.add('active');
                this.applyTheme(btn.dataset.theme);
                this.saveSettings();
            });
        });

        // Mute toggle
        const muteToggle = document.getElementById('mute-toggle');
        if (muteToggle) {
            muteToggle.addEventListener('change', () => this.saveSettings());
        }

        // Reset
        const resetBtn = document.getElementById('btn-reset-defaults');
        if (resetBtn) {
            resetBtn.addEventListener('click', () => this.resetDefaults());
        }
    }

    async onActivate() {
        await this.loadSettings();
    }

    async loadSettings() {
        try {
            const app = window.go?.main?.App;
            if (!app) return;

            const s = await app.GetSettings();
            if (!s) return;
            this.settings = s;

            // DB stores seconds — display as minutes
            document.getElementById('val-focus').textContent = Math.round((s.focusDuration || 1500) / 60);
            document.getElementById('val-short').textContent = Math.round((s.shortBreakDuration || 300) / 60);
            document.getElementById('val-long').textContent = Math.round((s.longBreakDuration || 900) / 60);
            document.getElementById('val-sessions').textContent = s.sessionsBeforeLongBreak || 4;

            // Theme
            const theme = s.theme || 'system';
            this.container.querySelectorAll('.theme-option').forEach(btn => {
                btn.classList.toggle('active', btn.dataset.theme === theme);
            });

            // Mute
            const muteToggle = document.getElementById('mute-toggle');
            if (muteToggle) muteToggle.checked = !!s.mute;
        } catch (err) {
            console.error('Failed to load settings:', err);
        }
    }

    async saveSettings() {
        try {
            const app = window.go?.main?.App;
            if (!app) return;

            const activeTheme = this.container.querySelector('.theme-option.active');
            const muteToggle = document.getElementById('mute-toggle');

            // Convert minutes (displayed) to seconds (stored)
            const settings = {
                focusDuration: parseInt(document.getElementById('val-focus').textContent) * 60,
                shortBreakDuration: parseInt(document.getElementById('val-short').textContent) * 60,
                longBreakDuration: parseInt(document.getElementById('val-long').textContent) * 60,
                sessionsBeforeLongBreak: parseInt(document.getElementById('val-sessions').textContent),
                notificationSound: this.settings?.notificationSound || 'bell',
                mute: muteToggle?.checked || false,
                theme: activeTheme?.dataset.theme || 'system',
            };

            await app.UpdateSettings(settings);
        } catch (err) {
            console.error('Failed to save settings:', err);
        }
    }

    async resetDefaults() {
        try {
            const app = window.go?.main?.App;
            if (!app) return;

            await app.ResetToDefaults();
            await this.loadSettings();
        } catch (err) {
            console.error('Failed to reset settings:', err);
        }
    }

    applyTheme(theme) {
        if (theme === 'system') {
            document.documentElement.removeAttribute('data-theme');
        } else {
            document.documentElement.dataset.theme = theme;
        }
    }
}
