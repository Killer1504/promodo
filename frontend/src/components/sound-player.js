/**
 * SoundPlayer — Web Audio API synthesized chimes.
 * No external files needed. Works inside WebView2.
 *
 * Usage:
 *   import { SoundPlayer } from './sound-player.js';
 *   const player = new SoundPlayer();
 *   player.play('focus');     // 3-note ascending bell
 *   player.play('break');     // soft single bell
 *   player.play('longBreak'); // 3-note descending bell
 *   player.muted = true;      // silence
 */
export class SoundPlayer {
    constructor() {
        this._ctx = null;
        this.muted = false;
    }

    _getContext() {
        if (!this._ctx) {
            this._ctx = new (window.AudioContext || window.webkitAudioContext)();
        }
        // Resume if suspended (browser autoplay policy)
        if (this._ctx.state === 'suspended') {
            this._ctx.resume();
        }
        return this._ctx;
    }

    /**
     * Play a bell note at a given frequency.
     * @param {number} freq - Hz
     * @param {number} startTime - AudioContext time offset (seconds)
     * @param {number} duration - seconds to sustain
     * @param {number} gain - 0..1 volume
     */
    _bellNote(freq, startTime, duration = 1.2, gain = 0.4) {
        const ctx = this._getContext();
        const now = ctx.currentTime + startTime;

        // Oscillator — sine wave for bell body
        const osc = ctx.createOscillator();
        osc.type = 'sine';
        osc.frequency.setValueAtTime(freq, now);
        // Slight frequency drop for natural bell decay
        osc.frequency.exponentialRampToValueAtTime(freq * 0.98, now + duration);

        // Gain envelope — sharp attack, exponential decay
        const gainNode = ctx.createGain();
        gainNode.gain.setValueAtTime(0, now);
        gainNode.gain.linearRampToValueAtTime(gain, now + 0.01);
        gainNode.gain.exponentialRampToValueAtTime(0.001, now + duration);

        // A tiny harmonic overtone for richness
        const osc2 = ctx.createOscillator();
        osc2.type = 'sine';
        osc2.frequency.setValueAtTime(freq * 2.0, now); // octave up
        const gainNode2 = ctx.createGain();
        gainNode2.gain.setValueAtTime(0, now);
        gainNode2.gain.linearRampToValueAtTime(gain * 0.15, now + 0.005);
        gainNode2.gain.exponentialRampToValueAtTime(0.001, now + duration * 0.5);

        osc.connect(gainNode);
        gainNode.connect(ctx.destination);
        osc2.connect(gainNode2);
        gainNode2.connect(ctx.destination);

        osc.start(now);
        osc.stop(now + duration);
        osc2.start(now);
        osc2.stop(now + duration);
    }

    /**
     * Play one of the predefined chime patterns.
     * @param {'focus'|'break'|'longBreak'} type
     */
    play(type) {
        if (this.muted) return;

        try {
            switch (type) {
                case 'focus':
                    // 3-note ascending bell: C5 → E5 → G5
                    this._bellNote(523, 0.0, 1.4, 0.35);
                    this._bellNote(659, 0.3, 1.4, 0.35);
                    this._bellNote(784, 0.6, 1.8, 0.40);
                    break;

                case 'break':
                    // Single warm mid bell: G4
                    this._bellNote(392, 0.0, 1.6, 0.35);
                    break;

                case 'longBreak':
                    // 3-note descending: G5 → E5 → C5
                    this._bellNote(784, 0.0, 1.2, 0.35);
                    this._bellNote(659, 0.3, 1.2, 0.35);
                    this._bellNote(523, 0.6, 1.8, 0.40);
                    break;

                default:
                    this._bellNote(523, 0.0, 1.2, 0.3);
            }
        } catch (err) {
            console.warn('Sound playback failed:', err);
        }
    }
}
