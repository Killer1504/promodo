/**
 * Session Dots Component
 * Displays cycle position as 4 dots below the progress ring.
 * Filled dots = completed focus sessions in current cycle.
 * Current dot = pulsing animation.
 */

export class SessionDots {
    constructor(container) {
        this.container = container;
        this.count = 4;
        this.render(0);
    }

    render(cyclePosition) {
        const dots = [];
        for (let i = 0; i < this.count; i++) {
            let cls = 'session-dot';
            if (i < cyclePosition) {
                cls += ' filled';
            } else if (i === cyclePosition) {
                cls += ' current';
            }
            dots.push(`<span class="${cls}" aria-label="Session ${i + 1}${i < cyclePosition ? ' completed' : i === cyclePosition ? ' current' : ''}"></span>`);
        }

        this.container.innerHTML = `
      <div class="session-dots" role="group" aria-label="Focus sessions progress (${cyclePosition} of ${this.count})">
        ${dots.join('')}
      </div>
    `;
    }

    /**
     * Update the dots display.
     * @param {number} cyclePosition - 0-based index of current cycle
     * @param {number} [total] - Optional total dots (default 4)
     */
    update(cyclePosition, total) {
        if (total && total !== this.count) {
            this.count = total;
        }
        this.render(cyclePosition);
    }
}
