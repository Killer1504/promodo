/**
 * Main Application Entry Point
 * Handles routing, theme detection, and event listener setup.
 */
import { TabBar } from './components/tab-bar.js';

class App {
    constructor() {
        this.currentView = null;
        this.views = {};
        this.init();
    }

    async init() {
        // Initialize theme
        await this.applyTheme();

        // Mount tab bar
        const tabBarEl = document.getElementById('tab-bar');
        this.tabBar = new TabBar(tabBarEl, (tabId) => this.navigate(tabId));

        // Navigate to Timer view (default)
        this.navigate('timer');

        // Listen for settings changes to update theme live
        if (window.runtime) {
            window.runtime.EventsOn('settings:changed', () => {
                this.applyTheme();
            });
        }
    }

    async applyTheme() {
        let theme = 'system';

        // Try to get theme from backend
        try {
            if (window.go && window.go.main && window.go.main.App) {
                theme = await window.go.main.App.GetTheme();
            }
        } catch {
            // Fallback to system detection
        }

        if (theme === 'system' || (!theme)) {
            // Auto-detect from OS
            const prefersDark = window.matchMedia('(prefers-color-scheme: dark)').matches;
            theme = prefersDark ? 'dark' : 'light';
        }

        document.documentElement.setAttribute('data-theme', theme);
    }

    navigate(viewId) {
        // Hide all views
        document.querySelectorAll('.view').forEach(v => v.classList.remove('active'));

        // Show target view
        const viewEl = document.getElementById(`view-${viewId}`);
        if (viewEl) {
            viewEl.classList.add('active');
        }

        this.currentView = viewId;

        // Lazy-load view module and call its onActivate if available
        this.activateView(viewId);
    }

    async activateView(viewId) {
        if (!this.views[viewId]) {
            try {
                switch (viewId) {
                    case 'timer': {
                        const mod = await import('./views/timer.js');
                        if (mod.TimerView) {
                            this.views[viewId] = new mod.TimerView(document.getElementById('view-timer'));
                        }
                        break;
                    }
                    case 'stats': {
                        const mod = await import('./views/stats.js');
                        if (mod.StatsView) {
                            this.views[viewId] = new mod.StatsView(document.getElementById('view-stats'));
                        }
                        break;
                    }
                    case 'settings': {
                        const mod = await import('./views/settings.js');
                        if (mod.SettingsView) {
                            this.views[viewId] = new mod.SettingsView(document.getElementById('view-settings'));
                        }
                        break;
                    }
                }
            } catch (err) {
                console.error(`Failed to load view: ${viewId}`, err);
            }
        }

        // Call onActivate hook if the view has one
        if (this.views[viewId]?.onActivate) {
            this.views[viewId].onActivate();
        }
    }
}

// Boot
document.addEventListener('DOMContentLoaded', () => {
    new App();
});
