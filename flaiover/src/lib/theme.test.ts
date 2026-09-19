// The token pairs the theme relies on, with their WCAG ratios (S-0044).
// Values mirror src/routes/layout.css; change both together.
import { describe, expect, it } from 'vitest';
import { readFileSync } from 'node:fs';

function lum(hex: string): number {
	const c = hex.replace('#', '');
	const [r, g, b] = [0, 2, 4].map((i) => parseInt(c.slice(i, i + 2), 16) / 255);
	const f = (v: number) => (v <= 0.03928 ? v / 12.92 : ((v + 0.055) / 1.055) ** 2.4);
	return 0.2126 * f(r) + 0.7152 * f(g) + 0.0722 * f(b);
}
export function contrast(a: string, b: string): number {
	const [hi, lo] = [lum(a), lum(b)].sort((x, y) => y - x);
	return (hi + 0.05) / (lo + 0.05);
}

/** Parse the token blocks out of layout.css so the test cannot drift from the stylesheet. */
function tokens(css: string, selector: string): Record<string, string> {
	const start = css.indexOf(selector + ' {');
	const end = css.indexOf('}', start);
	const out: Record<string, string> = {};
	for (const m of css.slice(start, end).matchAll(/--t-([a-z-]+):\s*(#[0-9a-f]{6})/g))
		out[m[1]] = m[2];
	return out;
}
const css = readFileSync(new URL('../routes/layout.css', import.meta.url), 'utf8');
const light = tokens(css, ':root');
const dark = tokens(css, "[data-theme='dark']");

// [foreground, background, minimum] per theme. 4.5 for text, 3 for large text and controls.
const textPairs: [string, string][] = [
	['ink', 'ground'],
	['ink', 'surface'],
	['ink-soft', 'ground'],
	['ink-soft', 'surface'],
	['muted', 'ground'],
	['muted', 'surface'],
	['accent', 'ground'],
	['accent', 'surface'],
	['on-primary', 'primary'],
	['good', 'good-soft'],
	['warn', 'warn-soft'],
	['danger', 'danger-soft'],
	['info', 'info-soft'],
	['good', 'ground'],
	['warn', 'ground'],
	['danger', 'ground'],
	['info', 'ground']
];
const controlPairs: [string, string][] = [
	['primary', 'ground'],
	['accent', 'ground'],
	['line-strong', 'ground'],
	['line-strong', 'surface']
];

// Board cards (S-0055): the title, the details, and the BLOCKED flag sit on the nature tint.
// Telling the pastels apart is waived, the card says both in text; reading on them is not.
// The type colours are an edge stripe and a legend swatch, and carry no text.
const natures = ['feature', 'improvement', 'remediation', 'research', 'experiment'];
const types = ['epic', 'story', 'task'];
const cardPairs: [string, string][] = natures.flatMap((n): [string, string][] => [
	['ink', `nature-${n}`],
	['muted', `nature-${n}`],
	['danger', `nature-${n}`]
]);

describe.each([
	['light', light],
	['dark', dark]
] as const)('%s theme', (_name, t) => {
	it('has every token', () => {
		for (const k of [
			'ground',
			'surface',
			'raised',
			'ink',
			'ink-soft',
			'muted',
			'line',
			'line-strong',
			'primary',
			'on-primary',
			'accent',
			'accent-strong',
			'good',
			'warn',
			'danger',
			'info',
			...natures.map((n) => `nature-${n}`),
			...types.map((t) => `type-${t}`)
		]) {
			expect(t[k], k).toMatch(/^#[0-9a-f]{6}$/);
		}
	});
	it.each(textPairs)('%s on %s reads at AA (4.5:1)', (fg, bg) => {
		expect(contrast(t[fg], t[bg])).toBeGreaterThanOrEqual(4.5);
	});
	it.each(cardPairs)('%s on %s reads at AA (4.5:1) on a board card', (fg, bg) => {
		expect(contrast(t[fg], t[bg])).toBeGreaterThanOrEqual(4.5);
	});
	it('gives every nature and every type a colour of its own', () => {
		expect(new Set(natures.map((n) => t[`nature-${n}`])).size).toBe(natures.length);
		expect(new Set(types.map((k) => t[`type-${k}`])).size).toBe(types.length);
	});
	it.each(controlPairs)('%s against %s clears 3:1 for controls', (fg, bg) => {
		expect(contrast(t[fg], t[bg])).toBeGreaterThanOrEqual(3);
	});
});

describe('brand colours keep their roles', () => {
	it('uses the five brand values where contrast allows', () => {
		expect(light.ground).toBe('#f7f3e3');
		expect(dark.ink).toBe('#f7f3e3');
		expect(light.muted).toBe('#645853');
		expect(light.primary).toBe('#054a91');
		expect(dark.accent).toBe('#23b5d3');
		expect(light['accent-strong']).toBe('#23b5d3');
		expect(light['good-strong']).toBe('#119822');
		expect(dark['good-strong']).toBe('#119822');
	});
	it('documents why the brand green is a fill, not text', () => {
		// #119822 reads at 3.4:1 on the cream ground; text uses #0c6e19, the brand green fills badges with ink text.
		expect(contrast('#119822', light.ground)).toBeLessThan(4.5);
		expect(contrast(light.ink, '#119822')).toBeGreaterThanOrEqual(4.5);
	});
	it('documents why the light accent text is a darker step of the cyan', () => {
		// #23b5d3 on the light ground is about 2.2:1, so text and focus use #17788c.
		expect(contrast('#23b5d3', light.ground)).toBeLessThan(3);
		expect(contrast(light.accent, light.ground)).toBeGreaterThanOrEqual(4.5);
	});
});
