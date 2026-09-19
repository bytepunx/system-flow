import { describe, expect, it } from 'vitest';
import { readFileSync } from 'node:fs';
import { natureTint, stripeClass, stripeFor, tintFor, typeStripe, typeSwatch } from './cardcolour';
import { CATEGORICAL, NATURE_SLOT } from './viz/palette';

describe('the card colour maps', () => {
	const css = readFileSync(new URL('../routes/layout.css', import.meta.url), 'utf8');

	it('use only colours that are theme tokens in both themes', () => {
		const names = [
			...Object.values(natureTint).map((c) => c.replace(/^bg-/, '')),
			...Object.values(typeStripe).map((c) => c.replace(/^border-l-/, '')),
			...Object.values(typeSwatch).map((c) => c.replace(/^bg-/, ''))
		];
		for (const name of names) {
			expect(css, name).toContain(`--color-${name}: var(--t-${name});`);
			expect(css.split(`--t-${name}: #`).length - 1, name).toBe(2);
		}
	});

	it('give a type the same colour on the card and in the legend', () => {
		for (const [type, stripe] of Object.entries(typeStripe))
			expect(typeSwatch[type]).toBe(stripe.replace(/^border-l-/, 'bg-'));
	});

	it('fall back to the plain card for values the schema does not know', () => {
		expect(tintFor('feature')).toBe('bg-nature-feature');
		expect(tintFor('chore')).toBe('bg-ground');
		expect(stripeFor('task')).toBe('border-l-4 border-l-type-task');
		expect(stripeFor('theme')).toBe('hover:border-l-line-strong');
		expect(stripeClass('epic')).toBe('border-l-4 border-l-type-epic');
		expect(stripeClass('theme')).toBe('');
	});
});

// A chart mark needs a saturated colour, which a pastel cannot be, so the charts keep their
// palette; what they share with the cards is the hue a nature is known by.
describe('a nature keeps its hue between the board and the charts', () => {
	const css = readFileSync(new URL('../routes/layout.css', import.meta.url), 'utf8');
	const hue = (hex: string) => {
		const [r, g, b] = [1, 3, 5].map((i) => parseInt(hex.slice(i, i + 2), 16) / 255);
		const max = Math.max(r, g, b);
		const d = max - Math.min(r, g, b);
		const h = max === r ? ((g - b) / d) % 6 : max === g ? (b - r) / d + 2 : (r - g) / d + 4;
		return (h * 60 + 360) % 360;
	};
	const apart = (a: number, b: number) => Math.min(Math.abs(a - b), 360 - Math.abs(a - b));
	const tints = (nature: string) =>
		[...css.matchAll(new RegExp(`--t-nature-${nature}: (#[0-9a-f]{6})`, 'g'))].map((m) => m[1]);

	it.each(Object.keys(natureTint))('%s', (nature) => {
		const [light, dark] = tints(nature);
		const slot = NATURE_SLOT[nature];
		expect(slot, 'every coded nature has a chart slot').toBeTypeOf('number');
		expect(apart(hue(light), hue(CATEGORICAL.light[slot]))).toBeLessThanOrEqual(25);
		expect(apart(hue(dark), hue(CATEGORICAL.dark[slot]))).toBeLessThanOrEqual(25);
	});
});
