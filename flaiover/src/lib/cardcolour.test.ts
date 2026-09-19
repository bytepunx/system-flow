import { describe, expect, it } from 'vitest';
import { readFileSync } from 'node:fs';
import { natureTint, stripeFor, tintFor, typeStripe, typeSwatch } from './cardcolour';

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
	});
});
