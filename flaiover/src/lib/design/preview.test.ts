import { describe, expect, it } from 'vitest';
import {
	contrast,
	contrastOnGround,
	dsDeclarations,
	luminance,
	previewStyle,
	toSRGB,
	unsafeVariable,
	type PreviewToken
} from './preview';

// A colour token as design.get answers it.
function colour(name: string, hex: string): PreviewToken {
	const n = parseInt(hex.slice(1), 16);
	const components = [n >> 16, (n >> 8) & 255, n & 255].map(
		(b) => Math.round((b / 255) * 1000) / 1000
	);
	return {
		path: `color.${name}`,
		type: 'color',
		value: { colorSpace: 'srgb', components, hex },
		css: hex,
		variables: [{ name: `--color-${name}`, value: hex }]
	};
}

// The colour roles of flai's starter tokens, as `flai design show --json` answered them for S-0382.
const STARTER: Record<'light' | 'dark', Record<string, string>> = {
	light: {
		background: '#ffffff',
		surface: '#f8fafc',
		text: '#0f172a',
		'text-muted': '#475569',
		border: '#848e9f',
		primary: '#2563eb',
		'on-primary': '#ffffff',
		secondary: '#475569',
		'on-secondary': '#ffffff',
		focus: '#2563eb',
		success: '#15803d',
		'success-soft': '#f0fdf4',
		warning: '#b45309',
		'warning-soft': '#fffbeb',
		danger: '#b91c1c',
		'danger-soft': '#fef2f2',
		info: '#1d4ed8',
		'info-soft': '#eff6ff'
	},
	dark: {
		background: '#0f172a',
		surface: '#1e293b',
		text: '#f1f5f9',
		'text-muted': '#94a3b8',
		border: '#64748b',
		primary: '#60a5fa',
		'on-primary': '#0f172a',
		secondary: '#cbd5e1',
		'on-secondary': '#0f172a',
		focus: '#60a5fa',
		success: '#4ade80',
		'success-soft': '#052e16',
		warning: '#fbbf24',
		'warning-soft': '#451a03',
		danger: '#f87171',
		'danger-soft': '#450a0a',
		info: '#60a5fa',
		'info-soft': '#172554'
	}
};

// The ratios flai reported for the same tokens' role pairs: foreground, background, ratio.
const FLAI_PAIRS: Record<'light' | 'dark', [string, string, number][]> = {
	light: [
		['text', 'background', 17.85],
		['text', 'surface', 17.06],
		['text-muted', 'background', 7.57],
		['text-muted', 'surface', 7.24],
		['success', 'background', 5.01],
		['success', 'surface', 4.79],
		['success', 'success-soft', 4.79],
		['warning', 'background', 5.02],
		['warning', 'surface', 4.79],
		['warning', 'warning-soft', 4.84],
		['danger', 'background', 6.47],
		['danger', 'surface', 6.18],
		['danger', 'danger-soft', 5.91],
		['info', 'background', 6.7],
		['info', 'surface', 6.4],
		['info', 'info-soft', 6.15],
		['on-primary', 'primary', 5.16],
		['on-secondary', 'secondary', 7.57],
		['border', 'background', 3.3],
		['focus', 'background', 5.16]
	],
	dark: [
		['text', 'background', 16.29],
		['text', 'surface', 13.35],
		['text-muted', 'background', 6.96],
		['text-muted', 'surface', 5.7],
		['success', 'background', 10.24],
		['success', 'surface', 8.39],
		['success', 'success-soft', 8.55],
		['warning', 'background', 10.69],
		['warning', 'surface', 8.76],
		['warning', 'warning-soft', 8.97],
		['danger', 'background', 6.45],
		['danger', 'surface', 5.28],
		['danger', 'danger-soft', 5.83],
		['info', 'background', 7.02],
		['info', 'surface', 5.75],
		['info', 'info-soft', 5.77],
		['on-primary', 'primary', 7.02],
		['on-secondary', 'secondary', 12.02],
		['border', 'background', 3.75],
		['focus', 'background', 7.02]
	]
};

const tokensOf = (theme: 'light' | 'dark') =>
	Object.entries(STARTER[theme]).map(([name, hex]) => colour(name, hex));

const oklch = (l: number, c: number, h: number | 'none') => ({
	colorSpace: 'oklch',
	components: [l, c, h]
});

describe('dsDeclarations and previewStyle (S-0382, ADR-0149)', () => {
	const typography: PreviewToken = {
		path: 'text.base',
		type: 'typography',
		value: {},
		css: '1rem',
		variables: [
			{ name: '--text-base', value: '1rem' },
			{ name: '--text-base--line-height', value: '1.5' }
		]
	};

	it("renames each token's variables to --ds-*, every variable of a token", () => {
		const d = dsDeclarations('light', [colour('primary', '#2563eb'), typography]);
		expect(d.declarations).toEqual([
			'--ds-color-primary: #2563eb;',
			'--ds-text-base: 1rem;',
			'--ds-text-base--line-height: 1.5;'
		]);
		expect(d.refused).toEqual([]);
	});

	it('keeps a font stack with quoted families', () => {
		const font: PreviewToken = {
			path: 'font.sans',
			type: 'fontFamily',
			value: [],
			css: '',
			variables: [{ name: '--font-sans', value: 'ui-sans-serif, "Segoe UI", sans-serif' }]
		};
		expect(dsDeclarations('light', [font]).declarations).toEqual([
			'--ds-font-sans: ui-sans-serif, "Segoe UI", sans-serif;'
		]);
	});

	it.each([
		['url(https://example.com/x.png)', 'url('],
		['URL(x)', 'url('],
		['red}body{color:red', 'a brace'],
		['red{', 'a brace'],
		['red; color: blue', 'a semicolon'],
		['</style><script>', '<'],
		['@import x', '@'],
		['\\72 ed', 'a backslash'],
		['red /* the rest', 'a comment']
	])('refuses, and reports, a value carrying %s', (value, what) => {
		const token: PreviewToken = {
			path: 'color.bad',
			type: 'color',
			value: {},
			css: value,
			variables: [{ name: '--color-bad', value }]
		};
		const d = dsDeclarations('dark', [colour('primary', '#2563eb'), token]);
		expect(d.declarations).toEqual(['--ds-color-primary: #2563eb;']);
		expect(d.refused).toEqual([
			{ theme: 'dark', name: '--color-bad', value, reason: `its value contains ${what}` }
		]);
	});

	it('refuses a name that is not a custom property', () => {
		expect(unsafeVariable('color: red; --x', 'red')).toMatch(/is not a custom property/);
		expect(unsafeVariable('--Color', 'red')).toMatch(/is not a custom property/);
		expect(unsafeVariable('--text-xs--line-height', '1.333')).toBeNull();
	});

	it("sets each theme's values on its panel, with its colour scheme, and gathers what it refused", () => {
		const bad: PreviewToken = {
			...colour('bad', '#000000'),
			variables: [{ name: '--color-bad', value: 'url(x)' }]
		};
		const { css, refused } = previewStyle({
			light: [colour('primary', '#2563eb')],
			dark: [colour('primary', '#60a5fa'), bad]
		});
		expect(css).toBe(
			"[data-ds-theme='light'] {\n\tcolor-scheme: light;\n\t--ds-color-primary: #2563eb;\n}\n\n" +
				"[data-ds-theme='dark'] {\n\tcolor-scheme: dark;\n\t--ds-color-primary: #60a5fa;\n}\n"
		);
		expect(refused).toEqual([
			{ theme: 'dark', name: '--color-bad', value: 'url(x)', reason: 'its value contains url(' }
		]);
	});
});

describe('toSRGB, as flai converts a colour (flai/internal/designsystem/color_test.go)', () => {
	it.each([
		[
			'srgb by its hex',
			{ colorSpace: 'srgb', components: [0, 0, 0], hex: '#F5F5F7' },
			'#f5f5f7',
			false
		],
		[
			'srgb by its components',
			{ colorSpace: 'srgb', components: [0.96, 0.96, 0.97] },
			'#f5f5f7',
			false
		],
		["CSS Color 4's oklch for red", oklch(0.628, 0.2577, 29.23), '#ff0000', false],
		['oklch white', oklch(1, 0, 0), '#ffffff', false],
		['oklch black', oklch(0, 0, 0), '#000000', false],
		['an achromatic oklch with no hue', oklch(0.5, 0, 'none'), '#636363', false],
		['an oklch green inside sRGB', oklch(0.5, 0.13, 150), '#137738', false],
		['an oklch blue just outside sRGB', oklch(0.52, 0.18, 255), '#0066cd', true],
		['an oklch green far outside sRGB', oklch(0.9, 0.4, 145), '#00ff00', true]
	])('%s', (_, value, hex, clipped) => {
		const c = toSRGB(value)!;
		expect(c.hex).toBe(hex);
		expect(c.clipped).toBe(clipped);
		for (const v of c.rgb) expect(v >= 0 && v <= 1).toBe(true);
	});

	it.each([
		['not an object', '#ffffff'],
		['another colour space', { colorSpace: 'display-p3', components: [1, 0, 0] }],
		['srgb with neither hex nor components', { colorSpace: 'srgb' }],
		['oklch without three numbers', { colorSpace: 'oklch', components: [0.5, 0.1] }]
	])('is null for %s', (_, value) => {
		expect(toSRGB(value)).toBeNull();
	});
});

describe('contrast, as flai reports it', () => {
	it('gives black on white 21 and white on white 1', () => {
		expect(luminance([1, 1, 1])).toBeCloseTo(1, 12);
		expect(luminance([0, 0, 0])).toBe(0);
		const white = { colorSpace: 'srgb', hex: '#ffffff' };
		const black = { colorSpace: 'srgb', hex: '#000000' };
		expect(contrast(black, white)!.ratio).toBe(21);
		expect(contrast(white, black)!.ratio).toBe(21);
		expect(contrast(white, white)!.ratio).toBe(1);
	});

	for (const theme of ['light', 'dark'] as const) {
		it.each(FLAI_PAIRS[theme])(`${theme}: %s on %s is %s, as flai reports it`, (fg, bg, ratio) => {
			const byName = Object.fromEntries(tokensOf(theme).map((t) => [t.path, t.value]));
			const c = contrast(byName[`color.${fg}`], byName[`color.${bg}`])!;
			expect(c.ratio).toBe(ratio);
			expect(c.foregroundHex).toBe(STARTER[theme][fg]);
			expect(c.backgroundHex).toBe(STARTER[theme][bg]);
			expect(c.clipped).toBe(false);
		});
	}

	it('truncates the ratio, never rounding a failing one up to pass, and keeps the exact one', () => {
		// #767676 on white is 4.54:1, just over 4.5; #777777 is 4.47:1, under it.
		const white = { colorSpace: 'srgb', hex: '#ffffff' };
		const over = contrast({ colorSpace: 'srgb', hex: '#767676' }, white)!;
		const under = contrast({ colorSpace: 'srgb', hex: '#777777' }, white)!;
		expect(over.ratio).toBe(4.54);
		expect(under.ratio).toBe(4.47);
		expect(under.exact).toBeGreaterThan(4.47);
		expect(under.exact).toBeLessThan(4.48);
	});

	it('says when either colour was clipped to sRGB', () => {
		const c = contrast(oklch(0.52, 0.18, 255), { colorSpace: 'srgb', hex: '#ffffff' })!;
		expect(c.clipped).toBe(true);
		expect(c.foregroundHex).toBe('#0066cd');
	});

	it("gives each colour token its contrast on the theme's ground, and leaves other tokens out", () => {
		const font: PreviewToken = {
			path: 'font.sans',
			type: 'fontFamily',
			value: ['system-ui'],
			css: 'system-ui',
			variables: [{ name: '--font-sans', value: 'system-ui' }]
		};
		const extra: PreviewToken = {
			path: 'color.brand-sun',
			type: 'color',
			value: oklch(0.82, 0.17, 85),
			css: '',
			variables: []
		};
		const rows = contrastOnGround([...tokensOf('light'), font, extra]);
		expect(rows.map((r) => r.path)).toEqual([
			...Object.keys(STARTER.light).map((n) => `color.${n}`),
			'color.brand-sun'
		]);
		expect(rows.find((r) => r.path === 'color.text')!.ratio).toBe(17.85);
		expect(rows.find((r) => r.path === 'color.background')!.ratio).toBe(1);
		const dark = contrastOnGround(tokensOf('dark'));
		expect(dark.find((r) => r.path === 'color.text')!.ratio).toBe(16.29);
		expect(dark.find((r) => r.path === 'color.border')!.ratio).toBe(3.75);
	});

	it('is empty for a theme with no ground', () => {
		expect(contrastOnGround([colour('text', '#000000')])).toEqual([]);
	});
});
