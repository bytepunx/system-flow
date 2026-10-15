// S-0382: the design page's token sheet lays out a project's tokens as design.get answers them:
// colours as swatches with their contrast on the ground in both themes, flai's role pairs, the
// type scale as sample text, and spacing, radii, and shadows drawn to size.
import { afterEach, describe, expect, it } from 'vitest';
import { flushSync, mount, unmount } from 'svelte';
import TokenSheet from './TokenSheet.svelte';
import type { DesignContrast, DesignTheme, DesignToken } from '$lib/server/repo';

// Tokens as flai's design.get answers its starter design system (flai design show --json).
const colour = (name: string, hex: string, line = 1): DesignToken => ({
	path: `color.${name}`,
	type: 'color',
	value: { colorSpace: 'srgb', hex },
	css: hex,
	variables: [{ name: `--color-${name}`, value: hex }],
	line
});
const text = (step: string, size: number, lineHeight: number): DesignToken => ({
	path: `text.${step}`,
	type: 'typography',
	value: { fontSize: { unit: 'rem', value: size }, lineHeight },
	css: `${size}rem`,
	variables: [
		{ name: `--text-${step}`, value: `${size}rem` },
		{ name: `--text-${step}--line-height`, value: String(lineHeight) }
	],
	line: 1
});
const length = (path: string, rem: number): DesignToken => ({
	path,
	type: 'dimension',
	value: { unit: 'rem', value: rem },
	css: `${rem}rem`,
	variables: [{ name: `--${path.replace('.', '-')}`, value: `${rem}rem` }],
	line: 1
});
const shadow = (alpha: string): DesignToken => {
	const css = `0px 4px 6px -1px #000000${alpha}, 0px 2px 4px -2px #000000${alpha}`;
	return {
		path: 'shadow.md',
		type: 'shadow',
		value: [],
		css,
		variables: [{ name: '--shadow-md', value: css }],
		line: 1
	};
};
const SANS = 'ui-sans-serif, system-ui, "Segoe UI", sans-serif';
const font: DesignToken = {
	path: 'font.sans',
	type: 'fontFamily',
	value: ['ui-sans-serif', 'system-ui', 'Segoe UI', 'sans-serif'],
	css: SANS,
	variables: [{ name: '--font-sans', value: SANS }],
	line: 1
};
// An extra colour in oklch, far outside sRGB, the same in both themes.
const neon: DesignToken = {
	path: 'color.brand-neon',
	type: 'color',
	value: { colorSpace: 'oklch', components: [0.7, 0.4, 150] },
	css: 'oklch(0.7 0.4 150)',
	variables: [{ name: '--color-brand-neon', value: 'oklch(0.7 0.4 150)' }],
	line: 1
};
const shared = [font, text('lg', 1.125, 1.556), text('2xl', 1.5, 1.333), length('spacing', 0.25)];

const LIGHT: DesignTheme = {
	name: 'light',
	tokens: [
		colour('background', '#ffffff'),
		colour('surface', '#f8fafc'),
		colour('text', '#0f172a'),
		colour('border', '#848e9f'),
		colour('primary', '#2563eb'),
		neon,
		...shared,
		length('radius.md', 0.375),
		shadow('1a')
	]
};
const DARK: DesignTheme = {
	name: 'dark',
	tokens: [
		colour('background', '#0f172a'),
		colour('surface', '#1e293b'),
		colour('text', '#f1f5f9'),
		colour('border', '#334155'),
		colour('primary', '#60a5fa'),
		neon,
		...shared,
		length('radius.md', 0.5),
		shadow('66')
	]
};
const pair = (
	theme: 'light' | 'dark',
	foreground: string,
	background: string,
	ratio: number,
	threshold: number
) => ({
	theme,
	foreground,
	background,
	foreground_hex: '#0f172a',
	background_hex: '#ffffff',
	ratio,
	threshold,
	passes: ratio >= threshold,
	clipped: false
});
const CONTRAST: DesignContrast[] = [
	{
		theme: 'light',
		pairs: [
			pair('light', 'color.text', 'color.background', 17.85, 4.5),
			pair('light', 'color.border', 'color.background', 3.3, 3)
		]
	},
	{
		theme: 'dark',
		pairs: [
			pair('dark', 'color.text', 'color.background', 16.29, 4.5),
			pair('dark', 'color.border', 'color.background', 1.74, 3)
		]
	}
];

describe('TokenSheet (S-0382)', () => {
	let c: ReturnType<typeof mount> | undefined;
	afterEach(() => {
		if (c) unmount(c);
		c = undefined;
		document.body.innerHTML = '';
	});

	function render(themes = [LIGHT, DARK], contrast = CONTRAST) {
		c = mount(TokenSheet, { target: document.body, props: { themes, contrast } });
		flushSync();
	}
	const q = <T extends Element = HTMLElement>(sel: string) => document.querySelector<T>(sel)!;
	const cell = (row: string, theme: string) =>
		q(`[data-testid="${row}"] [data-sheet-theme="${theme}"]`);
	const words = (el: Element) => el.textContent!.replace(/\s+/g, ' ').trim();

	it("shows each colour as a swatch in both themes, with its contrast on that theme's ground", () => {
		render();
		const light = cell('swatch-color.text', 'light');
		const dark = cell('swatch-color.text', 'dark');
		expect(
			light.querySelector<HTMLElement>('[data-testid="swatch"]')!.style.backgroundColor
		).not.toBe('');
		expect(light.getAttribute('data-sheet-theme')).toBe('light');
		expect(words(light)).toContain('#0f172a');
		expect(words(light.querySelector('[data-testid="contrast"]')!)).toBe('17.85:1 ≥ 4.5 text');
		expect(words(dark)).toContain('#f1f5f9');
		expect(words(dark.querySelector('[data-testid="contrast"]')!)).toBe('16.29:1 ≥ 4.5 text');
		// The border reaches 3 but not 4.5 in light, and neither in dark.
		expect(
			words(cell('swatch-color.border', 'light').querySelector('[data-testid="contrast"]')!)
		).toBe('3.30:1 ≥ 3 large text, UI');
		expect(
			words(cell('swatch-color.border', 'dark').querySelector('[data-testid="contrast"]')!)
		).toMatch(/^1\.\d\d:1 under 3$/);
		// The ground itself is named rather than measured against itself.
		expect(
			words(cell('swatch-color.background', 'dark').querySelector('[data-testid="contrast"]')!)
		).toBe('the ground');
	});

	it('names the sRGB colour an oklch value is measured as, and says when it was clipped', () => {
		render();
		const neonCell = cell('swatch-color.brand-neon', 'light');
		expect(words(neonCell)).toContain('oklch(0.7 0.4 150)');
		expect(words(neonCell)).toMatch(/#[0-9a-f]{6}/);
		expect(neonCell.querySelector('[data-testid="clipped"]')).not.toBeNull();
		expect(cell('swatch-color.text', 'light').querySelector('[data-testid="clipped"]')).toBeNull();
	});

	it("says when a theme has no ground to measure against, and still shows that theme's swatches", () => {
		render([
			LIGHT,
			{ name: 'dark', tokens: DARK.tokens.filter((t) => t.path !== 'color.background') }
		]);
		expect(document.body.textContent).toContain('Dark has no color.background');
		const dark = cell('swatch-color.text', 'dark');
		expect(words(dark)).toContain('#f1f5f9');
		expect(dark.querySelector('[data-testid="contrast"]')).toBeNull();
	});

	it("lists flai's role pairs with a pass or fail in each theme", () => {
		render();
		const verdict = (row: string, theme: string) =>
			words(cell(row, theme).querySelector('[data-testid="verdict"]')!);
		const border = 'pair-color.border-on-color.background';
		expect(words(cell(border, 'light'))).toContain('3.30:1');
		expect(verdict(border, 'light')).toBe('passes 3');
		expect(verdict(border, 'dark')).toBe('fails 3');
		expect(verdict('pair-color.text-on-color.background', 'dark')).toBe('passes 4.5');
	});

	it('sets each type step as sample text at its size and line height, in the sans stack', () => {
		render();
		const samples = document.querySelectorAll<HTMLElement>(
			'[data-testid="type-text.lg"] [data-testid="type-sample"]'
		);
		// The same in both themes: one sample, not two.
		expect(samples.length).toBe(1);
		expect(samples[0].style.fontSize).toBe('1.125rem');
		expect(samples[0].style.lineHeight).toBe('1.556');
		expect(samples[0].style.fontFamily).toContain('ui-sans-serif');
		expect(q('[data-testid="type-text.lg"]').textContent).toContain('Light and dark');
		expect(
			q<HTMLElement>('[data-testid="type-text.2xl"] [data-testid="type-sample"]').style.fontSize
		).toBe('1.5rem');
		expect(
			q<HTMLElement>('[data-testid="font-font.sans"] [data-testid="font-sample"]').style.fontFamily
		).toContain('system-ui');
	});

	it('draws spacing, radii, and shadows to size, light and dark side by side where they differ', () => {
		render();
		const bars = [
			...document.querySelectorAll<HTMLElement>(
				'[data-testid="spacing-light"] [data-testid="spacing-bar"]'
			)
		];
		expect(bars.length).toBe(7);
		// Four units wide: the style holds the product, which the DOM may simplify.
		expect(bars[2].style.width).toMatch(/^calc\((0\.25rem \* 4|1rem)\)$/);
		expect(bars[6].style.width).toMatch(/^calc\((0\.25rem \* 16|4rem)\)$/);
		expect(q('[data-testid="spacing-light"]').textContent).toContain('4 · 1rem');
		// The spacing unit is the same in both themes: one set of bars.
		expect(document.querySelector('[data-testid="spacing-dark"]')).toBeNull();

		const radii = [
			...document.querySelectorAll<HTMLElement>(
				'[data-testid="radius-radius.md"] [data-testid="radius-specimen"]'
			)
		];
		expect(radii.map((r) => r.style.borderRadius)).toEqual(['0.375rem', '0.5rem']);

		const shadows = [
			...document.querySelectorAll<HTMLElement>(
				'[data-testid="shadow-shadow.md"] [data-testid="shadow-specimen"]'
			)
		];
		expect(shadows.length).toBe(2);
		expect(shadows[0].getAttribute('style')).toContain('#0000001a');
		expect(shadows[1].getAttribute('style')).toContain('#00000066');
	});

	it('writes no value into a style that could not go into a stylesheet', () => {
		const bad = { ...colour('brand', '#000000'), css: 'red; background: url(https://x.test/a)' };
		render([{ name: 'light', tokens: [colour('background', '#ffffff'), bad] }], []);
		const swatch = cell('swatch-color.brand', 'light').querySelector('[data-testid="swatch"]')!;
		expect(swatch.getAttribute('style') ?? '').not.toContain('url(');
		expect(cell('swatch-color.brand', 'dark').textContent).toContain('not in this theme');
	});
});
