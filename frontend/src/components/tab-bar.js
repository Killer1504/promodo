/**
 * Tab Bar Component
 * Bottom navigation with 3 tabs: Timer, Stats, Settings
 */

const TAB_ITEMS = [
  { id: 'timer', label: 'Timer', icon: '⏱' },
  { id: 'stats', label: 'Stats', icon: '📊' },
  { id: 'settings', label: 'Settings', icon: '⚙️' },
];

export class TabBar {
  constructor(container, onChange) {
    this.container = container;
    this.onChange = onChange;
    this.activeTab = 'timer';
    this.render();
    this.bindEvents();
  }

  render() {
    this.container.innerHTML = `
      <nav class="tab-bar" role="tablist" aria-label="Main navigation">
        ${TAB_ITEMS.map(tab => `
          <button
            class="tab-item ${tab.id === this.activeTab ? 'active' : ''}"
            role="tab"
            id="tab-${tab.id}"
            aria-selected="${tab.id === this.activeTab}"
            aria-controls="view-${tab.id}"
            data-tab="${tab.id}"
            tabindex="${tab.id === this.activeTab ? '0' : '-1'}"
          >
            <span class="tab-icon">${tab.icon}</span>
            <span class="tab-label">${tab.label}</span>
          </button>
        `).join('')}
      </nav>
    `;
  }

  bindEvents() {
    this.container.addEventListener('click', (e) => {
      const btn = e.target.closest('[data-tab]');
      if (btn) this.setActive(btn.dataset.tab);
    });

    // Keyboard navigation: left/right arrows cycle tabs
    this.container.addEventListener('keydown', (e) => {
      const tabs = TAB_ITEMS.map(t => t.id);
      const idx = tabs.indexOf(this.activeTab);

      if (e.key === 'ArrowRight' || e.key === 'ArrowLeft') {
        e.preventDefault();
        const next = e.key === 'ArrowRight'
          ? (idx + 1) % tabs.length
          : (idx - 1 + tabs.length) % tabs.length;
        this.setActive(tabs[next]);
        this.container.querySelector(`[data-tab="${tabs[next]}"]`)?.focus();
      }
    });
  }

  setActive(tabId) {
    if (this.activeTab === tabId) return;
    this.activeTab = tabId;

    this.container.querySelectorAll('.tab-item').forEach(btn => {
      const isActive = btn.dataset.tab === tabId;
      btn.classList.toggle('active', isActive);
      btn.setAttribute('aria-selected', isActive);
      btn.setAttribute('tabindex', isActive ? '0' : '-1');
    });

    if (this.onChange) this.onChange(tabId);
  }

  getActive() {
    return this.activeTab;
  }
}
